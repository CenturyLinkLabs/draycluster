package provider

import (
	"encoding/base64"
	"errors"
	"github.com/CenturyLinkLabs/draycluster/deploy"
	"github.com/CenturyLinkLabs/draycluster/utils"
	"os"
)

type Amazon struct {
}

func NewAmazon() CloudProvider {
	c := Amazon{}
	return c
}

func (amz Amazon) ProvisionAgent() (deploy.CloudServer, error) {

	utils.LogInfo("\nDeploying Panamax remote agent in Amazon EC2")

	apiID := os.Getenv("AWS_ACCESS_KEY_ID")
	apiK := os.Getenv("AWS_SECRET_ACCESS_KEY")
	loc := os.Getenv("REGION")

	if apiID == "" || apiK == "" || loc == "" {
		return deploy.CloudServer{}, errors.New("\n\nMissing Params Or No Matching AMI found...Check Docs...\n\n")
	}

	var pk, puk, kn string
	if os.Getenv("MASTER_PUBLIC_KEY") != "" && os.Getenv("MASTER_PRIVATE_KEY") != "" {
		s1, _ := base64.StdEncoding.DecodeString(os.Getenv("MASTER_PRIVATE_KEY"))
		s2, _ := base64.StdEncoding.DecodeString(os.Getenv("MASTER_PUBLIC_KEY"))
		pk = string(s1)
		puk = string(s2)
		//kn = os.Getenv("AMAZON_SSH_KEY_NAME")
	} else {
		utils.LogInfo("\nCreating New Keys")
		pk, puk, _ = utils.CreateSSHKey()
	}

	c := deploy.Amazon{}
	c.AmiName = "ubuntu/images/hvm-ssd/ubuntu-trusty-14.04-amd64-server-20150123"
	c.AmiOwnerId = "099720109477"
	c.ApiAccessKey = apiK
	c.ApiKeyID = apiID
	c.Location = loc
	c.PrivateKey = pk
	c.PublicKey = puk
	c.ServerCount = 1
	c.TCPOpenPorts = []int{3001}
	c.ServerNames = []string{"Agent"}
	c.VMSize = "t2.micro"
	c.SSHKeyName = kn

	s, e := c.DeployVMs()
	if e != nil {
		return deploy.CloudServer{}, e
	}

	s[0].PublicSSHKey = puk
	s[0].PrivateSSHKey = pk

	utils.LogInfo("\nLogin Successful...Creating VMs...")

	return s[0], nil
}
