package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

var reg Regexp

var rule = map[string]interface{}{
	"required":map[string]interface{}{  // "required"
		"info":"参数[%s]不能为空",
		"fun":reg.Required,
	},
	"numeric":map[string]interface{}{  // "numeric"
		"info":"参数[%s]类型必须为数字",
		"fun":reg.Numeric,
	},
	"string":map[string]interface{}{
		"info":"参数[%s]类型必须为字符串",
		"fun":reg.String,
	},
	"between":map[string]interface{}{ // "between:5,10"
		"info":"参数[%s]的长度必须大于%d且小于%d",
		"fun":reg.Between,
	},
	"size":map[string]interface{}{ // "size:5,10"
		"info":"参数[%s]必须大于%d且小于%d",
		"fun":reg.Size,
	},
	"alpha":map[string]interface{}{
		"info":"参数[%s]必须是字母构成",
		"fun":reg.Alpha,
	},
	"alpha_num":map[string]interface{}{
		"info":"参数[%s]必须是字母和数字构成",
		"fun":reg.AlphaNum,
	},
	"alpha_dash":map[string]interface{}{
		"info":"参数[%s]必须是字母，数字或特殊字符其中一种构成",
		"fun":reg.AlphaDash,
	},
	"alpha_dash_all":map[string]interface{}{
		"info":"参数[%s]必须全部由字母、数字、特殊字符构成",
		"fun":reg.AlphaDashAll,
	},
	"password":map[string]interface{}{
		"info":"参数[%s]必须以字母开头，长度在6~18之间，只能包含字母、数字和下划线",
		"fun":reg.Password,
	},
	"date":map[string]interface{}{
		"info":"参数[%s]必须是日期格式",
		"fun":reg.Date,
	},
	"time":map[string]interface{}{
		"info":"参数[%s]必须是时间格式",
		"fun":reg.Time,
	},
	"date_time":map[string]interface{}{
		"info":"参数[%s]必须是日期时间格式",
		"fun":reg.DateTime,
	},
	"url":map[string]interface{}{
		"info":"参数[%s]必须是链接",
		"fun":reg.Url,
	},
	"email":map[string]interface{}{
		"info":"参数[%s]必须是正确的邮箱地址",
		"fun":reg.Email,
	},
	"mobile":map[string]interface{}{
		"info":"参数[%s]必须是正确的手机号",
		"fun":reg.Mobile,
	},
	"id_number":map[string]interface{}{
		"info":"参数[%s]必须是正确的身份证号",
		"fun":reg.IdNumber,
	},
}

// 参数对应中文解释
var attributes = map[string]string{
	"username":"用户名",
	"password":"密码",
	"mobile":"手机号",
	"title":"标题",
}

// 验证json参数 c,&v，第二个传结构体,规则绑定在tag上，只支持post
func VerifyParam(c *gin.Context,v interface{}) {
	if err := c.ShouldBindJSON(v); err != nil {
		ExitError("请求参数错误["+err.Error()+"]",-1)
	}
}

// 验证data参数 (data,{"date":"required|date"})
func VerifyData(data map[string]interface{},v map[string]string) {
	for param,item := range v { // 遍历请求参数列表
		var list []string
		list = strings.Split(item, "|") // 需要验证的规则列表
		for _,p := range list{
			var b,l []string
			b = strings.Split(p, ":") // 规则是否有子项
			var min,max int
			if len(b) > 1 {
				p = b[0]
				l = strings.Split(b[1], ",") // 子项的子项
				min,_ = strconv.Atoi(l[0])
				max,_ = strconv.Atoi(l[1])
			}
			if rule[p] == nil {
				continue
			}
			str := ParamToString(data[param]) // 需要验证的参数
			var source = param
			if attributes[param] != "" { // 请求参数中文对照
				source = attributes[param]
			}
			r := rule[p].(map[string]interface{})
			info := r["info"].(string)
			if p == "between" || p == "size" {
				fun := r["fun"].(func (str string,min,max int) bool)
				if !fun(str,min,max) { // 验证不通过，输出失败原因
					ExitError(fmt.Sprintf(info,source,min,max),-1)
				}
			} else if p == "string" {
				fun := r["fun"].(func (str interface{}) bool)
				if !fun(data[param]) {
					ExitError(fmt.Sprintf(info,source),-1)
				}
			} else {
				fun := r["fun"].(func (str string) bool)
				if !fun(str) {
					ExitError(fmt.Sprintf(info,source),-1)
				}
			}
		}
	}
}
