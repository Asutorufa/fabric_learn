export PATH=$PATH:../bin
export FABRIC_CFG_PATH=../
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org2MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/../crypto-config/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/../crypto-config/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051
ORDERER_CONTAINER=localhost:7050
CH_NAME=channel1
TLS_ROOT_CA="crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/tls/ca.crt"

# 获取配置块
peer channel fetch config config_block_2_2.pb -o $ORDERER_CONTAINER \
    -c $CH_NAME --tls --cafile $TLS_ROOT_CA

# 转换为json
configtxlator proto_decode --input config_block_2_2.pb --type common.Block \
    --output config_block_2_2.json

