package finder

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"

	"github.com/sahakkhotsanyan/adfind/pkg/config"
	"github.com/sahakkhotsanyan/adfind/pkg/fast"
)

type Finder interface {
	Find(uri, websiteType, wordlist string) ([]string, error)
}

type finder struct {
	fast fast.Client
	cfg  *config.Config
	stop bool
	base string
}

func NewFinder(fast fast.Client, cfg *config.Config, stop bool, basePath string) Finder {
	return &finder{
		fast: fast,
		cfg:  cfg,
		stop: stop,
		base: basePath,
	}
}

func (f *finder) Find(uri, websiteType, wordlist string) ([]string, error) {
	_, ok := f.cfg.WordLists[websiteType]
	if !ok && websiteType != "all" {
		return nil, errors.New("website type not found")
	}
	if websiteType != "all" {
		return f.processURI(uri, websiteType)
	}
	if wordlist != "native" {
		return f.processURIWithWordlist(uri, wordlist)
	}
	found := make([]string, 0)
	for wsT := range f.cfg.WordLists {
		tmp, err := f.processURI(uri, wsT)
		if err != nil {
			log.Errorf("error processing uri: %s with type %s", uri, wsT)
			continue
		}
		found = append(found, tmp...)
	}

	return found, nil
}

func (f *finder) processURIWithWordlist(uri, wordlist string) ([]string, error) {
	var err error
	var ok bool
	file, err := os.Open(wordlist)
	if err != nil {
		return nil, err
	}

	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			fmt.Println("Error closing file:", err)
		}
	}(file)

	scanner := bufio.NewScanner(file)

	found := make([]string, 0)
	uri = strings.TrimRight(uri, "/")
	for scanner.Scan() {

		tmpUri := uri + "/" + strings.TrimSpace(scanner.Text())
		ok, err = f.checkURI(tmpUri)
		if err != nil {
			log.Errorf("error checking uri: %s, err: %v", tmpUri, err)
			continue
		}
		if !ok {
			log.Debugf("uri not found: %s", tmpUri)
			continue
		}

		log.Infof("uri found: %s", tmpUri)
		found = append(found, tmpUri)
	}

	// Check for scanner errors
	if err = scanner.Err(); err != nil {
		log.Errorf("Error reading file: %v", err)
	}

	return found, err
}

func (f *finder) processURI(uri, websiteType string) ([]string, error) {
	fName := f.cfg.GetWordListFileName(websiteType)
	return f.processURIWithWordlist(uri, f.base+fName)
}

func (f *finder) checkURI(uri string) (bool, error) {
	statusCode, err := f.fast.CheckURL(uri)
	if err != nil {
		return false, err
	}

	if statusCode != fasthttp.StatusOK && statusCode >= fasthttp.StatusBadRequest {
		return false, nil
	}

	if f.stop {
		var yn string
		log.Infof("uri found: %s with status code %d", uri, statusCode)
		log.Infof("Do you want to continue? [y/n]")
		_, err = fmt.Scan(&yn)
		if err != nil {
			return false, err
		}
		if yn == "n" {
			os.Exit(0)
		}
	}

	return true, nil
}
