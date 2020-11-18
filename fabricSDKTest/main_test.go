package main

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	mspcommon "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/lifecycle"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/policydsl"
)

func TestSdk(t *testing.T) {
	sdk, err := fabsdk.New(config.FromFile("/mnt/shareSSD/code/Fabric/first/fabricSDKTest/config_test.yaml"))
	if err != nil {
		t.Error(err)
	}
	defer sdk.Close()

	sdkConfig, err := sdk.Config()
	if err != nil {
		t.Error(err)
	}
	t.Log(sdkConfig.Lookup("client"))
	t.Log(sdkConfig.Lookup("channels"))
	t.Log(sdkConfig.Lookup("organizations"))
	t.Log(sdkConfig.Lookup("orderers"))
	t.Log(sdkConfig.Lookup("peers"))
	t.Log(sdkConfig.Lookup("certificateAuthorities"))
	t.Log(sdkConfig.Lookup("entityMatchers"))
	t.Log(sdkConfig.Lookup("operations"))
	t.Log(sdkConfig.Lookup("metrics"))
}

/**
--------------------------Channel-----------------------------------
*/

func createChannel(resourceManager *resmgmt.Client) {

}

func CreateChannel(resourceManager *resmgmt.Client, channelID string,
	signed []mspcommon.SigningIdentity, channelConfigPath string) {
	req := resmgmt.SaveChannelRequest{
		ChannelID:         channelID,
		ChannelConfigPath: channelConfigPath,
		SigningIdentities: signed,
	}
	tx, err := resourceManager.SaveChannel(req, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(orderPeer))
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(tx)
}

func TestCreateChannels(t *testing.T) {
	var signed []mspcommon.SigningIdentity
	org1, err := org1MSPClient.GetSigningIdentity("Admin")
	if err != nil {
		log.Println(err)
		return
	}
	signed = append(signed, org1)
	org2, err := org2MSPClient.GetSigningIdentity("Admin")
	if err != nil {
		log.Println(err)
		return
	}
	signed = append(signed, org2)
	CreateChannel(orderRSM, channelID, signed, orgChannel)

	CreateChannel(org1RSM, channelID, []mspcommon.SigningIdentity{org1}, org1Anchor)

	CreateChannel(org2RSM, channelID, []mspcommon.SigningIdentity{org2}, org2Anchor)
}

func TestCreateChannel(t *testing.T) {
	userName := "Admin"
	orgName := "org1"

	sdk := createSDK("")
	resourceManagement := createResourceManagement(sdk, userName, orgName)
	if resourceManagement == nil {
		t.Error("resourceManagement create error")
	}

	// create Channel
	mspClient, err := msp.New(sdk.Context(), msp.WithOrg(orgName))
	if err != nil {
		t.Error(err)
	}

	adminIdentity, err := mspClient.GetSigningIdentity(userName)
	if err != nil {
		t.Error(err)
	}

	channelID := "channel1"
	req := resmgmt.SaveChannelRequest{
		ChannelID:         channelID,
		ChannelConfigPath: "/mnt/shareSSD/code/Fabric/fabric-samples/test-network/channel-artifacts_2/" + channelID + ".tx",
		SigningIdentities: []mspcommon.SigningIdentity{adminIdentity},
	}
	txID, err := resourceManagement.SaveChannel(req, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint("orderer.example.com"))
	if err != nil {
		t.Error(err)
	}
	t.Log(txID)
	// d1e58176d232494572f773c518b338ddae7a39241a455c02a32d96d48ae831fd
}

func JoinChannel(resourceManager *resmgmt.Client, channel, orderPeer string) {
	// join channel
	// 这里会将配置文件中的组织的peer全部加入
	// 所以要确认配置文件中的peer是真实存在的
	err := resourceManager.JoinChannel(channel, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(orderPeer))
	if err != nil {
		log.Println(err)
	}
}

func TestJoinChannel(t *testing.T) {
	JoinChannel(org1RSM, channelID, orderPeer)
	JoinChannel(org2RSM, channelID, orderPeer)
}

func QueryChannel(resourceManager *resmgmt.Client, peer string) {
	resp, err := resourceManager.QueryChannels(resmgmt.WithTargetEndpoints(peer))
	if err != nil {
		log.Println(err)
	}
	fmt.Println(resp)
}

func TestQueryChannel(t *testing.T) {
	QueryChannel(org1RSM, org1Peer0)
	QueryChannel(org2RSM, org2Peer0)
}

func FetchChannelBlock(resourceManager *resmgmt.Client) {
}

/**
-------------------------------chaincode-----------------------------
*/

