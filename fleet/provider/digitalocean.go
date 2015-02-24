package provider

import (
    "strconv"
    "fmt"
    "os"
    "github.com/CenturylinkLabs/draycluster/deploy"
    "github.com/CenturylinkLabs/draycluster/utils"
    "errors"
    "encoding/base64")

type DigitalOcean struct {}

func NewDigitalOcean() *DigitalOcean {
    cl := new(DigitalOcean)
    return cl
}

func (do *DigitalOcean) ProvisionCluster() ([]deploy.CloudServer, error) {
    utils.LogInfo("\nProvisioning cluster in Digital Ocean")

    apiT := os.Getenv("API_TOKEN")
    loc := os.Getenv("REGION")
    vmSize := os.Getenv("VM_SIZE")
    cnt, e := strconv.Atoi(os.Getenv("NODE_COUNT"))

    if apiT == "" || loc == "" || vmSize == "" {
        return nil, errors.New("\n\nMissing Params Or No Matching AMI found...Check Docs...\n\n")
    }

    pk, puk, _ := utils.CreateSSHKey()

    c := &deploy.DigitalOcean{}
    c.APIToken = apiT
    c.Location = loc
    c.PrivateKey = pk
    c.PublicKey = puk
    c.ServerCount = cnt
    c.VMSize = vmSize
    c.ImageName = "coreos-stable"
    c.UserData = createCloudConfigCluster(puk)
    c.SSHKeyName = "PanamaxKey-" + utils.RandSeq(4)

    for i := 0; i < c.ServerCount; i++ {
        c.ServerNames = append(c.ServerNames, fmt.Sprintf("Server-%d", i))
    }

    servers, e := c.DeployVMs()

    if e != nil {
        return nil, e
    }

    utils.LogInfo("\nCluster Creation Complete...")
    utils.SetKey("MASTER_PUBLIC_KEY", base64.StdEncoding.EncodeToString([]byte(puk)))
    utils.SetKey("UBUNTU_LOGIN_USER", "root")
    utils.SetKey("SSH_KEY_NAME", c.SSHKeyName)

    return servers, nil
}


