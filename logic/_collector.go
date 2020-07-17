package logic

import (
	"fmt"
	"log"
	"nms-controller/model"

	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/buger/jsonparser"
	"github.com/golang/glog"

	"github.com/savaki/jq"
)

func CallRunCollector(restart bool) {
	// CH = UserTicker() //搭配可以停止
	if restart {
		//清除上次線程
		i := 0
		for i < I {
			chanContainer[i] <- true
			close(chanContainer[i])
			i++
		}
		// time.Sleep(time.Duration(20) * time.Second)
		rc.FlushRedis()

		configList = pg.GetCustomConfigs()

		RunCollector(restart)
	} else {
		RunCollector(restart)
	}

	fmt.Println("run collector done")

}

var wg sync.WaitGroup

var I int //上次啟動的線程數量
func RunCollector(restart bool) {
	i := 0
	//one cfg is equal to one api target
	for _, config := range configList {
		rc.InsertQuerys(config) //insert querys to redis
		fmt.Println("========================================================")
		// fmt.Printf("%+v\n", o)

		//get all endpoints (including multi parameters-multi endpoints) of whole configFile
		endpoints := GetEndpoints(config)

		for _, endpoint := range endpoints {
			// wg.Add(1)
			chanContainer[i] = hit_api_and_collect_res_intervally(endpoint, config)
			fmt.Println("RunCollector endpoints index:", i)
			i++
			I = i

			// //測試用
			// hit_api_and_collect_res(endpoint, config)

		}
	}
	fmt.Printf("=====RunCollector load configFile finish=====\n\n")
	// wg.Wait()
}

func hit_api_and_collect_res_intervally(endpoint string, cfg model.CustomConfig) chan bool {
	intervalTime := time.Duration(cfg.Interval)
	uptimeTicker := time.NewTicker(intervalTime * time.Second)

	stopChan := make(chan bool)
	// done := make(chan bool, 1)
	go func(ticker *time.Ticker) {
		defer uptimeTicker.Stop()
		for {
			select {
			case <-uptimeTicker.C:
				fmt.Println("endpoint:", endpoint, "intervalTime:", intervalTime)
				hit_api_and_collect_res(endpoint, cfg)

			// case <-stopChan:
			// 	fmt.Println("[Stop] hit_api_and_collect_res_intervally ")
			// 	wg.Done()
			// 	return
			case stop := <-stopChan:
				if stop {
					fmt.Println("[Stop] hit_api_and_collect_res_intervally ")
					return
				}
			}
		}
	}(uptimeTicker)
	return stopChan
}

func hit_api_and_collect_res(endpoint string, cfg model.CustomConfig) {
	//delete redis list
	//....

	//打api，拿到api response
	jsonResponse, err := getResponseAndCheckJson(endpoint)
	if err != nil {
		glog.Error("GetJsonByUrl err: ", err)
	}
	// fmt.Print("Res:", string(jsonResponse))

	//舊版getResponseValueByMetricPath
	// _ = getResponseValueByMetricPath(cfg, jsonResponse) //由jsonpath取response值
	// fmt.Printf("input redis value:%+v\n", redisMetrics)

	//測試用
	//新版getResponseValueByMetricPath2
	// redisMetrics := getResponseValueByMetricPath2(cfg, jsonResponse) //由jsonpath取response值
	getResponseValueByMetricPath2(cfg, jsonResponse) //由jsonpath取response值
	// fmt.Printf("input redis value:%+v\n", redisMetrics)
}

//如果要改metric取path方式改這裡就好
func getResponseValueByMetricPath2(cfg model.CustomConfig, response []byte) []model.RedisMetric {
	//一個cfg可能會有多個metric
	var redisMetrics []model.RedisMetric // redisMetrics:=[]RedisMetric{}//不能醬寫

	metrics := cfg.Metrics

	// one cfg of metrics contains multi metric
	for _, metric := range metrics {
		rc.DeleteRedis(cfg, metric)
		// fmt.Println("Processing metric:\t", metric.Header)
		path := metric.Path
		// fmt.Println("path:", path)

		//如果回復陣列格式資料
		if strings.Contains(path, ".[].") { //-> response= [{},{}...]

			//# v1.3.2 fix bug
			if response == nil {
				response = []byte(`[]`)
			}

			path = strings.ReplaceAll(path, ".[].", "")
			//取[]裡面{}

			jsonparser.ArrayEach(response, func(responseEach []byte, dataType jsonparser.ValueType, offset int, err error) {
				v, err := jsonparser.GetUnsafeString(responseEach, path)
				if err != nil {
					log.Print("get [] type response err :", err)
					log.Print("responseEach:", string(responseEach))
					log.Print("path:", path)
				}
				// fmt.Println("metric path value:", v)
				value, err := strconv.ParseFloat(v, 64)
				if err != nil {
					log.Print("get [] value err :", err)
				}
				labels := GetLabel(metric, responseEach)
				uidValue := GetUidValue(responseEach, cfg.UidPath)
				redisMetric := model.RedisMetric{
					UID:    uidValue, //modifyyyyyyyyyyyyyyyyyyyyy
					Value:  value,
					Labels: labels,
				}
				rc.InsertRedis(cfg, metric, redisMetric)
				redisMetrics = append(redisMetrics, redisMetric) //delete
			})

			//如果回復物件格式資料
		} else { //取{}裡面 -> response= {}
			// path = strings.ReplaceAll(path, ".", "")

			//# v1.3.2 fix bug
			if response == nil {
				response = []byte(`{}`)
			}

			route, _ := jq.Parse(path)       //setting是取值路徑
			vb, err := route.Apply(response) //apply(json物件)
			if err != nil {
				log.Print("get {} type response err :", err)
				// log.Print("response:", string(response))
				log.Print("path:", path)
			}
			v := string(vb)
			// fmt.Println("metric path value:", v)
			value, err := strconv.ParseFloat(v, 64)
			if err != nil {
				log.Print("get {} value err :", err)
			}

			labels := GetLabel(metric, response)
			uidValue := GetUidValue(response, cfg.UidPath)
			redisMetric := model.RedisMetric{
				UID:    uidValue,
				Value:  value,
				Labels: labels,
			}
			rc.InsertRedis(cfg, metric, redisMetric)
			redisMetrics = append(redisMetrics, redisMetric) //delete
		}

	}
	return redisMetrics
}
