package deploy

import (
    "errors"
    "fmt"
    "github.com/CenturyLinkLabs/draycluster/utils"
    "github.com/mitchellh/goamz/aws"
    "github.com/mitchellh/goamz/ec2"
    "os"
    "strings"
    "time"
   )

type Amazon struct {
    VMSize       string
    Location     string
    ApiKeyID     string
    ApiAccessKey string
    AmiName      string
    AmiOwnerId   string
    SSHKeyName   string
    PrivateKey   string
    PublicKey    string
    amzClient    *ec2.EC2
    ServerCount  int
    TCPOpenPorts []int
    UDPOpenPorts []int
    ServerNames  []string
    UserData     []byte
}

func (amz *Amazon) DeployVMs() ([]CloudServer, error) {
    var e error
    e = amz.init()
    if e != nil {
        return nil, e
    }

    amz.AmiName, e = amz.getAmiID()
    if amz.AmiName == "" || e != nil {
        return nil, errors.New("AMI Not found for provisioning. Cannot proceed.!!!")
    }
    utils.LogInfo(fmt.Sprintf("AMI Used: %s", amz.AmiName))

    if amz.SSHKeyName == "" && (amz.PublicKey != "" && amz.PrivateKey != "") {
        amz.SSHKeyName, e = amz.importKey(amz.PublicKey)
        if e != nil {
            return nil, e
        }
    } else if amz.SSHKeyName == "" {
        return nil, errors.New("Please pass ssh keyname or a Private Key & Public Key to create vms.")
    }

    sg, e := amz.createFWRules()
    if e != nil {
        return nil, e
    }

    req := &ec2.RunInstances{
        ImageId:        amz.AmiName,
        InstanceType:   amz.VMSize,
        MinCount:       amz.ServerCount,
        MaxCount:       amz.ServerCount,
        KeyName:        amz.SSHKeyName,
        SecurityGroups: []ec2.SecurityGroup{sg},
        UserData:       amz.UserData,
    }

    resp, e := amz.amzClient.RunInstances(req)
    if e != nil {
        return nil, e
    }

    utils.LogInfo("\nWaiting for servers to provision....")

    var servers []CloudServer
    for i, inst := range resp.Instances {
        s, e := amz.waitForServer(inst)
        if e != nil {
            return nil, e
        }
        amz.amzClient.CreateTags([]string{inst.InstanceId}, []ec2.Tag{{Key: "Name", Value: amz.ServerNames[i]}})
        s.Name = amz.ServerNames[i]
        s.PrivateSSHKey = amz.PrivateKey
        servers = append(servers, s)
    }
    utils.LogInfo("\nProvisioning complete...")
    return servers, nil
}

func (amz *Amazon) getAmiID() (string, error) {
    f := ec2.NewFilter()
    f.Add("name", "*"+amz.AmiName+"*")
    f.Add("owner-id", amz.AmiOwnerId)
    im, _ := amz.amzClient.Images(nil, f)
    if im != nil && len(im.Images) > 0 {
        return im.Images[0].Id, nil
    }
    return "", errors.New("Image not found")
}

func (amz *Amazon) createFWRules() (ec2.SecurityGroup, error) {

    g := ec2.SecurityGroup{}
    g.Name = "pmx-security-group-" + utils.RandSeq(4)
    g.Description = "panamax security group"
    var ps []ec2.IPPerm

    amz.TCPOpenPorts = append(amz.TCPOpenPorts, 22)
    for _, p := range amz.TCPOpenPorts {
        ps = append(ps, ec2.IPPerm{Protocol: "tcp", SourceIPs: []string{"0.0.0.0/0"}, ToPort: p, FromPort: p})
    }
    for _, p := range amz.UDPOpenPorts {
        ps = append(ps, ec2.IPPerm{Protocol: "udp", SourceIPs: []string{"0.0.0.0/0"}, ToPort: p, FromPort: p})
    }


    _, e := amz.amzClient.CreateSecurityGroup(g)
    if e != nil {
        return ec2.SecurityGroup{}, e
    }
    _, e = amz.amzClient.AuthorizeSecurityGroup(g, ps)
    if e != nil {
        return ec2.SecurityGroup{}, e
    }
    return g, nil
}

func (amz *Amazon) waitForServer(inst ec2.Instance) (CloudServer, error) {
    for {
        if inst.State.Code == 16 {
            break
        }
        time.Sleep(30 * time.Second)
        resp, e := amz.amzClient.Instances([]string{inst.InstanceId}, &ec2.Filter{})
        if e != nil {
            panic(e)
            return CloudServer{}, e
        }
        inst = resp.Reservations[0].Instances[0]
    }
    utils.LogInfo(fmt.Sprintf("\nServer Provisioned: Public IP: %s, Private IP: %s", inst.PublicIpAddress, inst.PrivateIpAddress))
    return CloudServer{PublicIP: inst.PublicIpAddress, Name: inst.DNSName, PrivateIP: inst.PrivateIpAddress}, nil
}

func (amz *Amazon) init() error {

    if amz.ApiKeyID == "" || amz.ApiAccessKey == "" || amz.Location == "" || amz.VMSize == "" || len(amz.ServerNames) != amz.ServerCount {
        return errors.New("\n\nMissing Params Or No Matching AMI found...Check Docs...\n\n")
    }

    os.Setenv("AWS_ACCESS_KEY_ID", amz.ApiKeyID)
    os.Setenv("AWS_SECRET_ACCESS_KEY", amz.ApiAccessKey)

    auth, e := aws.EnvAuth()
    if e != nil {
        return e
    }

    var r aws.Region
    for _, r = range aws.Regions {
        if strings.Contains(amz.Location, r.Name) {
            break
        }
    }

    amz.amzClient = ec2.New(auth, r)

    return nil
}

func (amz *Amazon) importKey(puk string) (string, error) {
    kn := "pmx-keypair-" + utils.RandSeq(4)
    _, e := amz.amzClient.ImportKeyPair(kn, puk)

    if e != nil {
        panic(e)
        return "", e
    }
    return kn, nil
}

func (amz Amazon) ExecSSHCmd(publicIP string, privateKey string, command string) (string, error) {
    return utils.ExecSSHCmd(publicIP, privateKey, "ec2-user", command)
}

