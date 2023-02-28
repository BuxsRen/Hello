package net

import (
	"Hello/app/libs/utils"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type net struct {
	c           *http.Client
	req         *http.Request
	returnToMap bool
}

// 网络类
/**
 * @Example:
	c := net.New("http://127.0.0.1","POST","a=1&b=2")
	s := c.Do()
	fmt.Println(s)
*/
func New(url, method, data string) *net {
	var n = net{}
	n.c = &http.Client{}
	n.req, _ = http.NewRequest(method, url, strings.NewReader(data))
	n.req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return &n
}

// 设置请求地址
/**
 * @param scheme string 协议 http/https
 * @param host string 域名
 * @Example:
	c := net.New()
	c.SetRequestUrl("http","127.0.0.1")
*/
func (this *net) SetRequestUrl(scheme, host string) *net {
	this.req.URL.Scheme = scheme
	this.req.Host = host
	return this
}

// 设置请求头
/**
 * @param key string 键
 * @param val string 值
 * @Example:
	c := net.New()
	c.SetHeader("Content-Type","application/json")
*/
func (this *net) SetHeader(key, val string) *net {
	if this.req == nil {
		panic("请先初始化")
	}
	this.req.Header.Set(key, val)
	return this
}

// 设置请求方式.GET/POST 等
func (this *net) SetMethod(method string) *net {
	this.req.Method = method
	return this
}

// 设置返回数据是否json转map
func (this *net) SetReturnToMap(j bool) *net {
	this.returnToMap = j
	return this
}

// 发送请求，并返回请求数据
func (this *net) Do() interface{} {
	if this.req == nil {
		utils.ExitError("请先初始化:", -1)
	}
	res, err := this.c.Do(this.req)
	if err != nil {
		utils.ExitError("接口请求失败:"+err.Error(), -1)
	}
	body, e := io.ReadAll(res.Body)
	if e != nil {
		utils.ExitError("请求结果读取失败:"+e.Error(), -1)
	}
	defer res.Body.Close()
	if this.returnToMap {
		data := make(map[string]interface{})
		e := json.Unmarshal(body, &data)
		if e != nil {
			utils.ExitError("数据解析失败:"+e.Error(), -1)
		}
		return data
	}
	return string(body)
}
