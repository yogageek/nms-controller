package prom

import (
	"fmt"
	"nms-controller/model"
	"nms-controller/util"
	"strconv"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/golang/glog"
	"github.com/savaki/jq"
)

//map key is guageName
var PromMetric map[string]float64
var PromLabels map[string][]string

func init() {
	PromMetric = make(map[string]float64)
	PromLabels = make(map[string][]string)
}

//UpdatePrometheusData 根據configlist去更新metric, label值
func UpdatePrometheusData() {
	fmt.Println("UpdatePrometheusData...")
	for _, cfg := range configList {
		updateValueForMetricAndLabels(cfg)
	}
}

//將打api取到的值寫入map, key為guageName
func updateValueForMetricAndLabels(c model.CustomConfig) {

	url := c.Target

	apiRes, err := util.GetResponseAndCheckJson(url)
	if err != nil {
		return
	}

	//一個metric有一個value, 多個label
	for _, metric := range c.Metrics {
		guageName := metric.GetGuageName(c.QueryName)
		formattedPath, ifArray := ifConfigPathIsArray(metric.Path)
		metricValue := getValueInJSON(apiRes, formattedPath, ifArray)
		setPromMetricValue(guageName, metricValue)

		labels := []string{}
		for _, label := range metric.Labels {
			if label.Value != "null" { //# null之後要改掉
				formattedPath, ifArray := ifConfigPathIsArray(label.Path)
				labelValue := getValueInJSON(apiRes, formattedPath, ifArray)
				labels = append(labels, labelValue)
				//fmt.Sprintf("%f", labelValue)) //float to string
			}
		}
		setPromLabelsValue(guageName, labels)
	}
}

//判斷configPath是用來取json或json array的值
func ifConfigPathIsArray(configPath string) (string, bool) {
	if strings.Contains(configPath, ".[].") {
		formattedPath := strings.ReplaceAll(configPath, ".[].", "")
		return formattedPath, true
	}
	// formattedPath := strings.ReplaceAll(configPath, ".", "")
	formattedPath := configPath
	return formattedPath, false
}

func getValueInJSON(apiRes []byte, formattedPath string, ifArray bool) string {
	getValueInJSONArray := func() string {
		var value string
		jsonparser.ArrayEach(apiRes, func(resEach []byte, dataType jsonparser.ValueType, offset int, err error) {
			value, err = jsonparser.GetUnsafeString(resEach, formattedPath)
			if err != nil {
				glog.Error("取值錯誤-> get [] type from response err:", err, " 取值路徑為:", formattedPath, " responseEach:", string(resEach))
			}
		})
		return value
	}
	//需要用jq 因為son config有index順序關係
	getValueInJSONObject := func() string {
		route, err := jq.Parse(formattedPath) //setting是取值路徑
		if err != nil {
			glog.Error("jq.Parse err", err)
		}
		b, err := route.Apply(apiRes) //apply(json物件)
		if err != nil {
			glog.Error("取值錯誤-> get {} type from response err:", err, " 取值路徑為:", formattedPath)
		}
		return string(b)
	}

	if ifArray {
		return getValueInJSONArray()
	} else {
		return getValueInJSONObject()
	}
}

func setPromMetricValue(key string, value string) {
	value64, err := strconv.ParseFloat(value, 64)
	if err != nil {
		glog.Error("ParseFloat err:", err)
	}
	PromMetric[key] = value64
}

func setPromLabelsValue(key string, value []string) {
	PromLabels[key] = value
}
