#!/bin/bash

NAMESPACE=${2:-default}
MANIFESTS=../manifests
NFS_DIR=/data

function wait_pod() {
    if [[ -z $1 || -z $2 ]]; then
        echo "Usage: wait_pod <POD_NAME_PATTERN> <STATUS>" ; exit 1
    fi
    while sleep 0.5; do
        kubectl -n ${NAMESPACE:-default} get pods | awk 'BEGIN{i=1}{if($1~/'$1'/&&$3=="'$2'"){print $1;i=0}}END{exit i}' && break
    done
}

services=(
    orderer0-ordererorg:orderer0-ordererorg.yaml
    orderer1-ordererorg:orderer1-ordererorg.yaml
    orderer2-ordererorg:orderer2-ordererorg.yaml
    peer0-org1:peer0-org1.yaml
    peer0-org2:peer0-org2.yaml
    peer1-org1:peer1-org1.yaml
    peer1-org2:peer1-org2.yaml
    ca-org1:ca-org1.yaml
    ca-org2:ca-org2.yaml
    cli:cli.yaml
)

function start() {
    [[ $NAMESPACE != default ]] && cp -a $NFS_DIR/default $NFS_DIR/$NAMESPACE
    [[ $NAMESPACE != default ]] && kubectl apply -f <(sed s/default/$NAMESPACE/ $MANIFESTS/namespace.yaml)

    for service in $MANIFESTS/pvc/*; do
        kubectl apply -f <(sed s/default/$NAMESPACE/ ${service})
    done

    for service in ${services[*]}; do
       kubectl apply -f <(sed s/default/$NAMESPACE/ $MANIFESTS/${service#*:})
       wait_pod "${service%:*}.*" 'Running' | xargs -i echo {} is Running!
    done

    kubectl -n $NAMESPACE get pods -o wide
}

function delete() {
    for service in ${services[*]}; do
        kubectl delete -f <(sed s/default/$NAMESPACE/ $MANIFESTS/${service#*:})
    done

    for service in $MANIFESTS/pvc/*; do
        kubectl delete -f <(sed s/default/$NAMESPACE/ ${service})
    done

    [[ $NAMESPACE != default ]] && kubectl delete -f <(sed s/default/$NAMESPACE/ $MANIFESTS/namespace.yaml)
    [[ $NAMESPACE != default ]] && rm -rf $NFS_DIR/$NAMESPACE
}

case "$1" in
    start)
        start
    ;;
    delete)
        delete
    ;;
    *)
        echo "$0 start  启动 fabric 网络"
        echo "$0 delete 删除 fabric 网络"
    ;;
esac

