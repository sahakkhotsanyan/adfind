package main

import (
	"errors"
	"flag"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/sahakkhotsanyan/adfind/pkg/config"
	"github.com/sahakkhotsanyan/adfind/pkg/fast"
	"github.com/sahakkhotsanyan/adfind/pkg/finder"
)

var verbose, veryVerbose, help, stop bool
var url, adminType, basePath, wordlist string
var timeout int64

const (
	base       = "/usr/share/adfind/"
	configFile = "config.yaml"
)

// headerList implements flag.Value so we can pass -H multiple times like curl
type headerList map[string]string

func (h *headerList) String() string {
	var sb strings.Builder
	for k, v := range *h {
		sb.WriteString(k)
		sb.WriteString(": ")
		sb.WriteString(v)
		sb.WriteString("\n")
	}
	return sb.String()
}
func (h *headerList) Set(value string) error {
	parts := strings.SplitN(value, ":", 2)
	if len(parts) != 2 {
		return errors.New("invalid header format, use Key: Value")
	}
	key := strings.TrimSpace(parts[0])
	val := strings.TrimSpace(parts[1])
	if *h == nil {
		*h = make(map[string]string)
	}
	(*h)[key] = val
	return nil
}

var headers headerList

func main() {
	initInfo()
	flag.BoolVar(&verbose, "v", false, "verbose mode")
	flag.BoolVar(&veryVerbose, "vv", false, "very verbose verbose mode")
	flag.StringVar(&wordlist, "w", "native", "wordlist for admin panel")
	flag.BoolVar(&stop, "s", false, "stop when admin panel was found")
	flag.Int64Var(&timeout, "to", 1000, "timeout for request in milliseconds")
	flag.StringVar(&url, "u", "", "URL of site {example: adfind -u https://example.com -t php}")
	flag.StringVar(&basePath, "b", base, "base path of config files (default is /usr/share/adfind/)")
	flag.StringVar(&adminType, "t", "all", "type of admin panel (default is all) {types: php , asp, aspx, js, cfm, cgi, brf. example:adfind -u http://example.com -t php}")
	flag.Var(&headers, "H", "custom header to add to request, can be used multiple times")
	flag.BoolVar(&help, "h", false, "show this help")

	flag.Parse()

	if url == "" || help {
		flag.Usage()
		return
	}

	log.SetLevel(log.InfoLevel)

	if verbose {
		log.SetLevel(log.DebugLevel)
	}

	if veryVerbose {
		log.SetLevel(log.TraceLevel)
	}

	cfg, err := config.NewConfig(basePath + configFile)
	if err != nil {
		log.Fatal(err)
	}

	timeOut := time.Duration(timeout) * time.Millisecond

	fastClient := fast.NewClient(timeOut, timeOut, timeOut)

	findProcessor := finder.NewFinder(fastClient, cfg, stop, basePath)

	found, err := findProcessor.Find(url, adminType, wordlist)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(found)
}

func initInfo() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	log.Println("Starting adfind...")
	log.Println("by @sahakkhotsanyan (c) 2023")
}
