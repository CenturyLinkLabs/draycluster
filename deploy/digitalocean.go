package deploy

import (
    "code.google.com/p/goauth2/oauth"
    "errors"
    "github.com/digitalocean/godo"
    "time"
    "github.com/CenturylinkLabs/draycluster/utils"
    "strings")

type DigitalOcean struct {
    Location    string
    SSHKeyName  string
    APIToken    string
    ServerCount int
    doClient    *godo.Client
    PrivateKey  string
    PublicKey   string
    UserData    string
    VMSize      string
    ImageName   string
    ServerNames  []string
}

func (do *DigitalOcean) DeployVMs() ([]CloudServer, error) {
    e := do.init()
    if e != nil {
        return nil, e
    }

    utils.LogInfo("Deploying Server")

    var servers []CloudServer
    for i := 0; i < do.ServerCount; i++ {
        s, e := do.deployServer(do.ServerNames[i])
        if e == nil {
            servers = append(servers, s)
        } else {
            return servers, e
        }
    }
    return servers, nil
}

func (do *DigitalOcean) init() error {
    if do.APIToken == "" || do.ServerCount == 0 || do.Location == "" {
        return errors.New("\n\nMissing Params...Check Docs...\n\n")
    }
    t := &oauth.Transport{Token: &oauth.Token{AccessToken: do.APIToken}}
    do.doClient = godo.NewClient(t.Client())
    return nil
}

func (do *DigitalOcean) createSSHKey() (int, error) {
    kId := do.getSshKeyId()
    if (kId != -1 ) {
        return kId, nil
    }
    r := &godo.KeyCreateRequest{
        Name : do.SSHKeyName,
        PublicKey : do.PublicKey,
    }
    k, _, e := do.doClient.Keys.Create(r)
    if e != nil {
        return -1, e
    }
    return k.ID, e
}

func (do DigitalOcean) getSshKeyId() int {
    ks, _, e := do.doClient.Keys.List(&godo.ListOptions{Page: 1, PerPage: 100})
    if e != nil {
        return -1
    }
    for _, k := range ks {
        if strings.EqualFold(k.Name, do.SSHKeyName) {
            return k.ID
        }
    }
    return -1
}

func (do *DigitalOcean) deployServer(name string) (CloudServer, error) {
    k, e := do.createSSHKey()
    if e != nil {
        panic(e)
    }
    var kIds []interface{}
    kIds = append(kIds, k)

    cr := &godo.DropletCreateRequest{
        Name:              name,
        Region:            do.Location,
        Size:              do.VMSize,
        Image:             do.ImageName,
        PrivateNetworking: true,
        UserData:          do.UserData,
        SSHKeys:           kIds,
    }

    d, _, e := do.doClient.Droplets.Create(cr)

    if e != nil {
        return CloudServer{}, e
    }

    for {
        d, _, _ = do.doClient.Droplets.Get(d.Droplet.ID)
        if d.Droplet.Status == "active" {
            break
        }
        time.Sleep(60 * time.Millisecond)
    }

    return CloudServer{
        Name:          d.Droplet.Name,
        PrivateIP:     d.Droplet.Networks.V4[0].IPAddress,
        PublicIP:      d.Droplet.Networks.V4[1].IPAddress,
        PrivateSSHKey: do.PrivateKey,
        PublicSSHKey:  do.PublicKey,
    }, nil
}
