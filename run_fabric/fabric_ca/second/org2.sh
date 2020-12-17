c_dir=`dirname $(readlink -f $0)`
fabric_ca_server_host=192.168.9.103:7054
peerName=peer_g2
peerPass=peer_g2pw
userName=peer_g2_user
userPass=peer_g2_userpw
adminName=org2admin
adminPass=org2adminpw

export PATH=$PATH:${c_dir}/../../bin 

# 创建目录
mkdir -p ${c_dir}/peerOrganizations/org2.example.com/
export FABRIC_CA_CLIENT_HOME=${c_dir}/peerOrganizations/org2.example.com/

# 登录
echo -e "\nLOGIN"
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
    OrganizationalUnitIdentifier: orderer'>${c_dir}/peerOrganizations/org2.example.com/msp/config.yaml

# 向fabric-ca发起注册请求
echo -e "\nREGISTER"
fabric-ca-client register --id.name ${peerName} --id.secret ${peerPass} --id.type peer
fabric-ca-client register --id.name ${userName} --id.secret ${userPass} --id.type client 
fabric-ca-client register --id.name ${adminName} --id.secret ${adminPass} --id.type admin 

# peer0 msp
echo "\nPEER MSP"
mkdir -p ${c_dir}/peerOrganizations/org2.example.com/peers
mkdir -p ${c_dir}/peerOrganizations/org2.example.com/peers/peer0.org2.example.com
fabric-ca-client enroll -u http://${peerName}:${peerPass}@${fabric_ca_server_host} -M ${c_dir}/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/msp --csr.hosts peer0.org2.example.com --csr.hosts localhost
cp ${c_dir}/peerOrganizations/org2.example.com/msp/config.yaml ${c_dir}/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/msp/config.yaml

# peer0 tls
echo "\nPEER TLS"
fabric-ca-client enroll -u http://${peerName}:${peerPass}@${fabric_ca_server_host} -M ${c_dir}/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls --enrollment.profile tls --csr.hosts peer0.org2.example.com --csr.hosts localhost 
cp ${c_dir}/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/tlscacerts/* ${c_dir}/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
cp ${c_dir}/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/signcerts/* ${c_dir}/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/server.crt
cp ${c_dir}/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/keystore/* ${c_dir}/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/server.key

# peer0 tlscacerts
echo -e "\nPEER TLSCACERTS"
mkdir -p ${c_dir}/peerOrganizations/org2.example.com/msp/tlscacerts
cp ${c_dir}/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/tlscacerts/* ${c_dir}/peerOrganizations/org2.example.com/msp/tlscacerts/ca.crt

# peer0 tlsca
echo -e "\nTLSCA"
mkdir -p ${c_dir}/peerOrganizations/org2.example.com/tlsca
cp ${c_dir}/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/tlscacerts/* ${c_dir}/peerOrganizations/org2.example.com/tlsca/tlsca.org2.example.com-cert.pem

# peer0 ca
echo -e "\nPEER Ca"
mkdir -p ${c_dir}/peerOrganizations/org2.example.com/ca
cp ${c_dir}/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/msp/cacerts/* ${c_dir}/peerOrganizations/org2.example.com/ca/ca.org2.example.com-cert.pem

# org2 user dir
mkdir -p ${c_dir}/peerOrganizations/org2.example.com/users

# org2 user1
echo -e "\nORG USER1"
mkdir -p ${c_dir}/peerOrganizations/org2.example.com/users/User1@org2.example.com
fabric-ca-client enroll -u http://${userName}:${userPass}@${fabric_ca_server_host} -M ${c_dir}/peerOrganizations/org2.example.com/users/User1@org2.example.com/msp 
cp ${c_dir}/peerOrganizations/org2.example.com/msp/config.yaml ${c_dir}/peerOrganizations/org2.example.com/users/User1@org2.example.com/msp/config.yaml

# org2 Admin
echo -e "\nORG ADMIN"
mkdir -p ${c_dir}/peerOrganizations/org2.example.com/users/Admin@org2.example.com

fabric-ca-client enroll -u http://${adminName}:${adminPass}@${fabric_ca_server_host} -M ${c_dir}/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp 
cp ${c_dir}/peerOrganizations/org2.example.com/msp/config.yaml ${c_dir}/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp/config.yaml