var (
	org1SDK           = createSDK2()
	org2SDK           = createSDK2("org2")
	org1RSM           = createResourceManagement(org1SDK, "Admin", "org1")
	org2RSM           = createResourceManagement(org2SDK, "Admin", "org2")
	orderRSM          = createResourceManagement(org1SDK, "Admin", "ordererorg")
	org1ChannelClient = createChannelClient(org1SDK, channelID, "Admin", "org1")
	org1Peer0         = "peer0.org1.example.com"
	org2Peer0         = "peer0.org2.example.com"
	org1MSP           = "Org1MSP"
	org2MSP           = "Org2MSP"
	org1MSPClient     *msp.Client
	org2MSPClient     *msp.Client
	channelID         = "channel1"
	orderPeer         = "orderer.example.com"

	orgChannel = "/mnt/shareSSD/code/Fabric/fabric-samples/test-network/channel-artifacts_2/channel1.tx"
	org1Anchor = "/mnt/shareSSD/code/Fabric/fabric-samples/test-network/channel-artifacts_2/Org1MSPanchors.tx"
	org2Anchor = "/mnt/shareSSD/code/Fabric/fabric-samples/test-network/channel-artifacts_2/Org2MSPanchors.tx"
)

func init() {
	var err error
	org1MSPClient, err = msp.New(org1SDK.Context())
	if err != nil {
		log.Panicln(err)
	}

	org2MSPClient, err = msp.New(org2SDK.Context())
	if err != nil {
		log.Panicln(err)
	}
}

func InstallCC(resourceManager *resmgmt.Client, peer string) {
	chaincodePath := "blockchainsTest"
	fmt.Println(os.Getenv("GOPATH"))
	ccPkg, err := gopackager.NewCCPackage(chaincodePath, os.Getenv("GOPATH"))
	if err != nil {
		log.Println(err)
		return
	}

	req := resmgmt.InstallCCRequest{
		Name:    "sacc",
		Path:    chaincodePath,
		Version: "1.0",
		Package: ccPkg,
	}

	resp, err := resourceManager.InstallCC(req, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithTargetEndpoints(peer))
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(resp)
}

func TestInstallCC(t *testing.T) {
	sdk := createSDK("")
	resourceManager := createResourceManagement(sdk, "Admin", "org1")
	peerI := "peer0.org1.example.com"
_org2Install:
	chaincodePath := "blockchainsTest"
	t.Log(os.Getenv("GOPATH"))
	ccPkg, err := gopackager.NewCCPackage(chaincodePath, os.Getenv("GOPATH"))
	if err != nil {
		t.Error(err)
		return
	}

	req := resmgmt.InstallCCRequest{
		Name:    "sacc",
		Path:    chaincodePath,
		Version: "1.0",
		Package: ccPkg,
	}

	resp, err := resourceManager.InstallCC(req, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithTargetEndpoints(peerI))
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(resp)

	resourceManager = createResourceManagement(sdk, "Admin", "org2")
	goto _org2Install
}

func TestInstantiateCC(t *testing.T) {
	sdk := createSDK("")
	resourceManager := createResourceManagement(sdk, "Admin", "org1")
	//ccPolicy, err := policydsl.FromString("OR('Org1MSP.Member','Org2MSP.Member')")
	//if err != nil {
	//	t.Error(err)
	//	return
	//}

	chaincodePath := "blockchainsTest"
	req := resmgmt.InstantiateCCRequest{
		Name:    "sacc",
		Path:    chaincodePath,
		Version: "1.0",
		Lang:    peer.ChaincodeSpec_GOLANG,
		Args:    [][]byte{[]byte("init"), []byte("a"), []byte("b")},
		Policy:  policydsl.SignedByAnyMember([]string{"Org1MSP"}),
	}

	resp, err := resourceManager.InstantiateCC("channel1", req, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithTargetEndpoints("peer0.org1.example.com"))
	if err != nil {
		t.Error(err)
	}
	t.Log(resp)
	//0af43b9122da2df5403dc71319608953fecb80b568f42f32667643b87127028f
}

/**
-----------------------上面使用lscc安装链码貌似在fabric2.0已经禁止使用---------------
*/

func packageCC() (string, []byte) {
	desc := &lifecycle.Descriptor{
		Path:  "/mnt/shareSSD/code/Fabric/first/blockchainsTest",
		Type:  peer.ChaincodeSpec_GOLANG,
		Label: "sacc",
	}
	ccPkg, err := lifecycle.NewCCPackage(desc)
	if err != nil {
		panic(err)
	}
	return desc.Label, ccPkg
}

