package main

import (
	"IPCertServer/cert"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
)

type ReturnMessage struct {
	Statu string  `json:"statu"`
	Data  []byte  `json:"data"`
}

type ReturnIPList struct {
	Statu    string  `json:"statu"`
	IPList  []string  `json:"ip_list"`
}

type EdgeConf struct {
	Community string `json:"community"`
	IPAddress string `json:"ip_address"`
	Key       string `json:"key"`
	Supernode string `json:"supernode"`
}

type Config struct {
	Port       string  `yaml:"Port"`
	Community  string  `yaml:"Community"`
	Key        string  `yaml:"Key"`
	Supernode  string  `yaml:"Supernode"`
}

var Community  string
var Key        string
var Supernode  string
var Port       string

var ipm *IPManagement

func main() {
	cmd := flag.String("c", "", "start")
	init := flag.String("i", "", "init")
	flag.Parse()

	// 获取配置文件数据
	var setting Config
	config, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		fmt.Println("read config.yaml failed.")
		return
	}
	err1 := yaml.Unmarshal(config, &setting)
	if err1 != nil{
		fmt.Println("Get config from config.yaml failed.")
		return
	}

	Community = setting.Community
	Key = setting.Key
	Supernode = setting.Supernode
	Port = setting.Port

	if Community==""||Key==""||Supernode==""||Port==""{
		fmt.Println("Community or Key or Supernode or Port is null.")
		return
	}

	ipm, err = NewIPManagement("192.168.176",10,240)
	if err!=nil{
		fmt.Println("create IPManagement error.")
        return
	}

	if *init =="init"{
		Init()
	}
	if *cmd =="start"{
		start()
	}

	//switch *cmd{
	//	case "init":
	//		Init()
	//	case "start":
	//		start()
	//	default:
	//		fmt.Println("No this cmd. please input the right cmd")
	//		fmt.Println("IPCertServer -c init       生成ca根证书和私钥(会删除原有的)")
	//		fmt.Println("IPCertServer -c start      启动IP地址分配与证书获取服务器监听")
	//}
}

func start(){
	// 判断是否已生成根证书和私钥
	// 检测根证书和根私钥是否存在
	rootCertFilePath := "./crypto/ca.crt"
	rootPrivateKeyFilePath := "./crypto/ca.key"

	if cert.PathIsExist(rootCertFilePath) == false {
		fmt.Println("no root cert.")
		return
	}

	if cert.PathIsExist(rootPrivateKeyFilePath) == false {
		fmt.Println("no root private key.")
		return
	}

	fmt.Println("服务器启动监听,监听IP端口: 0.0.0.0:" + Port)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
	// 注册处理函数，证书与IP地址生成
	http.HandleFunc("/register", Register)
	// 注册处理函数，使用证书获取IP地址
	http.HandleFunc("/getEdgeConf", GetEdgeConf)
	// 注册处理函数，查询已分配IP地址
	http.HandleFunc("/getIPList", GetIPList)
	// 监听绑定
	http.ListenAndServe("0.0.0.0:" + Port,nil)
}


func Init(){

	// 生成根证书和私钥
	certInfo := cert.CertInformation{
		Country:            []string{"China"},
		Organization:       []string{"buaa"},
		OrganizationalUnit: []string{"www.buaa.edu.cn"},
		EmailAddress:       []string{"wlkjaq@buaa.edu.cn"},
		StreetAddress:      []string{"37"},
		Province:           []string{"Beijing"},
		Locality:           []string{"haidian"},
		SubjectKeyId:       []byte{6, 5, 4, 3, 2, 1},
	}
	rootCertFilePath := "./crypto/ca.crt"
	rootPrivateKeyFilePath := "./crypto/ca.key"

	// 清空原有的文件
	_ = os.RemoveAll("./crypto")

	err := cert.CreateRootCertAndRootPrivateKey(certInfo,rootCertFilePath,rootPrivateKeyFilePath)
	if err != nil{
		fmt.Println("Create Root Cert And Root Private Key Error: ",err)
		panic("Failed To Create Root Cert And Root Private Key.")
	}

	fmt.Println("=============== Create RootCert And RootPrivateKey Successful ===============")

	err = DBInit(ipm)
	if err!= nil{
		fmt.Println(err)
	}
}

func DBInit(ipm *IPManagement) error {

	// 创建表
	sql_table := `
    CREATE TABLE IF NOT EXISTS IPCert(
       id INTEGER PRIMARY KEY AUTOINCREMENT,
       ip VARCHAR(256) NULL,
       cert BLOB NULL
    );
    `

	ipm.DB.Exec(sql_table)

	return nil
}