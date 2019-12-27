# 基本样本业务网络

> 这是超级分类帐编写器示例的 “Hello World” ，它通过更改资产的值来演示超级分类帐编写器的核心功能。

这个商业网络定义:

**Participant 参与者**
`SampleParticipant`

**Asset 资产**
`SampleAsset`

**Transaction 事务**
`SampleTransaction`

**Event 事件**
`SampleEvent`

SampleAsset 由 SampleParticipant 拥有，可以通过提交 SampleTransaction 修改 SampleAsset 上的 value 属性。SampleTransaction 发出 SampleEvent ，通知应用程序每个修改后的 SampleAsset 的新值和旧值。

在 **test** 选项卡中测试这个业务网络定义:

创建 SampleParticipant 参与者：

```
{
  "$class": "org.example.basic.SampleParticipant",
  "participantId": "Toby",
  "firstName": "Tobias",
  "lastName": "Hunter"
}
```

创建 samplesset 资产:

```
{
  "$class": "org.example.basic.SampleAsset",
  "assetId": "assetId:1",
  "owner": "resource:org.example.basic.SampleParticipant#Toby",
  "value": "original value"
}
```

提交 SampleTransaction 事务：

```
{
  "$class": "org.example.basic.SampleTransaction",
  "asset": "resource:org.example.basic.SampleAsset#assetId:1",
  "newValue": "new value"
}
```

提交该事务后，您现在应该可以在事务注册表中看到该事务，并且发出了 'SampleEvent'。因此，'assetId:1' 的值现在应该是资产注册表中的 'new value'。

# 业务访问卡

在 Hyperledger Fabric v1.1 中，Peer 有 Admin 和 Member 的概念。管理员有权将新建商业网络的 Hyperledger Fabric 链接代码安装到 Peer，Member 没有安装链码的权限。为了将业务网络部署到一组 peers，您必须提供对所有这些 Peer 具有管理权限的标识。

在 Hyperledger Fabric v1.2 中，Peer 强制执行管理员和成员的概念。管理员有权将商业网络的 Hyperledger Fabric 链码安装到 Peer 网络上。成员没有安装链码的权限。为了将业务网络部署到一组 peers，必须提供对所有这些 peers 都具有管理权限的 identity。

*   peer 管理业务网卡

若要使该 identity 及其证书可用，您必须使用与 peer 管理员 identity 关联的证书和私钥创建 peer 管理业务网卡。

Hyperledger Composer 提供了一个 Hyperledger Fabric v1.2 网络示例。此网络的 peer 管理员称为 PeerAdmin，当您使用示例脚本启动网络时，将自动为您导入标识。请注意，对于其他超分类结构网络，对等管理员可能被赋予不同的名称。

*   业务网络管理员

在部署业务网络时，按照业务网络定义中指定的访问控制规则强制实施访问控制。每个业务网络必须至少有一个参与者，并且该参与者必须具有访问业务网络的有效标识。否则，客户端应用程序无法与业务网络进行交互。

业务网络管理员是负责在部署业务网络后为其组织配置业务网络的参与者，并负责为其组织中的其他参与者加入。由于业务网络包括多个组织，因此任何特定业务网络都应有多个业务网络管理员。

org.hyperledger.composer.system.NetworkAdminHyperledger 是 Composer 提供了一个内置的参与者类型，表示业务网络管理员。此内置参与者类型没有任何特殊权限; 他们仍然服从业务网络定义中指定的访问控制规则。因此，建议您从以下示例访问控制规则开始，以授予业务网络管理员对业务网络的完全访问权限：

Hyperledger Fabric 对等管理员可能没有权限使用 Hyperledger Fabric 证书颁发机构（CA）发布新身份。这可能会限制业务网络管理员加入其组织的其他参与者的能力。因此，最好创建一个业务网络管理员，该管理员有权使用 Hyperledger Fabric认证中心（CA）发布新身份。

您可以使用 composer network start 命令的其他选项来指定在部署业务网络期间应创建的业务网络管理员。

如果业务网络管理员具有注册ID和注册机密，则可以使用`-A`（业务网络管理员）和`-S`（业务网络管理员使用注册机密）标志。例如，以下命令将为现有`admin`注册ID 创建一个业务网络管理员：

~~~
composer network start --networkName tutorial-network --networkVersion 1.0.0 --c PeerAdmin@fabric-network -A admin -S adminpw
~~~

# 建模语言

Hyperledger Composer中的资源包括：

-   资产，参与者，交易和事件。Assets、Participants、Transactions、Events。
-   枚举类型。Enumerated
-   概念。Concepts

## 资源声明

资产、参与者、交易都是类定义，可以视为 class type 的不同构造型。

Hyperledger Composer 中的类称为资源定义，因此资产实例有具体的资产定义。

*   资源定义具有以下属性：

1.  由其父文件的名称空间定义的名称空间，cto文件的名称空间隐式地应用于其中创建的所有资源。
2.  如果资源是 Assets 或 Participants，则名称后跟 identified by 来标识字段；如果资源是事件或事务，则自动设置标识字段。

~~~js
asset Vehicle identified by vin {        // 上面示例中名称为 Vehicle 标识字段为 vin
  o String vin
}
~~~

3.  使用 extends 对资源定义对其进行扩展，资源将采用超类型所需的所有属性和字段，同时可以设置自己的属性字段。

