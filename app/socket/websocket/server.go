package websocket

import (
	"Hello/app/libs/utils"
	"Hello/app/models"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"io"
	"net/http"
	"sync"
	"time"
)

// 在线用户列表
var List sync.Map

// WebSocket 用户结构
type WebSocket struct {
	ws      *websocket.Conn
	err     error
	id      float64     // 用户Id
	msgchan chan []byte // 消息通道
	index   int         // 列表中的索引位置
}

func init() {
	go ping()
}

//CheckOrigin防止跨站点的请求伪造
var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	}, // 检测客户端请求是否合法
}

// 客户端接入
func (this *WebSocket) Chat(c *gin.Context) {
	w := &WebSocket{}
	w.ws, w.err = upGrader.Upgrade(c.Writer, c.Request, nil) //客户端连接，升级get请求为webSocket协议
	if w.err != nil {
		return
	}
	w.id = c.GetFloat64("_id")
	w.msgchan = make(chan []byte, 1024)
	w.login()
}

// 客户端加入
func (this *WebSocket) login() {
	list, ok := List.Load(this.id)
	if ok { // 异地登录，挤掉在线客户端
		ws := list.(*WebSocket)
		//onAll(this.buildMsg("logout",nil,0)) // 给所有在线的用户推送下线通知
		ws.msgchan <- this.buildMsg("logout", "您在另一处登录了，该设备被迫下线", 0)
		go func(ws *websocket.Conn) {
			time.Sleep(10 * time.Millisecond)
			_ = ws.Close()            // 另一地点登录，将在线的客户端踢下线
			List.Store(this.id, this) // 加入到在线列表
			go this.onMessage()
			go this.onSend()
		}(ws.ws)
	} else {
		List.Store(this.id, this) // 加入到在线列表
		go this.onMessage()
		go this.onSend()
	}
	//fmt.Println("[login]", this.id)
	onAll(this.buildMsg("login", nil, 0)) // 给所有在线的用户推送上线通知
	var user []float64
	f := func(k, v interface{}) bool {
		ws := v.(*WebSocket)
		user = append(user, ws.id)
		return true
	}
	List.Range(f)
	this.index = len(user)
	this.msgchan <- this.buildMsg("list", user, 0) // 发送在线用户列表
}

// 客户端离开
func (this *WebSocket) logout() {
	list, ok := List.Load(this.id)
	if ok {
		ws := list.(*WebSocket)
		if ws == this { // 正常离线，忽略被挤掉的离线
			List.Delete(this.id) // 从在线用户列表中移除该用户
			//fmt.Println("[logout]", this.id)
			onAll(this.buildMsg("logout", nil, 0)) // 给所有在线的用户推送下线通知
		}
	}
}

// 心跳
func ping() {
	for {
		time.Sleep(30 * time.Second)
		var msg = map[string]interface{}{
			"type": "ping",
			"time": utils.GetTime(),
		}
		b, _ := json.Marshal(msg)
		onAll(b)
	}
}

// 群发消息
func onAll(msg []byte) {
	f := func(k, v interface{}) bool {
		ws := v.(*WebSocket)
		if ws != nil && ws.msgchan != nil {
			ws.msgchan <- msg
		}
		return true
	}
	List.Range(f)
}

// 消息监听
func (this *WebSocket) onMessage() {
	defer this.ws.Close()
	var whiteList = map[string]int{"answer": 1, "hangUp": 1, "refuse": 1, "call": 1, "busy": 1} // type 白名单 接听，挂断，拒绝，拨打，通话忙碌
	for this.ws != nil {
		_, r, err := this.ws.NextReader() // 读取客户端的消息
		if err != nil || r == nil {
			this.logout()
			break
		}
		message, err := io.ReadAll(r)
		if this.err != nil || len(message) == 0 { // 客户端断开连接,关闭并结束对该客户端的服务
			this.logout()
			break
		}
		data := make(map[string]interface{})
		e := json.Unmarshal(message, &data)
		// 非json数据不处理
		if e != nil || data["type"] == nil {
			continue
		}
		//fmt.Println(data)
		switch data["type"].(string) {
		case "send": // 私聊
			if data["data"] == nil || fmt.Sprintf("%T", data["toId"]) != "float64" {
				continue
			}
			msg := &models.Message{}
			list, ok := List.Load(data["toId"].(float64))
			if ok { // 是否在在线线列表中
				ws := list.(*WebSocket)
				if ws != nil && ws.ws != nil { // 对方是否在线
					_ = msg.Push(data["toId"], this.id, data["toId"], data["data"], 0)
					ws.msgchan <- this.buildMsg("send", data["data"], 0) // 推送给对方
				}
			} else { // 不在线
				msg.Push(data["toId"], this.id, data["toId"], data["data"], 0)
			}
			msg.Push(this.id, this.id, data["toId"], data["data"], 1)
			this.msgchan <- this.buildMsg("send", data["data"], 0) // 推送给自己
			break
		case "list": // 获取在线用户列表
			var user []float64
			f := func(k, v interface{}) bool {
				ws := v.(*WebSocket)
				user = append(user, ws.id)
				return true
			}
			List.Range(f)
			this.msgchan <- this.buildMsg("list", user, 0) // 发送在线用户列表
			break
		case "all": // 群发
			if data["data"] == nil {
				continue
			}
			onAll(this.buildMsg("all", data["data"], 0))
			break
		default:
			break
		}

		if whiteList[data["type"].(string)] == 1 { // 通话相关
			if fmt.Sprintf("%T", data["toId"]) != "float64" {
				continue
			}
			list, ok := List.Load(data["toId"].(float64))
			if ok { // 是否在在线线列表中
				ws := list.(*WebSocket)
				if ws != nil { // 是否未断开
					ws.msgchan <- this.buildMsg(data["type"].(string), nil, 0)
				}
			}
		}
	}
}

// 写入消息
func (this *WebSocket) onSend() {
	for {
		select {
		case msg, ok := <-this.msgchan:
			if !ok { // 消息通道关闭
				fmt.Printf("[onSend] %f发送消息携程结束\n", this.id)
				return
			}
			_ = this.ws.WriteMessage(websocket.TextMessage, msg)
		}
	}
}

// 构建消息
func (this *WebSocket) buildMsg(types string, data interface{}, id int64) []byte {
	var msg = map[string]interface{}{
		"type": types,
		"id":   id,
		"data": data,
		"from": this.id,
		"time": utils.GetTime(),
	}
	b, _ := json.Marshal(msg)
	return b
}
