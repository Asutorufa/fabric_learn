./bin/cryptogen generate --config=./crypto-config.yaml 
./bin/configtxgen -configPath ./ -profile TwoOrgsOrdererGenesis -channelID system-channel -outputBlock ./system-genesis-block/genesis.block
IMAGE_TAG=latest COMPOSER_PROJECT_NAME=my docker-compose -f docker-compose-net.yaml up
