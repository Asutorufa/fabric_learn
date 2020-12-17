source $(dirname $(readlink -f $0))/bin.sh
export FABRIC_CFG_PATH=${c_dir}/../
export CHANNEL_NAME="channel1"
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${c_dir}/../crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/server.crt
export CORE_PEER_MSPCONFIGPATH=${c_dir}/../crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/
# export CORE_PEER_MSPCONFIGPATH=${PWD}/crypto-config/peerOrganizations/org1.example.com/msp/
export CORE_PEER_ADDRESS=localhost:7051 
