package middleware

import (
	"Hello/app/libs/encry"
	"github.com/gin-gonic/gin"
)

// WebSocket路由中间件
type WebSocket struct {
	c *gin.Context
}

// 入口方法
func (this *WebSocket) Handle(c *gin.Context) {
	ws := WebSocket{c: c}
	token := ws.c.Query("token")
	data := encry.DecryptToken(token)
	if data == nil {
		c.Header("http-msg", "login timeout")
		c.Header("http-code", "-99")
		c.JSON(-99, gin.H{"code": -99, "msg": "login timeout"})
		c.Abort()
	} else {
		c.Set("_id", data["id"])
		c.Next()
	}
}
