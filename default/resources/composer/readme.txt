# 数据目录需要持久化，cli、playground、reserver 都需要挂载它，需要写入数据注意写入权限
/home/composer/.composer

# 由于 composer 只能运行在 composer 用户下 (1000) 所以下面文件需赋予 composer 用户为属主
export ORG1_CERT=/data/composer/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/admincerts/Admin@org1.example.com-cert.pem
export ORG1_KEY=/data/composer/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/keystore/d9fcc4030479fc24e21480f8a087ce968ed80a0379aa9d4235281b2eaa47f80e_sk

export ORG2_CERT=/data/composer/crypto-config/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp/admincerts/Admin@org2.example.com-cert.pem
export ORG2_KEY=/data/composer/crypto-config/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp/keystore/63d4314eaf8b7145b32e59908be02b42ff8038dcf0057e269d30d4db3ea664c1_sk

# 生成 PeerAdmin*.card 这样的文件 
composer card create -p connection-org1.json -c $ORG1_CERT -k $ORG1_KEY -u PeerAdmin -r PeerAdmin -r ChannelAdmin
composer card create -p connection-org2.json -c $ORG2_CERT -k $ORG2_KEY -u PeerAdmin -r PeerAdmin -r ChannelAdmin

kubectl cp composer-cli-68f644d544-l9hvs:/data/composer/PeerAdmin@fabric-network-org1.card PeerAdmin@fabric-network-org1.card
kubectl cp composer-cli-68f644d544-l9hvs:/data/composer/PeerAdmin@fabric-network-org2.card PeerAdmin@fabric-network-org2.card

# 导入 card
composer card import --file ./PeerAdmin@fabric-network-org1.card
composer card import --file ./PeerAdmin@fabric-network-org2.card

# 安装 baas 应用包，会连接 peer 节点，安装一个智能合约
composer network install --card PeerAdmin@fabric-network-org1 --archiveFile ./baas-network@0.1.5.bna --option npmrcFile=./ali-npmrc
composer network install --card PeerAdmin@fabric-network-org2 --archiveFile ./baas-network@0.1.5.bna --option npmrcFile=./ali-npmrc

# 生成用户证书，会连接 fabric-ca 服务器，申请证书
composer identity request --card PeerAdmin@fabric-network-org1 --user admin --enrollSecret adminpw --path alice
composer identity request --card PeerAdmin@fabric-network-org2 --user admin --enrollSecret adminpw --path bob

# 启动业务网络，这一步会部署链码到 Fabric 网络中
composer network start --card PeerAdmin@fabric-network-org1 --networkName baas-network --networkVersion 0.1.5 --option endorsementPolicyFile=./endorsement-policy.json --networkAdmin alice --networkAdminCertificateFile ./alice/admin-pub.pem --networkAdmin bob --networkAdminCertificateFile ./bob/admin-pub.pem

# 生成 org1 的 card 测试部署网络
composer card create --connectionProfileFile ./connection-org1.json --user alice --businessNetworkName baas-network --certificate ./alice/admin-pub.pem --privateKey ./alice/admin-priv.pem
composer card import --file alice@baas-network.card
composer network ping --card alice@baas-network

# 生成 org2 的 card 测试部署网络
composer card create --connectionProfileFile ./connection-org2.json --user bob --businessNetworkName baas-network --certificate ./bob/admin-pub.pem --privateKey ./bob/admin-priv.pem
composer card import --file bob@baas-network.card
composer network ping --card bob@baas-network

