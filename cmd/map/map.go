package main

import "fmt"

func main() {

	m := map[int]int{1: 1, 2: 0, 3: 3}
	fmt.Println("start:", m)

	for key, val := range m {
		if val == 0 {
			delete(m, key)
		}
	}

	fmt.Println("stop:", m)
}
