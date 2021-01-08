# IPCertServer

## 编译说明

- 将项目放在GOPATH下
- 安装第三方库
```
go get github.com/mattn/go-sqlite3

go get gopkg.in/yaml.v2

go get github.com/pkg/errors
```

- go build进行编译


## 使用说明

创建config.yaml文件，内容如下：

```
Port:          "9876"
Community:     "p2pmain"
Key:           "p2pmainpswd"
Supernode:     "8.131.232.232:7777"
```

### 使用命令如下：

```
// 初始化并启动服务
IPCertServer -i init -c start

// 只启动服务
IPCertServer -c start   

// 只初始化
IPCertServer -i init

// -i init  生成ca根证书和私钥并初始化数据库(会删除原有的)
// -c start 启动监听
```
