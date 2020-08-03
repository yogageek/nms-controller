package main

import (
	"log"
)

func test() {
	// a channel to tell it to stop
	stopchann := make(chan struct{})
	// a channel to signal that it's stopped
	stoppedchan := make(chan struct{})
	go func() { // work in background
		// close the stoppedchan when this func
		// exits
		defer close(stoppedchan)
		// TODO: do setup work
		defer func() {
			// TODO: do teardown work
		}()
		for {
			select {
			default:
				// TODO: do a bit of the work
			case <-stopchann:
				// stop
				return
			}
		}
	}()

	log.Println("stopping...")
	close(stopchann) // tell it to stop
	<-stoppedchan    // wait for it to have stopped
	log.Println("Stopped.")
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

/*
func hit_api_and_collect_res_intervally(endpoint string, cfg CustomConfig) chan bool {
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

*/

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
