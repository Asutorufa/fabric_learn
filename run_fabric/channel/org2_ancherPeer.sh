c_dir=`dirname $(readlink -f $0)`
echo "DIR -> ${c_dir}"
source ${c_dir}/org2_env.sh
ORDERER_CA="${c_dir}/../crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/tls/ca.crt"
peer channel update -o localhost:7050 -c $CHANNEL_NAME -f ${c_dir}/../channel-artifacts/Org2MSPanchors.tx --tls --cafile $ORDERER_CA