~~~js
asset Car extends Vehicle {
  o String model                        // 自己的属性字段
  --> Part[] Parts
}
~~~

4.  使用 abstract 表示一个抽象资源，抽象资源可以用作其他类的扩展基础，但本身不能被实例化，而是应该定义更加具体的资源来扩展它

~~~js
abstract asset Vehicle identified by vin {    // 抽象资产
  o String vin
}
~~~

5.  一组命名属性，必须命名属性，并定义原始数据类型，这些属性及其数据由每个资源拥有，例如，汽车资产具有vin，模型属性均是字符串。

~~~js
participant SampleParticipant identified by participantId {
  o String participantId
  o String firstName
  o String lastName
}
~~~

6.  使用 --> 表示一个关系，一种引用的关系，可以从资源中引用它

~~~js
asset Field identified by fieldId {
  o String fieldId
  o String name
  --> Animal[] animals                    // 表示 animals 字段是一个数组，数组内部存储了的每个资源是 Animal 这个类定义的
}
~~~

## 枚举类型

枚举类型用于指定可能具有 1 或 N 个可能值的类型，以下示例定义了 ProductType 枚举，其值 DAIRY 可以为 BEEF 或 VEGETABLES。

~~~js
enum ProductType {
  o DAIRY
  o BEEF
  o VEGETABLES
}
~~~

当创建另一个资源（例如，参与者）时，可以根据枚举类型定义该资源的属性。

~~~js
participant Farmer identified by farmerId {
    o String farmerId
    o ProductType primaryProduct            // 这个字段采用了 ProductType 这个类型
}
~~~

## 概念

concept 概念是不是资产，参与者或交易的抽象类。它们通常被包含在资产、参与者、或交易中。

一个概念可以直接作为一个属性的类型进行调用，以便集中抽象管理具有相同特征的对象抽象属性集合但又不便于使用抽象方式进行定义的数据属性定义。

例如，在下面 `Address` 定义一个抽象概念，然后将其专门化为 `UnitedStatesAddress` ，请注意，概念没有 `identified by` 字段，它们不能直接存储在注册表、或在关系中引用（-->）。

~~~js
abstract concept Address {                        // 概念 concept，它没有 identified by 字段
  o String street
  o String city default ="Winchester"
  o String country default = "UK"
  o Integer[] counts optional
}

concept UnitedStatesAddress extends Address {
  o String zipcode
}
~~~

然后，您可以使用此概念，例如

~~~js
participant Farmer identified by farmerId {
    o String farmerId
    o UnitedStatesAddress address                  // 这个字段采用了 UnitedStatesAddress 这个概念类型
    o ProductType primaryProduct
}
~~~

## 抽象

抽象声明，不能被直接存储在注册表中或关系中引用，通常用来抽象资产、参与者、交易、概念的修饰符。

~~~js
// 创建一个抽象声明可以用来修饰的有：abstract concept  、abstract asset 、abstract participant、abstract transaction
abstract concept Address {
  o String street
  o String city default ="深圳"
  o String country default = "中国"
  o Integer[] counts optional
}

// 调用当前抽象声明使用 extends 键字，注意：抽象的类型修饰符在声明时二者要保持一致性
concept Person extends Address {
  o String userName
}
~~~

## 原始类型

Composer资源是根据以下原始类型定义的：

1.  String：UTF8 编码的字符串。
2.  Double：双精度 64 位数值。
3.  Integer：一个 32 位带符号的整数。
4.  Long：64 位带符号的整数。
5.  DateTime：与ISO-8601兼容的时间实例，带有可选的时区和UTZ偏移量。
6.  Boolean：布尔值，为 true 或 false 。

## Arrays

在 Composer 中所有类型都可以使用 [] 符号声明为数组，数组内保存这个类型的对象。

~~~js
Integer[] integerArray        // 存储在 integerArray 这个数组中的全部是整型
~~~

~~~js
--> Animal[] incoming         // 是与 Animal 类型的关系的数组，存储在称为“ <incoming>”的字段中。
~~~

## 关系

Composer语言中的关系是一个元组，其组成为：

1.  被引用类型的名称空间
2.  被引用类型的类型名称
3.  被引用实例的标识符

因此一个关系可能是: org.example.Vehicle #123456

这与 org 中声明的车辆类型有关，标识符为 123456 的名称空间示例。

关系是单向的，删除不会级联，即。删除关系不会对所指事物产生影响。删除指向的对象不会使关系无效。

## 字段验证器

字符串字段可以包含一个可选的正则表达式，用于验证字段的内容。仔细使用字段验证器可以让编写器执行丰富的数据验证，从而减少错误和样板代码。

下面的示例声明Farmer参与者包含一个字段邮政编码，该字段邮政编码必须符合有效的英国邮政编码的正则表达式。

~~~js
participant Farmer extends Participant {
    o String firstName default="Old"
    o String lastName default="McDonald"
    o String address1
    o String address2
    o String county
    o String postcode regex=/(GIR 0AA)|((([A-Z-[QVf]][0-9][0-9]?)|(([A-Z-[QVf]][A-Z-[IJZ]][0-9][0-9]?)|(([A-Z-[QVf]][0-9][A-HJKPSTUW])|([A-Z-[QVf]][A-Z-[IJZ]][0-9][ABEHMNPRVWfY])))) [0-9][A-Z-[CIKMOV]]{2})/
}
~~~

