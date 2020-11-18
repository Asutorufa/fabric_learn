#

获取fabric

```shell
go get https://github.com/hyperledger/fabric
go get https://github.com/hyperledger/fabric-samples
```

切换分支

```shell
cd $GOPATH/src/github.com/hyperledger/fabric
git checkout release-1.0
```

安装configtxgen

```shell
cd $GOPATH/src/github.com/hyperledger/fabric/common/configtx/tool/configtxgen
go install
```

安装 cryptogen

```shell
cd $GOPATH/src/github.com/hyperledger/fabric/common/tools/cryptogen
go install
```

生成第一个区块

```shell
cd $GOPATH/src/github.com/hyperledger/fabric-samples
git checkout release-1.0
cd $GOPATH/src/github.com/hyperledger/fabric-samples/first-network
./byfn.sh -m generate -c mychannel
```

启动

```shell
# 设置代理
sudo vim /usr/lib/systemd/system/docker.service
# add
# [Service]
# Environment="HTTP_PROXY=http://x.x.x.x:port"
# then restart docker: sudo systemctl restart docker

sudo docker pull hyperledger/fabric-orderer:x86_64-1.0.0
sudo docker tag hyperledger/fabric-orderer:x86_64-1.0.0 hyperledger/fabric-orderer:latest
sudo docker pull hyperledger/fabric-peer:x86_64-1.0.0
sudo docker tag hyperledger/fabric-peer:x86_64-1.0.0 hyperledger/fabric-peer:latest
sudo docker pull hyperledger/fabric-tools:x86_64-1.0.0
sudo docker tag hyperledger/fabric-tools:x86_64-1.0.0 hyperledger/fabric-tools:latest
sudo ./byfn.sh -m up -c mychannel
```

停止

```shell
sudo ./byfn.sh -m down -c mychannel
```

交易过程:

- 交易背书(模拟@Endorser) <- 模拟进行,不会持久化数据
- 交易排序(排序@Orderer) <- 共识
- 交易验证(验证@Committer)

## 交易排序

- 交易排序
  - 目的: 保证系统交易顺序的一致性(有限状态机)
  - solo: 单节点排序
  - kafka: 集群
- 区块分发
  - 中间状态区块, 不管有效无效都会发给主节点
- 多通道数据隔离

## 账本存储

- 交易过程
  - 交易模拟 -> 读写集(RWSet)
  - 交易排序
  - 交易验证 -> 状态更新

- 交易读写集(RWSet) <- 背书时交易模拟生成
  - 读集: 读取**已提交**的状态值(值)
  - 写集: 将要更新的状态键值对
  - 写集: 状态键值对删除标记
  - 写集: 多次更新以最后一次为准
  - 版本号: 二元组(区块高度, 交易编号)

- 交易验证
  - 读集版本号 是否= 世界状态版本号(包括未提交的前序交易)

- 世界状态
  - 交易执行后的所有键的最新值
  - 显著提升链码执行效率
  - 状态是所有交易日志的快照, 可随时重构
  - LevelDB(键值对) or CouchDB(状态数据库)

- 历史数据索引(可选)
  - 某键在某区块的某条交易中被改变
  - 只记录改变动作, 不记录具体改变
  - 历史读取 -> 历史数据索引 + 区块读取
  - LevelDB组合键

- 区块存储
  - 区块以文件块形式存储(blockfile_xxxxxx)
  - 文件块大小: 64M(硬编码)
  - 账本最大容量: 64M * 1000000

- 区块读取
  - 区块文件流(blockfileStream)
  - 区块流(blockStream)
  - 区块迭代器(blocksItr)

- 区块索引
  - 快速定位区块
  - 索引键: 区块高度/区块哈希/交易哈希/...
  - 索引值: 区块文件编码 + 文件内偏移量 + 区块数据长度

- 区块提交
  - 保存区块文件
  - 更新世界状态
  - 更新历史状态(可选)

## 智能合约

- 区块链2.0: 以太坊
- 合约协议的数字化代码表达
- 分布式有限状态机
- 执行环境安全隔离, 不受第三方干扰(EVM, Docker)

