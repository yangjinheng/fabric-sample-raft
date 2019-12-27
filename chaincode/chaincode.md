# 链码 hello world

链码是 Fabric 应用层基石，是应用程序与区块链交互的媒介，也就是说如果我们想基于 Fabric 做一个区块链相关的项目，我们就必须写一个链码，然后将链码部署在 Fabric 上，最后基于 Fabric 提供的 SDK 编写一个应用程序，与部署在区块链上的链码进行交互。

链码的执行环境是一个 Docker 容器，通过 Grpc 接口与背书节点建立连接，也就是只有背书节点才会运行链码。

因为链码是一个 docker 容器，所以它有一整套的生命周期，主要包括以下几个阶段：

1.  打包，将编写的链码整合为一个文件，也可以理解为智能合约的代码编译
2.  安装，这一步仅仅是将链码文件上传到了背书节点
3.  实例化，这一步会执行 Init 方法，主要对链码进行初始化，这个方法只会被执行一次
4.  升级，如果之前的智能合约存在问题，可以通过升级来解决
5.  交互，对区块链的查询和写入

## 链码交互流程

1.  应用程序发起交互请求到背书节点
2.  背书节点收到消息后，会调用容器管理模块查看链码是否正在运行，如果没有启动就编译并启动链码容器
3.  启动后的链码容器会与背书节点建立 Grpc 连接
4.  连接建立好以后背书节点会把应用程序发来的交互请求转发给链码进行执行
5.  链码执行完毕后会返回一个执行结果给背书节点
6.  背书节点收到执行结果后会调用 ESCC 对结果进行签名背书，（ESCC ）
7.  背书节点将签名背书后的结果返回给应用程序

*   ESCC

是一种系统链码，用来完成一些系统功能，虽然是链码但是它运行于节点进程中，目前 Fabric 存在五种：

1.  LSCC（Lifecycle System Chaincode）用于管理链码的生命周期，主要管理安装、实例化、升级
2.  CSCC（Configuraction System Chaincode）配置管理链码，用于管理某个链的配置，比如：新节点加入链
3.  QSCC（Query System Chaincode）可以理解为一个区块存储的 Web 服务
4.  ESCC（Endorsement System Chaincode）交易背书链码，主要用于将交易模拟执行后的结果进行封装然后签名
5.  VSCC（Validation System Chaincode）交易验证，一笔交易被记账节点记录之前，需要完成有效性校验，除了交易读写集验证以外，还需要通过 VSCC 验证，比如：当某些某些节点链码升级以后，使用老版本模拟的交易就是无效的

## 链码开发步骤

1.  构建 Fabric 网络为链码的执行提供载体
2.  根据业务逻辑编写链码
3.  编写单元测试代码，测试链码是否能够正常工作
4.  部署链码到 Fabric 网络进行测试

## 编写链码依赖

~~~go
import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)
~~~

*   github.com/hyperledger/fabric/core/chaincode/shim

~~~
提供了链码与账本交互的中间层，它负责与 peer 节点进行通信，来读取 state 状态，账本等等，并且将自身执行的结果通过 grpc 返回。
~~~

*   pb "github.com/hyperledger/fabric/protos/peer"

~~~
Init 和 Invoke 当方法返回的 pb.response，它提供一种输出的格式，用来传递链码执行的结果，当链码执行完成后，需要通过 protobuf 将执行结果序列化之后传递给 peer 节点。
~~~

## 链码编程接口

~~~go
type Chaincode interface {
    Init(stub ChaincodeStubInterface) pb.Response
    Invoke(stub ChaincodeStubInterface) pb.Response
}
~~~

1.  Init 当链码实例化或者升级的时候被调用

2.  Invoke 当链码收到查询或者调用请求的时候调用

它们接收的 ChaincodeStubInterface 就是 github.com/hyperledger/fabric/core/chaincode/shim 这个中间层。返回值就是 pb "github.com/hyperledger/fabric/protos/peer"

## 链码编程的禁忌

1.  链码的执行是在分布式的环境中，在多个节点中隔离执行的

也就是同一笔交易会被执行很多次，，执行次数取决于背书策略的选择，比如让这条链上的所有节点都执行，或者只有某个组织的某个节点执行。所以链码编程中需要注意：执行的结果不能因为节点的不同而产生不一样的结果，因为客户端会去比较从不同节点返回的交易模拟的结果，如果不同则这笔交易就不会被发送到排序节点进行排序，导致执行结果不同的原因可能有：

