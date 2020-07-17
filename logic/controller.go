package logic

import (
	"fmt"

	"nms-controller/model"
	"time"
)

//定時更新的全局變數
var GroupDatas []model.Group

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

	// t1 := time.NewTimer(time.Second * 1) //這裡設幾秒 就會等幾秒
	t2 := time.NewTimer(time.Second * 1)

	//run prom
	// exporter.RunProm()

	for {
		select {
		//consumer - consume the messages
		case msg := <-input: //take data from chan
			fmt.Println(msg) //will print helle world
		// case <-t1.C: //t1.C拿出channel
		// RunAdapter()
		// t1.Reset(time.Second * 10) //使t1重新開始計時
		//# 取消spec 根據config定義秒數來做

		case <-t2.C:
			// 	fmt.Println("doSon...")
			// 	GroupData1 := doSon()
			// 	fmt.Println("doAmf...")
			// 	GroupData2 := doAmf()
			// 	fmt.Println("doImec...")
			GroupData3 := doImec()
			GroupDatas = []model.Group{}
			GroupDatas = append(GroupDatas, GroupData3...)
			// GroupDatas = append(append(GroupData1, GroupData2...), GroupData3...)
			t2.Reset(time.Second * 5)
		}
	}
}
