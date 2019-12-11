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

# 等待 pod 运行后取得名称
CLI=`wait_pod 'cli.*' 'Running'`

# 安装智能合约
commands='
export FABRIC_LOGGING_SPEC=INFO
export CORE_PEER_TLS_ENABLED=true

# peer0-org1
export CORE_PEER_LOCALMSPID=Org1MSP
export CORE_PEER_ADDRESS=peer0-org1:7051
export CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
peer chaincode install -n mycc -v 2.0 -p github.com/hyperledger/fabric/peer/resources/chaincodes/chaincode_example02

# peer1-org1
export CORE_PEER_LOCALMSPID=Org1MSP
export CORE_PEER_ADDRESS=peer1-org1:7051
export CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/crypto-config/peerOrganizations/org1.example.com/peers/peer1.org1.example.com/tls/ca.crt
peer chaincode install -n mycc -v 2.0 -p github.com/hyperledger/fabric/peer/resources/chaincodes/chaincode_example02

# peer0-org2
export CORE_PEER_LOCALMSPID=Org2MSP
export CORE_PEER_ADDRESS=peer0-org2:7051
export CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/crypto-config/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
export CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/crypto-config/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
peer chaincode install -n mycc -v 2.0 -p github.com/hyperledger/fabric/peer/resources/chaincodes/chaincode_example02

# peer1-org2
export CORE_PEER_LOCALMSPID=Org2MSP
export CORE_PEER_ADDRESS=peer1-org2:7051
export CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/crypto-config/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
export CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/crypto-config/peerOrganizations/org2.example.com/peers/peer1.org2.example.com/tls/ca.crt
peer chaincode install -n mycc -v 2.0 -p github.com/hyperledger/fabric/peer/resources/chaincodes/chaincode_example02
'
kubectl -n $NAMESPACE exec -it $CLI -- bash -c "$commands"


# 实例化智能合约 => 100
commands=`cat <<'EOF'
export FABRIC_LOGGING_SPEC=INFO
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID=Org1MSP
export CORE_PEER_ADDRESS=peer0-org1:7051
export CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export ORDERER_CA=/etc/hyperledger/fabric/crypto-config/ordererOrganizations/example.com/orderers/orderer0.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
peer chaincode instantiate -C mychannel -o orderer0:7050 --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -n mycc -v 2.0 -c '{"Args": ["init","a","100","b","200"]}' -P "OR ('Org1MSP.member','Org2MSP.member')"
sleep 10
peer chaincode query -C mychannel -n mycc -c '{"Args":["query","a"]}'
EOF
`
kubectl -n $NAMESPACE exec -it $CLI -- bash -c "$commands"

# 交易和查询智能合约 => 70，现实是返回了 90
commands=`cat <<'EOF'
export FABRIC_LOGGING_SPEC=INFO
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID=Org1MSP
export CORE_PEER_ADDRESS=peer0-org1:7051
export CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt

# invoke env
export PEER0_ORG1_CA=/etc/hyperledger/fabric/crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export PEER0_ORG2_CA=/etc/hyperledger/fabric/crypto-config/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
export ORDERER0_CA=/etc/hyperledger/fabric/crypto-config/ordererOrganizations/example.com/orderers/orderer0.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
export ORDERER1_CA=/etc/hyperledger/fabric/crypto-config/ordererOrganizations/example.com/orderers/orderer1.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
export ORDERER2_CA=/etc/hyperledger/fabric/crypto-config/ordererOrganizations/example.com/orderers/orderer2.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

# invoke
peer chaincode invoke -C mychannel -o orderer0:7050 --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER0_CA --peerAddresses peer0-org1:7051 --tlsRootCertFiles $PEER0_ORG1_CA --peerAddresses peer0-org2:7051 --tlsRootCertFiles $PEER0_ORG2_CA -n mycc -c '{"Args":["invoke","a","b","10"]}'
peer chaincode invoke -C mychannel -o orderer1:7050 --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER1_CA --peerAddresses peer0-org1:7051 --tlsRootCertFiles $PEER0_ORG1_CA --peerAddresses peer0-org2:7051 --tlsRootCertFiles $PEER0_ORG2_CA -n mycc -c '{"Args":["invoke","a","b","10"]}'
peer chaincode invoke -C mychannel -o orderer2:7050 --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER2_CA --peerAddresses peer0-org1:7051 --tlsRootCertFiles $PEER0_ORG1_CA --peerAddresses peer0-org2:7051 --tlsRootCertFiles $PEER0_ORG2_CA -n mycc -c '{"Args":["invoke","a","b","10"]}'

# after invoke
sleep 5
peer chaincode query -C mychannel -n mycc -c '{"Args":["query","a"]}'
EOF
`
kubectl -n $NAMESPACE exec -it $CLI -- bash -c "$commands"
