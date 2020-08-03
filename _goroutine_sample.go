package main

import (
	"fmt"
	"time"
)

//step1: 接收done channel作为参数
func hello(done chan bool) {
	fmt.Println("Hello world goroutine")

	//step2: 然后把true写入done channel
	done <- true
}
func main() {

	//宣告了done := make(chan bool)後 執行go hello(done)之外還需要 <-done 才不會阻塞
	done := make(chan bool)

	go hello(done)

	//step3: <-done表示从done channel接收数据，如果没有任何变量使用或存储该数据，这是完全合法的。
	//注意! 我们正在从done channel接收数据。这行代码是阻塞的，这意味着在Goroutine将数据写入done channel之前将会一直阻塞
	<-done

	//step4: 结束阻塞打印了main函数的文本
	fmt.Println("main function") //main Goroutine被阻塞直到done channel有数据写入
}

//-----------------------------------------------------------------
func hello(done chan bool) {
	fmt.Println("hello go routine is going to sleep")
	time.Sleep(4 * time.Second)
	fmt.Println("hello go routine awake and going to write to done")
	done <- true //并写入数据到channel
}
func main() {
	done := make(chan bool)
	fmt.Println("Main going to call hello go goroutine")
	go hello(done) //就會直接進去方法
	<-done         //main Goroutine将被阻塞，因为它正在等待来自<-done的通道的数据
	fmt.Println("Main received data")
}

//-----------------------------------------------------------------
func UserTicker() chan bool {
	fmt.Println("....")
	ticker := time.NewTicker(2 * time.Second)

	stopChan := make(chan bool)
	stopChan2 := make(chan bool)
	// stopChan <- true //傳進去要有人接

	go func(ticker *time.Ticker) {
		defer ticker.Stop() //Stop會停止Ticker，停止後，Ticker不會再被髮送，但是Stop不會關閉通道

		for {
			select {
			case <-ticker.C:
				fmt.Println("Ticker2....")
			case stop2 := <-stopChan2: //取出stopChan值
				fmt.Println(stop2, "into stopChan2...")

			case stop := <-stopChan: //取出stopChan值
				fmt.Println("into stopChan...")
				if stop {
					fmt.Println("Ticker2 Stop")
					return
				}
			}
		}
	}(ticker) //(這裡是參數)

	fmt.Println("return stopChan") //這裡會在UserTicker跑完馬上回傳
	return stopChan
}
func main() {
	// UserTicker() 本身回傳是一個chan bool 也就是stopChan 所以可以向他發送true
	ch := UserTicker()          //一個通道 裡面一個go routine做無限迴圈
	time.Sleep(6 * time.Second) //這裡要sleep  上面go routine才會一直執行
	ch <- true                  //true會傳進stopChan
	time.Sleep(3 * time.Second) //這裡要sleep  才讓go routine有空檔印出
	close(ch)
}
