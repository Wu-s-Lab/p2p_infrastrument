package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetIPList(w http.ResponseWriter,  r *http.Request){

	ipList,err := ipm.getIPListFromDB()
	if err!=nil{
		rm := ReturnIPList{
			Statu:  "GetIPList From DB error." + err.Error(),
			IPList: ipList ,
		}
		rmBytes, _ := json.Marshal(rm)
		_, _ = w.Write(rmBytes)
		return
	}

	fmt.Println("ipList: ",ipList)

	rm := ReturnIPList{
		Statu:   "ok",
		IPList:  ipList,
	}
	rmBytes, _ := json.Marshal(rm)
	_, _ = w.Write(rmBytes)
}