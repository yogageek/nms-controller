package exporter

import (
	"encoding/json"
	"fmt"
	"nms-controller/db"
	"nms-controller/util"

	"nms-controller/model"

	"strings"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
)

// 存在記憶體中的 map, 提供初始化時結構的 init, 另外在每次 GET /metrics 時, 以 key 取出對應的 prometheus 結構進行賦值
//# 1.2.0

var (
	rc                  = db.Rc
	pg                  = db.Pg
	PrometheusCollector = map[string]prometheus.Collector{}
	redisQueryKeys      []string
)

//step1
func RunProm() {
	fmt.Println("Run prom...")

	// Register dummy exporter, not necessary

	// get all configs together
	configList := pg.GetCustomConfigs()

	RegisterPromByConfig(configList)

	//改使用pg自己拼湊
	//get queryKeys in redis
	redisQueryKeys = rc.RedisClient.SMembers("querys").Val()

	// exporter := NewExporter(metricsPrefix)
	// prometheus.MustRegister(exporter)
	// RegisterOrUnregister(PrometheusCollector)
}

/*
func RunProm() {
	fmt.Println("Run prom...")

	// Register dummy exporter, not necessary

	// get all configs together
	configList := pg.GetCustomConfigs()

	//set up PrometheusCollector
	analyzePgConfigAndBuild(configList)

	//get queryKeys in redis
	redisQueryKeys = rc.RedisClient.SMembers("querys").Val()

	// exporter := NewExporter(metricsPrefix)
	// prometheus.MustRegister(exporter)
	// RegisterOrUnregister(PrometheusCollector)
}*/

func RegisterPromByConfig(cl []model.CustomConfig) {
	// fmt.Println("分析pg config, 生成 MustRegister")
	for _, cfg := range cl {

		// constLabels := map[string]string{
		// 	"query_name": cfg.QueryName,
		// 	"target":     cfg.Target,
		// }

		for _, metric := range cfg.Metrics {
			queryName := strings.Replace(cfg.QueryName, "-", "_", -1) //label name 如果有-則替換成底線
			// fmt.Println(cfg.QueryName + ":" + metric.Header)
			labels := []string{}
			for _, label := range metric.Labels {
				ln := strings.Replace(label.Name, " ", "_", -1) //label name 如果有空格則替換成底線
				labels = append(labels, ln)
			}
			exporter := NewExporter(metric, queryName, labels)
			prometheus.MustRegister(exporter.gauge)
			// buildPrometheusCollector(metric, queryName, constLabels, labels)
		}
	}
}

// logic啟動先selectConfig 這裡用cfg拿出
// get config in pg, register
func analyzePgConfigAndBuild(cl []model.CustomConfig) {
	// fmt.Println("分析pg config, 生成 MustRegister")
	for _, cfg := range cl {

		constLabels := map[string]string{
			"query_name": cfg.QueryName,
			"target":     cfg.Target,
		}

		for _, metric := range cfg.Metrics {
			queryName := strings.Replace(cfg.QueryName, "-", "_", -1) //label name 如果有-則替換成底線
			// fmt.Println(cfg.QueryName + ":" + metric.Header)
			labels := []string{}
			for _, label := range metric.Labels {
				ln := strings.Replace(label.Name, " ", "_", -1) //label name 如果有空格則替換成底線
				labels = append(labels, ln)
			}
			buildPrometheusCollector(metric, queryName, constLabels, labels)
		}
	}
}

//#1.0.0 bug , PrometheusCollector變數不要放記憶體, 之後再處理這個問題

// 想要加入自己定義的metrics，需在啟動http service前註冊自己的exporter物件到prometheus

//get query from redis, compare if in pc
func RegisterOrUnregister(pc map[string]prometheus.Collector) {
	fmt.Println("Prometheus Registion...")

	for pcKey, collector := range pc {
		//如果prom key有在redis queryKey裡面 就註冊
		if util.StringInSlice(pcKey, redisQueryKeys) {
			if err := prometheus.Register(collector); err != nil {
				glog.Error("Register Fail:", pcKey, " err: ", err)
			} else {
				fmt.Println("Register Success:", pcKey) //collector SON:HealthManagement_Genie:gauge registered.
			}
		} else {
			//如果prom key 沒有在redis queryKey裡面 就註銷註冊
			if prometheus.Unregister(collector) {
				fmt.Println("Unregister Success:", pcKey)
			} else {
				fmt.Println("Unregister Fail:", pcKey)
			}
		}
	}
}

// 從 Redis 讀出來之後, 取得數值, 再去向 prometheusCollector 這個存在記憶體中的 map 取出對應的 prometheus 結構進行賦值
// https://yami.io/golang-prometheus/
// https://mojotv.cn/go/prometheus-client-for-go
func GetOperationByMiddleware() {
	fmt.Println("Prometheus Server 每次調用 GET /metrics 時, 執行的 function, 會啟動 Redis 連線")
	rc := rc.RedisClient

	var metric model.MetricE
	var j int64

	for _, metricKey := range redisQueryKeys {

		size := rc.LLen(metricKey).Val()

		for j = 0; j < size; j++ {
			m, err := rc.LIndex(metricKey, j).Result()
			if err != nil {
				panic(err)
			}
			data := []byte(m)
			json.Unmarshal(data, &metric)

			//現有記憶體變數PrometheusCollector
			collector, ok := PrometheusCollector[metricKey]

			labels := make(map[string]string)
			// labels["user"] = metric.UID
			for _, label := range metric.Labels {
				labelName := strings.Replace(label.Name, " ", "_", -1)
				labels[labelName] = label.Value
				//fmt.Println("name:" + label.Name + "\t" + "value:" + label.Value)
			}
			if ok {
				keyPath := strings.Split(metricKey, ":")
				metricType := keyPath[len(keyPath)-1]
				switch metricType {
				case "gauge":
					gauge := collector.(prometheus.GaugeVec)
					gauge.With(prometheus.Labels(labels)).Set(metric.Value)
				case "counter":
					counter := collector.(prometheus.CounterVec)
					if metric.Value == -1 {
						counter.Reset()
					}
					counter.With(prometheus.Labels(labels)).Add(metric.Value)
				default:
					fmt.Println("Type was not defined")
				}

			} else {
				fmt.Println("Can not found " + metricKey)
			}

		}
	}
}

//depre
// //打api更新exporter 之後可拿掉 直接透過controller的api
// func UpdateCongfig(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 	fmt.Println("prometheusCollector length:", len(PrometheusCollector))
// 	unregisPrometheusCollector() //清除所有已經register的
// 	analyzePgConfig()            //分析config
// 	doRegis()
// 	w.WriteHeader(200)
// }
