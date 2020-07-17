package main

import (
	"flag"
	"log"
	"net/http"
	"nms-controller/logic"
	"nms-controller/middleware"

	"time"

	"github.com/rs/cors"
)

// Add() # 添加计数
// Done() # 减掉计数，等价于Add(-1)，这样sync.WaitGroup只有两个API了
// Wait() # 阻塞直到计数为零
func main() {
	//Init the command-line flags for glog
	flag.Set("v", "5")
	flag.Parse()

	//讓此獨立為一個線程(裡面是一個無限迴圈)，不然會無法繼續下一行
	go logic.RunControllerLoop()

	// 緩衝
	time.Sleep(3 * time.Second)

	// move to run controller
	// //#1.0.0 bug 啟動多個註冊反註冊失敗 待修
	// logic.RunProm()

	//包含啟動prom metrics頁面 (每次進localhost:metrics 時,會執行GetOperationByMiddleware)
	//開啟exporter update, adpater crud API
	r := middleware.Router()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	})
	handler := c.Handler(r)

	// ch := RunCollector(false)
	// time.Sleep(20 * time.Second)
	// ch <- true
	// close(ch)

	//go最終會停再ListenAndServe無法下一行
	err := http.ListenAndServe(":8080", handler) //使用 http.ListenAndServe("localhost:8080", nil) 函数，如果成功会返回空，否则会返回一个错误
	//設8080為了配合ingress
	log.Fatal(err)
}

// func RunCollector() {
// 	var wg sync.WaitGroup
// 	//one cfg is equal to one api target
// 	for _, config := range customConfigs {
// 		Rdb.InsertQuerys(config) //insert querys to redis
// 		fmt.Println("========================================================")
// 		// fmt.Printf("%+v\n", o)

// 		//get all endpoints (including multi parameters-multi endpoints) of whole configFile
// 		endpoints, err := GetEndpoints(config)
// 		if err != nil {
// 			panic(err)
// 		}

// 		for _, endpoint := range endpoints {
// 			go hit_api_and_collect_res_intervally(endpoint, config)
// 			wg.Add(1)
// 		}
// 	}
// 	fmt.Printf("=====load configFile finish=====\n\n")
// 	wg.Wait()
// }

//------------------------------------

// func prettyPrintStruct(i interface{}) string {
// 	s, _ := json.MarshalIndent(i, "", "\t")
// 	return string(s)
// }

// note
// fmt.Printf("[hitAPI] return %d\n", result)
// fmt.Printf("%+v\n", r.Response())
// statusCode := r.Response().StatusCode.(string)

// fmt.Println(cfg[0].Metrics[0].Labels[0].Name)

// // Similarly, we can also fetch data from the database, and the driver
// // will call the Scan() method to unmarshal the data to an Attr struct.
// item := new(Item)
// err = db.QueryRow("SELECT id, attrs FROM items ORDER BY id DESC LIMIT 1").Scan(&item.ID, &item.Attrs)
// if err != nil {
// 	log.Fatal(err)
// }

// // You can then use the struct fields as normal...
// weightKg := item.Attrs.Dimensions.Weight / 1000
// log.Printf("Item: %d, Name: %s, Weight: %.2fkg", item.ID, item.Attrs.Name, weightKg)
