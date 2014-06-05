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

	// builds the full map
	var mymap *purple.ConfigMap
	mymap, err = purple.BuildConfigMap(url, 25)
	fmt.Println((*mymap)["mail.telenor.dk"].IncomingServers[0].Hostname)

}
