c_dir=`dirname $(readlink -f $0)`
echo ${c_dir}
export PATH=$PATH:${c_dir}/../../bin
export FABRIC_CFG_PATH=${c_dir}/../../
#export CORE_PEER_TLS_ENABLED=true
#export CORE_PEER_LOCALMSPID="Org1MSP"
#export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/../crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
#export CORE_PEER_MSPCONFIGPATH=${PWD}/../crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
#export CORE_PEER_ADDRESS=localhost:7051
#ORDERER_CONTAINER=localhost:7050
#CH_NAME=channel1
#TLS_ROOT_CA="${PWD}/../crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/tlscacerts/tls-192-168-9-103-7055.pem"

#export CORE_PEER_TLS_ENABLED=true
#export CORE_PEER_LOCALMSPID="OrdererMSP"
#export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/../crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/tls/tlscacerts/tls-192-168-9-103-7055.pem
#export CORE_PEER_MSPCONFIGPATH=${PWD}/../crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/msp/
#export CORE_PEER_MSPCONFIGPATH=${PWD}/../crypto-config/ordererOrganizations/example.com/users/Admin@example.com/msp/
#export CORE_PEER_ADDRESS=localhost:7051
#ORDERER_CONTAINER=localhost:7050
#CH_NAME=channel1
# TLS_ROOT_CA="crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/tlscacerts/tlsca.org1.example.com-cert.pem"
# TLS_ROOT_CA="crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/tlscacerts/tls-192-168-9-103-7055.pem"

export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="OrdererMSP"
# export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/../fabric_ca/third/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
#export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/../fabric_ca/for_change/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
#export CORE_PEER_MSPCONFIGPATH=${PWD}/../crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_MSPCONFIGPATH=${c_dir}/../../crypto-config/ordererOrganizations/example.com/users/Admin@example.com/msp
#export CORE_PEER_MSPCONFIGPATH=${PWD}/../crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/msp
export CORE_PEER_ADDRESS=localhost:7050
ORDERER_CONTAINER=localhost:7050
CH_NAME=channel1
CH_NAME=system-channel
TLS_ROOT_CA="${c_dir}/../../crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt"
TLS_ROOT_CA="${c_dir}/../../crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/tls/ca.crt"

# 获取配置块
peer channel fetch config ${c_dir}/config_block.pb -o $ORDERER_CONTAINER -c $CH_NAME --tls --cafile $TLS_ROOT_CA

# 转换为json
configtxlator proto_decode --input ${c_dir}/config_block.pb --type common.Block --output ${c_dir}/config_block.json

 
