package util

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func CheckAndPrettyJson(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")
	return out.Bytes(), err
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func ToMapList(keys []string, values []string) []map[string]string {
	var maplist []map[string]string

	if !(len(keys) == len(values)) {
		fmt.Println("if length len(keys) == len(values) ->", len(keys) == len(values))
	}

	m := make(map[string]string)
	for i, a := range keys {
		// glog.Info(a)
		// glog.Info(values[i])
		m[a] = values[i]
	}
	maplist = append(maplist, m)
	return maplist
}
