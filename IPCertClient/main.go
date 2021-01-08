package main

import (
	"IPCertClient/cert"
	"encoding/json"
	"encoding/pem"
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type ReturnMessage struct {
	Statu string  `json:"statu"`
	Data  []byte  `json:"data"`
}

type EdgeConf struct {
	Community string `json:"community"`
	IPAddress string `json:"ip_address"`
	Key       string `json:"key"`
	Supernode string `json:"supernode"`
}

var ServerUrl string

func main() {
	cmd := flag.String("c", "", "cmd: init start")
	path := flag.String("f", "", "file path to save cert")
	edge := flag.String("e", "", "file path to save edge conf")
	name := flag.String("n", "", "your name")
	phone := flag.String("p", "", "your phone")
	server := flag.String("s", "", "server url")

	flag.Parse()

	if *cmd==""{
		fmt.Println("please input cmd.")
		os.Exit(-1)
	}

	if *server==""{
		fmt.Println("please input server url")
		os.Exit(-1)
	}

	ServerUrl = "http://" + *server

	//time.Sleep(3*time.Second)

	switch *cmd{
	case "init":
		if *path==""||*name==""||*phone==""{
			fmt.Println("please input all of path, name, phone.")
			os.Exit(-1)
		}
		Init(*path,*name,*phone)
		os.Exit(0)

	case "getEdgeConf":
		if *path==""||*edge==""{
			fmt.Println("please input cert file path and edge conf file path .")
			os.Exit(-1)
		}

		err := GetEdgeConf(*path, *edge)
		if err!=nil{
			fmt.Println("err: ", err)
			os.Exit(-1)
		}
		os.Exit(0)

	default:
		fmt.Println("No this cmd. please input the right cmd")
		fmt.Println("IPCertClient -c init      生成公私钥和csr并向服务器申请cert")
		fmt.Println("IPCertClient -c getEdgeConf -f ./crypto  使用证书查询IP地址")
		os.Exit(-1)
	}
}

// 生成公私钥和csr+申请证书
func Init(filePath, name, phone string){
	// 判断是否已经注册

	err := cert.CreateKeyAndCSR(filePath, name, phone)
    if err!=nil {
    	fmt.Println("Create Key And CSR ERROR: ", err)
		os.Exit(-1)
	}

	err = GetCertFromCAByCsr(filePath)
	if err!=nil {
		fmt.Println("Get Cert From CA ByCsr ERROR: ", err)
		os.Exit(-1)
	}
}

// 使用证书获取EdgeConf
func GetEdgeConf(certFilePath, edgeConfFilePath string) error {

	// 获取本地证书文件
	if cert.PathIsExist(certFilePath) == false {
		fmt.Println("cert file not exist.")
		return errors.New("cert file not exist.")
	}

	certPEMBlock, _ := ioutil.ReadFile(certFilePath)
	certDERBlock, _ := pem.Decode(certPEMBlock)
	if certDERBlock == nil {
		return errors.New("Failed to read cert.")
	}
	certBytes := certDERBlock.Bytes

	rmBytes, err := SendRep(ServerUrl+"/getEdgeConf", string(certBytes))
	if err!=nil {
		return errors.New("get cert from ca server error: " + err.Error())
	}

	var rm ReturnMessage
	err = json.Unmarshal(rmBytes,&rm)
	if err!=nil {
		return errors.New("Failed to json unmarshal return message.")
	}
	if rm.Statu!="ok"{
		return errors.New("get EdgeConf from ca server error: " + rm.Statu)
	}

	var ec EdgeConf
	err = json.Unmarshal(rm.Data,&ec)
	if err!=nil {
		return errors.New("Failed to json unmarshal EdgeConf.")
	}

	// 写入EdgeConf配置文件
	dstFile, err := os.Create(edgeConfFilePath)
	if err!=nil {
		fmt.Println(err.Error())
		return err
	}
	defer dstFile.Close()
	dstFile.WriteString("-c=" + ec.Community + "\n")
	dstFile.WriteString("-k=" + ec.Key + "\n")
	dstFile.WriteString("-a=" + ec.IPAddress + "\n")
	dstFile.WriteString("-l=" + ec.Supernode + "\n")

	return nil
}

// 使用csr向ca申请证书
func GetCertFromCAByCsr(filePath string) error {
	csrFilePath := filePath + "/client.csr"
	certFilePath := filePath + "/client.crt"

	csrPEMBlock, _ := ioutil.ReadFile(csrFilePath)
	csrDERBlock, _ := pem.Decode(csrPEMBlock)
	if csrDERBlock == nil {
		return errors.New("Failed to read csr.")
	}

	csrBytes := csrDERBlock.Bytes
	rmBytes, err := SendRep(ServerUrl+"/register", string(csrBytes))
	if err!=nil {
		return errors.New("get cert from ca server error: " + err.Error())
	}
	var rm ReturnMessage
	err = json.Unmarshal(rmBytes, &rm)
	if err!=nil {
		return errors.New("Failed to json unmarshal return message.")
	}
	if rm.Statu!="ok"{
		return errors.New("get cert from ca server error: " + rm.Statu)
	}

	// 保存证书到本地文件
	err = cert.Write(certFilePath, "CERTIFICATE", rm.Data)
	if err != nil {
		return errors.New("write cert error: " + err.Error())
	}

	return nil
}

func SendRep(url string, data string) ([]byte, error) {
	// 发送post请求
	res, err := http.Post(url, "application/json;charset=utf-8", strings.NewReader(data))
	if err != nil || res == nil {
		fmt.Println("Request error: ", err)
		return nil, err
	}

	defer res.Body.Close()

	// 读取服务器返回消息
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Fatal error: ", err)
		return nil, err
	}
	return content, err
}