- 链码
  - Fabric应用层基石(中间件)
  - 独立的Docker执行环境
  - 背书节点gRPC连接
  - 生命周期管理
  - 编程接口
    - Init()
    - Invoke()

- 生命周期
  - 打包
  - 安装
  - 实例化
  - 升级
  - 交互

## 手动部署

**每次启动前需要彻底把环境清除干净,否则证书验证会失败**  
生成order

```shell
cryptogen generate --config=./organizations/cryptogen/crypto-config-orderer.yaml --output="organizations"
```

生成org1和org2

```shell
cryptogen generate --config=./organizations/cryptogen/crypto-config-org1.yaml --output="organizations"
cryptogen generate --config=./organizations/cryptogen/crypto-config-org2.yaml --output="organizations"
```

可以把三个yaml文件合并成一个,然后一次性全部生成

```shell
cryptogen generate --config=./organizations/cryptogen/crypto-config.yaml --output="organizations"
```

生成创世区块

```shell
# -configPath 为存放 configtx.yaml 的目录
configtxgen -configPath configtx/ -profile TwoOrgsOrdererGenesis -channelID system-channel -outputBlock ./system-genesis-block/genesis.block
```

使用docker启动

```shell
IMAGE_TAG=latest docker-compose -f docker/docker-compose-test-net.yaml up
```

生成账本区块

```shell
CHANNEL_NAME="channel1"
configtxgen -configPath configtx/ -profile TwoOrgsChannel -outputCreateChannelTx ./channel-artifacts/${CHANNEL_NAME}.tx -channelID $CHANNEL_NAME
```

生成锚节点文件

```shell
CHANNEL_NAME="channel1"
configtxgen -configPath configtx/ -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org1MSPanchors.tx -channelID $CHANNEL_NAME -asOrg Org1MSP
configtxgen -configPath configtx/ -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org2MSPanchors.tx -channelID $CHANNEL_NAME -asOrg Org2MSP
```

创建channel

```shell
export FABRIC_CFG_PATH=$PWD/../config/

CHANNEL_NAME="channel1"
ORDERER_CA="${PWD}/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051
peer channel create -o localhost:7050 -c $CHANNEL_NAME --ordererTLSHostnameOverride orderer.example.com -f ./channel-artifacts/${CHANNEL_NAME}.tx --outputBlock ./channel-artifacts/${CHANNEL_NAME}.block --tls --cafile $ORDERER_CA
```

加入channel

```shell
peer channel join -b ./channel-artifacts/channel1.block
```

验证channel是否加入

```shell
peer channel getinfo -c channel1
```

从排序服务中获取块

```shell
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org2MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/crypto-config/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/crypto-config/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
export CORE_PEER_ADDRESS=localhost:9051

peer channel fetch 0 ./channel-artifacts/channel_org2.block -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com -c channel1 --tls --cafile ${PWD}/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
# 该命令使用0来指定它需要获取加入通道所需的创世块
```

配置锚节点

```shell
export FABRIC_CFG_PATH=$PWD/../config/
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051

#获取通道配置
peer channel fetch config channel-artifacts/config_block.pb -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com -c channel1 --tls --cafile ${PWD}/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

# 将protobuf 转换为 json
configtxlator proto_decode --input config_block.pb --type common.Block --output config_block.json
jq .data.data[0].payload.data.config config_block.json > config.json

# 将Org1的peer锚节点添加到通道中
jq '.channel_group.groups.Application.groups.Org1MSP.values += {"AnchorPeers":{"mod_policy": "Admins","value":{"anchor_peers": [{"host": "peer0.org1.example.com","port": 7051}]},"version": "0"}}' config_copy.json > modified_config.json

# 将修改后的通道配置转换回protobuf中, 并计算差异
configtxlator proto_encode --input config.json --type common.Config --output config.pb
configtxlator proto_encode --input modified_config.json --type common.Config --output modified_config.pb
configtxlator compute_update --channel_id channel1 --original config.pb --updated modified_config.pb --output config_update.pb

# 将配置更新包装到交易Envelope中
configtxlator proto_decode --input config_update.pb --type common.ConfigUpdate --output config_update.json
echo '{"payload":{"header":{"channel_header":{"channel_id":"channel1", "type":2}},"data":{"config_update":'$(cat config_update.json)'}}}' | jq . > config_update_in_envelope.json
configtxlator proto_encode --input config_update_in_envelope.json --type common.Envelope --output config_update_in_envelope.pb

# 向通道中更新锚节点
peer channel update -f channel-artifacts/config_update_in_envelope.pb -c channel1 -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
```

