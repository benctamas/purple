package main

import (
	"fmt"
	"github.com/benctamas/purple/purple"
)

func main() {
	url := "https://autoconfig.thunderbird.net/v1.1/"

	// fetches the list of config files (as urls)
	urls, err := purple.FetchConfigFileList(url)
	if err != nil {
		panic(err)
	}

	for _, url := range urls {
		fmt.Println(*url)
	}

	// fetches and parses one config file
	conf, err := purple.DownloadAndParseConfig(*urls[0])
	if err != nil {
		panic(err)
	}
	fmt.Println(conf)

	var mymap *purple.ConfigMap

	// builds the full map
	// let's use 25 concurrent workers for donwnloading/parsing
	mymap, err = purple.BuildConfigMap(url, 25)
	if err != nil {
		panic(err)
	}

	// full config datablock
	fmt.Println((*mymap)["mail.telenor.dk"].ClientConfig)

	// source url of datablock
	fmt.Println((*mymap)["mail.telenor.dk"].ClientConfig.SourceUrl)

	// EmailProvider datablock containing domain
	fmt.Println((*mymap)["mail.telenor.dk"].EmailProvider)

}
