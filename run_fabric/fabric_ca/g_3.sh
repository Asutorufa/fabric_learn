fabric_ca_server_host=192.168.9.103:7055
export PATH=$PATH:/mnt/shareSSD/code/Fabric/run_fabric/bin 

# 创建目录
mkdir -p organizations/ordererOrganizations/example.com
export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/ordererOrganizations/example.com

# 登录
echo "LOGIN"
fabric-ca-client enroll -u http://admin:adminpw@${fabric_ca_server_host}

# OUs and NodeOUs
echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/192-168-9-103-7055.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/192-168-9-103-7055.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/192-168-9-103-7055.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/192-168-9-103-7055.pem
    OrganizationalUnitIdentifier: orderer'>${PWD}/organizations/ordererOrganizations/example.com/msp/config.yaml

# 向fabric-ca发起注册请求
echo "\nREGISTER"
fabric-ca-client register --id.name orderer --id.secret ordererpw --id.type orderer
fabric-ca-client register --id.name ordererAdmin --id.secret ordererAdminpw --id.type admin

# peer0 msp
echo "\nPEER MSP"
mkdir -p organizations/ordererOrganizations/example.com/orderers
mkdir -p organizations/ordererOrganizations/example.com/orderers/example.com
mkdir -p organizations/ordererOrganizations/example.com/orderers/orderer.example.com

fabric-ca-client enroll -u http://orderer:ordererpw@${fabric_ca_server_host} -M ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp --csr.hosts orderer.example.com --csr.hosts localhost
cp ${PWD}/organizations/ordererOrganizations/example.com/msp/config.yaml ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/config.yaml

# peer0 tls
echo "\nPEER TLS"
fabric-ca-client enroll -u http://orderer:ordererpw@${fabric_ca_server_host}  -M ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/tls --enrollment.profile tls --csr.hosts orderer.example.com --csr.hosts localhost
cp ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/tls/tlscacerts/* ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/tls/ca.crt
cp ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/tls/signcerts/* ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/tls/server.crt
cp ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/tls/keystore/* ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/tls/server.key

# peer0 tlscacerts
echo "\nPEER TLSCACERTS"
mkdir -p ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts
cp ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/tls/tlscacerts/* ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

mkdir -p ${PWD}/organizations/ordererOrganizations/example.com/msp/tlscacerts
cp ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/tls/tlscacerts/* ${PWD}/organizations/ordererOrganizations/example.com/msp/tlscacerts/tlsca.example.com-cert.pem
  
# org2 Admin
echo "\nORG ADMIN"
mkdir -p organizations/ordererOrganizations/example.com/users
mkdir -p organizations/ordererOrganizations/example.com/users/Admin@example.com

fabric-ca-client enroll -u http://ordererAdmin:ordererAdminpw@${fabric_ca_server_host} -M ${PWD}/organizations/ordererOrganizations/example.com/users/Admin@example.com/msp
cp ${PWD}/organizations/ordererOrganizations/example.com/msp/config.yaml ${PWD}/organizations/ordererOrganizations/example.com/users/Admin@example.com/msp/config.yaml

