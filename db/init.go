package db

import (
	"log"
)

var (
	//redis連線
	Rc *Redis
	//pg連線
	Pg *Postgres
)

func init() {
	Rc = NewRedis()
	Pg = NewPostgres()

	//清除redis資料
	Rc.FlushRedis()
	log.Println("flush redis...")
}
