package middleware

import (
	"Hello/app/libs/redis"
	"Hello/app/libs/utils"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

// 后台管理路由中间件
type Admin struct {
	c     *gin.Context
	id    int                    // 用户id
	token string                 // 用户token
	info  map[string]interface{} // 原始用户信息
}

// 入口方法
func (this *Admin) Handle(c *gin.Context) {
	admin := Admin{c: c}
	if !admin.checkUser() {
		admin.c.Abort()
	} else {
		admin.addParam()
		admin.c.Next()
	}
}

// 验证用户
func (this *Admin) checkUser() bool {
	return true
	token := this.c.Request.Header.Get("X-Token")
	if token == "" {
		token = this.c.Query("token")
	}
	rdb := redis.Redis{}
	info := rdb.Get("Admin_" + token)
	if info == "" {
		this.c.JSON(200, gin.H{
			"code": -99,
			"msg":  "还没登录或登录状态超时",
		})
		return false
	}
	data := make(map[string]interface{})
	_ = json.Unmarshal([]byte(info), &data)
	this.id = int(data["id"].(float64))
	this.info = data
	this.token = token
	return true
}

// 追加参数
func (this *Admin) addParam() {
	data := utils.GetAllData(this.c)
	data["_id"] = this.id
	data["_token"] = this.token
	data["_info"] = this.info
	this.c.Set("_data", data)
}
