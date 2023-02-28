package encry

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

// AES加解密类
// AES 有五种加密模式：
// 1.电码本模式（Electronic Codebook Book (ECB)） 出于安全考虑，golang默认并不支持ECB模式。
// 2.密码分组链接模式（Cipher Block Chaining (CBC)）
// 3.计算器模式（Counter (CTR)）
// 4.密码反馈模式（Cipher FeedBack (CFB)）
// 5.输出反馈模式（Output FeedBack (OFB)）
// https://mp.weixin.qq.com/s/jnsypJLiYIFK14rXDNPpSQ
type AES struct {}

// AES加密
func (this *AES) AesEcrypt(origData []byte, key []byte) ([]byte, error) {
	//创建加密算法实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//获取块的大小
	blockSize := block.BlockSize()
	//对数据进行填充，让数据长度满足需求
	origData = this.pkcs7Padding(origData, blockSize)
	//采用AES加密方法中CBC加密模式
	blocMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	//执行加密
	blocMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// AES解密
func (this *AES) AesDeCrypt(cypted []byte, key []byte) ([]byte, error) {
	//创建加密算法实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//获取块大小
	blockSize := block.BlockSize()
	//创建加密客户端实例
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(cypted))
	//这个函数也可以用来解密
	blockMode.CryptBlocks(origData, cypted)
	//去除填充字符串
	origData, err = this.pkcs7UnPadding(origData)
	if err != nil {
		return nil, err
	}
	return origData, err
}

// PKCS7 填充模式
func (this *AES) pkcs7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	//Repeat()函数的功能是把切片[]byte{byte(padding)}复制padding个，然后合并成新的字节切片返回
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// 填充的反向操作，删除填充字符串
func (this *AES) pkcs7UnPadding(origData []byte) ([]byte, error) {
	//获取数据长度
	length := len(origData)
	if length == 0 {
		return nil, errors.New("加密字符串错误！")
	} else {
		//获取填充字符串长度
		unpadding := int(origData[length-1])
		//截取切片，删除填充字节，并且返回明文
		return origData[:(length - unpadding)], nil
	}
}