*   随机函数
*   系统时间
*   其他不稳定的外部依赖

## 最小的智能合约

~~~go
package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type SimpleAsset struct {
    
}

func (t *SimpleAsset) Init(stub shim.ChaincodeStubInterface) peer.Response {
    
}

func (t *SimpleAsset) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
    
}

func main() {
    if err := shim.Start(new(SimpleAsset)); err != nil {
        // TODO: printf error
    }
}
~~~

## 编写转账链码

*   需要实现的功能

~~~
- 初始化双方数字资产
- 实现双方的转账
- 查询双方的数字资产
- 删除某一方的实体状态
~~~

*   链码代码

~~~go
package main

import (
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode 链码结构体，下面会为 SimpleChaincode 绑定不同的方法
type SimpleChaincode struct{}

// peer chaincode install -n mycc -v 2.0 -p github.com/hyperledger/fabric/peer/resources/chaincodes/chaincode_example02
// peer chaincode instantiate -C mychannel -o orderer0:7050 --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -n mycc -v 2.0 -c '{"Args": ["init","a","100","b","200"]}' -P "OR ('Org1MSP.member','Org2MSP.member')"
// peer chaincode invoke -C mychannel -o orderer0:7050 --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER0_CA --peerAddresses peer0-org1:7051 --tlsRootCertFiles $PEER0_ORG1_CA --peerAddresses peer0-org2:7051 --tlsRootCertFiles $PEER0_ORG2_CA -n mycc -c '{"Args":["invoke","a","b","10"]}'
// peer chaincode query -C mychannel -n mycc -c '{"Args":["query","a"]}'

// Init 初始化链码时候执行的函数
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	// 获得初始化时候的函数名称和参数，即: {"Args": ["init","a","100","b","200"]} 其中：函数名为: init 参数为: "a","100","b","200"
	_, args := stub.GetFunctionAndParameters()

	// 检查参数部分 "a","100","b","200" 数量是否为 4 个
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	// 将传递过来的参数赋值给变量
	var A, B string
	var Aval, Bval int
	var err error

	A = args[0]
	Aval, err = strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("金额错误")
	}

	B = args[2]
	Bval, err = strconv.Atoi(args[3])
	if err != nil {
		return shim.Error("金额错误")
	}

	fmt.Printf("Aval = %d, Bval = %d\n", Aval, Bval)

	// 提交到世界状态
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	// 返回一个成功的响应
	return shim.Success(nil)
}

// Invoke 对链码的增删改查在这个接口实现
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	// 获得 Invoke 时候的操作函数和参数
	fmt.Println("ex02 Invoke")
	function, args := stub.GetFunctionAndParameters()

	// 根据不同函数进行路由
	if function == "invoke" {
		return t.invoke(stub, args)
	} else if function == "delete" {
		return t.delete(stub, args)
	} else if function == "query" {
		return t.query(stub, args)
	}

	// 如果没有匹配到的操作函数，返回错误
	return shim.Error("Invalid invoke function name. Expecting \"invoke\" \"delete\" \"query\"")
}

// delete 实现了一个对实体的删除操作
func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 检查参数数量
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	// 取得要删除的目标
	A := args[0]

	// 从世界状态中删除目标
	err := stub.DelState(A)
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	// 返回删除成功
	return shim.Success(nil)
}

// query 方法实现了查询即 {"Args":["query","a"]}
func (t *SimpleChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var A string
	var err error

	// 检查参数数量 {"Args":["query","a"]} 其中 ["query","a"]
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	// 取得查询目标名称，赋值给 A
	A = args[0]

	// 从世界状态中取出目标
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	// 查询结果返回
	jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success(Avalbytes)
}

