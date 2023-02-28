package routes

import (
	"Hello/app/socket/websocket"
	"github.com/gin-gonic/gin"
)

// websocket路由
type WebSocket struct {
	Route *gin.RouterGroup // gin 路由
	// WebSocket 不需要绑定返回数据处理
}

// 入口方法
func (this *WebSocket) Handle() {
	this.Route.GET("chat", (&websocket.WebSocket{}).Chat) //
}
