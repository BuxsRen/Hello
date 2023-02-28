package qrcode

import (
	"Hello/app/libs/encry"
	"Hello/app/libs/utils"
	"github.com/skip2/go-qrcode"
	"image/color"
)

type qrCode struct {
	q      *qrcode.QRCode
	border bool // 是否展示边框
	size   int  // 二维码大小
}

// 二维码生成类: 初始化,传入需要生成的二维码字符串内容
/**
 * @Example:
	q := qrcode.New("xxx")
	q.SetSize(100).SetBackgroundColor(colornames.Red)
	fmt.Println(q.CreateBase64QrCode())
*/
func New(content string) *qrCode {
	var e error
	c := qrCode{border: true, size: 300}
	c.q, e = qrcode.New(content, qrcode.Medium)
	if e != nil {
		utils.ExitError(e.Error(), -1)
	}
	return &c
}

// 设置二维码内容
func (this *qrCode) SetContent(content string) *qrCode {
	this.q.Content = content
	return this
}

// 设置边框
func (this *qrCode) SetBorder(border bool) *qrCode {
	this.border = border
	return this
}

// 设置二维码大小
func (this *qrCode) SetSize(size int) *qrCode {
	this.size = size
	return this
}

// 设置二维码背景颜色。colornames.Red
func (this *qrCode) SetBackgroundColor(color color.RGBA) *qrCode {
	this.q.BackgroundColor = color
	return this
}

// 设置二维码前景颜色。colornames.Red
func (this *qrCode) SetForegroundColor(color color.RGBA) *qrCode {
	this.q.BackgroundColor = color
	return this
}

// 创建base64二维码
func (this *qrCode) CreateBase64QrCode() string {
	if this.border {
		this.q.DisableBorder = true //去掉边框
	}
	data, _ := this.q.PNG(this.size)
	return "data:image/png;base64," + encry.Base64Encode(data)
}