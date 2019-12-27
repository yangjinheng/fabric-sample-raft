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
