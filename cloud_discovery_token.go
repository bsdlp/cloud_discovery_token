package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/coreos/coreos-cloudinit/config"
)

func LogError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func GetToken(BaseUrl string) (string, error) {
	NewTokenURL := fmt.Sprintf("%s/new", BaseUrl)
	resp, err := http.Get(NewTokenURL)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("%s: %s", NewTokenURL, resp.Status)
	}
	defer resp.Body.Close()

	token, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(token), nil
}

func WriteCloudConfig(FilePath *string, cfg config.CloudConfig) error {
	WriteFile, err := os.Create(*FilePath)
	if err != nil {
		return err
	}
	defer WriteFile.Close()

	_, err = WriteFile.WriteString(cfg.String())
	return err
}

func main() {
	FilePath := flag.String("config", "./cloud-config.yaml", "Path to cloud-config yaml file")
	BaseUrl := flag.String("url", "https://discovery.etcd.io", "URL to cluster discovery service")
	Overwrite := flag.Bool("overwrite", false, "Overwrite config with new token")
	flag.Parse()

	file, err := ioutil.ReadFile(*FilePath)
	LogError(err)

	if config.IsCloudConfig(string(file)) {
		CloudConfig, err := config.NewCloudConfig(string(file))
		LogError(err)

		if len(CloudConfig.Coreos.Etcd.Discovery) == 0 {
			token, err := GetToken(*BaseUrl)
			LogError(err)

			fmt.Printf("discovery: %s\n", token)
			CloudConfig.Coreos.Etcd.Discovery = token

			if *Overwrite == true {
				err = WriteCloudConfig(FilePath, *CloudConfig)
				LogError(err)
			}
		} else {
			fmt.Printf("discovery: %s\n", CloudConfig.Coreos.Etcd.Discovery)
		}
	}
}
