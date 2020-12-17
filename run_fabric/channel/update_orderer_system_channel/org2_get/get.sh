c_dir=`dirname $(readlink -f $0)`
source ${c_dir}/org2_env.sh
ORDERER_CONTAINER=localhost:7050
TLS_ROOT_CA="${c_dir}/../../../crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/tls/server.crt"

cd ${c_dir}
# 获取配置块
peer channel fetch config config_block.pb -o $ORDERER_CONTAINER -c $CHANNEL_NAME --tls --cafile $TLS_ROOT_CA

# 转换为json
configtxlator proto_decode --input config_block.pb --type common.Block --output config_block.json

