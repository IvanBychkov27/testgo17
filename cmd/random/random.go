package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	a := []string{"a", "b", "c"}

	idx := random(len(a))

	fmt.Println("res =", a[idx])
}

func random(n int) int {
	return rand.Intn(n)
}
