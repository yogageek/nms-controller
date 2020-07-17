package logic

import (
	"fmt"
	"strings"
	"time"

	"github.com/buger/jsonparser"
	"github.com/golang/glog"
)

//----------------------------
//依照metric給的path拿出response的value
// func getResponseValueByMetricPath(cfg CustomConfig, response []byte) []RedisMetric {
// 	//一個cfg可能會有多個metric
// 	var redisMetrics []RedisMetric // redisMetrics:=[]RedisMetric{}//不能醬寫

// 	metrics := cfg.Metrics

// 	// one cfg of metrics contains multi metric
// 	for _, metric := range metrics {
// 		Rdb.DeleteRedis(cfg, metric)
// 		// fmt.Println("Processing metric:\t", metric.Header)
// 		path := metric.Path
// 		// fmt.Println("path:", path)

// 		if strings.Contains(path, ".[].") { //-> response= [{},{}...]
// 			path = strings.ReplaceAll(path, ".[].", "")
// 			//取[]裡面{}

// 			jsonparser.ArrayEach(response, func(responseEach []byte, dataType jsonparser.ValueType, offset int, err error) {
// 				v, err := jsonparser.GetUnsafeString(responseEach, path)
// 				if err != nil {
// 					log.Print(err)
// 				}
// 				// fmt.Println("metric path value:", v)
// 				value, err := strconv.ParseFloat(v, 64)
// 				if err != nil {
// 					// log.Print(err)
// 				}
// 				labels := GetLabel(metric, responseEach)
// 				uidValue := GetUidValue(responseEach, cfg.UidPath)
// 				redisMetric := RedisMetric{
// 					UID:    uidValue, //modifyyyyyyyyyyyyyyyyyyyyy
// 					Value:  value,
// 					Labels: labels,
// 				}
// 				Rdb.InsertRedis(cfg, metric, redisMetric)
// 				redisMetrics = append(redisMetrics, redisMetric) //delete
// 			})

// 		} else { //取{}裡面 -> response= {}
// 			path = strings.ReplaceAll(path, ".", "")
// 			v, err := jsonparser.GetUnsafeString(response, path)
// 			if err != nil {
// 				log.Print(err)
// 			}
// 			// fmt.Println("metric path value:", v)
// 			value, err := strconv.ParseFloat(v, 64)
// 			if err != nil {
// 				// log.Print(err)
// 			}
// 			labels := GetLabel(metric, response)
// 			uidValue := GetUidValue(response, cfg.UidPath)
// 			redisMetric := RedisMetric{
// 				UID:    uidValue,
// 				Value:  value,
// 				Labels: labels,
// 			}
// 			Rdb.InsertRedis(cfg, metric, redisMetric)
// 			redisMetrics = append(redisMetrics, redisMetric) //delete
// 		}

// 	}
// 	return redisMetrics
// }

//根據path取值
func getValueByPath(response []byte, path string) interface{} {
	//transfer path first
	switch strings.Contains(path, ".[].") {
	case true:
		path = strings.ReplaceAll(path, ".[].", "") //給jq取值用 不需要帶點
	case false:
		path = strings.ReplaceAll(path, ".", "") //給jq取值用 不需要帶點
	default:
		fmt.Println("unexpected path:", path)
	}

	value, err := jsonparser.GetUnsafeString(response, path)
	if err != nil {
		glog.Error(err)
		glog.Error("responseEach:", string(response))
	}
	return value
}

func CloseGoroutine() {
	doneChan := make(chan interface{})

	go func(done <-chan interface{}) {
		for {
			select {
			case <-done:
				return
			default:
			}
		}
	}(doneChan)

	// 父 goroutine 关闭子 goroutine
	close(doneChan)
}

//----------------------------
// ch := UserTicker()
// time.Sleep(20 * time.Second)
// ch <- true
// close(ch)
func UserTicker() chan bool {
	ticker := time.NewTicker(5 * time.Second)

	stopChan := make(chan bool)

	go func(ticker *time.Ticker) {
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				fmt.Println("Ticker2....")
			case stop := <-stopChan: //取出stopChan值
				if stop {
					fmt.Println("Ticker2 Stop")
					return
				}
			}
		}
	}(ticker)

	return stopChan
}

// fmt.Println("if length key list == value list ->", len(klist) == len(vlist))

// jsonparser.ArrayEach(response, func(responseEach []byte, dataType jsonparser.ValueType, offset int, err error) {
// 	value, err := jsonparser.GetUnsafeString(response, "sonHM", "name")
// 	if err != nil {
// 		glog.Error(err)
// 		glog.Error("responseEach:", string(response))
// 	}
// 	fmt.Println(value)
// })
// return 0

// key會重複  無法
// var m map[string]interface{}
// for i, k := range klist {
// }

// func paObj(response []byte, objs []string) []map[string]interface{} {
// 	var strlist []string
// 	// You can use `ObjectEach` helper to iterate objects { "key1":object1, "key2":object2, .... "keyN":objectN }
// 	jsonparser.ObjectEach(response, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
// 		fmt.Printf("Key: '%s'\n Value: '%s'\n Type: %s\n", string(key), string(value), dataType)
// 		return nil
// 	}, "bsHM")
// 	return strlist
// }

// func BuildSun(klist []string, vlist []string) model.Sun {
// 	var members []model.Member
// 	for i, k := range klist {
// 		member := model.Member{
// 			Key:   k,
// 			Value: vlist[i],
// 		}
// 		members = append(members, member)
// 	}
// 	sun := model.Sun{
// 		GroupKey:  "son",
// 		GroupName: "son",
// 		Members:   members,
// 	}

// 	fmt.Println(sun)

// 	return sun
// }

/*
	// data := []byte(`{
	// 	"person": {
	// 	  "name": {
	// 		"first": "Leonid",
	// 		"last": "Bugaev",
	// 		"fullName": "Leonid Bugaev"
	// 	  },
	// 	  "github": {
	// 		"handle": "buger",
	// 		"followers": 109
	// 	  },
	// 	  "avatars": [
	// 		{ "url": "https://avatars1.githubusercontent.com/u/14009?v=3&s=460", "type": "thumbnail" }
	// 	  ]
	// 	},
	// 	"company": {
	// 	  "name": "Acme"
	// 	}
	//   }`)
	// jsonparser.ObjectEach(response, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
	// 	fmt.Printf("Key: '%s'\n Value: '%s'\n Type: %s\n", string(key), string(value), dataType)
	// 	return nil
	// }, "sonHM")
*/
