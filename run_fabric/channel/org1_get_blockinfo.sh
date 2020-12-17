MYDIR=`dirname $(readlink -f $0)`
echo "DIR -> ${MYDIR}"
source ${MYDIR}/org1_env.sh
peer channel getinfo -c $CHANNEL_NAME
