package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"unicode/utf8"
)

// 正则匹配
type Regexp struct {}

// 验证必填
func (this *Regexp) Required(str string) bool {
	return str != ""
}

// 字符串是否是纯数字
func (this *Regexp) Numeric(str string) bool {
	_ ,e:= strconv.Atoi(str)
	return e == nil
}

// 是否是字符串
func (this *Regexp) String(str interface{}) bool {
	return fmt.Sprintf("%T",str) == "string"
}

// 字符串是否在指定长度之内
func (this *Regexp) Between(str string,min,max int) bool {
	return utf8.RuneCountInString(str) > min && utf8.RuneCountInString(str) < max
}

// 字符串(数字)是否在指定范围大小之间
func (this *Regexp) Size(str string,min,max int) bool {
	num ,_:= strconv.Atoi(str)
	return 	num > min && num < max
}

// 字符串是否是字母构成
func (this *Regexp) Alpha(str string) bool {
	reg ,_ := regexp.MatchString(`[a-zA-Z]+`, str)
	return reg
}

// 字符串是否是字母和数字构成
func (this *Regexp) AlphaNum(str string) bool {
	reg := regexp.MustCompile(`[a-zA-Z]+`)
	word := reg.FindAllStringSubmatch(str,1)
	reg = regexp.MustCompile(`\d+`)
	number:= reg.FindAllStringSubmatch(str,1)
	return len(word) > 0 && len(number) > 0
}

// 字符串只能是字母或数字或中划线或特殊字符构成
func (this *Regexp) AlphaDash(str string) bool {
	reg ,_ := regexp.MatchString(`^[a-zA-Z0-9-*/+.~!@#$%^&*()]+$`, str)
	return reg
}

// 字符串必须全部由字母、数字、特殊字符构成
func (this *Regexp) AlphaDashAll(str string) bool {
	word_1:= regexp.MustCompile(`[a-zA-Z]+`).FindAllStringSubmatch(str,1)
	word_2:= regexp.MustCompile(`[\d+]`).FindAllStringSubmatch(str,1)
	word_3:= regexp.MustCompile(`[-*/+.~!@#$%^&*()]`).FindAllStringSubmatch(str,1)
	return len(word_1) > 0 && len(word_2) > 0 &&len(word_3) > 0
}

// 密码 以字母开头，至少需要包含一个数字，6-20位之间，可以有特殊字符组成的密码
func (this *Regexp) Password(str string) bool {
	length ,_ := regexp.MatchString(`^.{6,20}$`, str)
	matched, _ := regexp.MatchString("(.*([a-zA-Z].*))(.*[0-9].*)[a-zA-Z0-9-*/+.~!@#$%^&*()]{0,20}$", str)
	return matched && length
}

// 匹配日期 2004-04-30 | 2004-02-29
func (this *Regexp) Date(str string) bool {
	reg ,_ := regexp.MatchString(`^[0-9]{4}-(((0[13578]|(10|12))-(0[1-9]|[1-2][0-9]|3[0-1]))|(02-(0[1-9]|[1-2][0-9]))|((0[469]|11)-(0[1-9]|[1-2][0-9]|30)))$`, str)
	return reg
}

// 匹配时间 10:23 | 12:00:00
func (this *Regexp) Time(str string) bool {
	reg ,_ := regexp.MatchString(`^(([0-1]?[0-9])|([2][0-3])):([0-5]?[0-9])(:([0-5]?[0-9]))?$`, str)
	return reg
}

// 匹配日期时间 2022-01-13 10:23:14
func (this *Regexp) DateTime(str string) bool {
	reg ,_ := regexp.MatchString(`^[0-9]{4}-(((0[13578]|(10|12))-(0[1-9]|[1-2][0-9]|3[0-1]))|(02-(0[1-9]|[1-2][0-9]))|((0[469]|11)-(0[1-9]|[1-2][0-9]|30))) (([0-1]?[0-9])|([2][0-3])):([0-5]?[0-9])(:([0-5]?[0-9]))?$`, str)
	return reg
}

// 匹配链接
func (this *Regexp) Url(str string) bool {
	reg ,_ := regexp.MatchString(`(https?|ftp|file)://[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]`, str)
	return reg
}

// 匹配邮箱
func (this *Regexp) Email(str string) bool {
	reg ,_ := regexp.MatchString(`^[A-Za-z0-9\\u4e00-\\u9fa5]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`, str)
	return reg
}

// 验证手机号
func (this *Regexp) Mobile(str string) bool {
	reg ,_ := regexp.MatchString(`^(((13[0-9]{1})|(14[0-9]{1})|(15[0-9]{1})|(16[0-9]{1})|(17[0-9]{1})|(19[0-9]{1})|(18[0-9]{1}))+\d{8})$`, str)
	return reg
}

// 验证身份证号
func (this *Regexp) IdNumber(str string) bool {
	reg1 := regexp.MustCompile(`^([1-6][1-9]|50)\d{4}(18|19|20)\d{2}((0[1-9])|10|11|12)(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx]$`) // 18 位
	reg2 := regexp.MustCompile(`^([1-6][1-9]|50)\d{4}\d{2}((0[1-9])|10|11|12)(([0-2][1-9])|10|20|30|31)\d{3}$`) // 15 位
	word1:= reg1.FindAllStringSubmatch(str,1)
	word2:= reg2.FindAllStringSubmatch(str,1)
	return len(word1) > 0 || len(word2) > 0
}

// 是否包含汉字
func (this *Regexp) Char(str string) bool {
	reg ,_ := regexp.MatchString(`^[\\u4e00-\\u9fa5]{0,}$`, str)
	return reg
}