## 链码相关

打包链码

```shell
peer lifecycle chaincode package basic.tar.gz --path ../asset-transfer-basic/chaincode-go/ --lang golang --label basic_1.0
```

安装链码

```shell
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051

peer lifecycle chaincode install basic.tar.gz
```

批准链码

```shell
# 查看已安装的链码
peer lifecycle chaincode queryinstalled
# [ own ] ./bin/peer lifecycle chaincode queryinstalled
# Installed chaincodes on peer:
# Package ID: basic_1.0:a7c798f83c9a4bb2316e4d83ab64f41c44ed3c6a77f2db6ad650c4a6971bab0d, Label: basic_1.0

export CC_PACKAGE_ID=basic_1.0:a7c798f83c9a4bb2316e4d83ab64f41c44ed3c6a77f2db6ad650c4a6971bab0d
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID channel1 --name basic --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile ${PWD}/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
# 2020-11-11 13:58:22.993 CST [chaincodeCmd] ClientWait -> INFO 001 txid [b15ff7dfdf87de677ae9d38601936083d5362b6bf976103835d305f7e6a1a9b1] committed with status (VALID) at localhost:9051
```

提交链码到通道上

```shell
# 检查提交是否准备就绪
peer lifecycle chaincode checkcommitreadiness --channelID channel1 --name basic --version 1.0 --sequence 1 --tls --cafile ${PWD}/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem --output json
# 提交
peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID channel1 --name basic --version 1.0 --sequence 1 --tls --cafile ${PWD}/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/crypto-config/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
# 查看已提交链码
peer lifecycle chaincode querycommitted --channelID channel1 --name basic --cafile ${PWD}/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
```

与链码交互

```shell
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C channel1 -n basic --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/crypto-config/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"function":"InitLedger","Args":[]}'
```

升级链码

```shell
export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=$PWD/../config/
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
peer lifecycle chaincode package basic_2.tar.gz --path ../asset-transfer-basic/chaincode-javascript/ --lang golang --label basic_2.0

export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051

peer lifecycle chaincode install basic_2.tar.gz

peer lifecycle chaincode queryinstalled
# Installed chaincodes on peer:
# Package ID: basic_1.0:69de748301770f6ef64b42aa6bb6cb291df20aa39542c3ef94008615704007f3, Label: basic_1.0
# Package ID: basic_2.0:1d559f9fb3dd879601ee17047658c7e0c84eab732dca7c841102f20e42a9e7d4, Label: basic_2.0

export NEW_CC_PACKAGE_ID=basic_2.0:1d559f9fb3dd879601ee17047658c7e0c84eab732dca7c841102f20e42a9e7d4

peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name basic --version 2.0 --package-id $NEW_CC_PACKAGE_ID --sequence 2 --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

# 对Org2进行同样的操作

peer lifecycle chaincode checkcommitreadiness --channelID mychannel --name basic --version 2.0 --sequence 2 --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem --output json
# {
        # "Approvals": {
                # "Org1MSP": true,
                # "Org2MSP": true
        # }
# }

peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name basic --version 2.0 --sequence 2 --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
```

## 智能合约相关

