package middleware

import (
	"Hello/app/libs/encry"
	"Hello/app/libs/redis"
	"Hello/app/libs/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

type Api struct {
	c        *gin.Context
	id       int    // 用户id
	token    string // 用户token
	username string //	 用户名
}

// 入口方法
func (this *Api) Handle(c *gin.Context) {
	api := Api{c: c}
	api.checkUser()
	//api.checkRedisUser()
	api.addParam()
	api.c.Next()
}

// 用户验证
func (this *Api) checkUser() {
	token := this.c.Request.Header.Get("Id")
	param := encry.DecryptToken(token)
	if param != nil {
		this.id = int(param["id"].(float64))
		this.token = token
		this.username = param["username"].(string)
	} else {
		utils.ExitError("登录超时了，请重新登录[错误代码:-99]", -99)
	}
}

// 用户验证 redis用户验证
func (this *Api) checkRedisUser() {
	rdb := redis.Redis{}
	day := strconv.Itoa(utils.GetNow().Day())
	info := rdb.Get("User_" + day + "_" + this.token)
	if info == "" {
		utils.ExitError("登录超时了，请重新登录哦!", -99)
	}
	id, _ := strconv.Atoi(info)
	this.id = id
}

// 追加参数
func (this *Api) addParam() {
	data := utils.GetAllData(this.c)
	data["_id"] = this.id
	data["_token"] = this.token
	data["_username"] = this.username
	this.c.Set("_data", data)
}

