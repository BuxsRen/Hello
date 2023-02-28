package routes

import (
	"Hello/app/http/middleware"
	"Hello/bootstrap/config"
	"Hello/routes"
	"github.com/gin-gonic/gin"
)

type Route struct {
	Router     *gin.Engine                    // gin 路由
	controller func(*gin.Context) interface{} // 路由绑定的控制器
}

// 路由配置
func (this *Route) Handle() {

	// 全局中间件
	this.Router.Use((&middleware.Global{}).Handle)

	// 404
	this.Router.NoRoute(func(c *gin.Context) {
		if config.App.Server.Template {
			c.HTML(404, "404.tpl", gin.H{})
		} else {
			c.JSON(404, gin.H{"code": 404, "msg": "url not found"})
		}
	})

	// web路由配置
	web := this.Router.Group("/")                       // 配置路由组
	web.Use((&middleware.Web{}).Handle)                 // 使用中间件
	(&routes.Web{Route: web, Bind: this.bind}).Handle() // 关联子路由

	// 接口路由配置
	api := this.Router.Group("/api")
	api.Use((&middleware.Api{}).Handle)
	(&routes.Api{Route: api, Bind: this.bind}).Handle()

	// 后台管理接口路由配置
	admin := this.Router.Group("/admin")
	admin.Use((&middleware.Admin{}).Handle)
	admin.Use((&middleware.AdminRoleAuth{}).Handle)
	(&routes.Admin{Route: admin, Bind: this.bind}).Handle()

	// websocket路由
	ws := this.Router.Group("/ws")
	ws.Use((&middleware.WebSocket{}).Handle)
	(&routes.WebSocket{Route: ws}).Handle()

	// ... 声明其他路由文件
}
