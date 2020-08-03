package middleware

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"nms-controller/db"
	"nms-controller/model"

	"github.com/golang/glog"

	_ "github.com/lib/pq" // postgres golang driver
)

//更新pg config
func postConfig(w http.ResponseWriter, r *http.Request) {
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
		json.NewEncoder(w).Encode(model.ErrRes{Message: "illegal config, please check again"})
		// http.Error(w, e.Message, http.StatusBadRequest) //也可簡寫成這樣
		return
	} else {
		//insert config into pg
		err = db.Pg.InsertConfig(reqBody)
		if err != nil {
			//成功返回
			res := "done"
			json.NewEncoder(w).Encode(res)
		}
	}

	// 失敗
	// Srv.Close()
	// time.Sleep(3 * time.Second)
	// Handler()

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

	fmt.Fprintf(w, "Success, Restarting Prometheus...\n")
	go restart()
}

func postConfigWrapper(hf http.HandlerFunc) http.HandlerFunc {
	// fmt.Println("PostConfigWrapper") //when init service run here
	count := 0
	return func(w http.ResponseWriter, r *http.Request) { //when hit api run here
		count++
		fmt.Printf("Handler Function called %d times\n", count)
		hf(w, r)
	}
}

//在k8s上會自動重啟
func restart() {
	time.Sleep(time.Duration(5) * time.Second)
	defer os.Exit(3) //exit code 3 means ERROR_PATH_NOT_FOUND
}

// func hello(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Hello World!\n")
// }

// func log(hf http.HandlerFunc) http.HandlerFunc {
// 	count := 0
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		count++
// 		fmt.Printf("Handler Function called %d times\n", count)
// 		hf(w, r)
// 	}
// }

//r.HandleFunc("/hello", log(hello)).Methods("GET", "OPTIONS")

// type MyHandler struct{}

// func (wh *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Hello World!\n")
// }

// func log(h http.Handler) http.Handler {
// 	count := 0
// 	f := func(w http.ResponseWriter, r *http.Request) {
// 		count++
// 		fmt.Printf("Handler Function called %d times\n", count)
// 		h.ServeHTTP(w, r)
// 	}
// 	return http.HandlerFunc(f)
// }

// http.Handle("/hello", log(&myHandler))
