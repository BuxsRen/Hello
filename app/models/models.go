package models

import (
	"Hello/app/libs/utils"
	"Hello/bootstrap/database"
	"fmt"
	"github.com/gohouse/gorose/v2"
	"strconv"
	"strings"
)

func DB(table interface{}) gorose.IOrm {
	return database.MySQL.NewOrm().Table(table)
}

func DBSESSION(table interface{}) gorose.ISession {
	return database.MySQL.NewSession().Bind(table)
}

var (
	page        = 1                                                                               // 分页默认开始页数
	limit       = 20                                                                              // 分页默认每页条数
	selectBlack = []string{"page", "limit", "file", "access_token", "token", "x-token", "images"} // 不参与筛选的参数 下划线开头的也不会参与
	selectVague = []string{"username", "name", "nickname", "remarks", "title", "school"}          // 参与模糊查询的参数
)

type Model struct {
	Table     gorose.IOrm
	Data      map[string]interface{} // 需要筛选的数据
	Order     string                 // 排序字段
	OrderType string                 // 排序方式 asc,desc
	Search    map[string]interface{} // 搜索关键词
	Model     gorose.IOrm            // 操作的model
}

// 数据筛选/分页
/**
 *@Example:
m := Model{
		Model:this,
		Table: DB(this),
		Data: data,
		Search: map[string]interface{}{
			"title":"测试",
		},
	}
	return m.PageSearch()
*/
func (this *Model) PageSearch() map[string]interface{} {
	count := this.Model
	search := this.Model

	for k, v := range this.Data { // 筛选
		if k[0:1] == "_" { // 忽略 _ 开头的参数
			continue
		}
		if !utils.InArray(k, selectBlack) && !utils.InArray(k, selectVague) { // 不在黑名单 且 不在 模糊搜索名单的键
			if fmt.Sprintf("%T", v) != "[]string" && fmt.Sprintf("%T", v) != "map[string]interface{}" { // 条件查询
				this.Table.Where(k, v)
				count.Where(k, v)
				search.Where(k, v)
			}
		} else if utils.InArray(k, selectVague) { // 模糊筛选
			item := utils.ParamToString(v)
			this.Table.Where(k, "like", "%"+item+"%")
			count.Where(k, "like", "%"+item+"%")
			search.Where(k, "like", "%"+item+"%")
		}
	}

	if fmt.Sprintf("%T", this.Search) != "map[]" { // 搜索 {"title":"关键字"};
		var result = false
		var key string
		for k, v := range this.Search {
			key = k
			_, res := search.Where(k, "like", "%"+utils.ParamToString(v)+"%").Get()
			if res == nil { // 搜索到结果直接退出，搜不到继续搜
				result = true
				this.Table.Where(k, "like", "%"+utils.ParamToString(v)+"%")
				count.Where(k, "like", "%"+utils.ParamToString(v)+"%")
				break
			}
		}
		if !result { // 搜索不到就把记录为空
			this.Table.Where(key, "@!->:)!<-@")
			count.Where(key, "@!->:)<-!@")
		}
	}

	if this.Data["_sort"] != nil { // 处理排序字段[接口专用] 格式 升序：_sort=+|id 或 降序：_sort=-|id
		sort := strings.Split(utils.ParamToString(this.Data["_sort"]), "|")
		if len(sort) == 2 {
			this.Order = sort[1]
			if sort[0] == "-" {
				this.OrderType = "desc"
			} else {
				this.OrderType = "asc"
			}
		}
	}

	num, _ := count.Count() // 统计

	if this.Data["page"] != nil {
		page, _ = strconv.Atoi(this.Data["page"].(string))
	} else {
		page = 1
	}
	if this.Data["limit"] != nil {
		limit, _ = strconv.Atoi(this.Data["limit"].(string))
	} else {
		limit = 20
	}
	this.Table.Offset((page - 1) * limit).Limit(limit) // 分页

	if this.Order != "" { // 排序
		if this.OrderType == "" {
			this.OrderType = "desc"
		}
		this.Table.Order(this.Order + " " + this.OrderType)
	}

	list, e := this.Table.Get()
	if e != nil {
		utils.ExitError(e.Error(), -1)
	}

	return map[string]interface{}{
		"count": num,
		"items": list,
	}
}
