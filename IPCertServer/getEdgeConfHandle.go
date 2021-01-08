package main

import (
	"IPCertServer/cert"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetEdgeConf(w http.ResponseWriter,  r *http.Request){
	certByte, err := ioutil.ReadAll(r.Body)
	if err!=nil {
		fmt.Println("ioutil.ReadAll(r.Body) err: ",err)
		rm := ReturnMessage{
			Statu:"read cert bytes error.",
			Data:[]byte(""),
		}
		rmBytes,_:=json.Marshal(rm)
		_, _ = w.Write(rmBytes)
		return
	}

	// 验证证书是否为根证书签发
	isRight := cert.VerifyCAChain(certByte)
	if !isRight{
		fmt.Println("cert is not from root cert.")
		rm := ReturnMessage{
			Statu:"cert is not from root cert.",
			Data:[]byte(""),
		}
		rmBytes,_:=json.Marshal(rm)
		_, _ = w.Write(rmBytes)
		return
	}

	// 获取分配的IP地址
	IP, err := ipm.getIPAddress(certByte)
	if err!=nil{
		fmt.Println("ip get err: ",err)
		rm := ReturnMessage{
			Statu: "get ip error: " + err.Error(),
			Data:[]byte(""),
		}
		rmBytes,_:=json.Marshal(rm)
		_, _ = w.Write(rmBytes)
		return
	}

	// 创建EdgeConf
	ec := EdgeConf{
		Community: Community,
		IPAddress: IP,
		Key: Key,
		Supernode: Supernode,
	}
	fmt.Println(ec)
	ecBytes,_:= json.Marshal(ec)

	rm := ReturnMessage{
		Statu: "ok",
		Data:  ecBytes,
	}
	rmBytes,_:=json.Marshal(rm)
	_, _ = w.Write(rmBytes)
}
