package logic

import (
	"nms-controller/db"
	"nms-controller/model"

	"strconv"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/golang/glog"

	"github.com/savaki/jq"
)

var (
	rc         = db.Rc
	pg         = db.Pg
	configList []model.CustomConfig
)

func RunAdapter() {
	configList = pg.GetCustomConfigs() // all configs together

	// just print info
	for _, c := range configList {
		c.PrintQueryNameAndUrl()
	}

	processAdapter(configList)
}

func processAdapter(configList []model.CustomConfig) {
	//one cfg is equal to one api target
	for _, config := range configList {
		//注意是一個config裡的所有endpoints
		urls := getEndpoints(config)

		for _, url := range urls {
			apiResponse, err := getResponseAndCheckJson(url)
			if err != nil {
				glog.Error("api or pgconfig err: ", err)
			} else {
				ParsingJsonAndInsertRedis(apiResponse, config)
			}
		}
		rc.InsertQueryNameKeys(config) //之後可刪
	}
}

//根據config取response的值再塞入redis (如果要改metric取path方式改這裡就好)
func ParsingJsonAndInsertRedis(apiResponse []byte, cfg model.CustomConfig) {
	//一個config可能會有多個metric

	// one cfg of metrics contains multi metric
	for _, metric := range cfg.Metrics {
		path := metric.Path

		//如果回復陣列格式資料
		if strings.Contains(path, ".[].") { //-> response= [{},{}...]
			path = strings.ReplaceAll(path, ".[].", "")
			//取[]裡面{}

			jsonparser.ArrayEach(apiResponse, func(responseEach []byte, dataType jsonparser.ValueType, offset int, err error) {
				v, err := jsonparser.GetUnsafeString(responseEach, path)
				if err != nil {
					glog.Error("取值錯誤-> get [] type from response err:", err, " 取值路徑為:", path, " responseEach:", string(responseEach))
				}
				// fmt.Println("metric path value:", v)
				value, err := strconv.ParseFloat(v, 64)
				if err != nil {
					glog.Error("ParseFloat err:", err)
				}
				labels := getLabel(metric, responseEach)
				uidValue := getUidValue(responseEach, cfg.UidPath)
				redisMetric := model.RedisMetric{
					UID:    uidValue, //modifyyyyyyyyyyyyyyyyyyyyy
					Value:  value,
					Labels: labels,
				}
				rc.InsertRedis(cfg, metric, redisMetric)
			})

			//如果回復物件格式資料

		} else { //取{}裡面 -> response= {}

			route, _ := jq.Parse(path)          //setting是取值路徑
			vb, err := route.Apply(apiResponse) //apply(json物件)
			if err != nil {
				glog.Error("取值錯誤-> get {} type from response err:", err, " 取值路徑為:", path)
				// glog.Error("response:", string(response))
			}
			v := string(vb)
			// fmt.Println("metric path value:", v)
			value, err := strconv.ParseFloat(v, 64)
			if err != nil {
				glog.Error("ParseFloat err:", err)
			}

			labels := getLabel(metric, apiResponse)
			uidValue := getUidValue(apiResponse, cfg.UidPath)
			redisMetric := model.RedisMetric{
				UID:    uidValue,
				Value:  value,
				Labels: labels,
			}
			rc.InsertRedis(cfg, metric, redisMetric)
		}
	}
}
