package controllers

import (
	"Hello/app/libs/code"
	"Hello/app/libs/redis"
	"Hello/app/libs/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

type ApiController struct{}

// 图形验证码
func (this *ApiController) GetCode(c *gin.Context) interface{} {
	m := &code.Code{}
	s := m.CreateStrVerifyCode()
	rdb := redis.Redis{}
	rdb.Setex("Verify_Code_"+s.Id, s.Str, 500)
	return utils.OK("获取成功，验证码分钟有效", gin.H{
		"id":    s.Id,
		"thumb": s.Thumb,
	})
}

// 应用检查更新
func (this *ApiController) AppUpdate(c *gin.Context) interface{} {
	data := utils.GetAllSourceData(c)
	utils.VerifyData(data, map[string]string{
		"version": "required|numeric",
	})
	num := 3
	version, _ := strconv.Atoi(data["version"].(string))
	if version < num {
		return utils.OK("检查到有新版本", `更新日志：
	1.可以修改个人资料了
	2.优化了语音通话的一些逻辑
	3.圈子、通知系统雏形展示`)
	}
	return utils.NO("已是最新版本", -1, nil)
}
