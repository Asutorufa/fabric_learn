fabric_ca_server_host=192.168.9.103:7054
peerName=peer_g
peerPass=peer_gpw
userName=peer_g_user
userPass=peer_g_userpw
adminName=org1admin
adminPass=org1adminpw

export PATH=$PATH:/mnt/shareSSD/code/Fabric/run_fabric/bin 

# 创建目录
mkdir -p organizations/peerOrganizations/org1.example.com/
export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/org1.example.com/

# 登录
echo "LOGIN"
fabric-ca-client enroll -u http://admin:adminpw@${fabric_ca_server_host}

# OUs and NodeOUs
echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/192-168-9-103-7054.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/192-168-9-103-7054.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/192-168-9-103-7054.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/192-168-9-103-7054.pem
    OrganizationalUnitIdentifier: orderer' >${PWD}/organizations/peerOrganizations/org1.example.com/msp/config.yaml

# 向fabric-ca发起注册请求
echo "\nREGISTER"
fabric-ca-client register --id.name ${peerName} --id.secret ${peerPass} --id.type peer
fabric-ca-client register --id.name ${userName} --id.secret ${userPass} --id.type client 
fabric-ca-client register --id.name ${adminName} --id.secret ${adminPass} --id.type admin 

# peer0 msp
echo "\nPEER MSP"
mkdir -p organizations/peerOrganizations/org1.example.com/peers
mkdir -p organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com
fabric-ca-client enroll -u http://${peerName}:${peerPass}@${fabric_ca_server_host} -M ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/msp --csr.hosts peer0.org1.example.com --csr.hosts localhost
cp ${PWD}/organizations/peerOrganizations/org1.example.com/msp/config.yaml ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/msp/config.yaml

# peer0 tls
echo "\nPEER TLS"
fabric-ca-client enroll -u http://${peerName}:${peerPass}@${fabric_ca_server_host} -M ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls --enrollment.profile tls --csr.hosts peer0.org1.example.com --csr.hosts localhost 
cp ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/tlscacerts/* ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
cp ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/signcerts/* ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/server.crt
cp ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/keystore/* ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/server.key

# peer0 tlscacerts
echo "\nPEER TLSCACERTS"
mkdir -p ${PWD}/organizations/peerOrganizations/org1.example.com/msp/tlscacerts
cp ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/tlscacerts/* ${PWD}/organizations/peerOrganizations/org1.example.com/msp/tlscacerts/ca.crt

# peer0 tlsca
echo "\nTLSCA"
mkdir -p ${PWD}/organizations/peerOrganizations/org1.example.com/tlsca
cp ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/tlscacerts/* ${PWD}/organizations/peerOrganizations/org1.example.com/tlsca/tlsca.org1.example.com-cert.pem

# peer0 ca
echo "\nPEER Ca"
mkdir -p ${PWD}/organizations/peerOrganizations/org1.example.com/ca
cp ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/msp/cacerts/* ${PWD}/organizations/peerOrganizations/org1.example.com/ca/ca.org1.example.com-cert.pem

# org1 user dir
mkdir -p organizations/peerOrganizations/org1.example.com/users

# org1 user1
echo "\nORG USER1"
mkdir -p organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com
fabric-ca-client enroll -u http://${userName}:${userPass}@${fabric_ca_server_host} -M ${PWD}/organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp 
cp ${PWD}/organizations/peerOrganizations/org1.example.com/msp/config.yaml ${PWD}/organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/config.yaml

# org1 Admin
echo "\nORG ADMIN"
mkdir -p organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com

fabric-ca-client enroll -u http://${adminName}:${adminPass}@${fabric_ca_server_host} -M ${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp 
cp ${PWD}/organizations/peerOrganizations/org1.example.com/msp/config.yaml ${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/config.yaml