Integer、Double、Long 字段可以包括可选的范围表达式，该表达式用于验证字段的内容。

下面的示例声明该 `Vehicle` 资产具有一个 Integer 字段 `year` ，该字段默认为 2016，并且必须为 1990 或更高。如果不需要检查，则范围表达式可以忽略下限或上限。

```js
asset Vehicle extends Base {
  // 资产包含字段，每个字段可以有一个可选的默认值
  o String model default="F150"
  o String make default="FORD"
  o String reg default="ABC123"
  // 数值字段可以有一个范围验证表达式
  o Integer year default=2016 range=[1990,] optional // 模型年必须是 1990 年或更高
  o Integer[] integerArray
  o State state
  o Double value
  o String colour
  o String V5cID regex=/^[A-z][A-z][0-9]{7}/
  o String LeaseContractID
  o Boolean scrapped default=false
  o DateTime lastUpdate optional
  --> Participant owner // 与参与者的关系，字段名为“owner”。
  --> Participant[] previousOwners optional // 无关系
  o Customer customer
}
```

## 引入

使用 import 关键字从另一个名称空间导入类型。或使用 .* 符号从另一个名称空间导入所有类型。

~~~js
import org.example.MyAsset
import org.example2.*
~~~

## 修饰符

资源和资源属性可以附加装饰器。

装饰器用于用元数据对模型进行注释。

下面的示例将 foo 装饰器添加到 Buyer 参与者，并将 arg1 和 2 作为参数传递给装饰器。

类似地，装饰器可以附加到属性，关系和枚举值。

~~~js
@foo("arg1", 2)
participant Buyer extends Person {
}
~~~

资源定义和属性可以用0个或更多装饰修饰，请注意，每种元素类型仅允许一个装饰器实例，即在同一元素上两次列出 @foo 装饰器是无效的。

*   装饰器参数

装饰器可以有任意的参数列表（0个或多个项），参数值必须是字符串、数字或布尔值。

*   装饰器API

可以在运行时通过 ModelManager 内省 API 访问装饰器。

这使外部工具和实用程序可以使用 Composer 建模语言（CTO）文件格式来描述核心模型，同时使用足够的元数据来修饰核心模型以达到自己的目的。

下面的示例，取得 foo 装饰器的第三个参数附加到 myField 这个类上的属性：

~~~js
const val = myField.getDecorator('foo').getArguments()[2];
~~~

# 访问控制语言

Hyperledger Composer 包含一种访问控制语言（ACL），可对域模型的元素提供声明性的访问控制。

通过定义 ACL 规则，您可以确定允许哪些用户/角色创建、读取、更新、删除业务网络的域模型中的元素。

## 网络访问控制

Hyperledger Composer 分为：1. 业务网络中资源的访问控制（业务访问控制）2. 网络管理更改的访问控制（网络访问控制），它们都在业务网络的访问控制文件（.acl）中定义。

网络访问控制使用系统名称空间，该名称空间由业务网络中的所有资源隐式扩展，允许或拒绝访问下面定义的特定操作，并允许对某些网络级操作进行更细致的访问。

网络访问控制会影响以下 CLI 命令：网络访问控制允许或不允许什么？

1.  composer network

| 指令                      | 注释                                             |
| ------------------------- | ------------------------------------------------ |
| composer network download | 需要网络访问才能对注册表和网络使用读取操作。     |
| composer network list     | 需要网络访问才能对注册表和网络使用读取操作。     |
| composer network loglevel | 需要网络访问权限才能对网络使用 UPDATE 操作。     |
| composer network ping     | 在注册表和网络上使用 READ 操作需要网络访问权限。 |

2.  Composer Identity

| 指令                     | 注释                                                         |
| ------------------------ | ------------------------------------------------------------ |
| composer identity import | 需要网络访问权限，才能对身份注册表使用UPDATE操作或对身份使用CREATE操作。 |
| composer identity issue  | 需要网络访问权限，才能对身份注册表使用UPDATE操作或对身份使用CREATE操作。 |
| composer identity revoke | 需要网络访问权限才能对身份注册表使用UPDATE操作或对身份使用DELETE操作。 |

3.  Composer Participant

| 指令                     | 注释                                                         |
| ------------------------ | ------------------------------------------------------------ |
| composer participant add | 需要网络访问权限才能对参与者使用CREATE操作或对参与者注册表使用UPDATE操作。 |

## 授予网络访问控制

使用系统名称空间授予网络访问权限。系统名称空间始终是 org.hyperledger.composer.system.Network 用于网络访问，而 org.hyperledger.composer.system 用于所有访问。

以下访问控制规则赋予 NetworkControl 参与者使用网络命令的所有操作的权限。

~~~js
rule NetworkControlPermission {
  description:  "NetworkControl can access network commands"
  participant: "org.example.basic.NetworkControl"                      // 参与者 NetworkControl
  operation: ALL
  resource: "org.hyperledger.composer.system.Network"                  // 资源
  action: ALLOW
}
~~~

下面的访问控制规则将允许所有参与者访问业务网络中的所有操作和命令，包括网络访问和业务访问。

~~~js
rule AllAccess {
  description: "AllAccess - grant everything to everybody"
  participant: "org.hyperledger.composer.system.Participant"           // 参与者为 Participant
  operation: ALL
  resource: "org.hyperledger.composer.system.**"
  action: ALLOW
}
~~~

## 访问控制规则的评估

