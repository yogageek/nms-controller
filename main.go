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

	// RunControllerLoop startProm會初始化regHandler 完成之後router裡面才拿的到 故需緩衝
	time.Sleep(3 * time.Second)

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

	//設8080為了配合ingress
	//使用 http.ListenAndServe("localhost:8080", nil) 函数，如果成功会返回空，否则会返回一个错误
	//go最終會停再ListenAndServe無法下一行
	err := http.ListenAndServe(":8080", handler)
	log.Fatal(err)
}
