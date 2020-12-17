c_dir=`dirname $(readlink -f $0)`
source ${c_dir}/org2_env.sh
peer channel getinfo -c $CHANNEL_NAME 
