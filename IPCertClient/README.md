# 使用说明

```
./IPCertClient -c init -s 8.131.232.232:9876 -f ./crypto -n lixiaopeng -p 132000000001

./IPCertClient -c getEdgeConf -s 8.131.232.232:9876 -f ./crypto/client.crt -e edge.conf

./IPCertClient -c init -s 127.0.0.1:9876 -f ./crypto -n lixiaopeng -p 132000000001

./IPCertClient -c getEdgeConf -s 127.0.0.1:9876 -f ./crypto/client.crt -e edge.conf

// - f 证书文件存放路径
// - n 名字
// - p 电话
// - c 操作命令     init: 初始化获取证书    getEdgeConf: 获取edge配置文件
// - e 配置文件存放路径
// - s 服务器IP和端口
```

