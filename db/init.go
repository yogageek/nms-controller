package db

import (
	"log"
	"nms-controller/model"
)

var (
	//Rc redis連線
	Rc *theRedis
	//Pg pg連線
	Pg *postgres
	//GlobalCustomConfigs 供所有pkg存取
	GlobalCustomConfigs []model.CustomConfig
)

func init() {
	Rc = newTheRedis()
	Pg = newPostgres()

	//清除redis資料
	Rc.FlushRedis()
	log.Println("flush redis...")

	//查詢pg
	GlobalCustomConfigs = Pg.GetCustomConfigs()
}
