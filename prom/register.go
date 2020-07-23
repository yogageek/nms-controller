package prom

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	//RegHandler 讓router去註冊
	RegHandler http.Handler
)

//StartPrometheus 根據configlist去註冊metric, label名
func StartPrometheus() {
	fmt.Println("StartPrometheus...")

	//new my own register
	reg := prometheus.NewRegistry()

	for _, cfg := range configList {
		//之後取消promDescs  直接寫在collector物件
		promDescs := cfg.GetPromNameAndConstLabelAndLabel()
		for _, promDesc := range promDescs {
			guageName := promDesc.GuageName
			constLabels := promDesc.ConstLabels
			labels := promDesc.Labels

			//方法1
			collector := newCollector(guageName, constLabels, labels)
			reg.MustRegister(collector)

			//方法2
			// exporter := newExporter(guageName, constLabels, labels)
			// reg.MustRegister(exporter) //reg.MustRegister(exporter.gauge) 這樣寫會進不到collector&describe方法
		}
	}
	RegHandler = promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
}

/*
	取消註冊系統預設metrics, 但還是會遺留http metrics, 改用newRegistry可完全減去預設metrics
	a := prometheus.ProcessCollectorOpts{}
	prometheus.Unregister(prometheus.NewProcessCollector(a))
	prometheus.Unregister(prometheus.NewGoCollector())
*/
