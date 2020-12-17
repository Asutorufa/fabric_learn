export PATH=$PATH:../../bin
export FABRIC_CFG_PATH=../../
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/../../crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/../../crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051
ORDERER_CONTAINER=localhost:7050
CH_NAME=channel1
TLS_ROOT_CA="${PWD}/../../crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/tlscacerts/tls-192-168-9-103-7055.pem"

# 获取配置块
peer channel fetch config config_block.pb -o $ORDERER_CONTAINER -c $CH_NAME --tls --cafile $TLS_ROOT_CA

# 转换为json
configtxlator proto_decode --input config_block.pb --type common.Block --output config_block.json

