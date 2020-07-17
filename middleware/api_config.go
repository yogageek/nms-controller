package middleware

import (
	"database/sql"
	"encoding/json" // package to encode and decode the json into struct and vice versa
	"fmt"
	"io/ioutil"
	"log"
	"net/http" // used to access the request and response object of the api

	// model package where Config schema is defined
	"nms-controller/model"
	"time"

	"os"      // used to read the environment variable
	"strconv" // package used to covert string into int type

	"github.com/golang/glog"
	"github.com/gorilla/mux" // used to get the params from the route

	"github.com/joho/godotenv" // package used to read the .env file
	_ "github.com/lib/pq"      // postgres golang driver
)

// response format
type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

//-----------------------------
type error interface {
	Error() string
}

// Consider a different name, too close to shadowing an inbuilt type, with chance of people using one expecting to be using the other
type errRes struct {
	Message     string `json:"error,omitempty"`
	Description string `json:"description,omitempty"`
}

func (e *errRes) Error() string {
	return fmt.Sprintf("%d:%d:", e.Message, e.Description)
}

//-----------------------------

//--------------------------
// create connection with postgres db
func createConnection() *sql.DB {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Open the connection
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

	if err != nil {
		panic(err)
	}

	// check the connection
	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	// return the connection
	return db
}

// ConfigPost create a config in the postgres db
func ConfigPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// create an empty config of type model.Config
	var reqestConfig []model.CustomConfig

	// decode the json request to config
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		glog.Error(err)
	}

	err = json.Unmarshal(body, &reqestConfig)
	if err != nil {
		glog.Error("Unable to decode the request body. err:", err)
		e := errRes{Message: "illegal config"}
		//設定返回header
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(e)
		// http.Error(w, e.Message, http.StatusBadRequest) //也可簡寫成這樣
		return
	}
	// fmt.Println(reqestConfig)

	// call insert config function and pass the config
	insertID := insertConfig(reqestConfig)

	//重啟collector
	// ch := UserTicker()
	// time.Sleep(20 * time.Second)
	// ch <- true
	// close(ch)

	//搭配可以停止
	// tickerCH := logic.CH //ticker的channel
	// tickerCH <- true     //把true傳送過去
	// close(tickerCH)      //只寫這行不會停止ticker...
	time.Sleep(time.Duration(2) * time.Second)
	// logic.CallRunCollector(true)

	//更新exporter
	// if err := logic.updateExporter(); err != nil {
	// 	glog.Error("update exporter error. err:", err)
	// 	e := errRes{
	// 		Message:     "update exporter error",
	// 		Description: err.Error(),
	// 	}
	// 	// send the response
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	json.NewEncoder(w).Encode(e)
	// 	return
	// }

	//成功返回
	// format a response object
	res := response{
		ID:      insertID,
		Message: "Config created successfully",
	}
	// send the response
	json.NewEncoder(w).Encode(res)

}

// ConfigIDGet will return a single config by its id
func ConfigIDGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// get the id from the request params, key is "id"
	params := mux.Vars(r)

	// convert the id type from string to int
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Println("Unable to convert the string into int.  %v", err)
		failRes := response{
			Message: "Unable to convert the string into int.",
		}
		json.NewEncoder(w).Encode(failRes)
	}

	// call the getConfig function with config id to retrieve a single config
	config, err := getConfig(int64(id))

	if err != nil {
		log.Fatalf("Unable to get config. %v", err)
	}

	// send the response
	json.NewEncoder(w).Encode(config)
}

// ConfigsGet will return all the configs
func ConfigsGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// get all the configs in the db
	configs, err := getAllConfigs()

	if err != nil {
		log.Println("Unable to get all config. %v", err)
		failRes := response{
			Message: "Unable to get all config.",
		}
		json.NewEncoder(w).Encode(failRes)
	}

	// send all the configs as response
	json.NewEncoder(w).Encode(configs)
}

