package exporter

import (
	"fmt"
	"net/http"
	"nms-controller/db"

	"nms-controller/model"

	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// 存在記憶體中的 map, 提供初始化時結構的 init, 另外在每次 GET /metrics 時, 以 key 取出對應的 prometheus 結構進行賦值
//# 1.2.0

var (
	rc                  = db.Rc
	pg                  = db.Pg
	PrometheusCollector = map[string]prometheus.Collector{}
	redisQueryKeys      []string
	reg                 *prometheus.Registry
	RegHandler          http.Handler
)

//step1
func RunProm() {
	fmt.Println("Run prom...")

	// Register dummy exporter, not necessary

	// get all configs together
	// configList := pg.GetCustomConfigs()

	reg = prometheus.NewRegistry()

	// RegisterPromByConfig(configList)

	//取消註冊系統預設metrics, 但還是會遺留http metrics, 改用newRegistry可完全減去預設metrics
	// a := prometheus.ProcessCollectorOpts{}
	// prometheus.Unregister(prometheus.NewProcessCollector(a))
	// prometheus.Unregister(prometheus.NewGoCollector())

	getex := func() *Exporter {
		gauge := prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "golang",
			Name:      "A",
			Help:      "This is a yoga help"})
		//vec not work
		gaugeVec := *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "golang2",
			Name:      "B",
			Help:      "This is a yoga help"},
			[]string{"myLabel"})

		return &Exporter{
			gauge:    gauge,
			gaugeVec: gaugeVec,
		}
	}

	exporter := getex()

	reg.MustRegister(exporter.gauge)
	// reg.MustRegister(exporter.gaugeVec)
	RegHandler = promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	return
	//改使用pg自己拼湊
	//get queryKeys in redis
	// redisQueryKeys = rc.RedisClient.SMembers("querys").Val()

	// exporter := NewExporter(metricsPrefix)
	// prometheus.MustRegister(exporter)
	// RegisterOrUnregister(PrometheusCollector)
}

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
			reg.MustRegister(exporter.gauge)
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

// 想要加入自己定義的metrics，需在啟動http service前註冊自己的exporter物件到prometheus
