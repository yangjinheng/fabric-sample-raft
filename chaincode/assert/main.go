package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// AssetExangeCC 链码结构体，下面会为 AssetExangeCC 绑定不同的方法
type AssetExangeCC struct{}

// 用户第一次登记资产时候，这个资产的原始所有者是不知道的，这里设置一个默认的值
const (
	OriginOwner = "originOwnerPlaceholder"
)

// User 用户
type User struct {
	Name   string   `json:"name"`
	ID     string   `json:"id"`
	Assets []string `json:"assets"`
	// Assets map[string]string `json:"assets"` // {资产id:资产name} // map 是无序的，会造成在不同节点执行失败，不能使用
}

// Asset 资产
type Asset struct {
	Name     string `json:"name"`
	ID       string `json:"id"`
	Metadata string `json:"metadata"`
	// Metadata map[string]string `json:"metadata"` // 特殊属性 // map 是无序的，会造成在不同节点执行失败，不能使用
}

// AssetHistory 资产变更历史
type AssetHistory struct {
	AssetID        string `json:"asset_id"`
	OriginOwnerID  string `json:"origin_owner_id"`  // 资产的原始拥有者
	CurrentOwnerID string `json:"current_owner_id"` // 资产的当前拥有者
}

// 组合键
func constructUserKey(userID string) string {
	return fmt.Sprintf("user_%s", userID)
}

// 组合键
func constructAssetKey(assetID string) string {
	return fmt.Sprintf("asset_%s", assetID)
}

// 组合键
// func constructAssetHistoryKey(originUserID, assetID, currentUserID string) string {
// 	return fmt.Sprintf("history_%s_%s_%s", originUserID, assetID, currentUserID)
// }

// peer chaincode install -n mycc -v 2.0 -p github.com/hyperledger/fabric/peer/resources/chaincodes/chaincode_example02
// peer chaincode instantiate -C mychannel -o orderer0:7050 --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -n mycc -v 2.0 -c '{"Args": ["init","a","100","b","200"]}' -P "OR ('Org1MSP.member','Org2MSP.member')"
// peer chaincode invoke -C mychannel -o orderer0:7050 --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER0_CA --peerAddresses peer0-org1:7051 --tlsRootCertFiles $PEER0_ORG1_CA --peerAddresses peer0-org2:7051 --tlsRootCertFiles $PEER0_ORG2_CA -n mycc -c '{"Args":["invoke","a","b","10"]}'
// peer chaincode query -C mychannel -n mycc -c '{"Args":["query","a"]}'

// 用户开户
func userRegister(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 1. 检查参数的个数
	if len(args) != 2 {
		return shim.Error("参数数量错误")
	}

	// 2. 验证参数的正确性
	name := args[0]
	id := args[1]
	if name == "" || id == "" {
		return shim.Error("参数值不能为空字符")
	}

	// 3. 验证数据是否存在，应该存在 or 不应该存在
	userBytes, err := stub.GetState(constructUserKey(id))
	if err == nil && len(userBytes) != 0 {
		return shim.Error("用户已经存在了")
	}

	// 4.写入状态
	user := &User{
		Name:   name,
		ID:     id,
		Assets: make([]string, 0),
	}
	userBytes, err = json.Marshal(user)
	if err != nil {
		return shim.Error(fmt.Sprintf("序列化出错 %s", err))
	}

	if err := stub.PutState(constructUserKey(id), userBytes); err != nil {
		return shim.Error(fmt.Sprintf("写入账本失败 %s", err))
	}

	return shim.Success(nil)
}

// 用户销户
func userDestroy(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 1. 检查参数的个数
	if len(args) != 1 {
		return shim.Error("参数数量错误")
	}
	// 2. 验证参数的正确性
	id := args[0]
	if id == "" {
		return shim.Error("参数值不能为空字符")
	}
	// 3. 验证数据是否存在，应该存在 or 不应该存在
	userBytes, err := stub.GetState(constructUserKey(id))
	if err != nil && len(userBytes) == 0 {
		return shim.Error("用户不存在")
	}

	// 4. 写入状态，也就是删除用户，和删除用户的资产
	if err = stub.DelState(constructUserKey(id)); err != nil {
		return shim.Error(fmt.Sprintf("删除用户失败 %s", err))
	}
	// 删除资产
	user := new(User)
	if err := json.Unmarshal(userBytes, user); err != nil {
		return shim.Error(fmt.Sprintf("反序列化用户错误 %s", err))
	}
	for _, assetid := range user.Assets {
		if err := stub.DelState(constructAssetKey(assetid)); err != nil {
			return shim.Error(fmt.Sprintf("删除资产失败 %s", err))
		}
	}
	return shim.Success(nil)
}

