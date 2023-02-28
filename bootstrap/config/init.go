package config

import (
	"Hello/bootstrap/helper"
	"fmt"
	"gopkg.in/yaml.v2" // go get gopkg.in/yaml.v2
	"os"
)

func init() {
	App = loadConfig()
}

// 加载 app.yaml 配置
func loadConfig() *app {
	var h = &helper.Helper{}
	file, e := os.ReadFile("./config/app.yaml")
	if e != nil {
		h.Exit("✘ Config File Read Failed!", 3)
	}

	var app app
	e = yaml.Unmarshal(file, &app)
	if e != nil {
		h.Exit("✘ Config Loading Failed!", 3)
	}
	fmt.Println("🔨 Config -> ./config/app.yaml")
	return &app
}