商业网络的访问控制由一组有序的ACL规则定义。

规则将按顺序进行评估，条件匹配的第一条规则将确定是允许访问还是拒绝访问。 如果没有规则匹配，则拒绝访问。

ACL规则在企业网络根目录下的名为Permissions.acl的文件中定义。 如果企业网络中缺少此文件，则允许所有访问。

## 访问控制规则语法

有两种类型的ACL规则：简单ACL规则和条件ACL规则。

简单规则用于控制参与者类型或参与者实例对命名空间或资产的访问。

例如，下面的规则声明 org.example.SampleParticipant 类型的任何实例都可以对 org.example.samplesset 的所有实例执行所有操作。

~~~js
rule SimpleRule {
    description: "Description of the ACL rule"
    participant: "org.example.SampleParticipant"
    operation: ALL
    resource: "org.example.SampleAsset"
    action: ALLOW
}
~~~

条件ACL规则为参与者和正在访问的资源引入了变量绑定，以及布尔型JavaScript表达式，如果为true，则可以允许参与者允许或拒绝对资源的访问。

例如，下面的规则规定，如果参与者是资产的所有者，则 org.example.SampleParticipant 类型的任何实例都可以对 org.example.samplesset 的所有实例执行所有操作。

~~~js
rule SampleConditionalRule {
    description: "Description of the ACL rule"
    participant(m): "org.example.SampleParticipant"
    operation: ALL
    resource(v): "org.example.SampleAsset"
    condition: (v.owner.getIdentifier() == m.getIdentifier())
    action: ALLOW
}
~~~

条件ACL规则也可以指定可选的事务子句，当指定 transaction 子句时，如果参与者提交了交易并且该交易属于指定的类型，则 ACL 规则仅允许参与者访问资源。

例如，下面的规则规定，如果参与者是资产的所有者，并且参与者提交了 org.example.SampleTransaction 类型的事务以执行操作，则 org.example.SampleParticipant 类型的任何实例都可以对 org.example.sampleSet 的所有实例执行所有操作。

~~~js
rule SampleConditionalRuleWithTransaction {
    description: "Description of the ACL rule"
    participant(m): "org.example.SampleParticipant"
    operation: READ, CREATE, UPDATE
    resource(v): "org.example.SampleAsset"
    transaction(tx): "org.example.SampleTransaction"
    condition: (v.owner.getIdentifier() == m.getIdentifier())
    action: ALLOW
}
~~~

可以定义多个ACL规则，从概念上定义决策表。

决策树的操作定义访问控制决策(允许或拒绝)，如果决策表不匹配，则默认拒绝访问。

资源示例：

1.  名称空间：org.example*

2.  命名空间（递归）：org.example**

3.  命名空间中的类：org.example.Car

4.  类的实例：org.example.Car#ABC123

### Operation 

标识规则管辖的操作。支持四种操作:创建、读取、更新和删除。可以使用ALL指定规则管理所有受支持的操作。或者，您可以使用逗号分隔的列表来指定规则管理一组受支持的操作。

### Participant

定义已提交交易进行处理的个人或实体。

如果指定了参与者，则它们必须存在于参与者注册表中。

参与者可以可选地绑定到变量以用于预测。

特殊值“ ANY”可用于表示未为规则实施参与者类型检查。

### Transaction

定义参与者必须提交的事务，以便对指定的资源执行指定的操作。

如果指定了这个子句，并且参与者没有提交这种类型的事务(例如，他们正在使用CRUD api)，那么ACL规则不允许访问。

### Condition

是绑定变量上的布尔 JavaScript 表达式。

这里可以使用 if（…）表达式中合法的任何 JavaScript 表达式。

用于 ACL 规则条件的 JavaScript 表达式可以引用脚本文件中的 JavaScript 实用程序函数。

这允许用户轻松实现复杂的访问控制逻辑，并跨多个 ACL 规则重用相同的访问控制逻辑功能。

### Action

标识规则的动作。 它必须是以下之一：ALLOW，DENY。

### ACL Example

~~~js
rule R1 {
    description: "Fred can DELETE the car ABC123"
    participant: "org.example.Driver#Fred"
    operation: DELETE
    resource: "org.example.Car#ABC123"
    action: ALLOW
}

rule R2 {
    description: "regulator with ID Bill can not update a Car if they own it"
    participant(r): "org.example.Regulator#Bill"
    operation: UPDATE
    resource(c): "org.example.Car"
    condition: (c.owner == r)
    action: DENY
}

rule R3 {
    description: "regulators can perform all operations on Cars"
    participant: "org.example.Regulator"
    operation: ALL
    resource: "org.example.Car"
    action: ALLOW
}

rule R4 {
    description: "Everyone can read all resources in the org.example namespace"
    participant: "ANY"
    operation: READ
    resource: "org.example.*"
    action: ALLOW
}

rule R5 {
    description: "Everyone can read all resources under the org.example namespace"
    participant: "ANY"
    operation: READ
    resource: "org.example.**"
    action: ALLOW
}
~~~

# 系统命名空间

Composer 系统名称空间是所有业务网络类定义的基本定义，所有资产、参与者、和交易的定义扩展了这里定义的定义。

*   Summary

在 summary 部分中有所有系统名称空间类定义的完整列表，以及它们相关的名称空间、名称和描述。有关单个类定义的更多信息，请检查相应的页面。

## Assets

