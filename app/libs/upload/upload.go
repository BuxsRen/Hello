package upload

import (
	"Hello/app/libs/utils"
	"Hello/bootstrap/config"
	"crypto/md5"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"strings"
)

// 文件上传类
/**
 * @Example:
	var c *gin.Context
	file, err := c.FormFile("file")
	yun := upload.Upload{File:file,Format:"images"}
	fmt.Println(yun.Upload())
*/
type Upload struct {
	File   *multipart.FileHeader // 上传的文件
	Format string                // 上传的类型
	suffix string                // 文件后缀
}

// 上传文件，返回上传成功后的地址
func (this Upload) Upload() string {
	this.checkFormat()
	this.checkSuffix()
	//this.checkType()
	this.checkSize()
	return this.save()
}

func (this *Upload) checkFormat() {
	if format[this.Format] == "" {
		utils.ExitError("不允许的类型", -1)
	}
}

// 验证文件后缀
func (this *Upload) checkSuffix() {
	this.suffix = strings.ToLower(path.Ext(this.File.Filename)) // 取文件后缀
	this.suffix = strings.ToLower(this.suffix)
	if suffixList[this.Format][this.suffix] == "" {
		utils.ExitError("文件格式错误", -1)
	}
}

// 验证文件类型
func (this *Upload) checkType() {
	types := this.File.Header.Values("Content-Type")
	fmt.Println(this.File)
	if len(types) == 0 || formatList[this.Format][types[0]] == "" {
		utils.ExitError("文件类型错误", -1)
	}
}

// 验证文件大小
func (this *Upload) checkSize() {
	if this.File.Size > sizeList[this.Format]*1024*1024 {
		utils.ExitError("上传的图片过大", -1)
	}
}

// 保存文件
func (this *Upload) save() string {
	// 打开文件
	src, e := this.File.Open()
	file, e := this.File.Open()
	if e != nil {
		utils.ExitError(e.Error(), -1)
	}
	defer src.Close()
	name := this.getFileMd5(file) + this.suffix // 获取文件的唯一名称
	fileName := this.getFilePath() + name       // 保存文件的绝对路径
	if utils.IsExist(fileName) {                // 避免文件重复保存
		return name
	}
	out, e := os.Create(fileName)
	if e != nil {
		utils.ExitError(e.Error(), -1)
	}
	defer out.Close()
	_, e = io.Copy(out, src)
	if e != nil {
		utils.ExitError(e.Error(), -1)
	}
	return name
}

// 获取文件md5
func (this *Upload) getFileMd5(f multipart.File) string {
	md5hash := md5.New()
	_, _ = io.Copy(md5hash, f)
	return fmt.Sprintf("%x", md5hash.Sum(nil))
}

// 获取文件保存的绝对路径
func (this *Upload) getFilePath() string {
	dir := config.App.Other.PublicDir
	if dir == "" {
		utils.ExitError("未设置 public_dir 文件上传静态目录", -1)
	}
	if !utils.IsDir(dir) {
		_, err := os.Create(dir)
		if err != nil {
			utils.ExitError(fmt.Sprintf("创建目录失败,%v", err), -1)
		}
	}
	return dir
}
