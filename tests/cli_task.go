package main

import (
	"Hello/app/console/task"
	"fmt"
	"github.com/urfave/cli"
	"os"
	"time"
)

// 计划任务独立版

// 注册任务，调用方法：console.Start("SyncPrize")
var fun = map[string]func(){
	"ContinueStock": task.ContinueStock,
}

func main() {
	app := cli.NewApp()
	app.Name = "运行任务(任务单独运行)"
	app.Version = "1.0.0"
	app.Commands = []cli.Command{
		{
			Name:    "task:start",               // 命令的名字
			Aliases: []string{"t:s"},            // 命令的缩写
			Usage:   "运行任务 命令： task:start 任务名称", // 命令的用法注释，这里会在输入 程序名 -help的时候显示命令的使用方法
			Action: func(c *cli.Context) error { // 命令的处理函数
				run(c.Args().Get(0))
				return nil
			},
		},
	}
	_ = app.Run(os.Args)
}

func run(name string) {
	if fun[name] == nil {
		fmt.Println("任务不存在")
		return
	}
	fmt.Println("运行任务...")
	for {
		fun[name]()
		time.Sleep(1 * time.Millisecond) // 避免任务中忘记写延时一直占用cpu
	}
}
