package middleware

import (
	"encoding/json" // package to encode and decode the json into struct and vice versa
	"io/ioutil"
	"net/http" // used to access the request and response object of the api

	// model package where Config schema is defined
	"nms-controller/db"
	"nms-controller/model"

	// used to read the environment variable
	// package used to covert string into int type
	"github.com/golang/glog"
	// used to get the params from the route
	// package used to read the .env file
	_ "github.com/lib/pq" // postgres golang driver
)

func PostConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var customConfigs []model.CustomConfig

	// decode request body (byte)
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		glog.Error("read request body err:", err)
	}

	//request(byte) to struct
	err = json.Unmarshal(reqBody, &customConfigs)
	//json.Marshal(configs) //struct to byte
	if err != nil {
		glog.Error("byte to struct err:", err)

		w.WriteHeader(http.StatusBadRequest) //設定返回header
		json.NewEncoder(w).Encode(model.ErrRes{
			Message: "illegal config"},
		) // http.Error(w, e.Message, http.StatusBadRequest) //也可簡寫成這樣
		return
	}

	//insert into pg
	err = db.Pg.InsertConfig(reqBody)
	if err != nil {
		//成功返回
		res := "done"
		json.NewEncoder(w).Encode(res)
	}

	//重啟collector
	// ch := UserTicker()
	// time.Sleep(20 * time.Second)
	// ch <- true
	// close(ch)

	//搭配可以停止
	// tickerCH := logic.CH //ticker的channel
	// tickerCH <- true     //把true傳送過去
	// close(tickerCH)      //只寫這行不會停止ticker...
	// time.Sleep(time.Duration(2) * time.Second)
}
