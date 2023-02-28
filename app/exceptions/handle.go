package exceptions

import (
	"Hello/app/libs/encry"
	"Hello/app/libs/push"
	"Hello/bootstrap/config"
	"Hello/bootstrap/helper"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/url"
	"runtime/debug"
)

// 异常处理中间件
type exception struct {
	c        *gin.Context
	e        interface{} // 原始错误
	str      string      // 原始错误转字符串
	debug    bool
	env      string
	template bool
	push     bool
	info     map[string]interface{}
}

// 入口方法
func Handle(c *gin.Context) {
	defer func() {
		e := recover() // 捕获异常
		if e != nil {
			ep := &exception{
				c:        c,
				e:        e,
				debug:    config.App.Server.Debug,
				env:      config.App.Server.Env,
				template: config.App.Server.Template,
				push:     config.App.Push.Use,
			}
			ep.str = ep.string()
			ep.start()
		}
		c.Set("_data", "")
		c.Set("_source_data", "")
		c.Abort()
	}()
	c.Next()
}

// 开始处理
func (this *exception) start() {
	err := json.Unmarshal([]byte(this.str), &this.info)
	if err != nil { // 系统抛出的异常
		this.system()
	} else {
		this.user()
	}
}

// 系统抛出的异常
func (this *exception) system() {
	if this.debug && this.env == "local" { // 开启debug，并且是开发环境
		fmt.Printf("\n\x1b[%dm [ERROR] [%s] \x1b[0m\n\n", 31, this.str)
		debug.PrintStack() // 输出到终端
		if this.template { // 使用了模板
			this.c.HTML(500, "500.tpl", gin.H{"msg": this.str})
		} else {
			this.c.JSON(500, gin.H{"code": 500, "msg": this.str, "data": nil})
		}
	} else {               // 写入到文件中
		if this.template { // 使用了模板
			this.c.HTML(500, "500.tpl", gin.H{"msg": "Server Error"})
		} else {
			this.c.JSON(500, gin.H{"code": 500, "msg": "Server Error", "data": nil})
		}
		text := "\n" + this.str + "\n" + string(debug.Stack()) + "\n"
		_ = helper.LOG.Output(1, text) // 错误输出到文件
	}
	if this.push && this.env == "production" { // 开启了错误信息推送
		data := this.c.GetStringMap("_data")
		data["error"] = this.str
		content, _ := json.Marshal(data)
		go push.New().Push(string(content)) // 执行错误信息推送
	}
}

// 用户抛出的异常
func (this *exception) user() {
	code := this.info["code"].(float64)
	if int(code) == 200 {
		str, _ := url.QueryUnescape(this.info["data"].(string))
		var list interface{}
		e := encry.JsonDecode(str, &list)
		if e != nil {
			this.c.JSON(500, gin.H{"code": 500, "msg": e.Error(), "data": nil})
		} else {
			this.c.JSON(200, gin.H{"code": 200, "msg": "获取成功", "data": list})
		}
	} else {
		str, _ := url.QueryUnescape(this.info["msg"].(string))
		this.c.JSON(200, gin.H{"code": int(code), "msg": str, "data": nil})
	}
}

// error 转 string
func (this *exception) string() string {
	switch v := this.e.(type) {
	case error:
		return v.Error()
	default:
		return this.e.(string)
	}
}
