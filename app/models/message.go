package models

import (
	"Hello/app/libs/utils"
)

// 聊天记录表
type Message struct{}

func (this *Message) TableName() string {
	return "message"
}

// 获取未读消息记录数
func (this *Message) GetUnReadMessageCount(uid, from interface{}) int64 {
	list, e := DB(this).Fields("id", "from_id", "users_id", "content", "create_at").Where("users_id", uid).Where("from_id", from).Where("is_read", 0).Count()
	if e != nil {
		return 0
	}
	return list
}

// 追加聊天记录
func (this *Message) Push(uid, from, to, content interface{}, read interface{}) int64 {
	id, e := DB(this).Data(map[string]interface{}{
		"users_id":  uid,
		"from_id":   from,
		"to_id":     to,
		"content":   content,
		"is_read":   read,
		"create_at": utils.GetTime(),
	}).InsertGetId()
	if e != nil {
		return 0
	}
	return id
}

// 获取聊天记录
func (this *Message) GetMessageList(data map[string]interface{}) map[string]interface{} {
	table := DB(this)
	id := data["id"]
	table = table.Where("users_id", data["_id"]).Where("create_at", ">=", utils.GetTime()-86400*7).Where(func() {
		table.Where("from_id", id).OrWhere("to_id", id)
	}).Fields("id", "is_read", "from_id", "users_id", "content", "create_at")
	delete(data, "id")
	m := Model{Model: DB(this), Table: table, Data: data, Order: "id", OrderType: "desc"}
	return m.PageSearch()
}

// 消息至已读
func (this *Message) MessageToRead(data map[string]interface{}) {
	_, _ = DB(this).Where("users_id", data["_id"]).Where("from_id", data["id"]).Data(map[string]interface{}{"is_read": 1}).Update()
}
