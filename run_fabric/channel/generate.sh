source $(dirname $(readlink -f $0))/bin.sh
export CHANNEL_NAME="channel1"

configtxgen -configPath ${c_dir}/../ -profile TwoOrgsChannel -outputCreateChannelTx ${c_dir}/../channel-artifacts/${CHANNEL_NAME}.tx -channelID $CHANNEL_NAME
configtxgen -configPath ${c_dir}/../ -profile TwoOrgsChannel -outputAnchorPeersUpdate ${c_dir}/../channel-artifacts/Org1MSPanchors.tx -channelID $CHANNEL_NAME -asOrg Org1MSP
configtxgen -configPath ${c_dir}/../ -profile TwoOrgsChannel -outputAnchorPeersUpdate ${c_dir}/../channel-artifacts/Org2MSPanchors.tx -channelID $CHANNEL_NAME -asOrg Org2MSP 


