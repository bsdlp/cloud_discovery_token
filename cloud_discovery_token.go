package main

import (
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
	file, err := ioutil.ReadFile("/tmp/cloud-config.yaml")
	if err != nil {
		log.Fatalln(err)
	}

	if config.IsCloudConfig(string(file)) {
		CloudConfig, err := config.NewCloudConfig(string(file))
		if err != nil {
			log.Fatalln(err)
		}
		if len(CloudConfig.Coreos.Etcd.Discovery) == 0 {
			token := GetToken("https://discovery.etcd.io")
			fmt.Printf("discovery: %s\n", token)
			CloudConfig.Coreos.Etcd.Discovery = token
		} else {
			fmt.Printf("discovery: %s\n", CloudConfig.Coreos.Etcd.Discovery)
		}
	}
}