// 资产注册
func assetEnroll(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 1. 检查参数的个数
	if len(args) != 4 {
		return shim.Error("参数数量错误")
	}
	// 2. 验证参数的正确性
	assetname := args[0]
	assetid := args[1]
	metadata := args[2]
	ownerid := args[3]
	if assetname == "" || assetid == "" || metadata == "" || ownerid != "" {
		return shim.Error("参数值不能为空字符")
	}

	// 3. 验证数据是否存在，应该存在 or 不应该存在
	userBytes, err := stub.GetState(constructUserKey(ownerid))
	if err == nil && len(userBytes) != 0 {
		return shim.Error("用户已经存在了")
	}
	assetBytes, err := stub.GetState(constructAssetKey(assetid))
	if err == nil && len(assetBytes) != 0 {
		return shim.Error("资产已经存在了")
	}

	// 写入状态，1. 写入资产对象，2. 更新用户对象 3. 插入资产变更记录
	// 写入资产对象
	asset := &Asset{
		Name:     assetname,
		ID:       assetid,
		Metadata: metadata,
	}
	assetBytes, err = json.Marshal(asset)
	if err != nil {
		return shim.Error(fmt.Sprintf("序列化资产失败 %s", err))
	}
	if err := stub.PutState(constructAssetKey(assetid), assetBytes); err != nil {
		return shim.Error(fmt.Sprintf("资产写入账本失败 %s", err))
	}
	// 更新用户对象
	user := new(User)
	if err := json.Unmarshal(userBytes, user); err != nil {
		return shim.Error(fmt.Sprintf("反序列化用户失败 %s", err))
	}
	user.Assets = append(user.Assets, assetid)
	if userBytes, err = json.Marshal(user); err != nil {
		return shim.Error(fmt.Sprintf("序列化用户失败 %s", err))
	}
	if err := stub.PutState(constructUserKey(user.ID), userBytes); err != nil {
		return shim.Error(fmt.Sprintf("更新用户写入账本失败 %s", err))
	}
	// 插入资产变更记录，OriginOwner 是一个静态变量用于表示第一次登记资产时候的原始拥有者
	history := &AssetHistory{
		AssetID:        assetid,
		OriginOwnerID:  OriginOwner,
		CurrentOwnerID: ownerid,
	}
	historyBytes, err := json.Marshal(history)
	if err != nil {
		return shim.Error(fmt.Sprintf("序列化资产历史失败 %s", err))
	}
	// 创建组合键
	historyKey, err := stub.CreateCompositeKey("history", []string{OriginOwner, assetid, ownerid})
	if err != nil {
		shim.Error(fmt.Sprintf("组合键制作失败%s", err))
	}
	if err = stub.PutState(historyKey, historyBytes); err != nil {
		shim.Error(fmt.Sprintf("资产历史写入账本失败 %s", err))
	}
	return shim.Success(nil)
}

// 资产转移
func assetExchange(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 1. 检查参数的个数
	if len(args) != 3 {
		return shim.Error("参数数量错误")
	}

	// 2. 验证参数的正确性
	ownerid := args[0]
	assetid := args[1]
	currentownerid := args[0]
	if ownerid == "" || assetid == "" || currentownerid == "" {
		return shim.Error("参数值不能为空字符")
	}
	// 3. 验证数据是否存在，应该存在 or 不应该存在
	OriginOwnerBytes, err := stub.GetState(constructUserKey(ownerid))
	if err != nil || len(OriginOwnerBytes) == 0 {
		return shim.Error("用户不存在")
	}
	currentOwnerBytes, err := stub.GetState(constructUserKey(currentownerid))
	if err != nil || len(currentOwnerBytes) == 0 {
		return shim.Error("用户不存在")
	}
	if assetBytes, err := stub.GetState(constructAssetKey(assetid)); err != nil || len(assetBytes) == 0 {
		return shim.Error("资产不存在")
	}
	// 校验原始拥有者确实拥有当前变更的资产
	originOwner := new(User)
	if err := json.Unmarshal(OriginOwnerBytes, originOwner); err != nil {
		return shim.Error(fmt.Sprintf("反序列化用户失败 %s", err))
	}
	aidexist := false
	for _, aid := range originOwner.Assets {
		if aid == assetid {
			aidexist = true
			break
		}
	}
	if !aidexist {
		return shim.Error("原始拥有者没有这个资产")
	}

	// 4. 写入状态，1. 原始拥有者删除资产，2. 新拥有者增加资产，3. 资产变更记录添加一条
	// 原始拥有者删除资产
	assetIds := make([]string, 0)
	for _, aid := range originOwner.Assets {
		if aid == assetid {
			continue
		}
		assetIds = append(assetIds, aid)
	}
	originOwner.Assets = assetIds
	OriginOwnerBytes, err = json.Marshal(originOwner)
	if err != nil {
		return shim.Error(fmt.Sprintf("序列化原始拥有者失败%s", err))
	}
	if err := stub.PutState(constructUserKey(ownerid), OriginOwnerBytes); err != nil {
		return shim.Error(fmt.Sprintf("原始拥有者写入账本失败%s", err))
	}

	// 新拥有者增加资产
	currentOwner := new(User)
	if err := json.Unmarshal(currentOwnerBytes, currentOwner); err != nil {
		return shim.Error(fmt.Sprintf("反序列化用户失败 %s", err))
	}
	currentOwner.Assets = append(currentOwner.Assets, assetid)
	currentOwnerBytes, err = json.Marshal(currentOwner)
	if err != nil {
		return shim.Error(fmt.Sprintf("序列化新拥有者失败%s", err))
	}
	if err := stub.PutState(constructUserKey(currentownerid), currentOwnerBytes); err != nil {
		return shim.Error(fmt.Sprintf("新拥有者写入账本失败%s", err))
	}

	// 插入资产变更记录
	history := &AssetHistory{
		AssetID:        assetid,
		OriginOwnerID:  ownerid,
		CurrentOwnerID: currentownerid,
	}
	historyBytes, err := json.Marshal(history)
	if err != nil {
		return shim.Error(fmt.Sprintf("序列化资产历史失败 %s", err))
	}
	historyKey, err := stub.CreateCompositeKey("history", []string{OriginOwner, assetid, ownerid})
	if err != nil {
		shim.Error(fmt.Sprintf("组合键制作失败%s", err))
	}
	if err = stub.PutState(historyKey, historyBytes); err != nil {
		shim.Error(fmt.Sprintf("资产历史写入账本失败 %s", err))
	}
	return shim.Success(nil)
}

