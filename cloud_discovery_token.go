package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/coreos/coreos-cloudinit/config"
)

func GetToken(BaseUrl string) string {
	NewTokenURL := fmt.Sprintf("%s/new", BaseUrl)
	resp, err := http.Get(NewTokenURL)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	token, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return string(token)
}

func main() {
	FilePath := flag.String("config", "./cloud-config.yaml", "Path to cloud-config yaml file")
	BaseUrl := flag.String("url", "https://discovery.etcd.io", "URL to cluster discovery service")
	flag.Parse()

	file, err := ioutil.ReadFile(*FilePath)
	if err != nil {
		log.Fatalln(err)
	}

	if config.IsCloudConfig(string(file)) {
		CloudConfig, err := config.NewCloudConfig(string(file))
		if err != nil {
			log.Fatalln(err)
		}
		if len(CloudConfig.Coreos.Etcd.Discovery) == 0 {
			token := GetToken(*BaseUrl)
			fmt.Printf("discovery: %s\n", token)
			CloudConfig.Coreos.Etcd.Discovery = token
		} else {
			fmt.Printf("discovery: %s\n", CloudConfig.Coreos.Etcd.Discovery)
		}
	}
}
