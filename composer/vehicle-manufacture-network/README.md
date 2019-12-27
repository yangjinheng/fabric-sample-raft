# Vehicle Manufacture Network

这个网络跟踪车辆的制造，从最初的订单请求到制造商的完成。监管者能够在整个过程中提供监督。

## 模型定义

### Participants（参与者）

| 参与者       | 说明       |
| ------------ | ---------- |
| PrivateOwner | 私人       |
| Manufacturer | 汽车制造商 |
| Regulator    | 监管机构   |

### Assets（资产）

| 资产    | 说明 |
| ------- | ---- |
| Order   | 订单 |
| Vehicle | 车辆 |

### Transactions（事务）

| 事务              | 说明            |
| ----------------- | --------------- |
| PlaceOrder        | 下订单          |
| UpdateOrderStatus | 升级订单状态    |
| SetupDemo         | 用于测试的 Demo |

### Events（事件）

| 事件                   | 说明         |
| ---------------------- | ------------ |
| PlaceOrderEvent        | 下订单事件   |
| UpdateOrderStatusEvent | 升级订单状态 |

## 使用说明

### 用户创建订单

1. Person 使用制造商的应用程序来制造他们想要的汽车并订购它。
2. 应用程序提交一个 PlaceOrder（下订单） 事务，该事务创建一个新的 Order（订单） 资产，其中包含 Person 希望为他们生成的车辆的详细信息。

### 制造商生产汽车

1.  Manufacturer（制造商）开始在汽车上工作，当它进入生产阶段时，Manufacturer（制造商）提交 UpdateOrderStatus（更新订单状态）事务来标记 Order（订单）状态的变化，例如将状态从 PLACED（放置）状态更新到 SCHEDULED_FOR_MANUFACTURE（计划生产）。 
2.  一旦 Manufacturer（制造商）完成了 Order（订单）的生产，他们通过使用 VIN_ASSIGNED 状态提交 UpdateOrderStatus（更新订单状态）事务来注册汽车(还提供要注册的VIN车辆识别号)，并使用 Order（订单）中指定的详细信息将 Vehicle（车辆）资产正式添加到注册表中。

### 注册车辆

1.  一旦汽车被注册，那么 Manufacturer（制造商）提交一个状态为 OWNER_ASSIGNED 的 UpdateOrderStatus（更新订单状态）事务，此时 Vehicle（车辆）的 OWNER 字段被设置为与 Order（订单）的 orderer（订购方） 字段匹配。监管机构将对整个过程进行监督。

## 参与者权限

在这个网络权限中，为参与者列出了他们能做什么和不能做什么，权限的分配是比较重要的。

Permissions.acl 文件中的规则显式地允许参与者执行操作。 该文件中没有为参与者编写的操作将被阻止。

### Regulator

RegulatorAdminUser - 允许监管机构对所有资源执行所有操作

### Manufacturer

ManufacturerUpdateOrder - 允许制造商仅使用 UpdateOrderStatus（更新订单状态）来更新订单资产的数据。还必须在 Order asset 中将制造商指定为 <vehicleDetails.make>。

ManufacturerUpdateOrderStatus - 允许制造商创建和读取 UpdateOrderStatus（更新订单状态）事务，这些事务引用指定为 vehicleDetails.make 的订单。

ManufacturerReadOrder - 允许制造商读取在其中指定为 vehicleDetails.make 的订单资产。

ManufacturerCreateVehicle - 允许制造商仅使用 UpdateOrderStatus 事务来创建车辆资产。 事务的 orderStatus 必须是 VIN_ASSIGNED（已分配车辆识别号），并且 Order（订单）资产的制造商必须指定为 vehicleDetails.make

ManufacturerReadVehicle - 允许制造商读取指定为 vehicleDetails.make 的车辆资产。

### Person

PersonMakeOrder - 允许 Person 仅使用 PlaceOrder 事务创建 Order（订单）资产。还必须在 Order（订单）资产中将此人指定为 orderer（订购方）

PersonPlaceOrder - 允许 Person 创建和读取 PlaceOrder 事务的权限，这些事务引用在 order 中指定为 orderer（订购者） 的 order（订单）

PersonReadOrder - 允许 Person 读取指定为 orderer（订购者）的 order（订单）资产

## 测试网络

在 playground 进行下面的测试步骤

### 提交 SetupDemo

导航到 **Test** 选项卡，然后提交一个 SetupDemo 事务，这一步仅仅是为了方便测试而封装了一写批量的操作，这一步可以自行去调用 API 去完成。

```js
{
  "$class": "org.acme.vehicle_network.SetupDemo"
}
```

这将产生：

~~~js
3个 Manufacturer（制造商）participants，
14个 Person（个人）participants，
1个 Regulator（监管者）participants，
13个 Vehicle（车辆）assets。
~~~

### 消费者订购汽车

接下来，通过提交 PlaceOrder transaction 来订购您的汽车（橙色的 Arium Gamora）：

