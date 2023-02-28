package udpsocket

import (
	"Hello/app/libs/encry"
	"Hello/app/libs/utils"
	"Hello/bootstrap/config"
	"Hello/bootstrap/helper"
	"encoding/json"
	"fmt"
	"net"
	"sync"
)

var List sync.Map

type User struct {
	id   float64
	addr *net.UDPAddr
}

// UDP Server端
func Run(wg *sync.WaitGroup) {
	server, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: config.App.Server.UdpPort,
	})
	fmt.Printf("📡 Listen Udp Server -> 0.0.0.0:%d\n", config.App.Server.UdpPort)
	if err != nil {
		fmt.Println("➦ " + err.Error())
		(&helper.Helper{}).Exit("✘ This Port Is Already In Use !", 3)
		return
	}
	fmt.Printf("► OK! Start Udp Service...\n\n")
	wg.Done()
	defer server.Close()
	var whiteList = map[string]int{"call": 1} // type 白名单 通话确认，通话中
	for {
		var d [10240]byte
		n, addr, err := server.ReadFromUDP(d[:]) // 接收数据
		if err != nil {
			fmt.Println("读取客户端数据失败：", err)
			continue
		}
		message := d[:n]
		data := make(map[string]interface{})
		e := json.Unmarshal(message, &data)
		if e != nil || data["type"] == nil || data["token"] == nil || fmt.Sprintf("%T", data["token"]) != "string" { // 非json数据不处理
			continue
		}
		token := data["token"].(string)
		info := encry.DecryptToken(token)
		if info == nil {
			continue
		}

		if data["type"].(string) == "login" { // 登录
			List.Store(info["id"].(float64), addr)
			err := send(server, buildMsg(info["id"], "login", nil), addr)
			if err != nil {
				continue
			}
		}
		if whiteList[data["type"].(string)] == 1 {
			if data["toId"] == nil || fmt.Sprintf("%T", data["toId"]) != "float64" {
				continue
			}
			res, ok := List.Load(data["toId"].(float64))
			if ok {
				ip := res.(*net.UDPAddr)
				err := send(server, message, ip)
				if err != nil {
					continue
				}
			}
		}
	}
}

func send(server *net.UDPConn, message []byte, ip *net.UDPAddr) error {
	_, err := server.WriteToUDP(message, ip)
	if err != nil {
		fmt.Println("写入数据失败: ", err)
		return err
	}
	return nil
}

// 构建消息
func buildMsg(id interface{}, types string, data interface{}) []byte {
	var msg = map[string]interface{}{
		"type": types,
		"data": data,
		"from": id,
		"time": utils.GetTime(),
	}
	b, _ := json.Marshal(msg)
	return b
}
