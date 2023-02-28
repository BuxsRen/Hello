package push

import (
	"Hello/bootstrap/config"
)

// 推送
type push struct {
	mode string // 推送方式
}

// 消息推送初始化
/**
 * @Example:
	p := push.New()
	p.Push("test")
*/
func New() *push {
	p := new(push)
	p.mode = config.App.Push.Mode
	if config.App.Push.Mode == "" {
		p.mode = "bark"
	}
	return p
}

// 推送信息。推送的内容
func (this *push) Push(message string) {
	if fun[this.mode] == nil {
		panic("不支持的推送方式")
	}
	fun[this.mode](message)
}

// 更改推送方式，只对当前实例有效。模式(bark，dingTalk，dingTalkMarkDown，wechat)
func (this *push) SetPushMode(mode string) {
	this.mode = mode
}
