#!/bin/bash
export FABRIC_CFG_PATH=$PWD
export SYS_CHANNEL=sys-channel
export CHANNEL_NAME=mychannel

# clean all
rm -rf ./channel-artifacts/*

# Generating Orderer Genesis block
configtxgen -profile SampleMultiNodeEtcdRaft -channelID $SYS_CHANNEL -outputBlock ./channel-artifacts/genesis.block

# Generating channel configuration transaction 'channel.tx'
configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID $CHANNEL_NAME

# Generating anchor peer update for Org1MSP
configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org1MSPanchors.tx -channelID $CHANNEL_NAME -asOrg Org1MSP

# Generating anchor peer update for Org2MSP
configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org2MSPanchors.tx -channelID $CHANNEL_NAME -asOrg Org2MSP
