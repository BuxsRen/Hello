package listen

import (
	"Hello/app/events/trigger"
	"Hello/app/libs/utils"
)

// 监听事件名称 : 多个触发事件。执行顺序是 事件名称 a-z 排序
var event = map[string]interface{}{
	// 用户每日首次登录
	"user_login": []map[string]interface{}{
		{
			"async": true,         // 使用异步，否则是同步
			"fun":   trigger.Test, // 触发事件
		},
	},
}

// 触发监听事件 listen.Event("user_login",data)
func Event(e string, data map[string]interface{}) {
	if event[e] == nil {
		utils.ExitError("触发事件不存在", -1)
	}
	for _, item := range event[e].(map[string]interface{}) {
		for _, v := range item.([]map[string]interface{}) {
			fun := v["fun"].(func(map[string]interface{}))
			if v["async"].(bool) {
				go fun(data)
			} else {
				fun(data)
			}
		}
	}
}
