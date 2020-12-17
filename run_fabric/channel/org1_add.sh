c_dir=`dirname $(readlink -f $0)`
echo "DIR -> ${c_dir}"
source ${c_dir}/org1_env.sh
peer channel join -b ${c_dir}/../channel-artifacts/${CHANNEL_NAME}.block
peer channel getinfo -c $CHANNEL_NAME 