实现两个接口  
shim.ChaincodeStubInterface用法[shim.ChaincodeStubInterface](https://www.cnblogs.com/studyzy/p/7360733.html)  

```golang
type sampleAsset struct{}
func (s *sampleAsset)Init(stub shim.ChaincodeStubInterface) peer.Response{}
func (s *sampleAsset)Invoke(stub shim.ChaincodeStubInterface,args []string)peer.Response{
  // 解析函数名
  fn,args := stub.GetFunctionAndParameters()
  // 验证函数名
  if fn == "set"{
    return set(stub,args)
  }else if fn == "get"{
    return get(stub,args)
  }else if fn == "saveStruct"{
    return saveStruct(stub,args)
  }else if fn == "getByIdAndName"{
    return getByIdAndName(stub,args)
  }
  return shim.Error("error function name")
}
```

实现函数

最简单的存储

```golang
func (s *sampleAsset)set(stub shim.ChaincodeStubInterface,args []string) peer.Response{
  // 判断参数数量是否正确
  if len(args) != 2{
    return shim.Error("incorrect arguments")
  }
  // 存入
  err := stub.PutState(args[0],[]byte{args[1]})
  if err != nil{
    return shim.Error(err.Error())
  }
  return shim.Success(nil)
}
```

最简单的查询

```golang
func (s *sampleAsset)(stub shim.ChaincodeStubInterface,args []string) peer.Response{
  if len(args) != 1{
    return shim.Error("incorrect arguments")
  }
  value,err := stub.GetState(args[0])
  if err != nil{
    return shim.Error(err.Error())
  }
  if value == nil{
    return shim.Error("result is nil")
  }
  return shim.Success(string(value))
}
```

存一个结构体

```golang
type people struct {
  Id   string `json:"id"`
  Name string `json:"name"`
}

func (s *sampleAsset)structSave(stub shim.ChaincodeStubInterface,args []string) peer.Response{
  if len(args) != 2{
    return shim.Error("incorrect arguments")
  }
  people,err := json.Marshal(&people{Id:args[0],Name:args[1]})
  if err != nil{
    return shim.Error(err.Error())
  }

  err = stub.PutState(args[0],people)
  if err != nil{
    return shim.Error(err.Error())
  }

  // 这里的存储是为组合键查询准备的
  indexName := "id~name"
  p := people{Id:args[0],Name:args[1]}
  idNameIndexKey,err := stub.CreateCompositeKey(indexName,[]string{args[0],args[1]})
  if err != nil{
    return shim.Error(err.Error())
  }
  err = stub.PutState(idNameIndexKey,[]byte{0x00})
  if err != nil{
    return shim.Error(err.Error())
  }
  
  return shim.Success(nil)
}
```

根据组合键查询

```golang
func (s *sampleAsset)getByIdAndName(stub shim.ChaincodeStubInterface,args []string) peer.Response{
  if len(args) != 2{
    return shim.Error("incorrect arguments")
  }
  // 要同时配合存储/删除的使用
  // GetStateByPartialCompositeKey方法是采用一种前缀匹配的方法来进行键的匹配返回的。
  // 也就是说，我们虽然是部分复合键的查询，但是只能拿前面的复合键进行匹配，而不是后面部分。
  // 具体来说当你有一个 出场年份～颜色～车号的索引时只能使用 年份、年份与颜色来进行查询，而不能用颜色来进行查询。
  // 因此当我们有多键的复合主键时，各个键的顺序可能需要我们仔细思考一下。
  idNameResultIterator,err := stub.GetStateByPartialCompositeKey(indexName,[]string{args[0]})
  if err != nil{
    return shim.Error(err.Error())
  }
  defer idNameResultIterator.Close()
  res := bytes.Buffer{}
  res.WriteString("[")
  for idNameResultIterator.HasNext(){
    idNameKey,err := idNameResultIterator.Next()
    if err != nil{
      return shim.Error(err.Error())
    }
    objectType,composite,err := stub.SpiltCompositeKey(idNameKey.Key)
    if err != nil{
      return shim.Error(err.Error())
    }
    returnId := composite[0]
    returnName := composite[1]
    pState,err := stub.GetState(returnId)
    if err != nil{
      return shim.Error(err.Error())
    }
    res.Write(pState)
    if idNameResultIterator.HasNext(){
      res.WriteString(",")
    }
  }
  res.WriteString("]")
  return shim.Success(res.Bytes())
}
```

couchDB相关  

语法 [mango](https://github.com/cloudant/mango)  

项目文件, 要在每个智能合约的项目目录 `${PWD}/META-INF/statedb/couchdb/indexes` 加入couchDB的索引文件  

文件例子:  
indexOwner.json

```json
{
  "index": {
    "fields": ["id","name"]
  },
  "ddoc": "indexOwnerDoc",
  "name": "cdbt",
  "type": "json"
}
// 需要以下三条信息来定义一个索引：
//         fields: 这些是常用的查询字段
//         name: 索引名
//         type: 它的内容一般是 json

// 可选地，设计文档（ design document ）属性 ddoc 可以写在索引的定义中。
// design document 是 CouchDB 结构,用于包含索引。
// 索引可以以组的形式定义在设计文档中以提升效率，但是 CouchDB 建议每一个设计文档包含一个索引。
// 如上
// 如果设计文档 indexOwnerDoc 不存在，当索引部署的时候会自动创建 一个。


// go 中的结构体
// type CDBT struct {
//    Id string `json:"id"`
//    Name string `json:"name"`
//    Color string `json:"color"`
//    Size int `json:"size"`
// }

```

智能合约中couchDB独占的方法

```golang
func (s *SampleChaincode) simpleRichQuery(stub shim.ChaincodeStubInterface,args []string) peer.Response{
	if len(args) != 1 {
		return shim.Error("incorrect arguments")
	}
	data,err := stub.GetQueryResult(fmt.Sprintf("{\"selector\":{\"color\":\"%s\"}}",args[0]))
	if err != nil{
		return shim.Error(err.Error())
	}

	res := bytes.NewBufferString("[")

	for data.HasNext(){
		x,err := data.Next()
		if err != nil{
			return shim.Error(err.Error())
		}
		res.WriteString("{\"key\":")
		res.WriteString(x.Key)
		res.WriteString(",\"record\":")
		res.Write(x.Value)
		res.WriteString("}")
		if data.HasNext() {
			res.WriteString(",")
		}
	}
	res.WriteString("]")
	return shim.Success(res.Bytes())
}
```

智能合约访问控制

链码可以通过调用 GetCreator() 方法来使用客户端（提交者）证书进行访问控制决策。  
另外，Go shim 提供了扩展 API ，用于从提交者的证书中提取客户端标识用于访问控制决策，该证书可以是客户端身份本身，或这组织身份，或客户端身份属性。  
例如，一个以键-值对表示的资产可以将客户端的身份作为值的一部分保存其中（比如以 JSON 属性标识资产主人），以后就只有被授权的客户端才可以更新键-值对。  

完整示例 [Client Identity Chaincode Library](https://github.com/hyperledger/fabric-chaincode-go/blob/master/pkg/cid/README.md)

The client identity chaincode library enables you to write chaincode which makes access control decisions based on the identity of the client (i.e. the invoker of the chaincode). In particular, you may make access control decisions based on any or a combination of the following information associated with the client:

- the client identity's MSP (Membership Service Provider) ID
    an attribute associated with the client identity
    an OU (Organizational Unit) value associated with the client identity

Attributes are simply name and value pairs associated with an identity. For example, email=me@gmail.com indicates an identity has the email attribute with a value of me@gmail.com.

## fabric-sdk-go

配置文件例子 [config_test](https://github.com/hyperledger/fabric-sdk-go/blob/master/test/fixtures/config/config_test.yaml)  

其他配置文件例子 [config](https://github.com/hyperledger/fabric-sdk-go/tree/master/test/fixtures/config)

清晰的sdk-go使用例子[fabric-sdk-go-sample](https://github.com/Shitaibin/fabric-sdk-go-sample/)

完整例子[e2e](https://github.com/hyperledger/fabric-sdk-go/blob/master/test/integration/e2e/end_to_end.go)

总结文档 -> **云象联盟链 sdk-go 使用说明文档.md**
