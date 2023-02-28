package main

import (
	"Hello/app/libs/utils"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var Min = 1
var Max = 100

// WebSocket 客户端
type WebSocket struct {
	conn    *websocket.Conn // ws连接
	message chan []byte     // 发送消息通道
	status  bool            // ws状态
	timeout int             // 重连超时时间
	token   string          // 登录token 通过登录接口获取
	id      int
}

func main() {
	for i := Min; i <= Max; i++ {
		ws := newClient(i, EncryptToken(i, "Break", "test"))
		ws.connect()
		//time.Sleep(10*time.Millisecond)
	}
	fmt.Println("所有人全部上线")
	time.Sleep(86400 * time.Second)
	//var token = "eyJhaWQiOjE2LCJpZCI6MjM1MjQxLCJwbGF0Zm9ybSI6MzEwMTAxLCJ0IjoxNjQ3Mzk5ODg0LCJzaWduIjoiNmQ1YWU0ZTdkZDA2M2QwYTk3ZTNkZmE2OTI0ODMwMjEifQ=="
	//args := os.Args
	//if len(args) >= 2 {
	//	token = args[1]
	//}
	//ws := newClient(token)
	//ws.run()
}

// 构建客户端
func newClient(id int, token string) *WebSocket {
	var sendChan = make(chan []byte, 1024)
	var conn *websocket.Conn
	return &WebSocket{
		conn:    conn,
		message: sendChan,
		timeout: 86400,
		token:   token,
		id:      id,
	}
}

// 启动客户端，断开自动重新启动
func (this *WebSocket) run() {
	for {
		if this.status == false {
			this.connect()
		}
		time.Sleep(time.Second * time.Duration(this.timeout))
	}
}

// 连接服务端
func (this *WebSocket) connect() {
	var err error
	var rep *http.Response
	//this.conn, rep, err = websocket.DefaultDialer.Dial("ws://172.31.36.46:60000/ws/chat?token="+this.token,nil) // 连接服务端
	this.conn, rep, err = websocket.DefaultDialer.Dial("ws://127.0.0.1:9310/ws/chat?token="+this.token, nil) // 连接服务端
	if err != nil {
		if rep != nil {
			fmt.Println(rep.Header.Get("Http-Code"), rep.Header.Get("Http-Msg"))
		}
		fmt.Printf("\x1b[%dm ☹ [%d]连接主程序失败 \x1b[0m\n", 31, this.id)
		return
	}
	fmt.Println(this.id)
	this.status = true
	//go this.onSend()
	go this.onMessage()
}

// 读取消息
func (this *WebSocket) onMessage() {
	for {
		_, message, err := this.conn.ReadMessage()
		if err != nil { // 出现错误，退出读取，尝试重连
			fmt.Printf("\x1b[%dm x %s \x1b[0m\n", 31, err.Error())
			_ = this.conn.Close()
			this.status = false
			fmt.Println(err)
			break
		}
		if len(message) == 0 {
			continue
		}
		data := make(map[string]interface{})
		e := json.Unmarshal(message, &data)
		if e != nil {
			continue
		}
		if data["type"] == "send" {
			log.Printf("%s\n", string(message))
		}
	}
}

// 发送消息
func (this *WebSocket) onSend() {
	go func() {
		for {
			//time.Sleep(time.Duration(utils.Rand(1000,2500))*time.Millisecond)
			this.message <- this.buildMsg("send", "100,999", int64(utils.Rand(Min, Max)))
		}
	}()
	for {
		select {
		case msg, ok := <-this.message:
			if !ok {
				return
			}
			_ = this.conn.WriteMessage(websocket.TextMessage, msg)
		}
	}
}

// 构建消息
func (this *WebSocket) buildMsg(types string, data interface{}, id int64) []byte {
	var msg = map[string]interface{}{
		"type": types,
		"toId": id,
		"data": data,
	}
	b, _ := json.Marshal(msg)
	return b
}

// 创建token。用户id，活动id，平台id
func EncryptToken(userId, username, avatar interface{}) string {
	param := map[string]interface{}{
		"id":       userId,
		"username": username,
		"avatar":   avatar,
		"t":        utils.GetTime(),
	}
	data, _ := json.Marshal(param) // map转json
	param["sign"] = MD5(string(data) + "DcV8JaEiCBefsTYI2")
	data, _ = json.Marshal(param) // map转json
	return Base64Encode(data)
}

// base64编码，需要编码的字节码
func Base64Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

// 生成md5
func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}