## 前提条件

1.  必须使用 composer(UID:1000 GID:1000) 这个普通用户完成以下操作
2.  composer 可能必须运行于 nodejs 8 这个版本，（待验证）
3.  cli、playground、reserver 数据持久化目录都为 /home/composer/.composer
4.  注意 /home/composer/.composer 目录的权限，为 1000.1000

## 管理员访问卡

1.  修改 composer 目录权限

~~~bash
chown -R 1000.1000 /data/default/composer/
~~~

2.  准备每个组织的管理员用户证书和私钥

~~~bash
export ORG1_CERT=/data/composer/org1/Admin@org1.example.com-cert.pem
export ORG1_KEY=/data/composer/org1/d9fcc4030479fc24e21480f8a087ce968ed80a0379aa9d4235281b2eaa47f80e_sk

export ORG2_CERT=/data/composer/org2/Admin@org2.example.com-cert.pem
export ORG2_KEY=/data/composer/org2/63d4314eaf8b7145b32e59908be02b42ff8038dcf0057e269d30d4db3ea664c1_sk
~~~

3.  生成管理员的访问卡

~~~bash
composer card create -p connection-org1.json -c $ORG1_CERT -k $ORG1_KEY -u PeerAdmin -r PeerAdmin -r ChannelAdmin
composer card create -p connection-org2.json -c $ORG2_CERT -k $ORG2_KEY -u PeerAdmin -r PeerAdmin -r ChannelAdmin
~~~

4.  导入管理员访问卡

~~~bash
composer card import --file ./PeerAdmin@fabric-network-org1.card
composer card import --file ./PeerAdmin@fabric-network-org2.card
~~~

## 安装智能合约

1.  将编写好的智能合约打包为 bna 格式的归档文件

~~~bash
composer archive create --sourceType dir --sourceName .
~~~

2.  使用管理员访问卡，安装智能合约到每个组织的 peer 节点上

~~~bash
composer network install --card PeerAdmin@fabric-network-org1 --archiveFile ./transfer-example@0.0.1.bna --option npmrcFile=./ali-npmrc
composer network install --card PeerAdmin@fabric-network-org2 --archiveFile ./transfer-example@0.0.1.bna --option npmrcFile=./ali-npmrc
~~~

3.  注册业务管理员用户，这一步会到 fabric ca 服务器上注册两个用户

~~~bash
composer identity request --card PeerAdmin@fabric-network-org1 --user admin --enrollSecret adminpw --path alice
composer identity request --card PeerAdmin@fabric-network-org2 --user admin --enrollSecret adminpw --path bob
~~~

4.  指定背书策略启动网络，这一步会实例化智能合约，在每个背书节点上使用 ccenv 作为基础镜像构建一个新的智能合约镜像并运行为容器

~~~bash
--networkName       # 这个值在 package.json 内部的 name 字段做出的规定
--networkVersion    # 这个值在 package.json 内部的 version 字段做出规定
~~~

~~~bash
composer network start --card PeerAdmin@fabric-network-org1 --networkName transfer-example --networkVersion 0.0.1 --option endorsementPolicyFile=./endorsement-policy.json --networkAdmin alice --networkAdminCertificateFile ./alice/admin-pub.pem --networkAdmin bob --networkAdminCertificateFile ./bob/admin-pub.pem
~~~

5.  生成 org1 的 card 测试部署网络

~~~bash
composer card create --connectionProfileFile ./connection-org1.json --user alice --businessNetworkName transfer-example --certificate ./alice/admin-pub.pem --privateKey ./alice/admin-priv.pem
composer card import --file alice@transfer-example.card
composer network ping --card alice@transfer-example
composer network list --card alice@transfer-example
~~~

6.  生成 org2 的 card 测试部署网络

~~~bash
composer card create --connectionProfileFile ./connection-org2.json --user bob --businessNetworkName transfer-example --certificate ./bob/admin-pub.pem --privateKey ./bob/admin-priv.pem
composer card import --file bob@transfer-example.card
composer network ping --card bob@transfer-example
composer network list --card bob@transfer-example
~~~

