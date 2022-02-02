package main

import (
	"fmt"
)

func main() {

	data := 65 // max 65 = 9 223 372 036 854 775 808

	n := uint64(data)
	fmt.Println("f_Naive(", n, ") = ", factorialNaive(n))
	fmt.Println("f_Tree (", n, ") = ", factorialTree(n))
	fmt.Println("f_Rec  (", n, ") = ", factorialRecursion(n))

	//u, err := user.Current()
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//fmt.Println(u.Username, u.HomeDir, u.Name, u.Uid, u.Gid)

}

func factorialNaive(n uint64) uint64 { // наивный алгоритм
	res := uint64(1)
	for i := uint64(2); i < n+1; i++ {
		res *= i
	}
	return res
}

func factorialRecursion(n uint64) uint64 { // алгоритм recursion
	if n == 0 || n == 1 {
		return 1
	}
	return n * factorialRecursion(n-1)
}

//------ алгоритм деревом ------------------
func factorialTree(n uint64) uint64 { // алгоритм деревом
	if n < 0 {
		return 0
	}
	if n == 0 {
		return 1
	}
	if n == 1 || n == 2 {
		return n
	}
	return prodTree(2, n)
}

func prodTree(l, r uint64) uint64 {
	if l > r {
		return 1
	}
	if l == r {
		return l
	}
	if (r - l) == 1 {
		return r * l
	}

	m := (l + r) / 2

	return prodTree(l, m) * prodTree(m+1, r)
}
