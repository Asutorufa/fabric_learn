c_dir=`dirname $(readlink -f $0)`
export FABRIC_CFG_PATH=${c_dir}/../
source ${c_dir}/org1_env.sh
export FABRIC_CFG_PATH=${c_dir}/../
ORDERER_CA="${c_dir}/../crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/tls/ca.crt"
peer channel create -o localhost:7050 -c $CHANNEL_NAME \
    -f ${c_dir}/../channel-artifacts/${CHANNEL_NAME}.tx \
    --outputBlock ${c_dir}/../channel-artifacts/${CHANNEL_NAME}.block \
    --tls --cafile $ORDERER_CA
