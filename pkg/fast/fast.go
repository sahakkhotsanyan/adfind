package fast

import (
	"fmt"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

type Client interface {
	Get(uri string, resp *fasthttp.Response) error
	CheckURL(uri string) (int, []byte, error)
	Post(uri string, contentType string, body []byte, resp *fasthttp.Response) error
	SetCustomHeaders(headers map[string]string)
}

type client struct {
	client        *fasthttp.Client
	readTimeout   time.Duration
	customHeaders map[string]string
	headersMutex  sync.Mutex
}

var contentTypeJSON = []byte("application/json")

func NewClient(readTimeout time.Duration, writeTimeout time.Duration, idleTimeout time.Duration) Client {
	c := &fasthttp.Client{
		ReadTimeout:                   readTimeout,
		WriteTimeout:                  writeTimeout,
		MaxIdleConnDuration:           idleTimeout,
		NoDefaultUserAgentHeader:      true, // Don't send: User-Agent: fast
		DisableHeaderNamesNormalizing: true, // If you set the case on your headers correctly you can enable this
		DisablePathNormalizing:        true,
		// increase DNS cache time to an hour instead of default minute
		Dial: (&fasthttp.TCPDialer{
			Concurrency:      4096,
			DNSCacheDuration: time.Hour,
		}).Dial,
	}

	return &client{client: c, readTimeout: readTimeout}
}

// Get
// Don't forget to acquire and release response
func (c *client) Get(uri string, resp *fasthttp.Response) error {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(uri)
	req.Header.SetMethod(fasthttp.MethodGet)
	req.Header.SetContentTypeBytes(contentTypeJSON)
	req.Header.Set("Accept", "application/json")
	c.addCustomHeaders(req)
	err := c.client.DoRedirects(req, resp, 10)
	fasthttp.ReleaseRequest(req)
	if err != nil {
		return fmt.Errorf("failed to get %s: %w", uri, err)
	}
	return nil
}

// CheckURL checks if url is reachable
func (c *client) CheckURL(uri string) (int, []byte, error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(uri)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	req.Header.SetMethod(fasthttp.MethodGet)
	req.Header.SetContentTypeBytes(contentTypeJSON)
	c.addCustomHeaders(req)
	req.Header.Set("Accept", "application/json")
	err := c.client.DoRedirects(req, resp, 10)
	fasthttp.ReleaseRequest(req)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to get %s: %w", uri, err)
	}
	return resp.StatusCode(), resp.Body(), nil
}

// Post
// Don't forget to acquire and release response
func (c *client) Post(uri string, contentType string, body []byte, resp *fasthttp.Response) error {
	// per-request timeout
	reqTimeout := c.readTimeout
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(uri)
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentTypeBytes([]byte(contentType))
	c.addCustomHeaders(req)
	req.SetBodyRaw(body)
	err := c.client.DoTimeout(req, resp, reqTimeout)
	fasthttp.ReleaseRequest(req)

	if err != nil {
		return fmt.Errorf("failed to post %s: %w", uri, err)
	}
	return nil
}

// SetCustomHeaders sets custom headers to be added to each request
func (c *client) SetCustomHeaders(headers map[string]string) {
	c.headersMutex.Lock()
	defer c.headersMutex.Unlock()
	c.customHeaders = headers
}

func (c *client) addCustomHeaders(req *fasthttp.Request) {
	c.headersMutex.Lock()
	defer c.headersMutex.Unlock()
	for k, v := range c.customHeaders {
		req.Header.Set(k, v)
	}
}
