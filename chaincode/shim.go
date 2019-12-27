/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package shim

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hyperledger/fabric/protos/ledger/queryresult"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// Chaincode 所有的 chaincodes 必须实现的接口 
// fabric 按照指定的方式调用这些函数来运行事务。
type Chaincode interface {
	// 首次建立链码容器后，在实例化事务期间调用 Init，从而允许链码初始化其内部数据
	Init(stub ChaincodeStubInterface) pb.Response

	// 调用 Invoke 来更新或查询 proposal transaction 中的分类帐。
	// 在提交事务之前，更新的状态变量不会提交到分类帐。
	Invoke(stub ChaincodeStubInterface) pb.Response
}

// ChaincodeStubInterface 可部署的 Chaincode 应用程序使用 ChaincodeStubInterface 访问和修改其分类帐
type ChaincodeStubInterface interface {
	// GetArgs 以字节数组数组的形式返回链代码初始化和调用的参数。
	GetArgs() [][]byte

	// GetStringArgs 以字符串数组的形式返回链代码初始化和调用的参数。
	// 只有在客户端传递了打算用作字符串的参数时才使用 GetStringArgs。
	GetStringArgs() []string

	// GetFunctionAndParameters 将第一个参数作为函数名返回，其余参数作为字符串数组中的参数返回。
	// Only use GetFunctionAndParameters if the client passes arguments intended to be used as strings.
	// 只有在客户端传递了打算用作字符串的参数时才使用 GetFunctionAndParameters。
	GetFunctionAndParameters() (string, []string)

	// GetArgsSlice 以字节数组的形式返回链代码 Init 和调用的参数
	GetArgsSlice() ([]byte, error)

	// GetTxID 返回 transaction proposal 的 tx_id，该 tx_id 对于每个事务和每个客户机都是惟一的。
	// 参见 protos/common/common 中的 ChannelHeader。更多细节请见 proto。
	GetTxID() string

	// GetChannelID 返回 proposal 发送到的用于 chaincode 处理的通道。
	// 这将是 transaction proposal 的 channel_id (参见 protos/common/common.proto 中的 ChannelHeader)，除非链码在另一个通道上调用另一个
	GetChannelID() string

	// InvokeChaincode 本地使用相同的事务上下文调用指定的链码 "Invoke"；也就是说，调用链码的链码不会创建新的事务消息。
	// 如果被调用的链码在同一个通道上，它只需将被调用的链码读取集和写入集添加到调用事务中。
	// 如果被调用的链码在不同的信道上，则只向调用链码返回响应；
	// 来自被调用链码的任何 PutState 调用都不会对分类账产生任何影响；
	// 也就是说，在不同的通道上调用的链码不会将其读集和写集应用于事务。
	// 只有调用链码的读集和写集将应用于事务。
	// 实际上，在不同通道上调用的链码是一个 "Query"，它在随后的提交阶段不参与状态验证检查。
	// 如果 "channel" 为空，则假定调用方的 channel 为空。
	InvokeChaincode(chaincodeName string, args [][]byte, channel string) pb.Response

	// GetState 从分类帐返回指定的 'key' 值。
	// 请注意，GetState 不会从 writeset 读取数据，因为 writeset 尚未提交到分类帐。
	// 换句话说，GetState 不考虑 PutState 修改的、尚未提交的数据。
	// 如果状态数据库中不存在 'key'，则返回(nil, nil)。
	GetState(key string) ([]byte, error)

	// PutState 将指定的 'key'和 'value'作为数据写入 proposal 放入 transaction 的 writeset 中。
	// 在确认并成功提交交易之前，PutState 不会影响分类账。
	// 简单键不能是空字符串，也不能以空字符 0x00 开头，以避免与组合键的范围查询冲突，在内部以 0x00 作为组合键命名空间作为前缀。
	// 另外，如果使用 CouchDB，键只能包含有效的 UTF-8 字符串，不能以下划线("_")开头。
	PutState(key string, value []byte) error

	// DelState 在 transaction proposal 的 writeset 中记录要删除的指定 'key'。
	// `key` 及其值将在验证并成功提交交易后从分类帐中删除。
	DelState(key string) error

	// SetStateValidationParameter 设置 `key` 的密钥级别背书策略。
	SetStateValidationParameter(key string, ep []byte) error

	// GetStateValidationParameter 检索 `key` 的密钥级别背书策略。
	// 注意，这将在事务的 readset 中引入对 'key' 的读依赖。
	GetStateValidationParameter(key string) ([]byte, error)

	// GetStateByRange 返回分类帐中一组 `key` 的范围迭代器。
	// 迭代器可用于迭代 startKey(包含) 和 endKey(排除) 之间的所有键。
	// 但是，如果 startKey 和 endKey 之间的键数大于 totalQueryLimit（在core.yaml中定义），则此迭代器不能用于获取所有键（结果将由 totalQueryLimit 限制）。
	// 迭代器按词法顺序返回 `key`。
	// 注意，startKey 和 endKey 可以是空字符串，这意味着开始或结束时的范围查询是无界的。
	// 完成后，对返回的 StateQueryIteratorInterface 对象调用 Close()。
	// 在验证阶段重新执行查询，以确保自事务背书以来结果集没有更改（幻读检测）。
	GetStateByRange(startKey, endKey string) (StateQueryIteratorInterface, error)

	// GetStateByRangeWithPagination 返回分类帐中一组键的范围迭代器。
	// 迭代器可用于获取 startKey（包含）和 endKey（排除）之间的键。
	// 将空字符串作为值传递给 bookmark 参数时，返回的迭代器可用于获取 startKey（包含）和 endKey（排除）之间的第一个 pageSize keys。
	// 当 bookmark 是非空字符串时，迭代器可用于获取 bookmark(包括) 和 endKey(排除) 之间的第一个 pageSize keys。
	// 请注意，只有在查询结果的前一页（ResponseMetadata）中存在的 bookmark 可以用作 bookmark 参数的值。
	// 否则，必须将一个空字符串作为 bookmark 传递。
	// 迭代器按词法顺序返回 `key`。
	// 请注意，startKey 和 endKey 可以为空字符串，这意味着在 start 或 end 上进行无限制范围查询。
	// 完成后，对返回的StateQueryIteratorInterface对象调用Close()。
	// 仅在只读事务中支持此调用。
	GetStateByRangeWithPagination(startKey, endKey string, pageSize int32, bookmark string) (StateQueryIteratorInterface, *pb.QueryResponseMetadata, error)

	// GetStateByPartialCompositeKey 根据给定的部分组合键查询分类账中的状态。
	// 此函数返回一个迭代器，可用于迭代前缀与给定的部分组合键匹配的所有组合键。
	// 但是，如果匹配组合键的数量大于 totalQueryLimit (在core.yaml中定义)，则不能使用此迭代器获取所有匹配键(结果将受到 totalQueryLimit 的限制)。
	// `objectType` 和属性应该只有有效的 utf8 字符串，不应该包含 U+0000 (nil字节) 和 U+10FFFF (最大且未分配的代码点)。
	// 参见相关函数 SplitCompositeKey 和 CreateCompositeKey。
	// 完成后，对返回的StateQueryIteratorInterface对象调用Close()。
	// 在验证阶段重新执行查询，以确保自事务背书以来结果集没有更改（幻读检测）。
	GetStateByPartialCompositeKey(objectType string, keys []string) (StateQueryIteratorInterface, error)

	// GetStateByPartialCompositeKeyWithPagination 根据给定的部分组合键查询分类账中的状态。
	// 此函数返回一个迭代器，该迭代器可用于迭代前缀与给定部分组合键匹配的组合键。
	// 当一个空字符串作为值传递给 bookmark 参数时，返回的迭代器可用于获取第一个 pageSize 组合键，其前缀与给定的部分组合键匹配。
	// 当 bookmark 是非空字符串时，迭代器可用于获取 bookmark(包括)和最后一个匹配的组合键之间的第一个 “pageSize” 键。
	// 注意，只有查询结果的前一页中的 bookmark (ResponseMetadata)可以用作 bookmark 参数的值。
	// 否则，必须将空字符串作为 bookmark 传递。
	// “objectType” 和属性应该只有有效的 utf8 字符串，并且不应该包含 U+0000（零字节）和 U+10FFFF（最大和未分配的代码点）。
	// 请参见相关函数SplitCompositeKey和CreateCompositeKey。
	// 完成后对返回的 StateQueryIteratorInterface 对象调用Close()。
	// 此调用仅在只读事务中受支持。
	GetStateByPartialCompositeKeyWithPagination(objectType string, keys []string, pageSize int32, bookmark string) (StateQueryIteratorInterface, *pb.QueryResponseMetadata, error)

	// CreateCompositeKey 组合给定的 `attributes` 以形成组合键。
	// objectType 和属性应该只有有效的 utf8 字符串，不应该包含 U+0000 (nil字节) 和 U+10FFFF (最大且未分配的代码点)。
	// 得到的组合键可以用作 PutState() 中的 key。
	CreateCompositeKey(objectType string, attributes []string) (string, error)

	// 将指定的键拆分为组成复合键的属性。
	// 因此，在范围查询或部分组合键查询期间找到的组合键可以被分割成它们的组合部分。
	SplitCompositeKey(compositeKey string) (string, []string, error)

	// GetQueryResult 对状态数据库执行“丰富”查询。
	// 只有支持丰富查询的状态数据库（例如CouchDB）才支持它。
	// 查询字符串是底层状态数据库的原生语法。
	// 返回一个迭代器，可用于迭代查询结果集中的所有键。
	// 但是，如果查询结果集中的键数大于 totalQueryLimit(在core.yaml中定义)，则不能使用此迭代器获取查询结果集中的所有键(结果将受到 totalQueryLimit 的限制)。
	// 在验证阶段不会重新执行查询，也不会检测到幻读，也就是说，其他提交的事务可能添加、更新或删除了影响结果集的键，而在验证/提交时不会检测到这些键。
	// 因此，易受此影响的应用程序不应将 GetQueryResult 用作更新分类帐的事务的一部分，而应将其使用范围限制为只读链码操作。
	GetQueryResult(query string) (StateQueryIteratorInterface, error)

	// GetQueryResultWithPagination 对状态数据库执行“丰富”查询。
	// 只有支持丰富查询的状态数据库（例如CouchDB）才支持它。
	// 查询字符串是底层状态数据库的原生语法。
	// 返回一个迭代器，可用于迭代查询结果集中的键。
	// 当将一个空字符串作为值传递给 bookmark 参数时，返回的迭代器可用于获取查询结果的第一个 “pageSize”。
	// 当 bookmark 是非空字符串时，迭代器可用于获取 bookmark 和查询结果中最后一个键之间的第一个 “pageSize” 键。
	// 注意，只有查询结果（ResponseMetadata）的前一页中存在的 bookmark 才能用作 bookmark 参数的值。
	// 否则，必须将空字符串作为 bookmark 传递。
	// 此调用仅在只读事务中受支持。
	GetQueryResultWithPagination(query string, pageSize int32,
		bookmark string) (StateQueryIteratorInterface, *pb.QueryResponseMetadata, error)

	// GetHistoryForKey 返回一段时间内键值的历史记录。
	// 对于每个历史密钥更新，将返回历史值以及关联的事务标识和时间戳。
	// 时间戳是客户端在 proposal 头中提供的时间戳。
	// GetHistoryForKey 需要 peer 配置 core.ledger.history。enableHistoryDatabase 为 true。
	// 在验证阶段不重新执行查询，不检测幻读，也就是说，其他已提交的事务可能同时更新了键，从而影响结果集，而这在验证/提交时不会被检测到。
	// 因此，易受此影响的应用程序不应将GetHistoryForKey用作更新分类帐的事务的一部分，而应将其使用限制为只读的链码操作。
	GetHistoryForKey(key string) (HistoryQueryIteratorInterface, error)

	// GetPrivateData 从指定的 “集合” 返回指定的 “键” 的值。
	// 请注意，GetPrivateData 不会从私有 writeset 读取数据，因为它尚未提交到 “collection” 中。
	// 换句话说，GetPrivateData 不考虑由尚未提交的 PutPrivateData 修改的数据。
	GetPrivateData(collection, key string) ([]byte, error)

	// GetPrivateDataHash 返回指定集合中指定 “key” 的值的哈希值`
	GetPrivateDataHash(collection, key string) ([]byte, error)

	// PutPrivateData 将指定的 “key” 和 “value”放入事务的私有 writeset。
	// 注意，只有 private writeset 的 hash 才会进入事务建议响应（发送到发出事务的客户端），而实际的 private writeset 会临时存储在一个临时存储中。 
	// 在验证并成功提交事务之前，PutPrivateData 不会影响 “collection”。
	// 简单键不能是空字符串，也不能以空字符（0x00）开头，以避免与组合键的范围查询冲突，组合键在内部以0x00作为组合键命名空间的前缀。
	// 此外，如果使用 CouchDB，则键只能包含有效的 UTF-8 字符串，并且不能以下划线（“\u”）开头。
	PutPrivateData(collection string, key string, value []byte) error

	// DelState 在事务的私有写集中记录要删除的指定 “key”。
	// 注意，只有 private writeset 的 hash 才会进入事务建议响应（发送到发出事务的客户端），而实际的 private writeset 会临时存储在一个临时存储中。
	// 验证并成功提交事务时，将从集合中删除 “key” 及其值。
	DelPrivateData(collection, key string) error

	// SetPrivateDataValidationParameter 为 “key” 指定的私有数据的键级背书策略。
	SetPrivateDataValidationParameter(collection, key string, ep []byte) error

	// GetPrivateDataValidationParameter 检索由' key '指定的私有数据的键级背书策略。
	// 注意，这在事务的 readset 中引入了对 'key' 的读依赖。
	GetPrivateDataValidationParameter(collection, key string) ([]byte, error)

	// GetPrivateDataByRange 在给定私有集合的一组键上返回范围迭代器。
	// 迭代器可用于迭代 startKey(包含) 和 endKey(排除) 之间的所有键。
	// 迭代器按词法顺序返回键。
	// 注意，startKey 和 endKey 可以是空字符串，这意味着开始或结束时的范围查询是无界的。
	// 完成后，对返回的 StateQueryIteratorInterface 对象调用 Close()。
	// 在验证阶段重新执行查询，以确保自事务背书之后结果集没有更改（幻读检测）。
	GetPrivateDataByRange(collection, startKey, endKey string) (StateQueryIteratorInterface, error)

	// GetPrivateDataByPartialCompositeKey 根据给定的部分组合键查询给定私有集合中的状态。
	// 此函数返回一个迭代器，可用于迭代前缀与给定的部分组合键匹配的所有组合键。
	// “objectType” 和属性应该只有有效的 utf8 字符串，不应该包含 U+0000 (nil字节)和 U+10FFFF(最大且未分配的代码点)。
	// 参见相关函数 SplitCompositeKey 和 CreateCompositeKey。
	// 完成后，对返回的 StateQueryIteratorInterface 对象调用Close()。
	// 在验证阶段重新执行查询，以确保自事务背书以来结果集没有更改（幻读检测）。
	GetPrivateDataByPartialCompositeKey(collection, objectType string, keys []string) (StateQueryIteratorInterface, error)

	// GetPrivateDataQueryResult 对给定的私有集合执行“富”查询。
	// 它只支持支持富查询的状态数据库，例如couchdb。
	// 查询字符串是底层状态数据库的原生语法。
	// 返回一个迭代器，该迭代器可用于迭代(下一个)查询结果集。
	// 在验证阶段不重新执行查询，不检测幻读。
	// 也就是说，其他提交的事务可能添加、更新或删除了影响结果集的键，而在验证/提交时不会检测到这些键。
	// 因此，易受此影响的应用程序不应将GetQueryResult用作更新分类帐的事务的一部分，而应将其使用限制为只读的链码操作。
	GetPrivateDataQueryResult(collection, query string) (StateQueryIteratorInterface, error)

	// GetCreator 返回 `SignedProposal` 的 `SignatureHeader.Creator`（例如身份）。
	// 这是提交事务的代理(或用户)的身份。
	GetCreator() ([]byte, error)

	// GetTransient 返回 `ChaincodeProposalPayload.Transient` 字段。
	// 它是一个包含数据（例如加密材料）的映射，这些数据可能用于实现某种形式的应用程序级机密性。
	// 由 ChaincodeProposalPayload 规定的该字段的内容应始终从交易中省略，并从分类帐中排除。
	GetTransient() (map[string][]byte, error)

	// 将事务绑定返回给提案本身，事务绑定用于将应用程序数据(如存储在上面的临时字段中的数据)之间的链接强制执行到提案本身。
	// 这对于避免可能的重放攻击很有用。
	GetBinding() ([]byte, error)

	// GetDecorations 返回关于来自对 peer 的 proposal 的附加数据（如果适用）. 
	// 此数据由 peer 端的装饰器设置，这些装饰器追加或更改传递给链码的链码输入。
	GetDecorations() map[string][]byte
	
	// GetSignedProposal 返回 SignedProposal 对象，该对象包含 transaction proposal 的所有数据元素。
	GetSignedProposal() (*pb.SignedProposal, error)

	// GetTxTimestamp 返回创建事务时的时间戳。
	// 这是从 transaction ChannelHeader 获取的，因此它将指示客户机的时间戳，并且在所有背书人中具有相同的值。
	GetTxTimestamp() (*timestamp.Timestamp, error)

	// SetEvent 允许链码在对建议的响应上设置一个事件，该事件将作为事务的一部分包括在内。
	// 无论事务的有效性如何，事件都将在提交的块中的事务内可用。
	SetEvent(name string, payload []byte) error
}

