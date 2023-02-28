package routes

import (
	"Hello/app/http/controllers/api"
	"github.com/gin-gonic/gin"
)

// 应用接口路由
type Api struct {
	Route *gin.RouterGroup                                                                           // gin 路由
	Bind  func(func(string, ...gin.HandlerFunc) gin.IRoutes, string, func(*gin.Context) interface{}) // 路由与控制器绑定：请求方式，路由，控制器
}

// 入口方法
func (this *Api) Handle() {
	this.Bind(this.Route.GET, "/user/list", (&api.UserController{}).List)
	this.Bind(this.Route.POST, "/user/star", (&api.UserController{}).Star)
	this.Bind(this.Route.GET, "/user/info", (&api.UserController{}).Info)
	this.Bind(this.Route.GET, "/message/list", (&api.UserController{}).GetMessageList)
	this.Bind(this.Route.POST, "/message/read", (&api.UserController{}).MessageToRead)
	this.Bind(this.Route.POST, "/user/update", (&api.UserController{}).UpdateInfo)
	this.Bind(this.Route.POST, "/upload", (&api.UserController{}).Upload)
}
