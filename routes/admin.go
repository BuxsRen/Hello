package routes

import (
	"github.com/gin-gonic/gin"
)

// 管理后台路由
type Admin struct {
	Route *gin.RouterGroup                                                                           // gin 路由
	Bind  func(func(string, ...gin.HandlerFunc) gin.IRoutes, string, func(*gin.Context) interface{}) // 路由与控制器绑定：请求方式，路由，控制器
}

// 入口方法
func (this *Admin) Handle() {

}

