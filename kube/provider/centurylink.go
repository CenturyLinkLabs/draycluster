package provider

import (
    "errors"
    "github.com/CenturyLinkLabs/clcgo"
    "github.com/CenturyLinkLabs/draycluster/deploy"
    "github.com/CenturyLinkLabs/draycluster/utils"
    "os"
    "strconv"
    "strings"
    "fmt"
    )

// Centurylink has the data that is used for provisioning a server. Most of the
// data is passed in environment variables. The following env vars are required
// for provisioning a server in Centurylink, USERNAME, PASSWORD, GROUP_ID, CPU,
// MEMORY_GB, OPEN_TCP_PORTS
type Centurylink struct {
    clcClient   *clcgo.Client
    gID         string
    cpu         int
    memGb       int
    masterPK    string
    masterPuK   string
    uname       string
    password    string
    miCount     int
    minionPorts []int
}

// NewCenturylink is used to create a new client for using Centurylink client to
// create RHEL 7 server cluster.
func NewCenturylink() *Centurylink {
    cl := new(Centurylink)
    return cl
}

// ProvisionCluster is used to provision a cluster of RHEL7 VMs (1 Master +
// n Minions).
func (clc Centurylink) ProvisionCluster() ([]deploy.CloudServer, error) {
    utils.LogInfo("\nProvisioning Server Cluster in Centurylink")

    e := clc.initProvider()
    utils.LogInfo("\nMinion Count: " + strconv.Itoa(clc.miCount))

    if e != nil {
        return nil, e
    }

    var sn []string
    sn = append(sn, "MASTER")
    for i := 0; i < clc.miCount; i++ {
        sn = append(sn, fmt.Sprintf("MIN%d", i))
    }

    c := deploy.Centurylink{
        PrivateSSHKey: clc.masterPK,
        PublicSSHKey:  clc.masterPuK,
        APIUsername:   clc.uname,
        APIPassword:   clc.password,
        GroupID:       clc.gID,
        CPU:           clc.cpu,
        MemoryGB:      clc.memGb,
        ServerNames:   sn,
        TCPOpenPorts:  clc.minionPorts,
        ServerCount : clc.miCount+1,
    }

    s, e := c.DeployVMs()
    if e != nil {
        return nil, e
    }


    for i := 0; i <= clc.miCount; i++ {
        fmt.Printf("\n\n%s\n\n", s[i].Name)
        if !strings.Contains(s[i].Name, "MASTER") {
            s[i].PrivateSSHKey = ""
        } else {
            s[i].PrivateSSHKey = clc.masterPK
            s[i].PublicSSHKey = clc.masterPuK
        }
    }

    return s, nil
}

func (clc *Centurylink) initProvider() error {
    clc.uname = os.Getenv("USERNAME")
    clc.password = os.Getenv("PASSWORD")
    clc.gID = os.Getenv("GROUP_ID")
    clc.cpu, _ = strconv.Atoi(os.Getenv("CPU"))
    clc.memGb, _ = strconv.Atoi(os.Getenv("MEMORY_GB"))
    ps := os.Getenv("OPEN_TCP_PORTS")
    clc.miCount, _ = strconv.Atoi(os.Getenv("MINION_COUNT"))

    if ps != "" {
        s := strings.Split(ps, ",")
        for _, p := range s {
            v, e := strconv.Atoi(p)
            if e == nil {
                clc.minionPorts = append(clc.minionPorts, v)
            }
        }
    }

    if clc.uname == "" || clc.password == "" || clc.gID == "" {
        return errors.New("\n\nMissing Params.. in cluster creation...Check Docs....\n\n")
    }

    if clc.cpu <= 0 || clc.memGb <= 0 || clc.miCount <= 0 {
        return errors.New("\n\nMake sure CPU, MemoryGB and MINION_COUNT values are greater than 0.\n\n")
    }

    pk, puk, err := utils.CreateSSHKey()
    clc.masterPK = pk
    clc.masterPuK = puk

    if err != nil {
        return err
    }

    return nil
}
