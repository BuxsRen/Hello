package main

import (
	"Hello/app/libs/utils"
	"os"
)

// 用于build时的解压缩

func main() {
	str := os.Args
	if len(str) > 2 {
		utils.ZipDir(str[1], str[2])
	}
}
