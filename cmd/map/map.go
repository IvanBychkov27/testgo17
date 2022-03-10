package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	//del()

	jsonMap()
}

func jsonMap() {
	data := []byte(`
{ 
  	"a":1,
	"b":2,
	"c":3
}`)

	res := map[string]int{}
	err := json.Unmarshal(data, &res)
	if err != nil {
		fmt.Println("error unmarshal", err.Error())
	}

	fmt.Println("res:", res)
	fmt.Println("a =", res["a"])
	fmt.Println("b =", res["b"])
	fmt.Println("c =", res["c"])

	fmt.Printf("1e9 =%d", int(1e9))

}

func del() {
	m := map[int]int{1: 1, 2: 0, 3: 3}
	fmt.Println("start:", m)

	for key, val := range m {
		if val == 0 {
			delete(m, key)
		}
	}

	fmt.Println("stop:", m)
}