// invoke 实现了转账方法 {"Args":["invoke","a","b","10"]}
func (t *SimpleChaincode) invoke(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	// 取得转账的两个人
	var A, B string
	var Aval, Bval int
	var X int
	var err error

	A = args[0]
	B = args[1]

	// 从世界状态中取得 A、B 当前的值
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Avalbytes == nil {
		return shim.Error("Entity not found")
	}
	Aval, _ = strconv.Atoi(string(Avalbytes))

	Bvalbytes, err := stub.GetState(B)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Bvalbytes == nil {
		return shim.Error("Entity not found")
	}
	Bval, _ = strconv.Atoi(string(Bvalbytes))

	// 取得转账的金额
	X, err = strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("Invalid transaction amount, expecting a integer value")
	}

	// 进行转账
	Aval = Aval - X
	Bval = Bval + X
	fmt.Printf("Aval = %d, Bval = %d\n", Aval, Bval)

	//写入世界状态
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
~~~

# 智能合约

在区块链上存储数据，是在每一个节点都有一份完整的数据备份，因此在区块链上存储数据对资源的开销是比较大的，所以我们必须选择性的将数据存储在区块链上，比如说：只是存储真实数据的一个 HASH 值，而且并不是所有的业务数据都可以存储在区块链上，比如说文件、PDF、图片、视频、基本上在现阶段的区块链中不适合存储。

什么样的数据适合存储区块链呢？，可以将区块链类比做数据库，从技术上来说能够被数据库存储的基本上能被区块链存储。

现在通行的做法是将文件存储在另外一个地方，然后在区块链上存储一个唯一性的标识，区块链只是保存一个存证的机制，比如说：将文件存储在 ipfs 上，ipfs 是一个分布式的区中心化的文件系统，也是区块链领域中一个比较火的项目。

*   开发流程

1.  需求整理：提炼出与区块链交互的动作，这与一般的需求分析也没有什么不同。

2.  编写链码：把上一步提炼的交互方法用链码编程接口去实现，实现以后把链码安装到区块链上。
3.  链码交互：安装实例化完成以后，官方提供 cli 只用于开发测试，我们一般会用官方提供的 go-sdk 接口来交互。

总的来说：我们需要编写两个东西，这两个加起来就可以说我们基于 Fabric 开发了一个应用程序。

1.  一个是链码就是智能合约本身，链码负责读取、更新区块链账本，它会产生读写集，去读写状态。
2.  一个是与链码进行交互的服务，它还可以直接读取账本的信息，应用程序与链码的每一次交互都会被包装成一个交易存储到区块链中，如果这比交易被区块链认为是有效的，他们还会对状态数据库产生影响。

## 链码调试

编写好的链码就是一个 go 程序，可以在本地启用链码调试环境，启用方式为：

~~~bash
peer node start --peer-chaincodedev=true    // 命令行方式
CORE_CHAINCODE_MODE=dev                     // 环境变量方式
~~~

## 阿里文档

https://www.alibabacloud.com/help/zh/doc-detail/141367.htm?spm=a2c63.p38356.b99.63.2aea1bc7QutVro

https://github.com/kevin-hf/kongyixueyuan.git（孔壹学院的学历认证智能合约）

## 需求整体

开发一个资产交易（转让）平台

### 资产

某个人拥有一个能够被转让的东西，比如房产、车辆等等，至于能够获得多少收益，这不是区块链负责的，区块链只是负责资产转让的记录。

### 平台功能

1.  用户的开户、销户
2.  资产的登记上链的过程，就是用户绑定资产

3.  资产的转让主要实现资产的所有权的转让，变更
4.  查询功能，用户的查询、资产的查询、用户资产的变更历史查询

### 业务实体

1.  用户

*   名字、

*   标识（身份证）

*   资产列表

2.  资产

* 名字
* 标识
* 特殊属性列表（比如车辆品牌、排量、座位数）

3.  资产变更记录
* 资产标识

* 资产原始拥有者（初次登记时为空）、

* 资产变更后的拥有者

### 交互方法

1.  用户开户
* 参数：1. 名字 2. 标识（身份证）
2.  用户销户

*   参数：1. 标识

3.  资产登录

*   参数：1. 名字 2. 标识 3. 特殊属性列表 4. 拥有者

4.  资产转让

*   参数：1. 拥有者 2. 资产标识 3. 受让者

5.  用户查询

*   参数：标识
*   返回值：用户实体

6.  资产查询

*   参数：1. 标识
*   返回值：资产实体

7.  资产的变更记录

*   参数：1. 资产的标识 2. 记录类型（登记/转让/全部）
*   返回值：资产变更列表

