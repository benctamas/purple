package purple

import (
	"code.google.com/p/go.net/html"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

// extracts urls from <a href> links (inner lambda called recursively)
func parseLinks(url string, node *html.Node) (urls []*string) {
	var process func(node *html.Node)
	process = func(node *html.Node) {
		var target string
		if node.Type == html.ElementNode && node.Data == "a" {
			for _, attr := range node.Attr {
				if attr.Key == "href" {
					if strings.Contains(attr.Val, ".") {
						target = url + attr.Val
						urls = append(urls, &target)
					}
				}
			}
		}
		for cur := node.FirstChild; cur != nil; cur = cur.NextSibling {
			process(cur)
		}
	}
	process(node)
	return urls
}

// extract list of links pointing to config files from index page
func FetchConfigFileList(index_url string) ([]*string, error) {
	response, err := http.Get(index_url)
	defer response.Body.Close()

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		err := errors.New(http.StatusText(response.StatusCode))
		fmt.Println(err)
		return nil, err
	}

	doc, err := html.Parse(response.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	urls := parseLinks(index_url, doc)
	return urls, nil
}

// download & parse config file at url
func DownloadAndParseConfig(url string) (*ClientConfig, error) {
	response, err := http.Get(url)
	defer response.Body.Close()

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		err := errors.New(http.StatusText(response.StatusCode))
		fmt.Println(err)
		return nil, err
	}

	raw, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	conf := ClientConfig{}
	err = xml.Unmarshal(raw, &conf)

	if err != nil {
		fmt.Println("error parsing file:", err)
		return nil, err
	}

	conf.SourceUrl = url
	return &conf, nil
}

func worker(chUrls chan string, chData chan *ClientConfig, wg *sync.WaitGroup) {
	defer wg.Done()
	// fmt.Println("worker start")

	for url := range chUrls {
		// fmt.Println("fetching %s", url)
		conf, err := DownloadAndParseConfig(url)
		if err != nil {
			fmt.Println("cannot get config at %s", url)
		} else {
			chData <- conf
		}
		// fmt.Println("done %s", url)
	}
	// fmt.Println("worker done")
}

// builds map of config files listed at index url, concurrently
func BuildConfigMap(index_url string, worker_cnt int) (*ConfigMap, error) {
	urls, err := FetchConfigFileList(index_url)
	if err != nil {
		fmt.Println("cannot get config list")
		return nil, err
	}

	chUrls := make(chan string)
	chData := make(chan *ClientConfig)
	wg := new(sync.WaitGroup)
	confmap := ConfigMap{}

	// start map builder
	go func(confmap *ConfigMap, chData chan *ClientConfig) {
		for conf := range chData {
			for _, ep := range conf.EmailProviders {
				for _, domain := range ep.Domains {
					domconf := DomainConfig{}
					domconf.Domain = domain
					domconf.EmailProvider = &ep
					domconf.ClientConfig = conf
					(*confmap)[domain] = &domconf
				}
			}
		}
	}(&confmap, chData)

	// create workers
	for i := 0; i < worker_cnt; i = i + 1 {
		wg.Add(1)
		go worker(chUrls, chData, wg)
	}

	// push urls to workers
	for _, url := range urls {
		chUrls <- *url
	}

	// wait till workers finish processing urls
	close(chUrls)
	wg.Wait()

	// stop map builder
	close(chData)

	return &confmap, nil
}
