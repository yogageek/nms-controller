package middleware

import (
	"fmt"
	"net/http"

	"github.com/rs/cors"
)

func Handler() {
	//開啟exporter update, adpater crud API
	r := Router()
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

	//新方法
	//拉出http物件改為可重啟 (postconfig裡面去close)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}
	fmt.Println("http ListenAndServe:", srv.Addr)
	srv.ListenAndServe()

	// 舊方法
	// 設8080為了配合ingress
	// 使用 http.ListenAndServe("localhost:8080", nil) 函数，如果成功会返回空，否则会返回一个错误
	// go最終會停再ListenAndServe無法下一行
	// err := http.ListenAndServe(":8080", handler)
	// log.Fatal(err)
}
