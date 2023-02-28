package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/url"
	"strconv"
	"time"
)

// 将 interface{} 类型转成 string 类型
func ParamToString(val interface{}) string {
	var str string
	T := fmt.Sprintf(`%T`, val)
	switch T {
	case "float64":str = strconv.FormatFloat(val.(float64),'f',-1,64)
	case "[]string":str = val.([]string)[0]
	case "string":str = val.(string)
	case "int":str = strconv.Itoa(val.(int))
	case "int64":str = strconv.Itoa(int(val.(int64)))
	case "[]uint8":str = string(val.([]uint8))
	case "time.Time":str = val.(time.Time).Format("2006-01-02 15:04:05")
	}
	return str
}

// map 转 url 请求参数
func HttpBuildQuery(data map[string]interface{}) string {
	var uri url.URL
	q := uri.Query()
	for k,v := range data {
		q.Add(k, ParamToString(v))
	}
	return q.Encode()
}

// 取请求参数 参数key为*取所有，参数source表示是否为原始参数
func GetInput(c *gin.Context, key string,source bool) interface{} {
	data := make(map[string]interface{})
	if source {
		data = GetAllSourceData(c)
	} else {
		data = GetAllData(c)
	}
	if key == "*" {
		return data
	} else {
		return data[key]
	}
}

// 取所有请求数据
func GetAllData(c *gin.Context) map[string]interface{} {
	return CloneMap(c.GetStringMap("_data"))
}

// 取所有请求原数据
func GetAllSourceData(c *gin.Context) map[string]interface{} {
	return CloneMap(c.GetStringMap("_source_data"))
}
