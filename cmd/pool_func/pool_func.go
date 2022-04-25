package main

import (
	"fmt"
	"time"
)

func main() {

	readCh(1)
	readCh(2)

	time.Sleep(time.Second)
	fmt.Println("Done.")
}

var ch = make(chan int, 1)

func a() {
	ch <- 1
	fmt.Println("a done")
}

func b() {
	ch <- 2
	fmt.Println("b done")
}

func readCh(d int) {
	select {
	case ch <- d:
		fmt.Println("ok d =", d)
	default:
		fmt.Println("def d =", d)
	}
}

var taskQueue = make(chan func(), 1)

func runInPool(task func()) {
	select {
	case taskQueue <- task:
		// submitted, everything is ok
		fmt.Println("ok")
		task()

	default:
		go func() {
			// do the given task
			fmt.Println("no ok")
			//task()

			const interval = 10 * time.Second
			ticker := time.NewTicker(interval)
			defer ticker.Stop()

			for {
				select {
				case t := <-taskQueue:
					t()
					ticker.Reset(interval)
				case <-ticker.C:
					return
				}
			}
		}()
	}
}

func processFunc() {
	for i := 1; i < 4; i++ {
		fmt.Printf("process func %d\n", i)
	}
}

func f() {
	fmt.Println("f start...")
	time.Sleep(time.Second * 2)
	fmt.Println("f done")
}