func packageID() string {
	l, c := packageCC()
	return lifecycle.ComputePackageID(l, c)
}

func LifecycleInstallCC(resourceManager *resmgmt.Client) {
	label, ccPkg := packageCC()
	req := resmgmt.LifecycleInstallCCRequest{
		Label:   label,
		Package: ccPkg,
	}

	packageID := lifecycle.ComputePackageID(label, ccPkg)
	resp, err := resourceManager.LifecycleInstallCC(req, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(packageID, resp)
}

func TestLifecycleInstallCC(t *testing.T) {
	LifecycleInstallCC(org1RSM)
	LifecycleInstallCC(org2RSM)
}

func LifecycleGetInstalledChaincode(resourceManager *resmgmt.Client, peer string) {
	resp, err := resourceManager.LifecycleGetInstalledCCPackage(packageID(), resmgmt.WithTargetEndpoints(peer), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		log.Println(err)
	}
	fmt.Println(len(resp)) // <- 返回的resp貌似是二进制文件
}

func TestLifecycleGetInstalledChaincode(t *testing.T) {
	LifecycleGetInstalledChaincode(org1RSM, org1Peer0)
	LifecycleGetInstalledChaincode(org2RSM, org2Peer0)
}

func LifecycleQueryChaincode(resourceManager *resmgmt.Client, peer string) {
	resp, err := resourceManager.LifecycleQueryInstalledCC(resmgmt.WithTargetEndpoints(peer), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		log.Println(err)
	}
	fmt.Println(resp)
}

func TestLifecycleQueryChaincode(t *testing.T) {
	LifecycleQueryChaincode(org1RSM, org1Peer0)
	LifecycleQueryChaincode(org2RSM, org2Peer0)
}

func LifecycleApproveChaincode(resourceManager *resmgmt.Client, msp, peer, orderPeer, channel string) {
	req := resmgmt.LifecycleApproveCCRequest{
		Name:              "sacc",
		Version:           "3.0",
		PackageID:         packageID(),
		Sequence:          3,
		EndorsementPlugin: "escc",
		ValidationPlugin:  "vscc",
		SignaturePolicy:   policydsl.SignedByAnyMember([]string{msp}),
		//ChannelConfigPolicy: "",
		//CollectionConfig:    nil,
		InitRequired: true,
	}

	resp, err := resourceManager.LifecycleApproveCC(channel, req, resmgmt.WithTargetEndpoints(peer), resmgmt.WithOrdererEndpoint(orderPeer), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(resp) // <- txnID 0122088e9eca7e99e487afdb4415ecde5f8632fa1a2a0ae91bd6b012fd17455
}

func TestLifecycleApproveChaincode(t *testing.T) {
	LifecycleApproveChaincode(org1RSM, org1MSP, org1Peer0, orderPeer, channelID)
	LifecycleApproveChaincode(org2RSM, org1MSP, org2Peer0, orderPeer, channelID)
	//LifecycleApproveChaincode(org2RSM, org2MSP, org2Peer0, orderPeer, channelID)
	//LifecycleApproveChaincode(org1RSM, org2MSP, org1Peer0, orderPeer, channelID)
}

func QueryApprovedCC(resourceManager *resmgmt.Client, channel, peer string) {
	req := resmgmt.LifecycleQueryApprovedCCRequest{
		Name:     "sacc",
		Sequence: 3,
	}

	resp, err := resourceManager.LifecycleQueryApprovedCC(channel, req, resmgmt.WithTargetEndpoints(peer), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(resp)
}

func TestLifecycleQueryApprovedCC(t *testing.T) {
	QueryApprovedCC(org1RSM, channelID, org1Peer0)
	QueryApprovedCC(org2RSM, channelID, org2Peer0)
}

func LifecycleCheckCommitCC(resourceManager *resmgmt.Client, channel, msp, peer string) {
	req := resmgmt.LifecycleCheckCCCommitReadinessRequest{
		Name:              "sacc",
		Version:           "3.0",
		PackageID:         packageID(),
		Sequence:          3,
		EndorsementPlugin: "escc",
		ValidationPlugin:  "vscc",
		SignaturePolicy:   policydsl.SignedByAnyMember([]string{msp}),
		//ChannelConfigPolicy: "",
		//CollectionConfig:    nil,
		InitRequired: true,
	}

	resp, err := resourceManager.LifecycleCheckCCCommitReadiness(channel, req, resmgmt.WithTargetEndpoints(peer), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(resp)
}

func TestLifecycleCheckCommitCC(t *testing.T) {
	LifecycleCheckCommitCC(org1RSM, channelID, org1MSP, org1Peer0)
	LifecycleCheckCommitCC(org2RSM, channelID, org2MSP, org2Peer0)
}

func LifecycleCommitChaincode(resourceManager *resmgmt.Client, channel, msp, peer, orderPeer string) {
	req := resmgmt.LifecycleCommitCCRequest{
		Name:              "sacc",
		Version:           "3.0",
		Sequence:          3,
		EndorsementPlugin: "escc",
		ValidationPlugin:  "vscc",
		SignaturePolicy:   policydsl.SignedByAnyMember([]string{org1MSP}),
		//ChannelConfigPolicy: "",
		//CollectionConfig:    nil,
		InitRequired: true,
	}

	resp, err := resourceManager.LifecycleCommitCC(
		channel,
		req,
		resmgmt.WithTargetEndpoints("peer0.org1.example.com", "peer0.org2.example.com"), // <- 这里注意因为背书策略的原因,必须同时提交两个背书节点,否则会因为不满足背书策略而失败
		resmgmt.WithOrdererEndpoint(orderPeer), resmgmt.WithRetry(retry.DefaultResMgmtOpts),
	)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(resp)
}

func TestLifecycleCommitChaincode(t *testing.T) {
	LifecycleCommitChaincode(org1RSM, channelID, org1MSP, org1Peer0, orderPeer)
	//LifecycleCommitChaincode(org2RSM, channelID, org2MSP, org2Peer0, orderPeer)
}

func LifecycleQueryCommittedCC(resourceManager *resmgmt.Client, channel, peer string) {
	req := resmgmt.LifecycleQueryCommittedCCRequest{
		Name: "sacc",
	}

	resp, err := resourceManager.LifecycleQueryCommittedCC(channel, req, resmgmt.WithTargetEndpoints(peer), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(resp)
}

func TestLifecycleQueryCommittedCC(t *testing.T) {
	LifecycleQueryCommittedCC(org1RSM, channelID, org1Peer0)
	LifecycleQueryCommittedCC(org2RSM, channelID, org2Peer0)
}

func LifecycleInitCC(channelClient *channel.Client) {
	resp, err := channelClient.Execute(
		channel.Request{
			ChaincodeID: "sacc",
			Fcn:         "init",
			Args:        [][]byte{[]byte("a"), []byte("xiaoxiao")},
			//TransientMap:    nil,
			//InvocationChain: nil,
			IsInit: true,
		},
		channel.WithRetry(retry.DefaultChannelOpts),
	)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(resp)
}

func TestLifecycleInitCC(t *testing.T) {
	LifecycleInitCC(org1ChannelClient)
}

func TestExecuteCC(t *testing.T) {
	resp, err := org1ChannelClient.Query(channel.Request{
		ChaincodeID: "sacc",
		Fcn:         "set",
		Args:        [][]byte{[]byte("1"), []byte("xiaoxiao")},
	},
		channel.WithRetry(retry.DefaultChannelOpts),
		channel.WithTargetEndpoints(org1Peer0, org2Peer0),
	)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(resp.TransactionID, resp.Responses)
}

func TestQueryCC(t *testing.T) {
	resp, err := org1ChannelClient.Query(
		channel.Request{
			ChaincodeID: "sacc",
			Fcn:         "get",
			Args:        [][]byte{[]byte("a")},
			//TransientMap:    nil,
			//InvocationChain: nil,
			//IsInit:          false,
		},
		channel.WithRetry(retry.DefaultChannelOpts),
		channel.WithTargetEndpoints(org1Peer0, org2Peer0),
	)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(resp.Payload), resp.TransactionID)
}

func UpgradeCC(resourceManager *resmgmt.Client) {

}

/**
--------------------------------ledger-------------------------------
*/
func TestLedgerQueryBlock(t *testing.T) {
	sdk := createSDK("")
	c := createLedgerClient(sdk, "channel1", "Admin", "org1")
	block, err := c.QueryBlock(0)
	if err != nil {
		t.Error(err)
	}
	t.Log(block.GetData())
	t.Log(block.GetHeader())
	t.Log(block.GetMetadata())
	t.Log(c.QueryBlock(1))
}

func TestQueryTransaction(t *testing.T) {
	sdk := createSDK("")
	c := createLedgerClient(sdk, "channel1", "Admin", "org1")
	t.Log(c.QueryTransaction("0"))
}

func TestQueryInfoAndConfig(t *testing.T) {
	sdk := createSDK("")
	c := createLedgerClient(sdk, "channel1", "Admin", "org1")
	t.Log(c.QueryInfo())
	t.Log(c.QueryConfig())
	t.Log(c.QueryConfigBlock())
}