// CommonIteratorInterface 允许链码检查是否要从迭代器获取更多结果，并在完成后关闭它。
type CommonIteratorInterface interface {
	// 如果范围查询迭代器包含其他键和值，则 HasNext 返回 true。
	HasNext() bool

	// Close关闭迭代器。
	// 当完成从迭代器读取以释放资源时，应该调用此函数。
	Close() error
}

// StateQueryIteratorInterface 允许链码在由 range 返回的一组键/值对上迭代并执行查询。
type StateQueryIteratorInterface interface {
	// 继承 HasNext() 和 Close()
	CommonIteratorInterface

	// Next 返回范围中的下一个键和值，并执行查询迭代器。
	Next() (*queryresult.KV, error)
}

// HistoryQueryIteratorInterface 允许链码在历史查询返回的一组键/值对上迭代。
type HistoryQueryIteratorInterface interface {
	// 继承 HasNext() 和 Close()
	CommonIteratorInterface

	// Next返回历史查询迭代器中的下一个键和值。
	Next() (*queryresult.KeyModification, error)
}

// MockQueryIteratorInterface 允许链码在范围查询返回的一组键/值对上迭代。
// TODO: 一旦在 MockStub 中实现了 execute 查询和 history 查询，我们需要更新这个接口
type MockQueryIteratorInterface interface {
	StateQueryIteratorInterface
}
