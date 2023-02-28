package push

import (
	"Hello/app/libs/net"
	"Hello/bootstrap/config"
	"encoding/json"
	"net/url"
)

// bark 推送
func bark(message string) {
	n := net.New(config.App.Push.BarkUrl+url.QueryEscape(message), "GET", "")
	n.Do()
}

// 钉钉推送
func dingTalk(message string) {
	data := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]interface{}{
			"content": message,
		},
	}
	param, _ := json.Marshal(data)
	n := net.New(config.App.Push.DingTalkUrl, "POST", string(param))
	n.SetHeader("Content-Type", "application/json;charset=utf-8").Do()
}

// 钉钉推送 markdown
func dingTalkMarkDown(message string) {
	data := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]interface{}{
			"title": "标题",
			"text":  message,
		},
	}
	param, _ := json.Marshal(data)
	n := net.New(config.App.Push.DingTalkUrl, "POST", string(param))
	n.SetHeader("Content-Type", "application/json;charset=utf-8").Do()
}

// 企业微信推送
func wechat(message string) {
	data := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]interface{}{
			"content": message,
		},
	}
	param, _ := json.Marshal(data)
	n := net.New(config.App.Push.WechatUrl, "POST", string(param))
	n.SetHeader("Content-Type", "application/json;charset=utf-8").Do()
}