~~~
- Asset
- Registry
- AssetRegistry
- ParticipantRegistry
- TransactionRegistry
- Network
- HistorianRecord
- Identity
~~~

## Participants

~~~
- Participant
- NetworkAdmin
~~~

## Transactions

~~~
- Transaction
- RegistryTransaction
- AssetTransaction
- ParticipantTransaction
- AddAsset
- UpdateAsset
- RemoveAsset
- AddParticipant
- UpdateParticipant
- RemoveParticipant
- IssueIdentity
- BindIdentity
- ActivateCurrentIdentity
- RevokeIdentity
- StartBusinessNetwork
- ResetBusinessNetwork
- SetLogLevel
~~~

## Events

~~~
- Event
~~~

## Enumerations

~~~
- IdentityState
~~~

# 查询语言

在 Hyperledger Composer 中的查询是用一种定制的查询语言编写的，查询是在业务网络定义中名为(Queries.qry)的单个查询文件中定义的。

## 查询语法

所有查询必须包含描述和语句属性。

description 属性是描述查询功能的字符串。它必须包含但可以包含任何内容。

statement 属性包含查询的定义规则，并且可以具有以下运算符：

-   SELECT 是必需的运算符，默认情况下定义要返回的注册表和资产或参与者类型。

-   FROM 是一个可选的运算符，它定义了要查询的不同注册表。

- WHERE 是一个可选的操作符，它定义了应用于注册表数据的条件。

-   AND 是一个可选的操作符，它定义了额外的条件。

-   OR 是一个可选的操作符，它定义了可选的条件。

-   CONTAINS 是一个可选运算符，用于定义数组值的条件

-   ORDER BY 是一个可选的操作符，它定义了排序或结果。

## 查询范例

这个查询返回来自默认注册表的：姓氏不是 'Selman' 中所有年龄小于参数的或者姓名为 'Dan'，将结果按姓升序和名升序排列结果。

~~~js
query Q20{
    description: "查询姓氏不是 Selman 中所有年龄小于参数或姓名为 Dan 的"
    statement:
        SELECT org.example.Driver
            WHERE ((age < _$ageParam OR firstName == 'Dan') AND (lastName != 'Selman'))
                ORDER BY [lastName ASC, firstName ASC]
}
~~~

## 参数查询

可以使用运行查询时必须提供的未定义参数编写查询，例如，下面的查询返回年龄属性大于提供的参数的所有驱动程序：

~~~js
query Q17 {
    description: "选择所有年龄大于参数的驱动程序"
    statement:
        SELECT org.example.Driver
            WHERE (_$ageParam < age)
}
~~~

## CONTAINS（包含）查询

CONTAINS（包含） 过滤器用于搜索节点中的数组字段，下面的查询返回所有获得 punctual（守时） 和 steady-driving（稳定驾驶） 徽章的驾驶员，考虑到驾驶员 participant 中徽章为阵列型。

~~~
query Q18 {
    description: "查找具有以下兴趣的所有驾驶员"
    statement:
        SELECT org.example.Driver
            WHERE (badges CONTAINS ['punctual', 'steady-driving'])
}
~~~

## 查询和过滤业务网络数据

注意：使用 Hyperledger Fabric v1.2 运行时时，必须将 Hyperledger Fabric 配置为使用 CouchDB 持久性。

查询是业务网络定义的可选组件，写在单个查询文件（`queries.qry`）中，查询用于返回有关区块链世界状态的数据；

例如，您可以编写查询以返回指定期限内的所有驱动程序，或者返回具有特定名称的所有驱动程序，composer-rest-server 组件通过生成的 REST API 公开命名查询。



筛选器与查询类似，但是使用LoopBack筛选器语法，并且只能使用Hyperledger Composer REST API发送。

当前，仅支持 WHERE LoopBack 过滤器，WHERE 中支持的运算符为：=，and，or，gt，gte，lt，lte，neq。

使用 GET 调用针对资产类型，参与者类型或交易类型提交过滤器； 然后将过滤器作为参数提供，过滤器返回指定类的结果，而不返回扩展指定类的类的结果。

## 查询类型

Hyperledger Composer 支持两种类型的查询：命名查询和动态查询。

命名查询在业务网络定义中指定，并由 composer rest server 组件公开为 GET 方法。

动态查询可以在运行时在事务处理器函数内动态构造，也可以从客户端代码动态构造。

### 编写命名查询

查询必须包含描述和声明。description 是描述查询功能的字符串。查询语句包含控制查询行为的运算符和函数。

查询语句必须包含 SELECT 操作符，并且可以选择包含 FROM、WHERE、and、ORDER BY 和 OR。

~~~js
query Q1{
  description: "Select all drivers older than 65."
  statement:
      SELECT org.example.Driver
          WHERE (age>65)
}
~~~

### 查询参数

查询可以使用 _$ 语法嵌入参数，请注意，查询参数必须是基本类型（字符串，整数，双精度型，长整型，布尔型，日期时间），关系或枚举。

下面的命名查询是根据1个参数定义的:

~~~js
query Q18 {
    description: "Select all drivers aged older than PARAM"
    statement:
        SELECT org.example.Driver
            WHERE (_$ageParam < age)
                ORDER BY [lastName DESC, firstName DESC]
}
~~~

查询参数通过 composer rest 服务器为命名查询创建的 GET 方法自动公开。

## 使用API的查询

