package main

import (
	"flag"
	"fmt"
	"nms-controller/logic"
	"nms-controller/middleware"
)

// Add() # 添加计数
// Done() # 减掉计数，等价于Add(-1)，这样sync.WaitGroup只有两个API了
// Wait() # 阻塞直到计数为零
func main() {
	fmt.Println("v1.0.2")

	//Init the command-line flags for glog
	flag.Set("v", "5")
	flag.Parse()

	//讓此獨立為一個線程(裡面是一個無限迴圈)，不然會無法繼續下一行
	go logic.RunControllerLoop()

	// RunControllerLoop startProm會初始化regHandler 完成之後router裡面才拿的到

	middleware.Handler()
}
