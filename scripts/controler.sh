#!/bin/bash
# Author: YangJinheng

# 部署到哪个名称空间
NAMESPACE=${2:-default}

# 部署文件目录
MANIFESTS=../fabric

# NFS 根目录位置
# NFS 配置：/data *(rw,fsid=0,sync,no_subtree_check,no_auth_nlm,insecure,no_root_squash)
NFS_DIR=/data

# 服务名:服务部署文件的部署顺序数组
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

# 等待 PVC 全部绑定
function wait_pvc() {
    printf "Wait for pvc to be ready..."
    while sleep 1; do
        kubectl -n ${NAMESPACE:-default} get --no-headers pvc 2>/dev/null | awk 'BEGIN{i=1}{if($2=="Bound"){i=0}}END{exit i}' && break
    done
    printf "\rWait for pvc to be ready... ok\n"
}

# 等待 POD 到达某状态返回名称
function wait_pod() {
    if [[ -z $1 || -z $2 ]]; then
        echo "Usage: wait_pod <POD_NAME_PATTERN> <STATUS>" ; exit 1
    fi
    while sleep 1; do
        kubectl -n ${NAMESPACE:-default} get --no-headers pods 2>/dev/null | awk 'BEGIN{i=1}{if($1~/'$1'/&&$3=="'$2'"){print $1;i=0}}END{exit i}' && break
    done
}

# 创建名称空间，创建 PVC
function pvc() {
    [[ $NAMESPACE != default ]] && cp -an $NFS_DIR/default $NFS_DIR/$NAMESPACE
    [[ $NAMESPACE != default ]] && kubectl apply -f <(sed s/default/$NAMESPACE/ $MANIFESTS/namespace.yaml)

    for service in $MANIFESTS/pvc/*; do
        kubectl apply -f <(sed s/default/$NAMESPACE/ ${service})
    done

    wait_pvc

    kubectl -n $NAMESPACE get pvc
}

# 启动各个服务，创建 PVC
function start() {
    for service in ${services[*]}; do
       kubectl apply -f <(sed s/default/$NAMESPACE/ $MANIFESTS/${service#*:})
       wait_pod "${service%:*}.*" 'Running' | xargs -i echo {} is Running!
    done

    kubectl -n $NAMESPACE get pods -o wide
}

# 停止各个服务
function stop() {
    for service in ${services[*]}; do
        kubectl delete -f <(sed s/default/$NAMESPACE/ $MANIFESTS/${service#*:})
    done
}

# 停止各个服务，并删除 PVC
function delete() {
    for service in $MANIFESTS/pvc/*; do
        kubectl delete -f <(sed s/default/$NAMESPACE/ ${service})
    done

    [[ $NAMESPACE != default ]] && kubectl delete -f <(sed s/default/$NAMESPACE/ $MANIFESTS/namespace.yaml)
    [[ $NAMESPACE != default ]] && rm -rf $NFS_DIR/$NAMESPACE
}

case "$1" in
    pvc)
        pvc
    ;;
    start)
        pvc
        start
    ;;
    stop)
        stop
    ;;
    delete)
        stop
        delete
    ;;
    *)
        echo "$0 start  启动 fabric 网络"
        echo "$0 stop   停止 fabric 网络"
        echo "$0 pvc    创建 pvc 存储"
        echo "$0 delete 删除 fabric 网络"
    ;;
esac
