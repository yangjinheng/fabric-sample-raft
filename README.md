## 端口规划

*   Orderer

| 服务     | 端口  |
| -------- | ----- |
| orderer0 | 30001 |
| orderer1 | 30002 |
| orderer2 | 30003 |

*   Org1

| 服务       | 端口                  |
| ---------- | --------------------- |
| ca-org1    | 30100                 |
| peer0-org1 | 30101 / 30102 / 30103 |
| peer1-org1 | 30104 / 30105 / 30106 |
| peer2-org1 | 30107 / 30108 / 30109 |
| peer3-org1 | 30110 / 30111 / 30112 |

*   Org2

| 服务       | 端口                  |
| ---------- | --------------------- |
| ca-org2    | 30200                 |
| peer0-org2 | 30201 / 30202 / 30203 |
| peer1-org2 | 30204 / 30205 / 30206 |
| peer2-org2 | 30207 / 30208 / 30209 |
| peer3-org2 | 30210 / 30211 / 30212 |

## 启动顺序

~~~bash
1.  Namespace
2.  PVC
3.  Orderer
4.  Peer
5.  Ca
6.  Cli
~~~

## 文件组织

## 启动网络

1.  生成证书
2.  生成通道
3.  准备NFS

~~~bash
/data *(rw,fsid=0,sync,no_subtree_check,no_auth_nlm,insecure,no_root_squash)
~~~

## 故障记录

*   关于实例化 node 链码速度很慢，由于 node 链码实例化时候需要下载 npm 包，而下载时候使用了默认的 npm 源，所以慢，这需要修改 peer 源码，如下：

~~~go
// GenerateDockerBuild Build Docker image
func (nodePlatform *Platform) GenerateDockerBuild(path string, code []byte, tw *tar.Writer) error {
	codepackage := bytes.NewReader(code)
	binpackage := bytes.NewBuffer(nil)
	registry := cutil.GetDockerfileFromConfig("chaincode.node.runtime.registry")
	err := util.DockerBuild(util.DockerBuildOptions{
		Cmd: fmt.Sprintf("cp -R /chaincode/input/src/. /chaincode/output && cd /chaincode/output && npm install --production --registry=%s", registry),
		InputStream:  codepackage,
		OutputStream: binpackage,
	})
	if err != nil {
		return err
	}
	return cutil.WriteBytesToPackage("binpackage.tar", binpackage.Bytes(), tw)
}

//GetMetadataProvider fetches metadata provider given deployment spec
func (nodePlatform *Platform) GetMetadataProvider(code []byte) platforms.MetadataProvider {
	return &ccmetadata.TargzMetadataProvider{Code: code}
}
~~~

~~~bash
源码修改内容为：在实例化链码的时候从环境变量中取得 npm 源，实例化时候向 npm 传参，环境变量设置为：

- name: CORE_CHAINCODE_NODE_RUNTIME_REGISTRY
  value: https://registry.npm.taobao.org
~~~

