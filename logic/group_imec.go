package logic

import (
	"encoding/json"
	"fmt"
	"nms-controller/model"
	"os"

	"github.com/buger/jsonparser"
	"github.com/golang/glog"
)

var uriImec = os.Getenv("URI_IMEC")
var prefix = "/status/"
var imec_management = []string{"gui", "edge", "notc"}
var imec_database = []string{"redis", "mongodb", "rabbitmq"}
var imec_system = []string{"dns", "gwc", "upf"}

//測試資料
var testResponse0 = []byte(`{"status":0}`)
var testResponse1 = []byte(`[{"status":1}]`)

func doImec() (groups []model.Group) {
	//# groupkey名稱先寫死 management database system
	group1 := buildGroupImec("management", imec_management)
	group2 := buildGroupImec("database", imec_database)
	group3 := buildGroupImec("system", imec_system)
	groups = append(groups, group1, group2, group3)
	// fmt.Printf("%+v", groups)
	return groups
}

//根據設定去打多支api
func buildGroupImec(groupkey string, groupsubs []string) model.Group {
	var members []model.Member
	for _, gsub := range groupsubs {
		v := getApi(gsub) //如果api錯誤的話v=nil
		member := model.Member{
			Key:   gsub,
			Value: v,
		}
		members = append(members, member)
	}

	group := model.Group{
		GroupKey: groupkey,
		Members:  members,
	}

	// fmt.Println(group)
	return group
}

func getApi(gsub string) *string {
	uri := uriImec + prefix + gsub
	fmt.Println(uri)
	response, err := getResponseAndCheckJson(uri)
	if err != nil {
		//如果api壞掉或沒回應 一律給2(=檢查中)
		v := "2"
		return &v
	}
	v := getImecValue(response)
	return &v
}

// //for testing
// func getApi(gsub string) (string, error) {
// 	rand.Seed(time.Now().UnixNano())
// 	var v string
// 	if rand.Intn(2) == 1 {
// 		v = getImecValue(testResponse0)
// 	} else {
// 		v = getImecValue(testResponse1)
// 	}
// 	return v, nil
// }

//根據返回json設定取值
//目前只取status
func getImecValue(response []byte) string {
	fmt.Println(string(response))
	keylist := []string{"status"}
	for _, k := range keylist {
		v, err := jsonparser.GetUnsafeString(response, k)
		if err != nil {
			glog.Error(err, ". key=", k, ". response:", string(response))
			//處理upf格式
			var v string
			jsonparser.ArrayEach(response, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
				v, _ = jsonparser.GetUnsafeString(value, k)
			})
			return v
		}
		return v
	}
	glog.Error("getImecValue return empty")
	return ""
}

func toJ(data interface{}) []byte {
	j, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		glog.Error(err)
	}
	fmt.Println(string(j))
	return j
}
