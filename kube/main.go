package main

import (
	"encoding/base64"
//	"fmt"
	"github.com/CenturylinkLabs/draycluster/kube/provider"
	"github.com/CenturylinkLabs/draycluster/utils"
//	"os"
	"strings"
    "fmt"
    "os")

func main() {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			os.Exit(1)
		}
	}()

    utils.CreateRootCerts()

	cfg, e := utils.LoadJsonConfig()

	p := cfg["provider"]

	cp := provider.New(p)

	s, e := cp.ProvisionCluster()
	if e != nil {
		panic(e.Error())
	}

	mPuIP := ""
	mPrIP := ""
	mPK := ""

	var miIP []string

	for _, v := range s {
		if v.PrivateSSHKey == "" || mPuIP != "" {
			miIP = append(miIP, v.PrivateIP)
		} else {
			mPuIP = v.PublicIP
			mPrIP = v.PrivateIP
			mPK = v.PrivateSSHKey
		}
	}

	utils.SetKey("CLOUD_PROVIDER", p)
	utils.SetKey("MASTER_PUBLIC_IP", mPuIP)
	utils.SetKey("MASTER_PRIVATE_IP", mPrIP)
	utils.SetKey("MASTER_PRIVATE_KEY", base64.StdEncoding.EncodeToString([]byte(mPK)))
	utils.SetKey("MINION_IPS", strings.Join(miIP, ","))
    utils.SetKey("CLUSTER_TYPE", "kubernetes")
}
