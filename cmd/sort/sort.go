package main

import (
	"fmt"
	"sort"
)

func main() {
	sortirovka()
}

func sortirovka() {
	a := []int{9, 1, 8, 2, 7, 3, 6, 4, 5}
	fmt.Println("orig a =", a)

	sortMas(a)
	fmt.Println("sort a =", a)
}

func sortMas(a []int) {
	sort.Ints(a)
}
