package prom

import (
	"net/http"
	"nms-controller/db"
	"nms-controller/model"
)

var (
	configs []model.CustomConfig
	//RegHandler 讓router去註冊(main.go)
	RegHandler http.Handler
)

func init() {
	// get latest configs from pg
	configs = db.GlobalCustomConfigs
	//初始化RegHandler
	StartPrometheus()
}
