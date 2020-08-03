package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"nms-controller/model"

	"os"

	"github.com/golang/glog"
	_ "github.com/lib/pq"
)

type postgres struct {
	sqlDB *sql.DB
}

func newPostgres() *postgres {
	return &postgres{
		sqlDB: createPGClient(),
	}
}

func createPGClient() *sql.DB {

	// Open the connection
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

	if err != nil {
		glog.Error("create pg connection err:", err)
		panic(err)
	}

	// check the connection
	err = db.Ping()

	if err != nil {
		glog.Error("ping postgres err:", err)
		panic(err)
	}

	fmt.Println("Successfully connected postgres!")

	return db
}

//查詢pg取得最新一筆config並轉換成[]CustomConfig
func (pg *postgres) GetCustomConfigs() []model.CustomConfig {

	doSelect := func() []byte {
		fmt.Println("select config from postgres...")

		//只找最新的一筆
		sql := `SELECT data FROM file WHERE db=$1 ORDER BY id desc limit 1`

		//根據環境變數redis設置 來查詢找哪筆config 也就是說pg可能存在不同最新config分別給不同套controller使用
		//# v1.0.2 拿掉redis改成放在postgres
		pgdb := os.Getenv("POSTGRES_DB")

		var b []byte
		err := pg.sqlDB.QueryRow(sql, pgdb).Scan(&b)
		if err != nil {
			glog.Error(err)
			// deprecated
			// if err == sql.ErrNoRows {
			// //#1.3.0 如果db為空就新增一筆空config
			// fmt.Println("insert empty config into table...")
			// _, err := pgClient.Exec(`INSERT INTO file (data, redisdb) VALUES ($1, $2)`, "[]", env_redisdb)
			// if err != nil {
			// 	glog.Fatal(err)
			// }
			// }
		}
		return b
	}

	doConvert := func(b []byte) []model.CustomConfig {
		var customConfigs []model.CustomConfig
		//檢查config是否合法 & 轉struct -> customConfigs
		if err := json.Unmarshal(b, &customConfigs); err != nil {
			fmt.Println("Unmarshal config in pg and see if legal...")
			if b == nil {
				glog.Error("config is empty")
			} else {
				glog.Error("convert config to struct err:", err)
			}
		}
		return customConfigs
	}

	return doConvert(doSelect())
}

// insert one Config in into DB
func (pg *postgres) InsertConfig(configsb []byte) error {

	sql := `INSERT INTO public.file (data, db) VALUES ($1, $2) RETURNING id`

	db := os.Getenv("POSTGRES_DB")

	var id int
	err := pg.sqlDB.QueryRow(sql, configsb, db).Scan(&id)
	if err != nil {
		glog.Error(err)
		return err
	}
	fmt.Printf("Inserted. id=%v", id)
	return nil
}
