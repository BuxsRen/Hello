package utils

import (
	"strings"
	"unicode/utf8"
)

//取字符串左边 原文本 起始位置 取出数量
func Left(str string,start int,end int) string {
	return str[start:start+end]
}

// 取字符串右边： 原文本，起始位置，取出数量
func Right(str string,start int,end int) string{
	return str[len(str)-start-end:len(str)-start]
}

// 取字符串中间：原文本，开始文本，结束文本
func GetSubstr(str, start, end string) string {
	n := strings.Index(str, start)
	if n == -1 {
		n = 0
	} else {
		n = n + len(start)  // 增加了else，不加的会把start带上
	}
	str = string([]byte(str)[n:])
	m := strings.Index(str, end)
	if m == -1 {
		m = len(str)
	}
	str = string([]byte(str)[:m])
	return str
}

// ascii转字符 65 => A
func Chr(ascii int) string {
	return string(rune(ascii))
}

// 字符转ascii A => 65
func Ord(chr string) int {
	c, _ := utf8.DecodeRune([]byte(chr))
	return int(c)
}

// 元素是否在数组中
func InArray(target string, array []string) bool {
	for _, item := range array{
		if target == item{
			return true
		}
	}
	return false
}