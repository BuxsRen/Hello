package routes

import (
	"Hello/app/http/controllers"
	"Hello/app/http/controllers/api"
	"Hello/bootstrap/config"
	"github.com/gin-gonic/gin"
)

// web路由
type Web struct {
	Route *gin.RouterGroup                                                                           // gin 路由
	Bind  func(func(string, ...gin.HandlerFunc) gin.IRoutes, string, func(*gin.Context) interface{}) // 路由与控制器绑定：请求方式，路由，控制器
}

// 入口方法
func (this *Web) Handle() {

	// 默认页
	this.Route.GET("/", func(c *gin.Context) {
		if config.App.Server.Template {
			c.HTML(200, "index.tpl", gin.H{})
		} else {
			c.String(200, "Hello World")
		}
	})

	this.Bind(this.Route.POST, "/api/login", (&api.UserController{}).Login)

	this.Bind(this.Route.POST, "/api/register", (&api.UserController{}).Register)

	this.Bind(this.Route.GET, "/api/get/code", (&controllers.ApiController{}).GetCode)

	this.Bind(this.Route.POST, "/api/emali/code", (&api.UserController{}).Emali)

	this.Bind(this.Route.GET, "/api/check/update", (&controllers.ApiController{}).AppUpdate)

}