```js
{
  "$class": "org.acme.vehicle_network.PlaceOrder",
  "orderId": "1234",
  "vehicleDetails": {
    "$class": "org.acme.vehicle_network.VehicleDetails",
    "make": "resource:org.acme.vehicle_network.Manufacturer#Arium",
    "modelType": "Gamora",
    "colour": "Sunburst Orange"
  },
  "options": {
    "trim": "executive",
    "interior": "red rum",
    "extras": ["tinted windows", "extended warranty"]
  },
  "orderer": "resource:org.acme.vehicle_network.Person#Paul"
}

// 关键部分注释
{
  "$class": "org.acme.vehicle_network.PlaceOrder",
  "orderId": "1234",                                                          // 订单 ID
  "vehicleDetails": {
    "$class": "org.acme.vehicle_network.VehicleDetails",
    "make": "resource:org.acme.vehicle_network.Manufacturer#Arium",           // 制造商 Arium
    "modelType": "Gamora",                                                    // 型号
    "colour": "Sunburst Orange"                                               // 颜色
  },
  "options": {
    "trim": "executive",                                                      // 装饰，管理人，有执行权的，总统的
    "interior": "red rum",                                                    // 内部，红色朗姆酒
    "extras": ["tinted windows", "extended warranty"]                         // 扩展信息，茶色的窗户，延长保修
  },
  "orderer": "resource:org.acme.vehicle_network.Person#Paul"                  // 订货人，Paul
}
```

这个 PlaceOrder transaction 在注册表中生成一个新的 Order（订单）资产，并发出 PlaceOrderEvent 事件。

### 制造商接受订单

现在，通过提交一个 UpdateOrderStatus transaction 来模拟被制造商接受的订单:

```js
{
  "$class": "org.acme.vehicle_network.UpdateOrderStatus",
  "orderStatus": "SCHEDULED_FOR_MANUFACTURE",
  "order": "resource:org.acme.vehicle_network.Order#1234"
}

// 关键参数注释
{
  "$class": "org.acme.vehicle_network.UpdateOrderStatus",
  "orderStatus": "SCHEDULED_FOR_MANUFACTURE",                                 // 订单状态修改为安排生产 
  "order": "resource:org.acme.vehicle_network.Order#1234"                     // 订单对象为 1234
}
```

这个 UpdateOrderStatus（更新订单状态）transaction 使用资产注册表中的 orderId 1234 更新 Order（订单）的 orderStatus（订单状态），并发出一个 UpdateOrderStatusEvent 事件。

### 制造商向监管机构注册车辆


通过提交 UpdateOrderStatus（更新订单状态）事务模拟 manufacturer（制造商）向监管机构注册车辆：

```js
{
  "$class": "org.acme.vehicle_network.UpdateOrderStatus",
  "orderStatus": "VIN_ASSIGNED",
  "order": "resource:org.acme.vehicle_network.Order#1234",
  "vin": "abc123"
}

// 关键参数注释
{
  "$class": "org.acme.vehicle_network.UpdateOrderStatus",
  "orderStatus": "VIN_ASSIGNED",                                               // 分配车辆识别号
  "order": "resource:org.acme.vehicle_network.Order#1234",                     // 订单为 1234
  "vin": "abc123"                                                              // VIN 车辆识别号
}
```

这个 UpdateOrderStatus（更新订单状态）transaction 使用资产注册表中的 orderId 1234 更新 Order（订单）的 orderStatus（订单状态），根据资产注册表中的 Order（订单）创建一个新的 Vehicle（车辆），并发出一个 UpdateOrderStatusEvent（更新订单状态事件）事件。在此阶段，由于没有为车辆分配所有者，因此将其状态声明为 OFF_THE_ROAD。

### 更新车辆

接下来使用 UpdateOrderStatus transaction 分配车辆的所有者

```js
{
  "$class": "org.acme.vehicle_network.UpdateOrderStatus",
  "orderStatus": "OWNER_ASSIGNED",
  "order": "resource:org.acme.vehicle_network.Order#1234",
  "vin": 'abc123'
}

// 关键参数注释
{
  "$class": "org.acme.vehicle_network.UpdateOrderStatus",
  "orderStatus": "OWNER_ASSIGNED",                                             // 分配车辆的所有者
  "order": "resource:org.acme.vehicle_network.Order#1234",                     // 订单为 1234
  "vin": 'abc123'                                                              // VIN 车辆识别号
}
```

这个 UpdateOrderStatus（更新订单状态）transaction 使用资产注册表中的 orderId 1234 更新 Order（订单）的 orderStatus（订单状态），
使用 VIN（车辆识别号）abc123 更新 Vehicle 资产，使其所有者为 Paul，状态为 ACTIVE，并发出 UpdateOrderStatusEvent（更新订单状态事件）事件。

### 完成车辆交付

最后，通过提交另一个 UpdateOrderStatus（更新订单状态）事务，将订单标记为 DELIVERED，从而完成订购过程：

