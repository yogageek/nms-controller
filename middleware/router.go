package middleware

import (
	"net/http"
	"nms-controller/exporter"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// Router is exported and used in main.go
func Router() *mux.Router {

	r := mux.NewRouter()

	//groups
	r.HandleFunc("/groups", GetGroups).Methods("GET", "OPTIONS")

	//deprecated
	//update exporter config
	// r.HandleFunc("/update", logic.UpdateCongfig).Methods("GET")

	//postgres config
	r.HandleFunc("/config/{id}", ConfigIDGet).Methods("GET", "OPTIONS")
	r.HandleFunc("/configs", ConfigsGet).Methods("GET", "OPTIONS")
	r.HandleFunc("/config", ConfigPost).Methods("POST", "OPTIONS")
	r.HandleFunc("/config/{id}", ConfigIDPut).Methods("PUT", "OPTIONS") //update adapter pg config
	r.HandleFunc("/config/{id}", ConfigIDDelete).Methods("DELETE", "OPTIONS")

	// 使用allMiddleware中间件处理
	r.Use(allMiddleware)

	//swagger doc
	sh := http.StripPrefix("/doc", http.FileServer(http.Dir("./docs/")))
	r.PathPrefix("/doc/").Handler(sh)
	//swagger
	r.PathPrefix("/swagger").Handler(httpSwagger.Handler(
		httpSwagger.URL("/doc/swagger.json"), //The url pointing to API definition
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("#swagger-ui"),
	))

	//普羅米修斯量測頁面
	//為了用use另創subrouter 會先promMiddleware再allMiddleware
	pr := r.PathPrefix("/metrics").Subrouter()
	pr.Use(promMiddleware) // 使用promMiddleware中间件处理

	//default prom register
	// pr.Handle("", promhttp.Handler())
	//custom prom register
	pr.Handle("", exporter.RegHandler)

	//not used
	// reg := prometheus.NewRegistry()
	// handler := promhttp.InstrumentMetricHandler(reg, promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	// router.Path("/metrics").Handler(handler)

	return r
}

func promMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//CurrentRoute返回当前请求的匹配路由（如果有）. 这仅在匹配路由的处理程序内部调用时有效，因为匹配的路由存储在请求上下文中，请求上下文在处理程序返回后清除.
		mux.CurrentRoute(r)
		//暫時拿掉
		// exporter.GetOperationByMiddleware()
		next.ServeHTTP(w, r)
	})
}

func allMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//CurrentRoute返回当前请求的匹配路由（如果有）. 这仅在匹配路由的处理程序内部调用时有效，因为匹配的路由存储在请求上下文中，请求上下文在处理程序返回后清除.
		mux.CurrentRoute(r)
		next.ServeHTTP(w, r)
	})
}