// ConfigIDPut update config's detail in the postgres db
func ConfigIDPut(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// get the id from the request params, key is "id"
	params := mux.Vars(r)

	// convert the id type from string to int
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Println("Unable to convert the string into int.  %v", err)
		failRes := response{
			Message: "Unable to convert the string into int.",
		}
		json.NewEncoder(w).Encode(failRes)
	}

	// create an empty config of type model.Config
	var data []model.CustomConfig

	// decode the json request to config

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	err = json.Unmarshal(body, &data)

	if err != nil {
		log.Println("Unable to decode the request body.  %v", err)
		failRes := response{
			Message: "Unable to decode the request body.",
		}
		json.NewEncoder(w).Encode(failRes)
	}

	// call update config to update the config
	updatedRows := updateConfig(int64(id), data)

	// format the message string
	msg := fmt.Sprintf("Config updated successfully. Total rows/record affected %v", updatedRows)

	// format the response message
	res := response{
		ID:      int64(id),
		Message: msg,
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

// ConfigIDDelete delete config's detail in the postgres db
func ConfigIDDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// get the id from the request params, key is "id"
	params := mux.Vars(r)

	// convert the id in string to int
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Println("Unable to convert the string into int.  %v", err)
		failRes := response{
			Message: "Unable to convert the string into int.",
		}
		json.NewEncoder(w).Encode(failRes)
	}

	// call the deleteConfig, convert the int to int64
	deletedRows := deleteConfig(int64(id))

	// format the message string
	msg := fmt.Sprintf("Config updated successfully. Total rows/record affected %v", deletedRows)

	// format the reponse message
	res := response{
		ID:      int64(id),
		Message: msg,
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

//------------------------- handler functions ----------------
// insert one Config in the DB
func insertConfig(configs []model.CustomConfig) int64 {

	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the insert sql query
	// returning id will return the id of the inserted config
	sqlStatement := `INSERT INTO public.file (data, redisdb) VALUES ($1, $2) RETURNING id`

	// the inserted id will store in this id
	var id int64

	// execute the sql statement
	// Scan function will save the insert id in the id
	configJ, _ := json.Marshal(configs)
	redisdb := os.Getenv("REDIS_DB")

	err := db.QueryRow(sqlStatement, configJ, redisdb).Scan(&id)

	if err != nil {
		log.Fatalf("Unable to execute the query.  %v", err)
	}

	fmt.Printf("Inserted a single record %v", id)

	// return the inserted id
	return id
}

// get one config from the DB by its id
func getConfig(id int64) (model.ConfigIdDetail, error) {
	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create a Config of model.Config type
	var config model.ConfigIdDetail

	// create the select sql query
	sqlStatement := `SELECT file.id, file.data as configs FROM public.file WHERE id=$1`

	// execute the sql statement
	row := db.QueryRow(sqlStatement, id)

	// unmarshal the row object to config
	var configFile []byte
	err := row.Scan(&config.ID, &configFile)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return config, nil
	case nil:
		if err := json.Unmarshal(configFile, &config.CustomConfig); err != nil {
			panic(err)
		}
		return config, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}

	// return empty config on error
	return config, err
}

// get all Configs
func getAllConfigs() ([]model.ConfigIdDetail, error) {
	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	var configs []model.ConfigIdDetail

	// create the select sql query
	sqlStatement := `SELECT file.id, file.data as configs FROM public.file`

	// execute the sql statement
	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// close the statement
	defer rows.Close()

	// iterate over the rows
	for rows.Next() {
		var config model.ConfigIdDetail

		// unmarshal the row object to config
		var configFile []byte
		err = rows.Scan(&config.ID, &configFile)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}
		if err := json.Unmarshal(configFile, &config.CustomConfig); err != nil {
			panic(err)
		}

		// append the config in the config slice
		configs = append(configs, config)

	}

	// return empty config on error
	return configs, err
}

// update config in the DB
func updateConfig(id int64, configs []model.CustomConfig) int64 {

	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the update sql query
	sqlStatement := `UPDATE public.file SET data=$2 WHERE id=$1`

	// execute the sql statement
	bytes, _ := json.Marshal(configs)
	res, err := db.Exec(sqlStatement, id, bytes)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// check how many rows affected
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}

// delete config in the DB
func deleteConfig(id int64) int64 {

	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the delete sql query
	sqlStatement := `DELETE FROM public.file WHERE id=$1`

	// execute the sql statement
	res, err := db.Exec(sqlStatement, id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// check how many rows affected
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}
