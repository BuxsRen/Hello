package email

import (
	"Hello/app/libs/utils"
	"Hello/bootstrap/config"
	"fmt"
	"gopkg.in/gomail.v2"
	"mime"
)

/* 邮箱类 */
type Email struct {
	name string
	user string
	pass string
	host string
	port int
	em   *gomail.Message
}

// 邮箱类，初始化方法
/**
 *@Example:
	e := email.New()
	e.SetTitle("测试")
	e.SetToEmail([]string{"buxsren@qq.com"})
	e.SetBody("<h1>哈哈哈</h1>")
	e.AddAttach("test.txt","/www/wwwroot/public/test.txt")
	fmt.Println(e.SendMail())
*/
func New() *Email {
	e := Email{
		name: config.App.Email.Name,
		user: config.App.Email.User,
		pass: config.App.Email.Pass,
		host: config.App.Email.Host,
		port: config.App.Email.Port,
	}
	if e.name == "" || e.user == "" || e.pass == "" || e.host == "" || e.port == 0 {
		utils.ExitError("请先配置邮箱设置", -1)
	}
	e.em = gomail.NewMessage()
	e.em.SetHeader("From", e.em.FormatAddress(e.user, e.name)) // 设置发件人名称
	return &e
}

// 发送邮件
func (this *Email) SendMail() error {
	d := gomail.NewDialer(this.host, this.port, this.user, this.pass)
	return d.DialAndSend(this.em)
}

// 设置/更改发件人名称
func (this *Email) SetFormName(name string) *Email {
	this.em.SetHeader("From", this.em.FormatAddress(this.user, name))
	return this
}

// 设置发件人
func (this *Email) SetToEmail(to []string) *Email {
	this.em.SetHeader("To", to...)
	return this
}

// 设置邮件标题
func (this *Email) SetTitle(title string) *Email {
	this.em.SetHeader("Subject", title)
	return this
}

// 设置邮件正文,支持html格式
func (this *Email) SetBody(body string) *Email {
	this.em.SetBody("text/html", body) //设置邮件正文
	return this
}

// 添加附件.附件名称(xxx.txt),附件路径(/www/wwwwroot/public/xxx.txt)
func (this *Email) AddAttach(name, path string) *Email {
	this.em.Attach(path, gomail.Rename(name),
		gomail.SetHeader(map[string][]string{ // 附件中文乱码解决方案
			"Content-Disposition": {fmt.Sprintf(`attachment; filename="%s"`, mime.QEncoding.Encode("UTF-8", name))},
		}),
	)
	return this
}