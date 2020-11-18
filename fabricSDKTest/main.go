package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/policydsl"
)

type Client struct {
	sdk           *fabsdk.FabricSDK
	rc            *resmgmt.Client
	cc            *channel.Client
	goPath        string
	chaincodePath string
	ChaincodeID   string
	channelID     string
}

func New(configPath string) (c *Client, err error) {
	c = &Client{}
	c.chaincodePath = ""
	c.ChaincodeID = ""
	c.channelID = "mychannel"
	c.goPath = os.Getenv("GOPATH")

	c.sdk, err = fabsdk.New(config.FromFile(configPath))
	if err != nil {
		return nil, err
	}

	rcp := c.sdk.Context(fabsdk.WithUser("Admin"), fabsdk.WithOrg("Org1"))
	c.rc, err = resmgmt.New(rcp)
	if err != nil {
		return nil, err
	}

	ccp := c.sdk.ChannelContext(c.channelID, fabsdk.WithUser("User1"))
	c.cc, err = channel.New(ccp)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Client) install(version string) error {
	ccPkg, err := gopackager.NewCCPackage(c.chaincodePath, c.goPath)
	if err != nil {
		return err
	}
	req := resmgmt.InstallCCRequest{
		Name:    c.ChaincodeID,
		Path:    c.chaincodePath,
		Version: version,
		Package: ccPkg,
	}
	reqPeers := resmgmt.WithTargetEndpoints("peer0.org1.example.com")
	resps, err := c.rc.InstallCC(req, reqPeers)
	if err != nil {
		return err
	}

	for index := range resps {
		if resps[index].Status != http.StatusOK {
			err = fmt.Errorf("%v\n %s", err, resps[index].Info)
		}
	}
	return err
}

func (c *Client) Instantiate() error {
	org1OrOrg2 := "OR('Org1MSP.Member','Org2MSP.Member')"
	ccPolicy, err := policydsl.FromString(org1OrOrg2)
	if err != nil {
		return err
	}
	instantiateReq := resmgmt.InstantiateCCRequest{
		Name:    c.ChaincodeID,
		Path:    c.chaincodePath,
		Version: "1.0",
		Args:    [][]byte{[]byte("init"), []byte("a"), []byte("100"), []byte("b"), []byte("200")},
		Policy:  ccPolicy,
	}
	reqPeers := resmgmt.WithTargetEndpoints("peer0.org1.example.com")
	resp, err := c.rc.InstantiateCC(c.channelID, instantiateReq, reqPeers)
	if err != nil {
		return err
	}
	fmt.Println(resp.TransactionID)
	return nil
}

func (c *Client) Upgrade() error {
	org1OrOrg2 := "OR('Org1MSP.Member','Org2MSP.Member')"
	ccPolicy, err := policydsl.FromString(org1OrOrg2)
	if err != nil {
		return err
	}
	upgradeReq := resmgmt.UpgradeCCRequest{
		Name:    c.ChaincodeID,
		Path:    c.chaincodePath,
		Version: "1.0",
		Args:    [][]byte{[]byte("init"), []byte("a"), []byte("100"), []byte("b"), []byte("200")},
		Policy:  ccPolicy,
	}
	reqPeers := resmgmt.WithTargetEndpoints("peer0.org1.example.com")
	upgradeResp, err := c.rc.UpgradeCC(c.channelID, upgradeReq, reqPeers)
	if err != nil {
		return err
	}
	fmt.Println(upgradeResp.TransactionID)
	return nil
}

func (c *Client) fabricCa() {
}

var (
	configPath     = "/mnt/shareSSD/code/Fabric/first/fabricSDKTest/config_test.yaml"
	configPathOrg2 = "/mnt/shareSSD/code/Fabric/first/fabricSDKTest/config_test_org2.yaml"
)

func createSDK(path string) *fabsdk.FabricSDK {
	if len(path) == 0 {
		path = configPath
	}
	sdk, err := fabsdk.New(config.FromFile(path))
	if err != nil {
		log.Println(err)
		return nil
	}

	return sdk
}

func createSDK2(organization ...string) *fabsdk.FabricSDK {
	if len(organization) != 1 {
		return createSDK("")
	}

	if strings.ToLower(organization[0]) == "org2" {
		//fmt.Println("org2 SDK")
		return createSDK(configPathOrg2)
	}
	return nil
}

func createResourceManagement(sdk *fabsdk.FabricSDK, username, organization string) *resmgmt.Client {
	if sdk == nil {
		return nil
	}
	rc, err := resmgmt.New(sdk.Context(fabsdk.WithUser(username), fabsdk.WithOrg(organization)))
	if err != nil {
		log.Println(err)
		return nil
	}
	return rc
}

func createLedgerClient(sdk *fabsdk.FabricSDK, channelID, username, organization string) *ledger.Client {
	if sdk == nil {
		return nil
	}
	c, err := ledger.New(sdk.ChannelContext(channelID, fabsdk.WithUser(username), fabsdk.WithOrg(organization)))
	if err != nil {
		log.Println(err)
		return nil
	}
	return c
}

func createChannelClient(sdk *fabsdk.FabricSDK, channelID, username, organization string) *channel.Client {
	if sdk == nil {
		return nil
	}

	c, err := channel.New(sdk.ChannelContext(channelID, fabsdk.WithOrg(organization), fabsdk.WithUser(username)))
	if err != nil {
		return nil
	}
	return c
}
func main() {}
