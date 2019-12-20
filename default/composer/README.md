## 前提条件

1.  必须使用 composer(UID:1000 GID:1000) 这个普通用户完成以下操作
2.  composer 必须运行于 node 8.9 及以上版本，但不支持 node 9
3.  cli、playground、reserver 数据持久化目录都为 /home/composer/.composer
4.  注意 /home/composer/.composer 目录的权限，为 1000.1000
5.  由于 cli、playground、reserver 与 fabric 网络通信时候需要建立 tls 连接，所以需要 tls 证书的访问权限

## 安装开发环境

*   安装环境

~~~bash
yum -y install gcc+ gcc-c++
npm config set registry https://registry.npm.taobao.org
npm install -g --unsafe-perm composer-cli@0.20
npm install -g --unsafe-perm composer-playground@0.20
npm install -g --unsafe-perm composer-rest-server@0.20
npm install -g --unsafe-perm generator-hyperledger-composer@0.20
npm install -g --unsafe-perm yo
~~~

*   创建 composer 用户

~~~bash
useradd -u 1000 composer
~~~

## 创建业务网络

*   一个业务网络定义有如下布局

~~~bash
models/ (optional)            # 模型文件，定义了业务网络的业务领域，模型文件通常由业务分析师创建，因为它们定义模型元素之间的结构和关系：资产，参与者和事务。
lib/                          # JavaScript 文件，包含事务处理器功能，事务处理器功能在Hyperledger Fabric上运行，并可访问存储在Hyperledger Fabric区块链世界状态中的资产注册表。
permissions.acl (optional)    # 访问控制文件，包含一组访问控制规则，用于定义业务网络中不同参与者的权限。
package.json
README.md (optional)
~~~

*   创建一个项目，按照菜单填写

~~~bash
yo hyperledger-composer
==========================================================================
We're constantly looking for ways to make yo better! 
May we anonymously report usage statistics to improve the tool over time? 
More info: https://github.com/yeoman/insight & http://yeoman.io
========================================================================== No
Welcome to the Hyperledger Composer project generator
? Please select the type of project: Business Network
You can run this generator using: 'yo hyperledger-composer:businessnetwork'
Welcome to the business network generator
? Business network name: mynetwork
? Description: This is my test network
? Author name:  Yangjinheng
? Author email: 420123641@qq.com
? License: Apache-2.0
? Namespace: org.example.biznet
? Do you want to generate an empty template network? No: generate a populated sample network
   create package.json
   create README.md
   create models/org.example.biznet.cto
   create permissions.acl
   create .eslintrc.yml
   create features/sample.feature
   create features/support/index.js
   create test/logic.js
   create lib/logic.js
~~~

*   升级应用

~~~bash
composer network upgrade --card PeerAdmin@fabric-network-org1 --archiveFile ./baas-network@0.1.5.bna --option npmrcFile=./ali-npmrc
~~~

## 管理员访问卡

1.  准备文件

| 文件                            | 描述                                       |
| ------------------------------- | ------------------------------------------ |
| connection.json                 | 描述 Fabric 网络的配置文件，具体在下面介绍 |
| Admin@org1.example.com-cert.pem | 组织管理员的证书                           |
| xxxxxxxxxxxxxxxxxxxxxxxxxxxx_sk | 组织管理员证书的私钥                       |

2.  准备每个组织的管理员用户证书和私钥

~~~bash
export ORG1_CERT=/data/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/admincerts/Admin@org1.example.com-cert.pem
export ORG1_KEY=/data/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/keystore/d9fcc4030479fc24e21480f8a087ce968ed80a0379aa9d4235281b2eaa47f80e_sk

export ORG2_CERT=/data/crypto-config/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp/admincerts/Admin@org2.example.com-cert.pem
export ORG2_KEY=/data/crypto-config/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp/keystore/63d4314eaf8b7145b32e59908be02b42ff8038dcf0057e269d30d4db3ea664c1_sk
~~~

3.  准备 connection 文件

~~~bash
https://github.com/hyperledger/composer/blob/master/packages/composer-tests-integration/profiles/tls-connection-org1.json    # 参考格式
~~~

4.  生成管理员的访问卡

~~~bash
composer card create -p connection-org1.json -c $ORG1_CERT -k $ORG1_KEY -u PeerAdmin -r PeerAdmin -r ChannelAdmin
composer card create -p connection-org2.json -c $ORG2_CERT -k $ORG2_KEY -u PeerAdmin -r PeerAdmin -r ChannelAdmin
~~~

5.  导入管理员访问卡

~~~bash
composer card import --file ./PeerAdmin@fabric-network-org1.card
composer card import --file ./PeerAdmin@fabric-network-org2.card
~~~

## 安装业务网络

1.  将编写好的智能合约打包为 bna 格式的归档文件

~~~bash
composer archive create --sourceType dir --sourceName .
~~~

2.  使用管理员访问卡，安装智能合约到每个组织的 peer 节点上

~~~bash
composer network install --card PeerAdmin@fabric-network-org1 --archiveFile ./transfer-example@0.0.2.bna --option npmrcFile=./ali-npmrc
composer network install --card PeerAdmin@fabric-network-org2 --archiveFile ./transfer-example@0.0.2.bna --option npmrcFile=./ali-npmrc
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
composer network start --card PeerAdmin@fabric-network-org1 --networkName transfer-example --networkVersion 0.0.2 --option endorsementPolicyFile=./endorsement-policy.json --networkAdmin alice --networkAdminCertificateFile ./alice/admin-pub.pem --networkAdmin bob --networkAdminCertificateFile ./bob/admin-pub.pem
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

7.  升级智能合约

~~~bash
composer network install --card PeerAdmin@fabric-network-org1 --archiveFile ./transfer-example@0.0.6.bna --option npmrcFile=./ali-npmrc
composer network install --card PeerAdmin@fabric-network-org2 --archiveFile ./transfer-example@0.0.6.bna --option npmrcFile=./ali-npmrc

composer network upgrade --card PeerAdmin@fabric-network-org1 --networkName transfer-example --networkVersion 0.0.6
~~~

