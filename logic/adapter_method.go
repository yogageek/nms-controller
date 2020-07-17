package logic

import (
	"fmt"
	"log"
	"nms-controller/model"
	"nms-controller/util"
	"strings"
	"time"

	"github.com/buger/jsonparser"
	"github.com/golang/glog"
	"github.com/imroc/req"
)

func getResponseAndCheckJson(url string) ([]byte, error) {

	// fmt.Println("Request URL:" + url)
	timeoutSeconds := time.Duration(5) * time.Second
	req.SetTimeout(timeoutSeconds)
	r, err := req.Get(url)
	if err != nil {
		glog.Error("URL: ", url, " FAIL!\nerr:=", err)
		return nil, err
	}

	rcode := r.Response().StatusCode
	if rcode != 200 {
		err := fmt.Errorf("URL: ", url, " FAIL! Return status code=%d", rcode)
		glog.Error(err)
		return nil, err
	}

	bytes, err := r.ToBytes()
	if err != nil {
		glog.Error("response converts to bytes err:", err)
		return nil, err
	}

	json, err := util.CheckAndPrettyJson(bytes)
	if err != nil {
		glog.Error("URL: ", url, " Response converting to Json Fail! err:=", err)
		return nil, err
	}

	// fmt.Printf("Response Body:%s\n", json)

	return json, nil
}

//取得一個CustomConfig的所有endpoint
func getEndpoints(cfg model.CustomConfig) []string {

	var endpoints []string

	if !cfg.HaveParameters { // target = one endpoint
		fmt.Println("query無參數" + cfg.QueryName)
		endpoint := cfg.Target
		endpoints = append(endpoints, endpoint)
	} else { // target = query parameters x endpoint
		fmt.Println("query有參數:" + cfg.QueryName)
		endpoint := cfg.Target
		queryParameters := cfg.QueryParameters
		//拼湊endpoint
		for _, qp := range queryParameters {
			qpName := qp.Name
			qpValue := qp.Value
			for _, replaceValue := range qpValue {
				endpointValue := strings.ReplaceAll(endpoint, qpName, replaceValue)
				endpoints = append(endpoints, endpointValue)
			}
		}
	}
	// fmt.Printf("%+v\n", endpoints)
	return endpoints
}

//取config中metrics(多個)的path
func getMetrics(cfg model.CustomConfig) []string {
	var jsonPaths []string
	metrics := cfg.Metrics
	for _, metric := range metrics {
		jsonPath := metric.Path
		jsonPaths = append(jsonPaths, jsonPath)
		// for _, label := range metric.Labels {
		// 	labelName := label.Name
		// 	labelPath := label.Path
		// }
	}
	return jsonPaths
}

func getLabel(metric model.Metric, response []byte) []model.Label {
	labels := []model.Label{}
	for _, cfgLabel := range metric.Labels {
		//取label value優先, 如果value為空則取path
		if cfgLabel.Value != "" {
			label := model.Label{
				Name:  cfgLabel.Name,
				Value: cfgLabel.Value,
			}
			labels = append(labels, label)
		} else {
			if strings.Contains(cfgLabel.Path, ".[].") { //-> response= [{},{}...]

				path := strings.ReplaceAll(cfgLabel.Path, ".[].", "")
				//取[]裡面{}

				v, err := jsonparser.GetUnsafeString(response, path)
				if err != nil {
					log.Print(err)
				}
				// fmt.Println("label path value:", v)
				label := model.Label{
					Name:  cfgLabel.Name,
					Value: v,
				}
				labels = append(labels, label)

			} else { //取{}裡面 -> response= {}
				path := strings.ReplaceAll(cfgLabel.Path, ".", "")
				v, err := jsonparser.GetUnsafeString(response, path)
				if err != nil {
					log.Print(err)
				}
				// fmt.Println("label path value:", v)
				label := model.Label{
					Name:  cfgLabel.Name,
					Value: v,
				}
				labels = append(labels, label)
			}
		}

	}
	return labels
}

func getUidValue(response []byte, uidPath string) string {
	if strings.Contains(uidPath, ".[].") { //-> response= [{},{}...]
		uidPath = strings.ReplaceAll(uidPath, ".[].", "")
		fmt.Println(string(response))
	} else {
		uidPath = strings.ReplaceAll(uidPath, ".", "")
	}
	uidValue, _ := jsonparser.GetUnsafeString(response, uidPath)
	return uidValue
}
