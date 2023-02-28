package middleware

import (
	"Hello/app/libs/encry"
	"Hello/app/libs/utils"
	"Hello/bootstrap/config"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"strconv"
	"strings"
)

// 全局中间件
type Global struct {
	c          *gin.Context
	data       map[string]interface{} // 所有数据
	sourceData map[string]interface{} // 请求原始数据
}

// 入口方法
func (this *Global) Handle(c *gin.Context) {
	global := Global{c: c}
	global.header()
	global.getAllParam()
	if global.data == nil {
		global.c.Abort()
	} else {
		global.addParam()
		if global.sign() {
			global.c.Next()
		} else {
			global.c.Abort()
		}
	}
}

// 构建所有的请求参数
func (this *Global) getAllParam() {
	this.data = this.getParam()
	this.sourceData = utils.CloneMap(this.data)
}

// 取所有请求参数
func (this *Global) getParam() map[string]interface{} {
	data := make(map[string]interface{})
	for k, v := range this.c.Request.URL.Query() { // 收集Url上面的参数
		data[k] = v[0]
	}
	if this.c.Request.Method != "GET" { // 收集非GET请求的 json 参数
		if strings.Contains(this.c.Request.Header.Get("Content-Type"), "json") { // this.c.ContentType() == "application/json"
			param, _ := io.ReadAll(this.c.Request.Body)
			this.c.Request.Body = io.NopCloser(bytes.NewBuffer(param)) // 数据重新写回 body
			e := json.Unmarshal(param, &data)
			if e != nil {
				if config.App.Server.Template {
					this.c.HTML(400, "400.tpl", gin.H{"msg": e.Error()})
				} else {
					this.c.JSON(400, gin.H{"msg": e.Error(), "data": nil})
				}
				return nil
			}
		} else {
			if strings.Contains(this.c.Request.URL.Path, "/api") && !strings.Contains(this.c.Request.URL.Path, "upload") { // api，且不是上传接口，只能传json
				if config.App.Server.Template {
					this.c.HTML(405, "405.tpl", gin.H{"msg": "Only Support Json Content Type"})
				} else {
					this.c.JSON(405, gin.H{"code": 405, "msg": "Only Support Json Content Type", "data": nil})
				}
				return nil
			} else { // 收集非GET请求的 form-data 参数
				_ = this.c.Request.ParseMultipartForm(128)
				for k, v := range this.c.Request.PostForm {
					data[k] = v[0]
				}
			}
		}
	}
	return data
}

// 设置头部信息
func (this *Global) header() {
	this.c.Header("Go-Server", "gin/"+gin.Version)
}

// 追加参数
func (this *Global) addParam() {
	var ip string
	if this.c.Request.Header.Get("ALI-CDN-REAL-IP") != "" {
		ip = this.c.Request.Header.Get("ALI-CDN-REAL-IP")
	} else {
		ip = this.c.ClientIP()
	}
	var user_agent string
	if len(this.c.Request.Header["User-Agent"]) != 0 {
		user_agent = this.c.Request.Header["User-Agent"][0]
	}
	t, ts, te := utils.GetTodayTime()
	this.data["_t"] = t           // 当前时间戳
	this.data["_ts"] = ts         // 今日开始时间戳
	this.data["_te"] = te         // 今日结束时间戳
	this.data["_ip"] = ip         // 请求IP
	this.data["_ua"] = user_agent // 浏览器标识
	this.c.Set("_data", this.data)
	this.c.Set("_source_data", this.sourceData)
}

// 签名效验
func (this *Global) sign() bool {
	data := this.sourceData
	clientSign := data["sign"]
	if clientSign != nil {
		delete(data, "sign")
		for k, v := range data {
			if k[0:1] == "_" { // 忽略 _ 开头的参数
				delete(data, k)
			}
			if k == "num" {
				T := fmt.Sprintf(`%T`, v) // 取变量类型
				if T == "float64" {       // json 类型
					data[k] = string(rune(v.(float64))) // 字母转ascii
				} else {
					iv, _ := strconv.Atoi(v.([]string)[0])
					data[k] = string(rune(iv)) // 字母转ascii
				}
			}
		}
		t := utils.ParamToString(data["t"])
		sign := encry.MD5(utils.HttpBuildQuery(data))
		it, _ := strconv.ParseInt(t, 10, 64)
		if utils.GetTime()-it > 60 {
			this.c.JSON(200, gin.H{"code": -1, "msg": "sign timeout"})
			return false
		}
		if sign != clientSign {
			this.c.JSON(200, gin.H{"code": -1, "msg": "sign error"})
			return false
		}
		data["sign"] = sign
	}
	return true
}