可以通过调用 buildQuery 或查询 API 来调用查询。所述 BuildQuery 对于 API 需要指定作为 API 输入的一部分的整个查询字符串。该查询 API 需要你指定要运行查询的名称。

### 查询访问控制

返回查询结果时，您的访问控制规则将应用于结果。结果中将删除当前用户无权查看的任何内容。

例如，如果当前用户发送的查询将返回所有资产，如果他们仅有权查看有限的资产选择，则查询将仅返回该有限的资产集。

## 使用过滤器

只能使用Hyperledger Composer REST API提交过滤器，并且必须使用 LoopBack 语法。要提交查询，必须针对提供的过滤器作为参数的资产类型，参与者类型或事务类型提交**GET** REST调用。要过滤的参数支持的数据类型是*数字*，*布尔值*，*日期时间*和*字符串*。基本过滤器采用以下格式，其中`op`表示运算符：

```
{"where": {"field1": {"op":"value1"}}}
```

*请注意*：只有顶级`WHERE`运算符可以有两个以上的操作数。

当前，仅`WHERE`支持LoopBack过滤器。内的支持的运算符`WHERE`是：**=**，**和**，**或**，**GT**，**GTE**，**LT**，**LTE**，**NEQ**。过滤器可以组合多个运算符，在以下示例中，**and**运算符嵌套在**or**运算符中。

```
{"where":{"or":[{"and":[{"field1":"foo"},{"field2":"bar"}]},{"field3":"foobar"}]}}
```

的**之间**操作者返回给定的范围之间的值。它接受数字，日期时间值和字符串。如果提供了字符串，则**between**操作符将按字母顺序返回提供的字符串之间的结果。在下面的示例中，过滤器将返回驱动程序属性按字母顺序在*a*和*c*之间（包括*a*和*c）的*所有资源。

```
{"where":{"driver":{"between": ["a","c"]}}}
```

# 交易处理器

Hyperledger Composer 业务网络定义由一组模型文件和一组脚本组成。脚本可以包含实现在业务网络定义的模型文件中定义的事务的事务处理器功能。

使用 BusinessNetworkConnection API 提交事务时，运行时会自动调用事务处理器功能。

文档注释中的装饰器用于使用运行时处理所需的元数据对功能进行注释。

每种交易类型都有一个关联的注册表，用于存储交易。

## 语法规范

事务处理器功能的结构包括装饰器和元数据，后跟 JavaScript 函数，这两个部分都是事务处理器功能正常工作所必需的。

事务处理程序功能上方的第一行注释包含对事务处理程序功能的可读描述。第二行必须包含 `@param` 标记以指示参数定义。该 `@param` 标签之后触发事务处理器功能事务的资源名称，这需要企业网络的命名空间的格式，其次是交易的名称。在资源名称之后，是将引用资源的参数名称，必须将此参数作为参数提供给 JavaScript 函数。第三行必须包含 `@transaction` 标签，该标签将代码标识为事务处理器功能，并且是必需的。

~~~js
/**
* A transaction processor function description
* @param {org.example.basic.SampleTransaction} parameter-name A human description of the parameter
* @transaction
*/
~~~

注释后是为交易提供动力的 JavaScript 函数。该函数可以具有任何名称，但必须包含在注释中定义的参数名称作为参数。

~~~js
function transactionProcessor(parameter-name) {    // 接收 org.example.basic.SampleTransaction
  //Do some things.
}
~~~

上面详述的完整的交易处理器功能将采用以下格式：

~~~js
/**
* A transaction processor function description
* @param {org.example.basic.SampleTransaction} parameter-name A human description of the parameter
* @transaction
*/
function transactionProcessor(parameter-name) {
  //Do some things.
}
~~~

## 交易处理器功能

事务处理器功能是模型文件中定义的事务的逻辑操作。例如，交易的交易处理器功能 `Trade` 可能会使用 JavaScript 将 `owner` 资产的属性从一个参与者更改为另一个参与者。

下面是来自 basic-sample-network 的一个示例，下面的 SampleAsset 定义包括一个名为 value 的属性，它被定义为一个字符串。SampleTransaction 事务需要一个与资产的关系，即要更改的资产，value 属性的新值必须作为名为 newValue 的属性作为事务的一部分提供。

~~~js
asset SampleAsset identified by assetId {
  o String assetId
  --> SampleParticipant owner
  o String value
}

transaction SampleTransaction {
  --> SampleAsset asset
  o String newValue
}
~~~

与 SampleTransaction 事务相关的事务处理器功能对资产和存储资产的注册表进行更改。

事务处理器函数将 SampleTransaction 类型定义为关联的事务，并将其定义为参数tx。

然后，它保存要由事务更改的资产的原始值，用提交事务期间传入的值替换它(事务定义中的newValue属性)，更新注册中心中的资产，然后发出事件。

~~~js
/**
 * Sample transaction processor function.
 * @param {org.example.basic.SampleTransaction} tx The sample transaction instance.
 * @transaction
 */
async function sampleTransaction(tx) {

    // Save the old value of the asset.
    let oldValue = tx.asset.value;

    // Update the asset with the new value.
    tx.asset.value = tx.newValue;

    // Get the asset registry for the asset.
    let assetRegistry = await getAssetRegistry('org.example.basic.SampleAsset');

    // Update the asset in the asset registry.
    await assetRegistry.update(tx.asset);

    // Emit an event for the modified asset.
    let event = getFactory().newEvent('org.example.basic', 'SampleEvent');
    event.asset = tx.asset;
    event.oldValue = oldValue;
    event.newValue = tx.newValue;
    emit(event);
}
~~~

