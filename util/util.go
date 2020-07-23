package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang/glog"
	"github.com/imroc/req"
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

func GetResponseAndCheckJson(url string) ([]byte, error) {
	ts, err := strconv.Atoi(os.Getenv("TIMEOUT_SEC"))
	if err != nil {
		glog.Error(err)
	}
	// fmt.Println("Request URL:" + url)
	timeoutSeconds := time.Duration(ts) * time.Second
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

	json, err := CheckAndPrettyJson(bytes)
	if err != nil {
		glog.Error("URL: ", url, " Response converting to Json Fail! err:=", err)
		return nil, err
	}

	// fmt.Printf("Response Body:%s\n", json)

	return json, nil
}
