package message

import (
	"Hello/app/libs/utils"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi" // 手动 go get 一下
)

var aliyun = map[string]interface{}{
	"text": map[string]string{
		"accessKeyId":     "test",
		"accessKeySecret": "test",
		"signName":        "测试",
		"templateCode":    "SMS_1111111",
	},
}

// 阿里云短信
type ALiYun struct{}

// 发送阿里云短信
/**
 * @param mobile string 手机号
 * @param code string 验证码
 * @param templet string 模板
 * @Example:
	yun := message.ALiYun{}
	yun.Send("13800138000","123456","test")
*/
func (this *ALiYun) Send(mobile, code, templet string) {
	if aliyun[templet] == "" {
		utils.ExitError("短信模板不存在", -1)
	}
	config := aliyun[templet].(map[string]string)
	client, err := dysmsapi.NewClientWithAccessKey("ap-northeast-1", config["accessKeyId"], config["accessKeySecret"])

	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = mobile                     //接收短信的手机号码
	request.SignName = config["signName"]             //短信签名名称
	request.TemplateCode = config["templateCode"]     //短信模板ID
	request.TemplateParam = `{"code":"` + code + `"}` //短信模板变量对应的实际值，JSON格式

	_, err = client.SendSms(request)
	if err != nil {
		utils.ExitError(err.Error(), -1)
	}
}
