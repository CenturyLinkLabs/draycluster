package provider

import (
	"fmt"
	"github.com/CenturylinkLabs/draycluster/deploy"
	"strings"
    "github.com/CenturylinkLabs/draycluster/utils")

type CloudProvider interface {
	ProvisionAgent() (deploy.CloudServer, error)
}

func New(providerType string) CloudProvider {
	pt := strings.TrimSpace(strings.ToLower(providerType))
    utils.LogInfo(fmt.Sprintf("\n\nProvider:%s",pt))
	switch string(pt) {
	case "centurylink":
		return NewCenturylink()
	case "amazon":
		return NewAmazon()
    case "digitalocean":
  		return NewDigitalOcean()
	}
	fmt.Printf("\nNIL Provider:%s", pt)
	return nil
}
