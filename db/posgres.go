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

// var PgClient *sql.DB
// var PostgresDB Postgres

type Postgres struct {
	SqlDB *sql.DB
}

func NewPostgres() *Postgres {
	return &Postgres{
		SqlDB: createPGClient(),
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

func (pg *Postgres) selectPg() []byte {

	fmt.Println("select config from postgres...")

	var data []byte

	//只找最新的一筆
	sql := `SELECT data FROM file WHERE redisdb=$1 ORDER BY id desc limit 1`

	//根據環境變數redis設置 來查詢找哪筆config 也就是說pg可能存在不同最新config分別給不同套controller使用
	envRedis := os.Getenv("REDIS_DB")

	err := pg.SqlDB.QueryRow(sql, envRedis).Scan(&data)
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

	// glog.Info(string(data))

	return data
}

//查詢pg取得最新一筆config並轉換成[]CustomConfig
func (pg *Postgres) GetCustomConfigs() []model.CustomConfig {

	bytes := pg.selectPg()

	var customConfigs []model.CustomConfig

	//檢查config是否合法 & 轉struct -> customConfigs
	if err := json.Unmarshal(bytes, &customConfigs); err != nil {
		fmt.Println("Unmarshal config in pg and see if legal...")
		if bytes == nil {
			glog.Error("config is empty")
		} else {
			glog.Error("convert config to struct err:", err)
		}
	}
	return customConfigs
}
