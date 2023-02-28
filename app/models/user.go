package models

import (
	"Hello/app/libs/encry"
	"Hello/app/libs/utils"
	"fmt"
	"github.com/gohouse/gorose/v2"
)

// 活动表
type User struct{}

func (this *User) TableName() string {
	return "users"
}

// 根据用户名和密码查询用户
func (this *User) CheckUser(data map[string]interface{}) gorose.Data {
	res, e := DB(this).Where("username", data["username"]).Where("password", encry.Password(data["password"].(string))).Where("delete_at", 0).First()
	if e != nil {
		return nil
	}
	if data["product"] != nil && data["manufacturer"] != nil {
		device := utils.ParamToString(data["manufacturer"]) + " " + utils.ParamToString(data["product"])
		_, _ = DB(this).Where("id", res["id"]).Data(map[string]interface{}{"device": device, "last_login": utils.GetTime()}).Update()
	}
	return res
}

// 根据用户名
func (this *User) CheckUserName(username string) gorose.Data {
	res, e := DB(this).Where("username", username).Where("delete_at", 0).First()
	if e != nil {
		return gorose.Data{}
	}
	return res
}

// 注册用户
func (this *User) Create(data map[string]interface{}) int64 {
	ins := map[string]interface{}{
		"username":  data["username"],
		"nickname":  "Hello" + fmt.Sprintf("%v", utils.Rand(100000, 999999)),
		"password":  encry.Password(data["password"].(string)),
		"birthday":  utils.GetTime(),
		"create_at": utils.GetTime(),
	}
	id, e := DB(this).Data(ins).Insert()
	if e != nil {
		utils.ExitError(e.Error(), -1)
	}
	return id
}

// 用户列表
func (this *User) List(data map[string]interface{}) map[string]interface{} {
	table := DB(this).Fields("id", "username", "nickname", "avatar", "is_ban", "cover", "star", "sex", "birthday", "identity", "address")
	m := Model{Model: DB(this), Table: table, Data: data}
	res := m.PageSearch()
	msg := Message{}
	list := res["items"].([]gorose.Data)
	for k, v := range list {
		list[k]["un_read_count"] = msg.GetUnReadMessageCount(data["_id"], v["id"])
	}
	res["items"] = list
	return res
}

// 点赞
func (this *User) Star(id interface{}) {
	_, _ = DB(this).Execute("update h_users set star=star+1 where id=?", id)
}

// 用户信息
func (this *User) Info(data map[string]interface{}) gorose.Data {
	info, e := DB(this).Where("id", data["id"]).Fields("id", "username", "nickname", "create_at", "avatar", "cover", "star", "sex", "birthday", "identity", "is_ban", "info", "device", "last_login", "address").First()
	if data["id"] == utils.ParamToString(data["_id"]) && data["product"] != nil && data["manufacturer"] != nil {
		device := utils.ParamToString(data["manufacturer"]) + " " + utils.ParamToString(data["product"])
		_, _ = DB(this).Where("id", data["id"]).Data(map[string]interface{}{"device": device, "last_login": utils.GetTime()}).Update()
	}
	if e != nil {
		return gorose.Data{}
	}
	return info
}

// 更新资料
func (this *User) Update(id interface{}, data map[string]interface{}) error {
	delete(data, "id")
	delete(data, "username")
	delete(data, "password")
	delete(data, "star")
	delete(data, "create_at")
	delete(data, "delete_at")
	fmt.Println(data)
	_, e := DB(this).Where("id", id).Data(data).Update()
	return e
}
