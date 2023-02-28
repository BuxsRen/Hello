package utils

import (
	"math/rand"
	"net/url"
	"strconv"
	"sync"
)

var mu sync.Mutex

// 输出成功json数据
func JSON(data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"code":200,
		"msg":"获取成功",
		"data":data,
	}
}

// 发送列表数据。信息，数据，额外数据
func LIST(message string,data map[string]interface{}, arr map[string]interface{}) map[string]interface{} {
	var param = make(map[string]interface{})
	param["code"] = 200
	param["msg"] = message
	param["data"] = data["items"]
	param["count"] = data["count"]
	for k,v := range arr {
		param[k] = v
	}
	return param
}

// 发送成功数据。信息，数据
func OK(message string,data interface{}) map[string]interface{} {
	var param = make(map[string]interface{})
	param["code"] = 200
	param["msg"] = message
	param["data"] = data
	return param
}

// 发送失败数据。gin.Context，信息，数据
func NO(message string,code int,data map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"code": code,
		"msg":  message,
		"data": data,
	}
}

// 取随机数。最小值，最大值
func Rand(min, max int) int {
	return rand.Intn(max-min) + min
}

// 输出错误信息并退出。信息，错误码
func ExitError(message string,code int) {
	message = url.QueryEscape(message)
	panic(`{"code":`+strconv.Itoa(code)+`,"msg":"`+message+`"}`)
}

// 直接输出成功json信息并退出，传递字符串json
func ExitJson(list string) {
	list = url.QueryEscape(list)
	panic(`{"code":200,"data":"`+list+`"}`)
}

// 克隆map
func CloneMap(data map[string]interface{}) map[string]interface{} {
	clone := make(map[string]interface{})
	for k, v := range data {
		clone[k] = v
	}
	return clone
}

// 互斥锁，并发编程同步
func WithLockContext(fn func()) {
	mu.Lock()
	defer mu.Unlock()
	fn()
}