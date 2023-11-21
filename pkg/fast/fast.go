package fast

import (
	"fmt"
	"time"

	"github.com/valyala/fasthttp"
)

type Client interface {
	Get(uri string, resp *fasthttp.Response) error
	CheckURL(uri string) (int, error)
	Post(uri string, contentType string, body []byte, resp *fasthttp.Response) error
}

type client struct {
	client      *fasthttp.Client
	readTimeout time.Duration
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
	err := c.client.Do(req, resp)
	fasthttp.ReleaseRequest(req)
	if err != nil {
		return fmt.Errorf("failed to get %s: %w", uri, err)
	}
	return nil
}

// CheckURL checks if url is reachable
func (c *client) CheckURL(uri string) (int, error) {
	var body []byte
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(uri)
	req.Header.SetMethod(fasthttp.MethodGet)
	req.Header.SetContentTypeBytes(contentTypeJSON)
	req.Header.Set("Accept", "application/json")
	code, _, err := c.client.Get(body, uri)
	if err != nil {
		return 0, err
	}
	fasthttp.ReleaseRequest(req)
	if err != nil {
		return 0, fmt.Errorf("failed to get %s: %w", uri, err)
	}
	return code, nil
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
	req.SetBodyRaw(body)
	err := c.client.DoTimeout(req, resp, reqTimeout)
	fasthttp.ReleaseRequest(req)

	if err != nil {
		return fmt.Errorf("failed to post %s: %w", uri, err)
	}
	return nil
}
