package api

import (
	"Hello/app/libs/code"
	"Hello/app/libs/email"
	"Hello/app/libs/encry"
	"Hello/app/libs/redis"
	"Hello/app/libs/upload"
	"Hello/app/libs/utils"
	"Hello/app/models"
	"Hello/bootstrap/config"
	"github.com/gin-gonic/gin"
	"strings"
)

// 用户相关接口控制器
type UserController struct{}

// 用户登录
func (this *UserController) Login(c *gin.Context) interface{} {
	data := utils.GetAllSourceData(c)
	utils.VerifyData(data, map[string]string{
		"username": "required|string|between:3,18",
		"password": "required|string|alpha_dash|between:5,18",
	})
	u := models.User{}
	user := u.CheckUser(data)
	if user == nil {
		return utils.NO("用户名或密码错误", -1, nil)
	}
	if user["is_ban"].(int64) == 1 {
		return utils.NO("该用户禁止登录", -38, nil)
	}
	user["token"] = encry.EncryptToken(user["id"], data["username"], data["avatar"])
	delete(user, "delete_at")
	delete(user, "update_at")
	delete(user, "is_ban")
	return utils.OK("登录成功", user)
}

// 注册
func (this *UserController) Register(c *gin.Context) interface{} {
	data := utils.GetAllSourceData(c)
	utils.VerifyData(data, map[string]string{
		"username": "required|string|email",
		"code":     "required|string|between:1,5",
		"password": "required|string|alpha_dash|between:5,19",
		"repeat":   "required|string|alpha_dash|between:5,19",
	})
	if data["password"] != data["repeat"] {
		return utils.NO("两次密码不一致", -1, nil)
	}
	rdb := redis.Redis{}
	if strings.ToUpper(rdb.Get("Emali_Code_"+data["username"].(string))) != strings.ToUpper(data["code"].(string)) {
		return utils.NO("验证码不正确", -1, nil)
	}
	u := models.User{}
	user := u.CheckUserName(data["username"].(string))
	if user != nil {
		return utils.NO("该账号已被注册", -1, nil)
	}
	u.Create(data)
	return utils.OK("注册成功", nil)
}

// 发送邮箱验证码
func (this *UserController) Emali(c *gin.Context) interface{} {
	data := utils.GetAllSourceData(c)
	utils.VerifyData(data, map[string]string{
		"username": "required|string|email",
		"code":     "required|string|between:1,5",
		"id":       "required|string",
	})
	rdb := redis.Redis{}
	verify := rdb.Get("Verify_Code_" + data["id"].(string))
	if strings.ToUpper(verify) != strings.ToUpper(data["code"].(string)) {
		return utils.NO("验证码不正确", -1, nil)
	}
	e := email.New()
	e.SetTitle("Hello，注册")
	e.SetToEmail([]string{data["username"].(string)})
	s := (&code.Code{}).CreateStrVerifyCode()
	e.SetBody("<h1>您的验证码是:" + strings.ToUpper(s.Str) + "<br>5分钟内有效</h1>")
	rdb.Setex("Emali_Code_"+data["username"].(string), s.Str, 500)
	err := e.SendMail()
	if err != nil {
		return utils.OK(err.Error(), err.Error())
	}
	return utils.OK("邮箱验证码发送成功", nil)
}

// 用户列表
func (this *UserController) List(c *gin.Context) interface{} {
	data := utils.GetAllData(c)
	u := models.User{}
	list := u.List(data)
	return utils.OK("获取成功", list)
}

// 点赞
func (this *UserController) Star(c *gin.Context) interface{} {
	data := utils.GetAllSourceData(c)
	utils.VerifyData(data, map[string]string{
		"id": "required|numeric",
	})
	u := models.User{}
	u.Star(data["id"])
	return utils.OK("点赞成功", nil)
}

// 获取用户信息
func (this *UserController) Info(c *gin.Context) interface{} {
	data := utils.GetAllData(c)
	utils.VerifyData(data, map[string]string{
		"id": "required|numeric",
	})
	u := models.User{}
	info := u.Info(data)
	if info["is_ban"] != nil && info["is_ban"].(int64) == 1 {
		return utils.NO("该用户禁止登录", -38, nil)
	}
	return utils.OK("获取成功", info)
}

// 消息至已读
func (this *UserController) MessageToRead(c *gin.Context) interface{} {
	data := utils.GetAllData(c)
	utils.VerifyData(data, map[string]string{
		"id": "required|numeric",
	})
	m := models.Message{}
	m.MessageToRead(data)
	return utils.OK("操作成功", nil)
}

// 获取聊天记录
func (this *UserController) GetMessageList(c *gin.Context) interface{} {
	data := utils.GetAllData(c)
	utils.VerifyData(data, map[string]string{
		"id": "required|numeric",
	})
	m := models.Message{}
	list := m.GetMessageList(data)
	return utils.OK("获取成功", list)
}

// 更新资料
func (this *UserController) UpdateInfo(c *gin.Context) interface{} {
	id := utils.GetInput(c, "_id", false)
	data := utils.GetAllSourceData(c)
	u := models.User{}
	return utils.OK("资料更新成功", u.Update(id, data))
}

// 上传
func (this *UserController) Upload(c *gin.Context) interface{} {
	data := utils.GetAllData(c)
	utils.VerifyData(data, map[string]string{
		"format": "required|string",
	})
	file, err := c.FormFile("file")
	if err != nil {
		utils.ExitError("请选择文件上传", -1)
	}
	f := upload.Upload{File: file, Format: data["format"].(string)}
	path := f.Upload()
	prefix := config.App.Other.PublicPrefix
	if prefix == "" {
		utils.ExitError("未设置 public_prefix", -1)
	}
	url := config.App.Server.Url
	return utils.OK("上传成功", url+prefix+path)
}
