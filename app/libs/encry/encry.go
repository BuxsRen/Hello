package encry

import (
	"Hello/app/libs/utils"
	"Hello/bootstrap/config"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	uuid "github.com/satori/go.uuid" // go get github.com/satori/go.uuid
	"net/url"
	"strconv"
	"time"
)

var tokenKey = "441479573" // token签名签发密钥

// base64解码。需要解码的字符串
func Base64Decode(str string) ([]byte, string) {
	b, _ := base64.StdEncoding.DecodeString(str)
	return b, string(b)
}

// base64编码，需要编码的字节码
func Base64Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

// URL编码
func UrlEncode(str string) string {
	return url.QueryEscape(str)
}

// URL解码
func UrlDecode(str string) (string, error) {
	return url.QueryUnescape(str)
}

// map/结构体 转 json
func JsonEncode(src interface{}) ([]byte, error) {
	return json.Marshal(src)
}

// json 转 map/结构体
func JsonDecode(str string, dst interface{}) error {
	return json.Unmarshal([]byte(str), &dst)
}

// 生成md5
func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// 创建token。用户id，活动id，平台id
func EncryptToken(userId, username, avatar interface{}) string {
	param := map[string]interface{}{
		"id":       userId,
		"username": username,
		"avatar":   avatar,
		"t":        utils.GetTime(),
	}
	data, _ := json.Marshal(param) // map转json
	if config.App.Other.TokenKey != "" {
		tokenKey = config.App.Other.TokenKey
	}
	param["sign"] = MD5(string(data) + tokenKey)
	data, _ = json.Marshal(param) // map转json
	return Base64Encode(data)
}

// 解析token，返回token解析后的内容。用户token
func DecryptToken(token string) map[string]interface{} {
	baseToken, _ := Base64Decode(token) //base64解码
	param := make(map[string]interface{})
	err := json.Unmarshal(baseToken, &param)
	if err != nil {
		return nil
	}
	var sign = ""
	if param["sign"] != nil {
		sign = param["sign"].(string)
		delete(param, "sign")
	}
	data, _ := json.Marshal(param) // map转json
	if config.App.Other.TokenKey != "" {
		tokenKey = config.App.Other.TokenKey
	}
	if MD5(string(data)+tokenKey) == sign { // 效验sign
		timestamp := int64(param["t"].(float64))
		day, _ := strconv.Atoi(time.Unix(timestamp, 0).In(time.FixedZone("CST", 8*3600)).Format("02"))
		if utils.GetNow().Day() != day { // token不是当天的 1 = 1
			//return nil
		}
		return param
	}
	return nil
}

// 生成UUID
func UUID() string {
	u := uuid.NewV4()
	return u.String()
}

// 解析UUID
func UUIDDecode(str string) (bool, error) {
	_, err := uuid.FromString(str)
	if err != nil {
		return false, err
	}
	return true, nil
}

// 密码
func Password(pwd string) string {
	if config.App.Other.TokenKey != "" {
		tokenKey = config.App.Other.TokenKey
	}
	return MD5(tokenKey + pwd)
}