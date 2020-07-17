package exporter

import (
	"fmt"

	"nms-controller/model"

	"github.com/prometheus/client_golang/prometheus"
)

//建立guagec和counter量測
func promContainer(metric model.Metric, queryName string, constLabels map[string]string, labels []string) map[string]prometheus.Collector {

	switch metricType := metric.Type; metricType {
	case "gauge":
		gauge := prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name:        queryName + "_" + metric.Header,
				Help:        metric.Help,
				ConstLabels: constLabels,
			},
			labels,
		)
		PrometheusCollector[queryName+":"+metric.Header+":gauge"] = *gauge
	case "counter":
		counter := prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        queryName + "_" + metric.Header,
				Help:        metric.Help,
				ConstLabels: constLabels,
			},
			labels,
		)
		PrometheusCollector[queryName+":"+metric.Header+":counter"] = *counter
	default:
		fmt.Println("Type was not defined")
	}
	return PrometheusCollector
}

//建立guagec和counter量測
func buildPrometheusCollector(metric model.Metric, queryName string, constLabels map[string]string, labels []string) {
	switch metricType := metric.Type; metricType {
	case "gauge":
		gauge := prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name:        queryName + "_" + metric.Header,
				Help:        metric.Help,
				ConstLabels: constLabels,
			},
			labels,
		)
		PrometheusCollector[queryName+":"+metric.Header+":gauge"] = *gauge
	case "counter":
		counter := prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        queryName + "_" + metric.Header,
				Help:        metric.Help,
				ConstLabels: constLabels,
			},
			labels,
		)
		PrometheusCollector[queryName+":"+metric.Header+":counter"] = *counter
	default:
		fmt.Println("Type was not defined")
	}
	// switch metricType := metric.Type; metricType {
	// case "gauge":
	// 	gauge := prometheus.NewGaugeVec(
	// 		prometheus.GaugeOpts{
	// 			Name:        queryName + "_" + metric.Header,
	// 			Help:        metric.Help,
	// 			ConstLabels: constLabels,
	// 		},
	// 		labels,
	// 	)
	// 	PrometheusCollector[queryName+":"+metric.Header+":gauge"] = *gauge
	// case "counter":
	// 	counter := prometheus.NewCounterVec(
	// 		prometheus.CounterOpts{
	// 			Name:        queryName + "_" + metric.Header,
	// 			Help:        metric.Help,
	// 			ConstLabels: constLabels,
	// 		},
	// 		labels,
	// 	)
	// 	PrometheusCollector[queryName+":"+metric.Header+":counter"] = *counter
	// default:
	// 	fmt.Println("Type was not defined")
	// }
}
