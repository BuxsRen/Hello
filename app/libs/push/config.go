package push

var fun = map[string]func(message string){
	"bark": bark,
	"dingTalk": dingTalk,
	"dingTalkMarkDown": dingTalkMarkDown,
	"wechat": wechat,
}