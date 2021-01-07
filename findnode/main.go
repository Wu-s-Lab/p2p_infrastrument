package main

import (
	"encoding/json"
	"fmt"
	"errors"
	"net"
	"net/http"
	"strings"
	"time"
)

type Peer struct{
	ID string               `json:"id"`
	InetIPAddress string    `json:"inet_ip_address"`
	MacAddress    string    `json:"mac_address"`
	EnetIPAddress string    `json:"enet_ip_address"`
	HInt string             `json:"h_int"`
	LastSeen string         `json:"last_seen"`
}

type ReturnMessage struct {
	Code string   `json:"code"` // error  ok
	Data string   `json:"data"`
}

func main(){

	fmt.Println("服务器启动监听,监听IP：8.131.232.232:443")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
	// 注册处理函数，用户连接，自动调用指定的处理函数
	http.HandleFunc("/GetIPAddressList",GetIPAddressList)
	// 监听绑定
	err := http.ListenAndServe(":443",nil)
	if err!=nil{
		fmt.Println("服务器监听失败：",err)
	}
}

// post接口接收json数据
func GetIPAddressList(w http.ResponseWriter,  r *http.Request)  {
    // 获取节点信息
	peers,err := GetIPAddressListFromSuperNode()
	if err!=nil {
		fmt.Println("GetIPAddressList From SuperNode ERROR: ",err)
		message := "GetIPAddressList From SuperNode ERROR: " + err.Error()
		rm := ReturnMessage{Code:"error",Data: message}
		rmByte,_:=json.Marshal(rm)
		_, _ = w.Write(rmByte)
		return
	}
	peersBytes,err := json.Marshal(peers)
	if err != nil {
		fmt.Println("json.Marshal peers error: ",err)
		message := "json.Marshal peers error: " + err.Error()
		rm := ReturnMessage{Code:"error",Data: message}
		rmByte,_:=json.Marshal(rm)
		_, _ = w.Write(rmByte)
		return
	}
	// 返回数据
	rm := ReturnMessage{Code:"ok",Data: string(peersBytes)}
	rmByte,_:=json.Marshal(rm)
	_, _ = w.Write(rmByte)
}

// 从本地supernode获取节点信息
func GetIPAddressListFromSuperNode() ([]Peer, error) {
	var peers []Peer

	// 创建连接
	socket, err := net.DialUDP("udp4", nil, &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 5645,
	})
	if err != nil {
		fmt.Println("连接失败!", err)
		return peers, err
	}
	defer socket.Close()

	// 发送数据
	senddata := []byte("\r")
	_, err = socket.Write(senddata)
	if err != nil {
		fmt.Println("发送数据失败!", err)
		return peers, err
	}

	// 设置数据读取的结束时间 2s
	socket.SetReadDeadline(time.Now().Add(2*time.Second))
	var tmpPeers []Peer
	count := 0

    // 接收数据
	for{
		data := make([]byte, 4096)
		_, _, err := socket.ReadFromUDP(data)
		if err != nil {
			fmt.Println("read over.")
			if count == 0 {
				return peers, errors.New("Connect To SuperNode Error.")
			}
			break
		}

		// 从第三个开始接收
		count++
		if count <= 2 {
			continue
		}
        // 处理字符串
		arr := strings.Fields(string(data))
		peer := Peer{ID:arr[0],InetIPAddress:arr[1],MacAddress:arr[2],EnetIPAddress:arr[3],HInt:arr[4],LastSeen:arr[5]}
		tmpPeers = append(tmpPeers,peer)
        fmt.Println("arr: ",arr)
	}

	// 删除最后一行无用数据，如果元素数量大于0
	length := len(tmpPeers)
	if length > 0 {
		peers = append(peers,tmpPeers[:length - 1]...)
	}

	return peers, nil
}


