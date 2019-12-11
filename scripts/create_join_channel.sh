#!/bin/bash
# Author: YangJinheng

# ------------------------------------------------------------ 变量定义 ------------------------------------------------------------
# Fabric 所在名称空间
NAMESPACE=default

# 接收正则表达式的 POD 名称，等待 POD 状态，返回 POD 名称，用法：pod=`wait_pod '.*hl-composer-cli.*' 'Running'`
function wait_pod() {
    if [[ -z $1 || -z $2 ]]; then
        echo "Usage: wait_pod <POD_NAME_PATTERN> <STATUS>" ; exit 1
    fi
    while sleep 0.5; do
        kubectl -n ${NAMESPACE:-default} get pods | awk 'BEGIN{i=1}{if($1~/'$1'/&&$3=="'$2'"){print $1;i=0}}END{exit i}' && break
    done
}

# get cli name
CLI=`wait_pod 'cli.*' 'Running'`

CORE_PEER_TLS_ENABLED=true

# ------------------------------------------------------------ create-channel ------------------------------------------------------------
# Channel creation
commands='
export FABRIC_LOGGING_SPEC=INFO
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID=Org1MSP
export CORE_PEER_ADDRESS=peer0-org1:7051
export CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export ORDERER_CA=/etc/hyperledger/fabric/crypto-config/ordererOrganizations/example.com/orderers/orderer0.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
peer channel create --channelID mychannel --orderer orderer0:7050 --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA --file resources/channel-artifacts/channel.tx --outputBlock resources/channel-artifacts/channel.block
'
kubectl -n $NAMESPACE exec $CLI -- bash -c "$commands"


# peer0.org1 channel join
commands='
export FABRIC_LOGGING_SPEC=INFO
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID=Org1MSP
export CORE_PEER_ADDRESS=peer0-org1:7051
export CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export ORDERER_CA=/etc/hyperledger/fabric/crypto-config/ordererOrganizations/example.com/orderers/orderer0.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
peer channel join --blockpath resources/channel-artifacts/channel.block
peer channel update --channelID mychannel --orderer orderer0:7050 --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA --file resources/channel-artifacts/Org1MSPanchors.tx
'
kubectl -n $NAMESPACE exec $CLI -- bash -c "$commands"


# peer1.org1 channel join
commands='
export FABRIC_LOGGING_SPEC=INFO
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID=Org1MSP
export CORE_PEER_ADDRESS=peer1-org1:7051
export CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/crypto-config/peerOrganizations/org1.example.com/peers/peer1.org1.example.com/tls/ca.crt
ORDERER_CA=/etc/hyperledger/fabric/crypto-config/ordererOrganizations/example.com/orderers/orderer0.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
peer channel join --blockpath resources/channel-artifacts/channel.block
'
kubectl -n $NAMESPACE exec $CLI -- bash -c "$commands"


# peer0.org2 channel join
commands='
export FABRIC_LOGGING_SPEC=INFO
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID=Org2MSP
export CORE_PEER_ADDRESS=peer0-org2:7051
export CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/crypto-config/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
export CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/crypto-config/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
export ORDERER_CA=/etc/hyperledger/fabric/crypto-config/ordererOrganizations/example.com/orderers/orderer0.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
peer channel join --blockpath resources/channel-artifacts/channel.block
peer channel update --channelID mychannel --orderer orderer0:7050 --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA --file resources/channel-artifacts/Org2MSPanchors.tx
'
kubectl -n $NAMESPACE exec $CLI -- bash -c "$commands"


# peer0.org2 channel join
commands='
export FABRIC_LOGGING_SPEC=INFO
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID=Org2MSP
export CORE_PEER_ADDRESS=peer1-org2:7051
export CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/crypto-config/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
export CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/crypto-config/peerOrganizations/org2.example.com/peers/peer1.org2.example.com/tls/ca.crt
export ORDERER_CA=/etc/hyperledger/fabric/crypto-config/ordererOrganizations/example.com/orderers/orderer0.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
peer channel join --blockpath resources/channel-artifacts/channel.block
'
kubectl -n $NAMESPACE exec $CLI -- bash -c "$commands"
