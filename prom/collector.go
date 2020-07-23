package prom

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
)

//定義各種量測類型
type collector struct {
	gaugeName   string
	gaugeMetric *prometheus.Desc
}

func newCollector(guageName string, constLabels map[string]string, labels []string) *collector {
	return &collector{
		gaugeName: guageName,
		gaugeMetric: prometheus.NewDesc(
			guageName,
			fmt.Sprintf("This is a %s help", guageName),
			labels,
			constLabels), //這裡只賦予label名稱, label值在collect()方法裡給
	}
}

//做收集數據時進來
func (e *collector) Collect(ch chan<- prometheus.Metric) {
	//根據key(=e.gaugeName)取值
	PromMetricValue := PromMetric[e.gaugeName]
	PromLabelsValue := PromLabels[e.gaugeName]

	metric, err := prometheus.NewConstMetric(e.gaugeMetric, prometheus.GaugeValue, PromMetricValue, PromLabelsValue...)
	if err != nil {
		glog.Error(e.gaugeName, ":", err)
		return
	}
	//只要有error, 則metric為nil, 會錯
	ch <- metric
}

//這裡控制要註冊哪些metrics
func (e *collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.gaugeMetric
}
