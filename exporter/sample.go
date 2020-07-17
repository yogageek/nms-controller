package exporter

import (
	"nms-controller/model"

	"github.com/prometheus/client_golang/prometheus"
)

//定義兩個metric輸出
type Exporter struct {
	gauge    prometheus.Gauge    //Gauge: 重點方法set，自己設定各種value 最常用
	gaugeVec prometheus.GaugeVec //GaugeVec: Gauge支援Label
}

//建立量測物件
func NewExporter(metric model.Metric, queryName string, lb []string) *Exporter {

	guageName := queryName + "_" + metric.Header

	gauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "golang",
		Name:      guageName,
		Help:      "This is a yoga help"})

	gaugeVec := *prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "golang",
		Name:      guageName,
		Help:      "This is a yoga help"},
		[]string{"myLabel"})

	return &Exporter{
		gauge:    gauge,
		gaugeVec: gaugeVec,
	}
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.gauge.Set(float64(0))
	e.gaugeVec.WithLabelValues("hello").Set(float64(0))
	e.gauge.Collect(ch)
	e.gaugeVec.Collect(ch)
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	e.gauge.Describe(ch)
	e.gaugeVec.Describe(ch)
}
