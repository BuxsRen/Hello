package utils

import (
	"Hello/bootstrap/helper"
	"archive/zip"
	"io"
	"log"
	"os"
	"path/filepath"
)

// 向文本文本中追加内容(文件地址，内容)，不存在则创建
func AppendFile(filePath, content string) {
	f, e := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if e != nil {
		log.Println(e)
	}
	defer f.Close()
	_, _ = io.WriteString(f, content)
}

// 判断文件是否存在
func IsExist(fileAddr string) bool {
	// 读取文件信息，判断文件是否存在
	_, err := os.Stat(fileAddr)
	if err != nil {
		if os.IsExist(err) { // 根据错误类型进行判断
			return true
		}
		return false
	}
	return true
}

// 判断目录是否存在
func IsDir(fileAddr string) bool {
	s, err := os.Stat(fileAddr)
	if err != nil {
		log.Println(err)
		return false
	}
	return s.IsDir()
}

// 写入文件 文件名，文件字节
func WriteFile(fileName string, content []byte) error {
	return os.WriteFile(fileName, content, 0644)
}

//读取文件 文件名/绝对路径
func ReadFile(name string) []byte {
	contents, err := os.ReadFile(name)
	if err != nil {
		ExitError(err.Error(), -1)
		return nil
	}
	return contents
}

// 写日志，保存在错误日志文件中，提供给未开启debug模式时使用
func Log(str string) {
	if helper.LOG == nil {
		log.Println("✘ log write failed,debug is off")
	} else {
		_ = helper.LOG.Output(2, "\n"+str+"\n\n")
	}
}

//zip 压缩 目录，压缩后的文件路径
func ZipDir(dir, zipFile string) {
	fz, err := os.Create(zipFile)
	if err != nil {
		log.Printf("Create zip file failed: %s\n", err.Error())
		return
	}
	defer fz.Close()

	w := zip.NewWriter(fz)
	defer w.Close()

	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			fDest, err := w.Create(path[len(dir):])
			if err != nil {
				log.Printf("Create failed: %s\n", err.Error())
				return nil
			}
			fSrc, err := os.Open(path)
			if err != nil {
				log.Printf("Open failed: %s\n", err.Error())
				return nil
			}
			defer fSrc.Close()
			_, err = io.Copy(fDest, fSrc)
			if err != nil {
				log.Printf("Copy failed: %s\n", err.Error())
				return nil
			}
		}
		return nil
	})
}

// 解压缩 zip 文件路径，解压缩之后的路径
func Unzip(zipFile string, destDir string) error {
	zipReader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	for _, f := range zipReader.File {
		fpath := filepath.Join(destDir, f.Name)
		if f.FileInfo().IsDir() {
			_ = os.MkdirAll(fpath, os.ModePerm)
		} else {
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}

			inFile, err := f.Open()
			if err != nil {
				return err
			}
			defer inFile.Close()

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, inFile)
			if err != nil {
				return err
			}
		}
	}
	return nil
}