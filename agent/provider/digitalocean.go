package provider

import (
    "os"
    "github.com/CenturyLinkLabs/draycluster/deploy"
    "github.com/CenturyLinkLabs/draycluster/utils"
    "errors"
    "encoding/base64")

type DigitalOcean struct {}

func NewDigitalOcean() *DigitalOcean {
    cl := new(DigitalOcean)
    return cl
}

func (do *DigitalOcean) ProvisionAgent() (deploy.CloudServer, error) {
    utils.LogInfo("\nProvisioning agent in Digital Ocean")

    apiT := os.Getenv("API_TOKEN")
    loc := os.Getenv("REGION")

    if apiT == "" || loc == "" {
        return deploy.CloudServer{}, errors.New("\n\nMissing Params Or No Matching AMI found...Check Docs...\n\n")
    }

    utils.LogInfo("\nParams Found")

    var pk, puk, kn string
    if os.Getenv("SSH_KEY_NAME") != "" {
        kn = os.Getenv("SSH_KEY_NAME")
    }
    if os.Getenv("MASTER_PUBLIC_KEY") != "" && os.Getenv("MASTER_PRIVATE_KEY") != "" {
        utils.LogInfo("Keys Found")
        s1, _ := base64.StdEncoding.DecodeString(os.Getenv("MASTER_PRIVATE_KEY"))
        s2, _ := base64.StdEncoding.DecodeString(os.Getenv("MASTER_PUBLIC_KEY"))
        pk = string(s1)
        puk = string(s2)
    } else {
        utils.LogInfo("\nCreating New Keys")
        pk, puk, _ = utils.CreateSSHKey()
    }

    c := &deploy.DigitalOcean{}
    c.APIToken = apiT
    c.Location = loc
    c.PrivateKey = pk
    c.PublicKey = puk
    c.SSHKeyName = kn
    c.ServerCount = 1
    c.VMSize = "512mb"
    c.ImageName = "ubuntu-14-10-x64"
    c.ServerNames  = []string{"Agent"}

    utils.LogInfo("\nDeploying Agent Server")
    s, e := c.DeployVMs()
    if e != nil {
        return deploy.CloudServer{}, e
    }

    utils.LogInfo("\nServer Created")
    s[0].PublicSSHKey = puk
    s[0].PrivateSSHKey = pk

    return s[0], nil
}


