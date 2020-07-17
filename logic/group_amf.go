package logic

import (
	"nms-controller/model"
	"os"

	"github.com/buger/jsonparser"
	"github.com/golang/glog"
)

var uriAmf = os.Getenv("URI_AMF")

func doAmf() (groups []model.Group) {
	response, err := getApiForAll(uriAmf)
	if err != nil {
		glog.Error(err)
	}
	reponseMaplist := getAmfValue(response)
	groups = buildGroupAmf(reponseMaplist)

	// fmt.Printf("%+v", groups)
	return groups
}

func getApiForAll(uri string) ([]byte, error) {
	response, err := getResponseAndCheckJson(uri)
	if err != nil {
		return nil, err
	}
	return response, nil
}

//取值並放成map後再拼裝成group
func buildGroupAmf(maplists []map[string]string) []model.Group {

	var groups []model.Group

	// map list
	for _, m := range maplists {
		//# AMF groupkey名稱有特殊規格 是取決於response json key=IMSI的值
		amfGroupkey := m["IMSI"]
		//map
		var members []model.Member //相當於一個group裡面的jsons
		for k, v := range m {
			member := model.Member{
				Key:   k,
				Value: v,
			}
			members = append(members, member)
		}
		group := model.Group{
			GroupKey: amfGroupkey,
			Members:  members,
		}
		groups = append(groups, group)
	}

	// fmt.Println(group)
	return groups
}

//根據返回json設定取值
func getAmfValue(response []byte) []map[string]string {
	klist := []string{"IMSI", "gNBIP", "ULAMBR", "DLAMBR", "ECMState", "EMMState"}

	maplists := []map[string]string{}
	maplist := paForAmf(response, klist)
	maplists = append(maplists, maplist...)

	return maplists
}

func paForAmf(response []byte, keylist []string) []map[string]string {
	var maplist []map[string]string
	// jsonparser.EachKey(data, func(idx int, value []byte, vt jsonparser.ValueType, err error) {
	// 	fmt.Println(value)
	// 	// fmt.Println(string(value))
	// 	// key := paths[idx][1]
	// 	// m := make(map[string]string)
	// 	// m[key] = string(value)
	// 	// maplist = append(maplist, m)
	// }, paths...)

	// // You can use `ArrayEach` helper to iterate items [item1, item2 .... itemN]
	jsonparser.ArrayEach(response, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		m := make(map[string]string)
		for _, k := range keylist {
			v, _ := jsonparser.GetUnsafeString(value, k)
			// fmt.Println(string(v))
			m[k] = string(v)
		}
		maplist = append(maplist, m)
	})

	// // // You can use `ObjectEach` helper to iterate objects { "key1":object1, "key2":object2, .... "keyN":objectN }
	// jsonparser.ObjectEach(data, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
	// 	fmt.Printf("Key: '%s'\n Value: '%s'\n Type: %s\n", string(key), string(value), dataType)
	// 	return nil
	// })

	// fmt.Println(maplist)

	return maplist
}

// // You can use `ArrayEach` helper to iterate items [item1, item2 .... itemN]
// jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
// 	fmt.Println(jsonparser.Get(value, "url"))
// }, "person", "avatars")

// // Or use can access fields by index!
// jsonparser.GetString(data, "person", "avatars", "[0]", "url")

// // You can use `ObjectEach` helper to iterate objects { "key1":object1, "key2":object2, .... "keyN":objectN }
// jsonparser.ObjectEach(data, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
//         fmt.Printf("Key: '%s'\n Value: '%s'\n Type: %s\n", string(key), string(value), dataType)
// 	return nil
// }, "person", "name")
