package logic

import (
	"fmt"

	"nms-controller/logic/group"
	"nms-controller/logic/prom"
	"nms-controller/model"
	"time"
)

//定時更新的全局變數
var GroupDatas []model.Group

var (
	t1         = time.NewTimer(time.Second * 1) //這裡設幾秒 就會等幾秒才開始
	t2         = time.NewTimer(time.Second * 5)
	t3         = time.NewTimer(time.Second * 3)
	t1Duration = time.Duration(10)
	t2Duration = time.Duration(5)
	t3Duration = time.Duration(5)
)

func RunControllerLoop() {
	//make a chan first
	input := make(chan interface{})
	//producer - produce the messages
	go func() {
		for i := 0; i < 5; i++ {
			input <- i //put data into chan
		}
		input <- "RunControllerLoop..."
	}()

	// run prom
	// prom.StartPrometheus()

	for {
		select {
		//consumer - consume the messages
		case msg := <-input: //take data from chan
			fmt.Println(msg) //will print helle world
		case <-t1.C: //t1.C拿出channel
			// t1.Reset(time.Second * 10) //使t1重新開始計時
			//# 拿掉這個規格: 根據config定義秒數來打api
			prom.UpdatePrometheusData()
			t1.Reset(time.Second * t1Duration) //五秒所有api打一次

		case <-t2.C:
			/*
				fmt.Println("do groups...")
				// fmt.Println("doSon...")
				GroupData1 := group.DoSon()
				// fmt.Println("doAmf...")
				GroupData2 := group.DoAmf()
				// fmt.Println("doImec...")
				GroupData3 := group.DoImec()
				GroupDatas = []model.Group{}
				GroupDatas = append(GroupDatas, GroupData3...)
				GroupDatas = append(append(GroupData1, GroupData2...), GroupData3...)
				t2.Reset(time.Second * t2Duration)
			*/

			//目前只放imec的
			GroupData3 := group.DoImec()
			GroupDatas = []model.Group{} //注意這裡要放在do之後 不然會空很久
			GroupDatas = append(GroupDatas, GroupData3...)
			t2.Reset(time.Second * t2Duration)

			// 由於程式啟動後router只註冊一次Reghandler 故之後就算更新reghandler也不會有變化
			// case <-t3.C:
			// 	prom.StartPrometheus()
			// 	t3.Reset(time.Second * t3Duration)
		}
	}
}
