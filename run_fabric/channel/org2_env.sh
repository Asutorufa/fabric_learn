source $(dirname $(readlink -f $0))/bin.sh
export FABRIC_CFG_PATH=${c_dir}/../
export CHANNEL_NAME="channel1"
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org2MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${c_dir}/../crypto-config/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${c_dir}/../crypto-config/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp/
export CORE_PEER_ADDRESS=localhost:9051
