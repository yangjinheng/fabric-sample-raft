#!/bin/bash

export NAMESPACE=baas

NFS_ADDR=$(ifconfig eth0 | awk -F '[: ]+' 'NR==2{print $4}')

kubectl apply -f <(sed "s/{{clusterName}}/$NAMESPACE/ ; s/{{nfsServer}}/$NFS_ADDR/" fabric-explorer.yaml)