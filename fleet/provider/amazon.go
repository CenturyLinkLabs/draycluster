package provider

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/CenturylinkLabs/draycluster/deploy"
	"github.com/CenturylinkLabs/draycluster/utils"
	"os"
	"strconv"
	"strings"
)

type Amazon struct {
}

// NewAmazon is used to create a new client for using Amazon client to
// create RHEL 7 server cluster.
func NewAmazon() *Amazon {
    cl := new(Amazon)
	return cl
}

func (amz *Amazon) ProvisionCluster() ([]deploy.CloudServer, error) {

	utils.LogInfo("\nProvisioning cluster in Amazon EC2")

	apiID := os.Getenv("AWS_ACCESS_KEY_ID")
	apiK := os.Getenv("AWS_SECRET_ACCESS_KEY")
	loc := os.Getenv("REGION")
	vmSize := os.Getenv("VM_SIZE")
	cnt, e := strconv.Atoi(os.Getenv("NODE_COUNT"))

	if apiID == "" || apiK == "" || loc == "" || vmSize == "" {
		return nil, errors.New("\n\nMissing Params Or No Matching AMI found...Check Docs...\n\n")
	}

	pk, puk, _ := utils.CreateSSHKey()

	c := &deploy.Amazon{}
	c.ApiAccessKey = apiK
	c.ApiKeyID = apiID
	c.Location = loc
	c.PrivateKey = pk
	c.PublicKey = puk
	c.ServerCount = cnt
	c.VMSize = vmSize
	c.AmiName = "CoreOS-stable-557.2.0-hvm"
	c.AmiOwnerId = "595879546273"
    c.UserData = []byte(createCloudConfigCluster(puk))

	c.TCPOpenPorts = []int{4001, 7001}
	for _, p := range strings.Split(os.Getenv("OPEN_TCP_PORTS"), ",") {
		v, e := strconv.Atoi(p)
		if e == nil {
			c.TCPOpenPorts = append(c.TCPOpenPorts, v)
		}
	}


	for i := 0; i < c.ServerCount; i++ {
		c.ServerNames = append(c.ServerNames, fmt.Sprintf("Server-%d", i))
	}

	servers, e := c.DeployVMs()

	if e != nil {
		return nil, e
	}

	utils.LogInfo("\nCluster Creation Complete...")
	utils.SetKey("AMAZON_SSH_KEY_NAME", c.SSHKeyName)
	utils.SetKey("MASTER_PUBLIC_KEY", base64.StdEncoding.EncodeToString([]byte(puk)))
	utils.SetKey("UBUNTU_LOGIN_USER", "ubuntu")

	return servers, nil
}
