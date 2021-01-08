package main

import (
	"IPCertServer/cert"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func Register(w http.ResponseWriter,  r *http.Request){

	csrBytes, err := ioutil.ReadAll(r.Body)
	if err!=nil {
		fmt.Println("ioutil.ReadAll(r.Body) err: ",err)
		rm := ReturnMessage{
			Statu:"read csr bytes error.",
			Data:[]byte(""),
		}
		rmBytes,_:=json.Marshal(rm)
		_, _ = w.Write(rmBytes)
		return
	}

	// 注册获取证书
	certBytes, err := cert.CreateCertWithCsr(csrBytes)
	if err!=nil {
		fmt.Println("create cert err: ", err)
		rm := ReturnMessage{
			Statu:"create cert error.",
			Data:[]byte(""),
		}
		rmBytes,_:=json.Marshal(rm)
		_, _ = w.Write(rmBytes)
		return
	}

	// 分配IP地址
	ip, err := ipm.iPAddressAssignment()
	if err!=nil {
		fmt.Println("IPAddress Assignment error: ", err)
		rm := ReturnMessage{
			Statu: "IPAddress Assignment error.",
			Data:[]byte(""),
		}
		rmBytes,_:=json.Marshal(rm)
		_, _ = w.Write(rmBytes)
		return
	}

	fmt.Println("分配的ip: ", ip)

	// 保存ip和cert
	err = ipm.saveIPAddress(ip, certBytes)
	if err!=nil{
		fmt.Println("SaveIPAddress error: ", err)
		rm := ReturnMessage{
			Statu: "Save IPAddress error.",
			Data:  []byte(""),
		}
		rmBytes,_:=json.Marshal(rm)
		_, _ = w.Write(rmBytes)
		return
	}

	rm := ReturnMessage{
		Statu:  "ok",
		Data:   certBytes,
	}
	rmBytes,_:=json.Marshal(rm)
	_, _ = w.Write(rmBytes)
}