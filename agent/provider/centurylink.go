package provider

import (
    "github.com/CenturyLinkLabs/draycluster/deploy"
    "github.com/CenturyLinkLabs/draycluster/utils"
    "os"
)

type Centurylink struct {
}

func NewCenturylink() CloudProvider {
    c := Centurylink{}
    return c
}

func (clc Centurylink) ProvisionAgent() (deploy.CloudServer, error) {

    pk, puk, _ := utils.CreateSSHKey()

    c := deploy.Centurylink{
        PublicSSHKey:   puk,
        APIUsername:    os.Getenv("USERNAME"),
        APIPassword:    os.Getenv("PASSWORD"),
        GroupID:        os.Getenv("GROUP_ID"),
        CPU:            1,
        MemoryGB:       1,
        TCPOpenPorts:   []int{3001},
        ServerNames:     []string{"AGENT"},
        ServerTemplate: "UBUNTU-14-64-TEMPLATE",
        ServerCount:    1,
    }

    utils.LogInfo("\nWaiting for server creation")

    s, e := c.DeployVMs()

    s[0].PrivateSSHKey = pk
    s[0].PublicSSHKey = puk

    if e != nil {
        return deploy.CloudServer{}, e
    }
    return s[0], nil
}
