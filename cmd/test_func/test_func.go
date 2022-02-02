package main

import (
	"fmt"
	"time"
)

func main() {

	process(timeFormat)

}

func timeFormat() (string, error) {
	return time.Now().Format("15:04:05"), nil
}

type getTimeFunc func() (string, error)

func process(t getTimeFunc) {
	d, _ := t()
	fmt.Println(d)
}
