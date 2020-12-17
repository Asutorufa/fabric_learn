c_dir=`dirname $(readlink -f $0)`
${c_dir}/generate.sh
${c_dir}/create_channel.sh
${c_dir}/org1_add.sh
${c_dir}/org2_add.sh
${c_dir}/org1_ancherPeer.sh
${c_dir}/org2_ancherPeer.sh
sleep 2s
${c_dir}/org1_get_blockinfo.sh
${c_dir}/org2_get_blockinfo.sh
