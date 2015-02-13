package provider

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func createCloudConfigCluster(pubkey string) string {
	println("Create Cloud Config Cluster")
	response, _ := http.Get("https://discovery.etcd.io/new")
	defer response.Body.Close()
	contents, _ := ioutil.ReadAll(response.Body)
	cloudConfig, _ := ioutil.ReadFile("cloud-config-init.yaml")
	discoveryUrl := fmt.Sprintf("discovery: %s", string(contents))
	cloudConfigNew := strings.Replace(string(cloudConfig), "discovery_url", discoveryUrl, -1)
    cloudConfigNew = strings.Replace(string(cloudConfigNew), "ssh-rsa", pubkey, -1)
	return string(cloudConfigNew)
}

func createCloudConfigAgent(pubkey string) string {
	println("Create Cloud Config Agent")
	cloudConfig, _ := ioutil.ReadFile("cloud-config-agent.yaml")
	cloudConfigNew := strings.Replace(string(cloudConfig), "ssh-rsa", pubkey, -1)
	return string(cloudConfigNew)
}
