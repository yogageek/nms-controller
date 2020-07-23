package prom

import (
	"nms-controller/db"
	"nms-controller/model"
)

var (
	configList []model.CustomConfig
)

func init() {
	// get latest configs from pg
	configList = db.GlobalCustomConfigs
}
