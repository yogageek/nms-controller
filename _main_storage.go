package main

import (
	"time"

	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(3)
	m := make(map[string]chan bool)

	m["func1"] = make(chan bool)
	m["func2"] = make(chan bool)
	m["func3"] = make(chan bool)

	go func() {
		for {
			select {
			case <-m["func1"]:
				println("done 1")
				wg.Done()
				return
			default:
				// Do other stuff
				println("-1-")
				time.Sleep(3 * time.Second)
			}
		}
	}()

	go func() {
		for {
			select {
			case <-m["func2"]:
				println("done 2")
				wg.Done()
				return
			default:
				// Do other stuff
				println("-2-")
				time.Sleep(3 * time.Second)
			}
		}
	}()

	go func() {
		for {
			select {
			case <-m["func3"]:
				println("done 3")
				wg.Done()
				return
			default:
				// Do other stuff
				println("-3-")
				time.Sleep(3 * time.Second)
			}
		}
	}()
	close(m["func2"])
	close(m["func3"])
	wg.Wait()
}
