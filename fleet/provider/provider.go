package provider

import "strings"
import (
    "github.com/CenturylinkLabs/draycluster/deploy"
)

// CloudProvider is used to deploy kubernetes cluster on any of the supported
// cloud providers.
type CloudProvider interface {
    ProvisionCluster() ([]deploy.CloudServer, error)
}

// New is used to instantiate a CloudProvider to use to provision the cluster.
func New(providerType string) CloudProvider {
    providerType = strings.ToLower(providerType)
    switch providerType {
        case "amazon":
        return NewAmazon()
        case "digitalocean":
        return NewDigitalOcean()
    }
    return nil
}
