package prom

import (
	"github.com/prometheus/client_golang/prometheus"
)

type fooCollector struct {
	fooMetric *prometheus.Desc
	barMetric *prometheus.Desc
}

//You must create a constructor for you collector that
//initializes every descriptor and returns a pointer to the collector
func newFooCollector() *fooCollector {
	return &fooCollector{
		fooMetric: prometheus.NewDesc("foo_metric",
			"Shows whether a foo has occurred in our cluster",
			nil, nil,
		),
		barMetric: prometheus.NewDesc("bar_metric",
			"Shows whether a bar has occurred in our cluster",
			nil, nil,
		),
	}
}

func (collector *fooCollector) Collect(ch chan<- prometheus.Metric) {
	//定時去更新全局變數 這裡再去拿全局變數的值
	//Implement logic here to determine proper metric value to return to prometheus
	//for each descriptor or call other functions that do so.
	var metricValue float64
	if 1 == 1 {
		metricValue = 1
	}

	//Write latest value for each metric in the prometheus metric channel.
	//Note that you can pass CounterValue, GaugeValue, or UntypedValue types here.
	// ch <- prometheus.MustNewConstMetric(collector.fooMetric, prometheus.CounterValue, metricValue)
	// ch <- prometheus.MustNewConstMetric(collector.barMetric, prometheus.CounterValue, metricValue)
	ch <- prometheus.MustNewConstMetric(collector.fooMetric, prometheus.GaugeValue, metricValue)
	ch <- prometheus.MustNewConstMetric(collector.barMetric, prometheus.GaugeValue, metricValue)

	// ch <- prometheus.MustNewConstMetric(collector.barMetric, prometheus.GaugeValue, metricValue, "labeltest") //if originally no label, will panic
	/*
	   func MustNewConstMetric(desc *Desc, valueType ValueType, value float64, labelValues ...string)
	*/
}

func (collector *fooCollector) Describe(ch chan<- *prometheus.Desc) {

	//Update this section with the each metric you create for a given collector
	ch <- collector.fooMetric
	ch <- collector.barMetric
}
