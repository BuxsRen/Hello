package routes

import (
	"Hello/app/libs/utils"
	"Hello/bootstrap/config"
	"Hello/bootstrap/helper"
	"fmt"
	"github.com/gin-gonic/gin"
	"reflect"
	"runtime"
	"time"
)

// 逐一将路由绑定到控制器 (gin路由(GET,POST...)，路由地址，指向控制器)
func (this *Route) bind(method func(string, ...gin.HandlerFunc) gin.IRoutes, url string, controller func(*gin.Context) interface{}) {
	if config.App.Server.Debug {
		m := utils.GetSubstr(runtime.FuncForPC(reflect.ValueOf(method).Pointer()).Name(), "ithub.com/gin-gonic/gin.(*RouterGroup).", "-fm")
		c := utils.GetSubstr(runtime.FuncForPC(reflect.ValueOf(controller).Pointer()).Name(), "/", "-fm")
		defer fmt.Printf("\u001B[%dm[RouteBing] %s    %s   ->   /%s \u001B[0m\n\n", 34, m, url, c)
	}
	r := &Route{controller: controller}
	method(url, r.handler) // 绑定路由地址和控制器到gin路由
}

// 执行控制器，获取结果并响应
func (this *Route) handler(c *gin.Context) {
	s := utils.GetNow()
	data := this.controller(c)
	c.Header("Rep-Time", time.Since(s).String()) // 控制器处理消耗时间
	t := fmt.Sprintf("%T", data)
	switch t {
	case "string": // 字符串处理 return "test"
		c.String(200, data.(string))
	case "helper.Views": // 模板处理 return helper.View("404.tpl",nil)
		temp := data.(helper.Views)
		c.HTML(200, temp.Template, temp.Data)
	case "<nil>":
	default: // 默认json处理 return map[string]interface{}{"code":200}
		c.JSON(200, data)
	}
}