```js
{
  "$class": "org.acme.vehicle_network.UpdateOrderStatus",
  "orderStatus": "DELIVERED",
  "order": "resource:org.acme.vehicle_network.Order#1234"
}

{
  "$class": "org.acme.vehicle_network.UpdateOrderStatus",
  "orderStatus": "DELIVERED",                                                 // 更新订单状态为已交付
  "order": "resource:org.acme.vehicle_network.Order#1234"                     // 订单为 1234
}
```

这个 UpdateOrderStatus（更新订单状态）transaction 使用资产注册表中的 orderId 1234 更新 Order（订单）的 orderStatus（订单状态），并发出一个 UpdateOrderStatusEvent（更新订单状态事件）事件。

此业务网络定义已用于创建模拟上述场景的演示应用程序。
你可以在网站上找到更多细节：https://github.com/hyperledger/composer-sample-applications/tree/master/packages/vehicle-manufacture

## 测试结果

### Order

过程中产生了一个订单

~~~js
{
  "$class": "org.acme.vehicle_network.Order",
  "orderId": "1234",                                                       // 订单编号
  "vehicleDetails": {
    "$class": "org.acme.vehicle_network.VehicleDetails",
    "make": "resource:org.acme.vehicle_network.Manufacturer#Arium",        // 汽车制造商为 Arium
    "modelType": "Gamora",                                                 // 汽车型号为 Gamora
    "colour": "Sunburst Orange"                                            // 颜色为橙色
  },
  "orderStatus": "DELIVERED",                                              // 订单的状态为：已经交付给订购者
  "options": {
    "$class": "org.acme.vehicle_network.Options",
    "trim": "executive",                                                   // 管理者
    "interior": "red rum",                                                 // 内部信息
    "extras": [                                                            // 一些更详细的信息，例如车的配置，内饰
      "tinted windows",
      "extended warranty"
    ]
  },
  "orderer": "resource:org.acme.vehicle_network.Person#Paul"               // 订购者为 Paul
}
~~~

### Vehicle

过程中生产了一辆汽车

~~~js
{
  "$class": "org.acme.vehicle_network.Vehicle",
  "vin": "abc123",                                                         // 汽车的车辆识别号
  "vehicleDetails": {
    "$class": "org.acme.vehicle_network.VehicleDetails",
    "make": "resource:org.acme.vehicle_network.Manufacturer#Arium",        // 汽车制造商为 Arium
    "modelType": "Gamora",                                                 // 汽车型号为 Gamora
    "colour": "Sunburst Orange"                                            // 颜色为橙色
  },
  "vehicleStatus": "ACTIVE",                                               // 车辆状态为活动
  "owner": "resource:org.acme.vehicle_network.Person#Paul"                 // 车辆的所有者是 Paul
}
~~~

# 部署网络

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
composer network install --card PeerAdmin@fabric-network-org1 --archiveFile ./vehicle-manufacture-network@0.2.6.bna --option npmrcFile=./ali-npmrc
composer network install --card PeerAdmin@fabric-network-org2 --archiveFile ./vehicle-manufacture-network@0.2.6.bna --option npmrcFile=./ali-npmrc
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
composer network start --card PeerAdmin@fabric-network-org1 --networkName vehicle-manufacture-network --networkVersion 0.2.6 --option endorsementPolicyFile=./endorsement-policy.json --networkAdmin alice --networkAdminCertificateFile ./alice/admin-pub.pem --networkAdmin bob --networkAdminCertificateFile ./bob/admin-pub.pem
~~~

5.  生成 org1 的 card 测试部署网络

~~~bash
composer card create --connectionProfileFile ./connection-org1.json --user alice --businessNetworkName vehicle-manufacture-network --certificate ./alice/admin-pub.pem --privateKey ./alice/admin-priv.pem
composer card import --file alice@vehicle-manufacture-network.card
composer network ping --card alice@vehicle-manufacture-network
composer network list --card alice@vehicle-manufacture-network
~~~

6.  生成 org2 的 card 测试部署网络

~~~bash
composer card create --connectionProfileFile ./connection-org2.json --user bob --businessNetworkName vehicle-manufacture-network --certificate ./bob/admin-pub.pem --privateKey ./bob/admin-priv.pem
composer card import --file bob@vehicle-manufacture-network.card
composer network ping --card bob@vehicle-manufacture-network
composer network list --card bob@vehicle-manufacture-network
~~~

7.  升级智能合约

~~~bash
composer network install --card PeerAdmin@fabric-network-org1 --archiveFile ./vehicle-manufacture-network@0.2.6.bna --option npmrcFile=./ali-npmrc
composer network install --card PeerAdmin@fabric-network-org2 --archiveFile ./vehicle-manufacture-network@0.2.6.bna --option npmrcFile=./ali-npmrc

composer network upgrade --card PeerAdmin@fabric-network-org1 --networkName vehicle-manufacture-network --networkVersion 0.0.6
~~~

