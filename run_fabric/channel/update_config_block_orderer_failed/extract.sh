c_dir=`dirname $(readlink -f $0)`

jq .data.data[0].payload.data.config ${c_dir}/config_block.json > ${c_dir}/config.json 
cp ${c_dir}/config.json ${c_dir}/modified_config.json