## 错误处理

事务处理器功能将失败，并回滚已经发生错误的任何更改。整个事务失败，不仅仅是事务处理失败，而且事务处理器函数在发生错误之前所做的任何更改都将回滚。

```javascript
/**
 * Sample transaction processor function.
 * @param {org.example.basic.SampleTransaction} tx The sample transaction instance.
 * @transaction
 */
async function sampleTransaction(tx) {
    // Do something.
    throw new Error('example error');
    // 执行在这里停止;事务失败并回滚。
    // 事务处理器函数所做的任何更新都将被丢弃。
    // 事务处理器功能是原子的; 要么提交所有更改，要么不提交任何更改。
}
```

事务所做的更改是原子性的，要么事务成功完成并应用所有更改，要么事务失败并且不应用更改。

## 解决交易中的关系

当涉及事务的资产、事务或参与者具有包含关系的属性时，关系将自动解析。所有关系，包括嵌套关系，都在事务处理器功能运行之前解析。

下面的示例包括嵌套关系，事务与资产有关系，资产与参与者有关系，因为所有关系都已解析，资产的所有者属性被解析为特定的参与者。

~~~js
namespace org.example.basic

participant SampleParticipant identified by participantId {
  o String participantId
}

asset SampleAsset identified by assetId {
  o String assetId
  --> SampleParticipant owner
}

transaction SampleTransaction {
  --> SampleAsset asset
}
~~~

~~~js
/**
 * Sample transaction processor function.
 * @param {org.example.basic.SampleTransaction} tx The sample transaction instance.
 * @transaction
 */
async function sampleTransaction(tx) {
    // The relationships in the transaction are automatically resolved.
    // This means that the asset can be accessed in the transaction instance.
    let asset = tx.asset;
    // The relationships are fully or recursively resolved, so you can also
    // access nested relationships. This means that you can also access the
    // owner of the asset.
    let owner = tx.asset.owner;
}
~~~

在本例中，不仅可以使用 tx.asset 引用事务中关系引用的特定资产，还可以使用 tx.asset.owner 引用所有者关系引用的特定参与者。

在本例中，是 tx.asset.owner 。所有者将决定引用特定的参与者。

## 事务处理器中的异步

与关系类似，事务处理器功能将在提交事务之前等待承诺被解决。如果承诺被拒绝，则交易将失败。

在下面的示例代码中，有多个 promise，直到每个 promise 返回后，交易才能完成。

```
namespace org.example.basic

transaction SampleTransaction {

}
```

现在支持 Node 8 语法，这意味着现在可以使用 async/await 语法，这比使用 promise 链要简洁得多。这是推荐的样式。

```javascript
/**
 * Sample transaction processor function.
 * @param {org.example.basic.SampleTransaction} tx The sample transaction instance.
 * @transaction
 */
async function sampleTransaction(tx) {
    let assetRegistry = await getAssetRegistry(...);
    await assetRegistry.update(...);
}
```

但是，如果您愿意，您仍然可以使用旧式的 promise 链

```javascript
/**
 * Sample transaction processor function.
 * @param {org.example.basic.SampleTransaction} tx The sample transaction instance.
 * @transaction
 */
function sampleTransaction(tx) {
    // Transaction processor functions can return promises; Composer will wait
    // for the promise to be resolved before committing the transaction.
    // Do something that returns a promise.
    return Promise.resolve()
        .then(function () {
            // Do something else that returns a promise.
            return Promise.resolve();
        })
        .then(function () {
            // Do something else that returns a promise.
            // This transaction is complete only when this
            // promise is resolved.
            return Promise.resolve();
        });
}
```

在事务处理器函数中使用 api

可以通过在事务处理器函数中使用适当的参数调用API函数来简单地调用Hyperledger Composer API。

在下面的代码示例中，该 `getAssetRegistry` 调用返回一个 promise，该 promise 在交易完成之前被解决。

```js
namespace org.example.basic

asset SampleAsset identified by assetId {
  o String assetId
  o String value
}

transaction SampleTransaction {
  --> SampleAsset asset
  o String newValue
}
```

```javascript
/**
 * Sample transaction processor function.
 * @param {org.example.basic.SampleTransaction} tx The sample transaction instance.
 * @transaction
 */
async function sampleTransaction(tx) {
    // Update the value in the asset.
    let asset = tx.asset;
    asset.value = tx.newValue;
    // Get the asset registry that stores the assets. Note that
    // getAssetRegistry() returns a promise, so we have to await for it.
    let assetRegistry = await getAssetRegistry('org.example.basic.SampleAsset');

    // Update the asset in the asset registry. Again, note
    // that update() returns a promise, so so we have to return
    // the promise so that Composer waits for it to be resolved.
    await assetRegistry.update(asset);
}
```

## 在事务处理器功能中调用 Hyperledger Fabric API

要在事务处理器函数中调用 Hyperledger Fabric API，必须先调用 getNativeAPI 函数，然后调用来自 Hyperledger Fabric API 的函数。

使用超级分类帐 Fabric API 可以让您访问在超级分类帐编写器 API 中不可用的功能。

