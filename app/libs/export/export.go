package export

import (
	"Hello/app/libs/utils"
	"encoding/csv"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
	"os"
	"strconv"
)

// 数据导出表格类
/**
 *@Example:
	header := []string{"ID","标题"}
	body := []string{"1","哈哈哈"}
	var list [][]string{header,body}
	ep := &export.Export{C: c,Name:"vote_record",Data:list}
	return ep.Export()
*/
type Export struct {
	C        *gin.Context
	Name     string     // 导出文件标题
	Type     string     // 导出类型 down 文件流，file 文件名 默认 文件流
	Format   string     // 保存格式 xlsx，csv 默认 xlsx
	Data     [][]string // 导出数据(标题+内容)
	filename string     //文件名
	path     string     // 文件绝对路径+文件名
}

// 数据导出到表格，返回表格地址或表格文件
func (this *Export) Export() interface{} {
	if this.Format == "csv" {
		return this.csv().export()
	} else {
		return this.xlsx().export()
	}
}

// 数据写入到CSV
func (this *Export) csv() *Export {
	this.filename = "export_" + fmt.Sprintf("%s_%s.csv", this.Name, utils.GetNow().Format("20060102150405"))
	this.path = "./storage/cache/" + this.filename
	f, e := os.Create(this.path)
	if e != nil {
		utils.ExitError("csv文件创建失败", -1)
	}
	defer f.Close()
	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM，防止中文乱码
	w := csv.NewWriter(f)
	w.WriteAll(this.Data) // W写入多条记录
	w.Flush()             // 清空缓存
	return this
}

// 数据写入到XLSX
func (this *Export) xlsx() *Export {
	this.filename = "export_" + fmt.Sprintf("%s_%s.xlsx", this.Name, utils.GetNow().Format("20060102150405"))
	this.path = "./storage/cache/" + this.filename
	xlsx := excelize.NewFile()
	var index = 1
	for _, v := range this.Data {
		col := strconv.Itoa(index)
		for key, item := range v { // 写入每列数据
			xlsx.SetCellValue("Sheet1", string(rune(65+key))+col, item)
		}
		index++
	}
	if err := xlsx.SaveAs(this.path); err != nil {
		utils.ExitError("xlsx文件创建失败", -1)
	}
	return this
}

// 导出
func (this *Export) export() interface{} {
	if this.Type == "file" { // 文件形式，获取文件名，使用 下载历史导出文件 方法进行下载
		return utils.OK("导出成功", this.filename)
	} else { // 文件流形式
		this.setHeader()
		return string(utils.ReadFile(this.path))
	}
}

// 设置响应头
func (this *Export) setHeader() {
	this.C.Header("Content-Type", "application/octet-stream")
	this.C.Header("Content-Type", "application/download")
	this.C.Header("Content-Disposition", "application/download")
	this.C.Header("Content-Disposition", `inline;filename="`+this.filename+`"`)
	this.C.Header("Content-Transfer-Encoding", "binary")
	this.C.Header("Cache-Control", "must-revalidate, post-check=0, pre-check=0")
	this.C.Header("Pragma", "no-cache")
}

// 下载历史导出文件
/**
 *@Example:
	ep := &export.Export{C: c}
	return ep.GetFile("export_vote_record_20220215163147.csv")
*/
func (this *Export) GetFile(filename string) string {
	this.filename = filename
	this.setHeader()
	return string(utils.ReadFile("./storage/cache/" + filename))
}
