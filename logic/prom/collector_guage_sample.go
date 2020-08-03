package prom

import (
	"github.com/prometheus/client_golang/prometheus"
)

//newExporter or newCollector difference?
//定義各種量測類型
type exporter struct {
	gauge       prometheus.Gauge    //Gauge: 重點方法set，自己設定各種value 最常用
	gaugeVec    prometheus.GaugeVec //GaugeVec: Gauge支援Label
	gaugeMetric *prometheus.Desc
}

func newExporter(guageName string, constLabels map[string]string, labels []string) *exporter {
	//各種量測類型統一寫在這, describe()方法裡去決定註冊哪些

	//right now not used in describe()
	// gauge := prometheus.NewGauge(prometheus.GaugeOpts{
	// 	Name:        guageName,
	// 	Help:        fmt.Sprintf("This is a %s help", guageName),
	// 	ConstLabels: constLabels,
	// })

	// gaugeVec := *prometheus.NewGaugeVec(prometheus.GaugeOpts{
	// 	Name:        guageName,
	// 	Help:        fmt.Sprintf("This is a %s help", guageName),
	// 	ConstLabels: constLabels,
	// }, labels) //這裡只賦予label名稱, label值在collect()方法裡給

	return &exporter{
		// gaugeName: guageName,
		// gauge:     gauge,
		// gaugeVec:  gaugeVec,
	}
}

//做收集數據時進來
func (e *exporter) Collect(ch chan<- prometheus.Metric) {

	// e.gauge.Set(float64(0))
	// e.gaugeVec.set

	//如果label錯會panic
	// e.gaugeVec.WithLabelValues("hello").Set(float64(0))
	// e.gauge.Collect(ch)
	// e.gaugeVec.Collect(ch)

	// var metricValue float64
	// command := string("date +%s")
	// cmdResult := exeCmd(command)
	// metricValue = cmdResult
	// ch <- prometheus.MustNewConstMetric(collector.cmdMetric, prometheus.GaugeValue, metricValue)

}

//做註冊動作時進來
//這裡控制要註冊哪些metrics
func (e *exporter) Describe(ch chan<- *prometheus.Desc) {
	//如果啟動註冊panic直接註解一行
	// e.gauge.Describe(ch)
	// e.gaugeVec.Describe(ch)

	/*
		ch <- e.gaugeMetric這樣寫的話不能搭配
		//方法1
		exporter := newExporter(guageName, constLabels, labels)
		reg.MustRegister(exporter) //reg.MustRegister(exporter.gauge) 這樣寫會進不到collector&describe方法
	*/

	ch <- e.gaugeMetric
}
