package model

import "strings"

//之後直接併入collector

type PromDesc struct {
	GuageName   string
	ConstLabels map[string]string
	Labels      []string
}

func (m *Metric) GetGuageName(queryName string) string {
	formattedQueryName := strings.ReplaceAll(queryName, "-", "_") //// queryName := strings.Replace(c.QueryName, "-", "_", -1) //label name 如果有-則替換成底線
	guageName := formattedQueryName + "_" + m.Header
	return guageName
}

//GetPromNameAndLabel 設定guageName string, constLabels map[string]string, labels []string規格
//這裡要改 一個cfg可能return多個guageName string, constLabels map[string]string, labels []string 因為有多組metric
func (c *CustomConfig) GetPromNameAndConstLabelAndLabel() []PromDesc {
	pds := []PromDesc{}
	for _, metric := range c.Metrics {
		// old
		//guageName
		// formattedQueryName := strings.ReplaceAll(c.QueryName, "-", "_") //// queryName := strings.Replace(c.QueryName, "-", "_", -1) //label name 如果有-則替換成底線
		// guageName = formattedQueryName + "_" + metric.Header
		//new
		guageName := metric.GetGuageName(c.QueryName)

		//constLabels
		constLabels := map[string]string{
			"query_name": c.QueryName,
			// "target":     c.Target, //取消
		}
		//labels
		labels := []string{}
		for _, label := range metric.Labels {
			formattedLabel := strings.ReplaceAll(label.Name, " ", "_") //strings.Replace(label.Name, " ", "_", -1) //label name 如果有空格則替換成底線
			labels = append(labels, formattedLabel)
		}
		// glog.Info(guageName, constLabels, labels)
		pd := PromDesc{
			GuageName:   guageName,
			ConstLabels: constLabels,
			Labels:      labels,
		}
		pds = append(pds, pd)
	}
	return pds
}
