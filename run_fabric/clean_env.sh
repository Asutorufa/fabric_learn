
# Obtain CONTAINER_IDS and remove them
# TODO Might want to make this optional - could clear other containers
# This function is called when you bring a network down
function clearContainers() {
  CONTAINER_IDS=$(docker ps -a | awk '($2 ~ /dev-peer.*/) {print $1}')
  if [ -z "$CONTAINER_IDS" -o "$CONTAINER_IDS" == " " ]; then
    echo "No containers available for deletion"
  else
    docker rm -f $CONTAINER_IDS
  fi
}

# Delete any images that were generated as a part of this setup
# specifically the following images are often left behind:
# This function is called when you bring the network down
function removeUnwantedImages() {
  DOCKER_IMAGE_IDS=$(docker images | awk '($1 ~ /dev-peer.*/) {print $3}')
  if [ -z "$DOCKER_IMAGE_IDS" -o "$DOCKER_IMAGE_IDS" == " " ]; then
    echo "No images available for deletion"
  else
    docker rmi -f $DOCKER_IMAGE_IDS
  fi
} 

clearContainers
removeUnwantedImages

# remove orderer block and other channel configuration transactions and certs
docker run --rm -v $(pwd):/data busybox sh -c 'cd /data && rm -rf system-genesis-block/*.block crypto-config/peerOrganizations crypto-config/ordererOrganizations'
## remove fabric ca artifacts
docker run --rm -v $(pwd):/data busybox sh -c 'cd /data && rm -rf crypto-config/fabric-ca/org1/msp crypto-config/fabric-ca/org1/tls-cert.pem crypto-config/fabric-ca/org1/ca-cert.pem crypto-config/fabric-ca/org1/IssuerPublicKey crypto-config/fabric-ca/org1/IssuerRevocationPublicKey crypto-config/fabric-ca/org1/fabric-ca-server.db'
docker run --rm -v $(pwd):/data busybox sh -c 'cd /data && rm -rf crypto-config/fabric-ca/org2/msp crypto-config/fabric-ca/org2/tls-cert.pem crypto-config/fabric-ca/org2/ca-cert.pem crypto-config/fabric-ca/org2/IssuerPublicKey crypto-config/fabric-ca/org2/IssuerRevocationPublicKey crypto-config/fabric-ca/org2/fabric-ca-server.db'
docker run --rm -v $(pwd):/data busybox sh -c 'cd /data && rm -rf crypto-config/fabric-ca/ordererOrg/msp crypto-config/fabric-ca/ordererOrg/tls-cert.pem crypto-config/fabric-ca/ordererOrg/ca-cert.pem crypto-config/fabric-ca/ordererOrg/IssuerPublicKey crypto-config/fabric-ca/ordererOrg/IssuerRevocationPublicKey crypto-config/fabric-ca/ordererOrg/fabric-ca-server.db'
docker run --rm -v $(pwd):/data busybox sh -c 'cd /data && rm -rf addOrg3/fabric-ca/org3/msp addOrg3/fabric-ca/org3/tls-cert.pem addOrg3/fabric-ca/org3/ca-cert.pem addOrg3/fabric-ca/org3/IssuerPublicKey addOrg3/fabric-ca/org3/IssuerRevocationPublicKey addOrg3/fabric-ca/org3/fabric-ca-server.db'
# remove channel and script artifacts
docker run --rm -v $(pwd):/data busybox sh -c 'cd /data && rm -rf channel-artifacts log.txt *.tar.gz'

docker container prune
docker network prune
docker volume prune
