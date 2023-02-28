package qrcode

//import (
//	"Hello/app/libs/utils"
//	"github.com/tuotoo/qrcode" // 这里 go mod tidy 安装失败的话。就手动安装一下
//	"os"
//)
//
//// 解析二维码。传入图片的绝对地址
//func DeQrCode(filePath string) string {
//	f, err := os.Open(filePath)
//	if err != nil {
//		utils.ExitError(err.Error(),-1)
//	}
//	q, err := qrcode.Decode(f)
//	if err != nil {
//		utils.ExitError(err.Error(),-1)
//	}
//	defer f.Close()
//	return q.Content
//}