// 用户查询
func queryUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 1. 检查参数的个数
	if len(args) != 1 {
		return shim.Error("参数数量错误")
	}
	// 2. 验证参数的正确性
	id := args[0]
	if id == "" {
		return shim.Error("参数值不能为空字符")
	}
	// 3. 验证数据是否存在，应该存在 or 不应该存在
	userBytes, err := stub.GetState(constructUserKey(id))
	if err != nil && len(userBytes) == 0 {
		return shim.Error("用户不存在")
	}
	return shim.Success(userBytes)
}

// 资产查询
func queryAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 1. 检查参数的个数
	if len(args) != 1 {
		return shim.Error("参数数量错误")
	}
	// 2. 验证参数的正确性
	assetid := args[0]
	if assetid == "" {
		return shim.Error("参数值不能为空字符")
	}
	// 3. 验证数据是否存在，应该存在 or 不应该存在
	assetBytes, err := stub.GetState(constructAssetKey(assetid))
	if err != nil && len(assetBytes) == 0 {
		return shim.Error("资产不存在")
	}
	return shim.Success(assetBytes)
}

// 资产历史查询
func queryAssetHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 1. 检查参数的个数
	if len(args) > 2 || len(args) < 1 {
		return shim.Error("参数数量错误")
	}
	// 2. 验证参数的正确性
	assetid := args[0]
	querytype := args[1] // 可选参数
	if assetid == "" {
		return shim.Error("参数值不能为空字符")
	}
	// 3. 验证数据是否存在，应该存在 or 不应该存在
	assetBytes, err := stub.GetState(constructAssetKey(assetid))
	if err != nil && len(assetBytes) == 0 {
		return shim.Error("资产不存在")
	}

	// 4. 查询相关数据
	keys := make([]string, 0)
	keys = append(keys, assetid)
	switch querytype {
	case "enroll":
		keys = append(keys, OriginOwner)
	case "exchange", "all": // 不添加任何附加

	default:
		return shim.Error(fmt.Sprintf("不支持的查询类型 %s", querytype))
	}
	// constructUserKey、constructAssetKey、constructAssetHistoryKey 是我们自己写的字符串拼接组合键
	// Fabric 内部实现了组合键 stub.CreateCompositeKey 这个方法接收
	result, err := stub.GetStateByPartialCompositeKey("history", keys)
	if err != nil {
		return shim.Error(fmt.Sprintf("查询资产历史失败%s", err))
	}
	defer result.Close()
	histories := make([]*AssetHistory, 0)
	for result.HasNext() {
		historyVal, err := result.Next()
		if err != nil {
			return shim.Error(fmt.Sprintf("查询失败%s", err))
		}
		history := new(AssetHistory)
		if err := json.Unmarshal(historyVal.GetValue(), history); err != nil {
			return shim.Error(fmt.Sprintf("反序列化历史失败%s", err))
		}
		// 过滤掉不是资产转让的记录
		if querytype == "exchange" && history.OriginOwnerID == OriginOwner {
			continue
		}
		histories = append(histories, history)
	}
	historyBytes, err := json.Marshal(histories)
	if err != nil {
		return shim.Error(fmt.Sprintf("序列化历史失败%s", err))
	}
	return shim.Success(historyBytes)
}

// Init 初始化链码时候执行的函数
func (t *AssetExangeCC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke 对链码的增删改查在这个接口实现
func (t *AssetExangeCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	// 获取函数名称和参数
	funcName, args := stub.GetFunctionAndParameters()
	switch funcName {
	case "userRegister":
		return userRegister(stub, args)
	case "userDestroy":
		return userDestroy(stub, args)
	case "assetEnroll":
		return assetEnroll(stub, args)
	case "assetExchange":
		return assetExchange(stub, args)
	case "queryUser":
		return queryUser(stub, args)
	case "queryAsset":
		return queryAsset(stub, args)
	case "queryAssetHistory":
		return queryAssetHistory(stub, args)
	default:
		return shim.Error(fmt.Sprintf("不支持的方法%s", funcName))
	}
}

func main() {
	err := shim.Start(new(AssetExangeCC))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
