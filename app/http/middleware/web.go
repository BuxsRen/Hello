package middleware

import "github.com/gin-gonic/gin"

// Web路由中间件
type Web struct {
	c *gin.Context
}

// 入口方法
func (this *Web)Handle(c *gin.Context){
	web := Web{c:c}
	web.c.Next()
}