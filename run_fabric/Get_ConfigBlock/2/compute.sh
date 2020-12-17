export PATH=$PATH:../../bin
export FABRIC_CFG_PATH=../../
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/../../crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
#export CORE_PEER_MSPCONFIGPATH=${PWD}/../crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/msp
export CORE_PEER_MSPCONFIGPATH=${PWD}/../../crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051
ORDERER_CONTAINER=localhost:7050
CH_NAME=channel1
# TLS_ROOT_CA="../crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/tlscacerts/tlsca.org1.example.com-cert.pem"
# TLS_ROOT_CA="../crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/tlscacerts/tls-192-168-9-103-7055.pem"


# export CORE_PEER_TLS_ENABLED=true
# export CORE_PEER_LOCALMSPID="OrdererMSP"
# export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/../../crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/tls/tlscacerts/tls-192-168-9-103-7055.pem
# export CORE_PEER_MSPCONFIGPATH=${PWD}/../../crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/msp/
# export CORE_PEER_MSPCONFIGPATH=${PWD}/../../crypto-config/ordererOrganizations/example.com/users/Admin@example.com/msp/
# export CORE_PEER_ADDRESS=localhost:7050
# ORDERER_CONTAINER=localhost:7050
# CH_NAME=channel1
# TLS_ROOT_CA="crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/tlscacerts/tlsca.org1.example.com-cert.pem"
#TLS_ROOT_CA="crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/tlscacerts/tls-192-168-9-103-7055.pem"


configtxlator proto_encode --input config.json --type common.Config --output config.pb
configtxlator proto_encode --input modified_config.json --type common.Config --output modified_config.pb
configtxlator compute_update --channel_id $CH_NAME --original config.pb --updated modified_config.pb --output config_update.pb 

configtxlator proto_decode --input config_update.pb --type common.ConfigUpdate --output config_update.json
echo '{"payload":{"header":{"channel_header":{"channel_id":"'$CH_NAME'", "type":2}},"data":{"config_update":'$(cat config_update.json)'}}}' | jq . > config_update_in_envelope.json
configtxlator proto_encode --input config_update_in_envelope.json --type common.Envelope --output config_update_in_envelope.pb
peer channel update -f config_update_in_envelope.pb -c $CH_NAME -o $ORDERER_CONTAINER --tls --cafile $CORE_PEER_TLS_ROOTCERT_FILE
