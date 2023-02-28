package middleware

import "github.com/gin-gonic/gin"

// 后台用户鉴权中间件
type AdminRoleAuth struct {
	c *gin.Context
}

// 入口方法
func (this *AdminRoleAuth) Handle(c *gin.Context) {
	auth := AdminRoleAuth{c:c}
	auth.c.Next()
}