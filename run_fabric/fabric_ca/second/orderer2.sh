c_dir=`dirname $(readlink -f $0)`
fabric_ca_server_host=192.168.9.103:7054
export PATH=$PATH:${c_dir}/../../bin 

# 创建目录
# mkdir -p ${c_dir}/ordererOrganizations/example.com
# export FABRIC_CA_CLIENT_HOME=${c_dir}/ordererOrganizations/example.com

# 登录
echo -e "\nLOGIN"
# fabric-ca-client enroll -u http://admin:adminpw@${fabric_ca_server_host}

# OUs and NodeOUs
# echo 'NodeOUs:
#   Enable: true
#   ClientOUIdentifier:
#     Certificate: cacerts/192-168-9-103-7054.pem
#     OrganizationalUnitIdentifier: client
#   PeerOUIdentifier:
#     Certificate: cacerts/192-168-9-103-7054.pem
#     OrganizationalUnitIdentifier: peer
#   AdminOUIdentifier:
#     Certificate: cacerts/192-168-9-103-7054.pem
#     OrganizationalUnitIdentifier: admin
#   OrdererOUIdentifier:
#     Certificate: cacerts/192-168-9-103-7054.pem
#     OrganizationalUnitIdentifier: orderer'>${c_dir}/ordererOrganizations/example.com/msp/config.yaml

# 向fabric-ca发起注册请求
echo -e "\nREGISTER"
# fabric-ca-client register --id.name orderer --id.secret ordererpw --id.type orderer
# fabric-ca-client register --id.name ordererAdmin --id.secret ordererAdminpw --id.type admin

# peer0 msp
echo -e "\nPEER MSP"
# mkdir -p ${c_dir}/ordererOrganizations/example.com/orderers
# mkdir -p ${c_dir}/ordererOrganizations/example.com/orderers/example.com
mkdir -p ${c_dir}/ordererOrganizations/example.com/orderers/orderer2.example.com

fabric-ca-client enroll -u http://orderer:ordererpw@${fabric_ca_server_host} -M ${c_dir}/ordererOrganizations/example.com/orderers/orderer2.example.com/msp --csr.hosts orderer2.example.com --csr.hosts localhost
cp ${c_dir}/ordererOrganizations/example.com/msp/config.yaml ${c_dir}/ordererOrganizations/example.com/orderers/orderer2.example.com/msp/config.yaml

# peer0 tls
echo -e "\nPEER TLS"
fabric-ca-client enroll -u http://orderer:ordererpw@${fabric_ca_server_host}  -M ${c_dir}/ordererOrganizations/example.com/orderers/orderer2.example.com/tls --enrollment.profile tls --csr.hosts orderer2.example.com --csr.hosts localhost
cp ${c_dir}/ordererOrganizations/example.com/orderers/orderer2.example.com/tls/tlscacerts/* ${c_dir}/ordererOrganizations/example.com/orderers/orderer2.example.com/tls/ca.crt
cp ${c_dir}/ordererOrganizations/example.com/orderers/orderer2.example.com/tls/signcerts/* ${c_dir}/ordererOrganizations/example.com/orderers/orderer2.example.com/tls/server.crt
cp ${c_dir}/ordererOrganizations/example.com/orderers/orderer2.example.com/tls/keystore/* ${c_dir}/ordererOrganizations/example.com/orderers/orderer2.example.com/tls/server.key

# peer0 tlscacerts
echo -e "\nPEER TLSCACERTS"
mkdir -p ${c_dir}/ordererOrganizations/example.com/orderers/orderer2.example.com/msp/tlscacerts
cp ${c_dir}/ordererOrganizations/example.com/orderers/orderer2.example.com/tls/tlscacerts/* ${c_dir}/ordererOrganizations/example.com/orderers/orderer2.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

# mkdir -p ${c_dir}/ordererOrganizations/example.com/msp/tlscacerts
# cp ${c_dir}/ordererOrganizations/example.com/orderers/orderer.example.com/tls/tlscacerts/* ${c_dir}/ordererOrganizations/example.com/msp/tlscacerts/tlsca.example.com-cert.pem
  
# org2 Admin
echo -e "\nORG ADMIN"
# mkdir -p ${c_dir}/ordererOrganizations/example.com/users
# mkdir -p ${c_dir}/ordererOrganizations/example.com/users/Admin@example.com

# fabric-ca-client enroll -u http://ordererAdmin:ordererAdminpw@${fabric_ca_server_host} -M ${c_dir}/ordererOrganizations/example.com/users/Admin@example.com/msp
# cp ${c_dir}/ordererOrganizations/example.com/msp/config.yaml ${c_dir}/ordererOrganizations/example.com/users/Admin@example.com/msp/config.yaml

