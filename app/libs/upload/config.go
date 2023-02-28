package upload

var (
	// 允许上传的后缀
	suffixList = map[string]map[string]string{
		"images": {
			".png":"png",
			".jpg":"jpg",
			".gif":"gif",
		},
		"excel": {
			".xlsx":"xlsx",
			".xls":"xls",
		},
	}
	// 允许上传的类型
	formatList = map[string]map[string]string{
		"images": {
			"image/jpeg":"image/jpeg",
			"image/png":"image/png",
			"image/gif":"image/gif",
		},
		"excel": {
			"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": "sheet",
		},
	}
	// 允许上传的大小 单位 M
	sizeList = map[string]int64{
		"images":2,
		"excel":2,
	}
	// 支持上传的文件类型
	format = map[string]string{
		"images":"images",
		"excel":"excel",
	}
)

type Api interface {
	Upload() string
}