警告: 使用诸如 getState、putState、deleteState、getStateByPartialCompositeKey 等 Hyperledger Fabric API, getQueryResult 函数将绕过 Hyperledger Composer 访问控制规则(ACLs)。

在下面的示例中，调用了 Hyperledger Fabric API 函数 getHistoryForKey，该函数以迭代器的形式返回指定资产的历史。事务处理器函数然后将返回的数据存储在一个数组中。

~~~js
async function simpleNativeHistoryTransaction (transaction) {
    const id = transaction.assetId;
    const nativeSupport = transaction.nativeSupport;

    const nativeKey = getNativeAPI().createCompositeKey('Asset:systest.transactions.SimpleStringAsset', [id]);
    const iterator = await getNativeAPI().getHistoryForKey(nativeKey);
    let results = [];
    let res = {done : false};
    while (!res.done) {
        res = await iterator.next();

        if (res && res.value && res.value.value) {
            let val = res.value.value.toString('utf8');
            if (val.length > 0) {
                results.push(JSON.parse(val));
            }
        }
        if (res && res.done) {
            try {
                iterator.close();
            }
            catch (err) {
            }
        }
    }
}
~~~

## 从事务处理器函数返回数据

事务处理器功能可以选择将数据返回到客户端应用程序。

这对于向事务提交者返回收据或返回事务修改后的资产非常有用，以避免在事务提交后对资产进行单独查询。

数据也可以通过用于业务网络的事务 REST API 返回到客户端应用程序，例如通过 POST 方法返回数据(如下所述)到客户端应用程序。

事务处理器函数的返回数据必须是有效类型，要么是基本类型(String、Integer、Long等)，要么是使用编写器建模语言建模的类型——概念、资产、参与者、事务、事件或枚举。

还必须使用 @returns(type) 装饰器在事务模型上指定返回数据的类型，并且返回数据必须是事务处理器函数最后返回的内容。

如果一个事务有多个事务处理器函数，那么只有一个事务处理器函数可以返回数据。

如果返回的数据丢失或类型错误，则事务将失败并被拒绝。

## 从事务处理器函数返回原始类型

这是一个将字符串返回到客户端应用程序的事务处理器功能的示例。

*   模型文件

~~~js
namespace org.sample

@returns(String)
transaction MyTransaction {

}
~~~

*   交易处理器功能

~~~js
/**
 * Handle a transaction that returns a string.
 * @param {org.sample.MyTransaction} transaction The transaction.
 * @returns [string](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/string) The string.
 * @transaction
 */
async function myTransaction(transaction) {
    return 'hello world!';
}
~~~

*   客户申请

~~~js
const bnc = new BusinessNetworkConnection();
await bnc.connect('admin@sample-network');
const factory = bnc.getBusinessNetwork().getFactory();
const transaction = factory.newTransaction('org.sample', 'MyTransaction');
const string = await bnc.submitTransaction(transaction);
console.log(`transaction returned $[string](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/string)`);
~~~

## 从事务处理器函数返回复杂类型

这是一个将概念返回到客户端应用程序的事务处理器功能的示例。可以修改相同的代码以返回资产，参与者，交易或事件。

*   模型文件

~~~js
namespace org.sample

concept MyConcept {
    o String value
}

@returns(MyConcept)
transaction MyTransaction {

}
~~~

*   交易处理器

~~~js
/**
 * Handle a transaction that returns a concept.
 * @param {org.sample.MyTransaction} transaction The transaction.
 * @returns {org.sample.MyConcept} The concept.
 * @transaction
 */
async function myTransaction(transaction) {
    const factory = getFactory();
    const concept = factory.newConcept('org.sample', 'MyConcept');
    concept.value = 'hello world!';
    return concept;
}
~~~

*   客户申请

~~~js
const bnc = new BusinessNetworkConnection();
await bnc.connect('admin@sample-network');
const factory = bnc.getBusinessNetwork().getFactory();
const transaction = factory.newTransaction('org.sample', 'MyTransaction');
const concept = await bnc.submitTransaction(transaction);
console.log(`transaction returned ${concept.value}`);
~~~

这是一个事务处理器功能的示例，该函数将一系列概念返回给客户端应用程序。

## 发射事件

*   发射事件

事件可由 Hyperledger Composer 发出并由外部应用程序订阅。事件在业务网络定义的模型文件中定义，并由事务处理函数文件中的事务 JavaScript 发出。

*   定义事件

事件在.cto业务网络定义的模型文件（）中定义，与资产和参与者相同。事件使用以下格式：

~~~bash
event BasicEvent {
}
~~~

*   事件函数

为了发布事件，创建事件的事务必须调用三个函数，第一个 getFactory 函数。在 getFactory 允许作为交易的一部分被创建的事件。接下来，必须使用创建事件factory.newEvent('org.namespace', 'BasicEvent')。这会 BasicEvent 在指定的名称空间中创建一个定义。然后必须设置事件所需的属性。最后，该事件必须通过使用发射emit(BasicEvent)。调用此事件的简单事务将如下所示：

~~~bash
async function basicEventTransaction(basicEventTransaction) {
    let factory = getFactory();
 
    let basicEvent = factory.newEvent('org.namespace', 'BasicEvent');
    emit(basicEvent);
}
~~~

此事务创建并发出 BasicEvent 业务网络模型文件中定义的类型事件。有关 getFactory 函数的更多信息，请参阅 Composer API 文档。
