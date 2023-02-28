package console

import (
	"Hello/app/console/task"
	"Hello/app/libs/utils"
	"fmt"
	"time"
)

// 注册任务，调用方法：console.Start("任务名称")
var fun = map[string]func(){
	"ContinueStock": task.ContinueStock,
}

// 后台任务
type Task struct {
	Name   string // 任务名称
	status bool   // 任务状态
	time   int64  // 任务启动时间
}

// 在线的任务列表
var list map[string]*Task

func init() {
	list = make(map[string]*Task)
}

// 开始任务 返回：信息，运行结果
func (this *Task) Start() map[string]interface{} {
	if fun[this.Name] == nil {
		return map[string]interface{}{"info": "任务[" + this.Name + "]不存在"}
	}
	if list[this.Name] == nil {
		this.status = true
		this.time = utils.GetTime()
		list[this.Name] = this
		go this.run()
		fmt.Println("开始任务[" + this.Name + "]...")
		return map[string]interface{}{"info": "开始任务[" + this.Name + "]"}
	} else {
		return map[string]interface{}{"info": "任务[" + this.Name + "]已在运行中"}
	}
}

// 停止任务 返回：信息，停止结果
func (this *Task) Stop() map[string]interface{} {
	if list[this.Name] == nil {
		return map[string]interface{}{"info": "任务[" + this.Name + "]不存在"}
	}
	list[this.Name].status = false
	delete(list, this.Name)
	return map[string]interface{}{"info": "停止任务[" + this.Name + "]信号已发送"}
}

// 查询任务 返回：信息，运行状态，运行时间，查询结果
func (this *Task) Info() map[string]interface{} {
	data := map[string]interface{}{"name": this.Name, "info": "查询成功"}
	if list[this.Name] == nil {
		data["info"] = "任务不存在"
	} else {
		data["time"] = utils.Date("2006-01-02 15:04:05", list[this.Name].time)
	}
	return data
}

func (this *Task) run() {
	for this.status {
		fun[this.Name]()
		time.Sleep(1 * time.Millisecond) // 避免任务中忘记写延时一直占用cpu
	}
	fmt.Println("任务[" + this.Name + "]已停止")
}
