package logic

import (
	"nms-controller/model"
	"nms-controller/util"
	"os"

	"github.com/buger/jsonparser"
	"github.com/golang/glog"
)

var uriSon = os.Getenv("URI_SON")

func doSon() (groups []model.Group) {
	response, err := getApiForAll(uriSon)
	if err != nil {
		glog.Error(err)
	}
	reponseMaplist := getSonValue(response)
	groups = buildGroupSon("", reponseMaplist)
	// fmt.Printf("%+v", groups)
	return groups
}

//根據返回json設定取值
func getSonValue(response []byte) map[string][]map[string]string {

	objlist1 := []string{"sonHM", "sonHMConnectivity"}
	objlist2 := []string{"bsHM", "sonPM"}

	mapmaplist := map[string][]map[string]string{}
	// maplists := []map[string]string{}

	klist := []string{} // var klist []string
	vlist := []string{} // var vlist []string

	for _, str := range objlist1 {
		kPaths := [][]string{
			//抓不到的會直接跳過
			{str, "[0]", "name"},
			{str, "[1]", "name"},
			{str, "[2]", "name"},
			{str, "[3]", "name"},
			{str, "[4]", "name"},
			{str, "[5]", "name"},
			{str, "[6]", "name"},
			{str, "[7]", "name"},
		}
		vPaths := [][]string{
			{str, "[0]", "status"},
			{str, "[1]", "status"},
			{str, "[2]", "status"},
			{str, "[3]", "status"},
			{str, "[4]", "status"},
			{str, "[5]", "status"},
			{str, "[6]", "status"},
			{str, "[7]", "status"},
		}
		klist = paForSon1(response, kPaths)
		vlist = paForSon1(response, vPaths)
		maplist := util.ToMapList(klist, vlist)
		/*
			maplists = append(maplists, maplist...)
			//map[db:2 dc:1 dp:0 field:1 genie:1 oam:0 routine:1 son:1]
		*/
		mapmaplist[str] = maplist
	}

	for _, str := range objlist2 {
		kPaths := [][]string{
			{str, "offlineDevices"},
			{str, "onlineDevices"},
			{str, "failedDevices"},
			{str, "fieldNumber"},
			{str, "deviceNumber"},
		}
		maplist := paForSon2(response, kPaths)
		// maplists = append(maplists, maplist...)
		mapmaplist[str] = maplist
	}
	// fmt.Printf("maplists:%+v", maplists)
	return mapmaplist
}

//son回應格式有兩種類型，故有兩種不同parse法
func paForSon1(response []byte, paths [][]string) []string {
	var strlist []string
	jsonparser.EachKey(response, func(idx int, value []byte, vt jsonparser.ValueType, err error) {
		// fmt.Println(paths[idx])
		// fmt.Println(paths[idx][0])
		strlist = append(strlist, string(value))
		// switch idx {
		// case 0: // []string{"person", "name", "fullName"}
		// 	fmt.Println(string(value))
		// case 1: // []string{"person", "avatars", "[0]", "url"}
		// 	fmt.Println(string(value))
		// case 2: // []string{"company", "url"},
		// 	fmt.Println(string(value))
		// }
	}, paths...)
	return strlist
}

//son回應格式有兩種類型，故有兩種不同parse法
func paForSon2(response []byte, paths [][]string) []map[string]string {
	var maplist []map[string]string
	jsonparser.EachKey(response, func(idx int, value []byte, vt jsonparser.ValueType, err error) {
		key := paths[idx][1]
		m := make(map[string]string)
		m[key] = string(value)
		maplist = append(maplist, m)
	}, paths...)
	return maplist // [map[offlineDevices:43] map[failedDevices:59]]
}

func buildGroupSon(name string, mapmaplist map[string][]map[string]string) []model.Group {

	var groups []model.Group
	// map list
	for g, m := range mapmaplist {
		var members []model.Member
		for _, v := range m {
			for k, v := range v {
				member := model.Member{
					Key:   k,
					Value: v,
				}
				members = append(members, member)
			}
		}
		group := model.Group{
			GroupKey: g,
			Members:  members,
		}
		// fmt.Println(group)
		groups = append(groups, group)
	}
	return groups
}

// func buildGroupSon(name string, mapmaplist map[string][]map[string]string) []model.Group {

// 	var groups []model.Group
// 	// map list
// 	for _, m := range mapmaplist {
// 		var members []model.Member
// 		//map
// 		for k, v := range m {
// 			member := model.Member{
// 				Key:   k,
// 				Value: v,
// 			}
// 			members = append(members, member)
// 		}
// 		group := model.Group{
// 			GroupKey: name,
// 			Members:  members,
// 		}
// 		fmt.Println(group)
// 		groups = append(groups, group)
// 	}
// 	return groups
// }
