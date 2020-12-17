c_dir=$(dirname $(readlink -f $0))
source ${c_dir}/orderer_env.sh
ORDERER_CA="${c_dir}/../../crypto-config/ordererOrganizations/example.com/msp/tlscacerts/tlsca.example.com-cert.pem"
#ORDERER_CA=${c_dir}/../../crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/tls/server.crt
ORDERER_CONTAINER=localhost:7050

cd $c_dir

configtxlator proto_encode --input config.json --type common.Config --output config.pb
configtxlator proto_encode --input modified_config.json --type common.Config --output modified_config.pb
configtxlator compute_update --channel_id $CHANNEL_NAME --original config.pb --updated modified_config.pb --output config_update.pb 

configtxlator proto_decode --input config_update.pb --type common.ConfigUpdate --output config_update.json
echo '{"payload":{"header":{"channel_header":{"channel_id":"'$CHANNEL_NAME'", "type":2}},"data":{"config_update":'$(cat config_update.json)'}}}' | jq . > config_update_in_envelope.json
configtxlator proto_encode --input config_update_in_envelope.json --type common.Envelope --output config_update_in_envelope.pb
peer channel update -f config_update_in_envelope.pb -c $CHANNEL_NAME -o $ORDERER_CONTAINER --tls --cafile $ORDERER_CA
#$CORE_PEER_TLS_ROOTCERT_FILE
