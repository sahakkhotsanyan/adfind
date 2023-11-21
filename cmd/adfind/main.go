package main

import (
	"flag"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/sahakkhotsanyan/adfind/pkg/config"
	"github.com/sahakkhotsanyan/adfind/pkg/fast"
	"github.com/sahakkhotsanyan/adfind/pkg/finder"
)

var verbose, help, stop bool
var url, adminType, basePath string
var timeout int64

const (
	base       = "/usr/share/adfind/"
	configFile = "config.yaml"
)

func main() {
	initInfo()
	flag.BoolVar(&verbose, "v", false, "verbose mode")
	flag.BoolVar(&stop, "s", false, "stop when admin panel was found")
	flag.Int64Var(&timeout, "to", 1000, "timeout for request in milliseconds")
	flag.StringVar(&url, "u", "", "URL of site {example: adfind -u https://example.com -t php}")
	flag.StringVar(&basePath, "b", base, "base path of config files (default is /usr/share/adfind/)")
	flag.StringVar(&adminType, "t", "all", "type of admin panel (default is all) {types: php , asp, aspx, js, cfm, cgi, brf. example:adfind -u http://example.com -t php}")
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

	cfg, err := config.NewConfig(basePath + configFile)
	if err != nil {
		log.Fatal(err)
	}

	timeOut := time.Duration(timeout) * time.Millisecond

	fastClient := fast.NewClient(timeOut, timeOut, timeOut)

	findProcessor := finder.NewFinder(fastClient, cfg, stop, basePath)

	found, err := findProcessor.Find(url, adminType)
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
