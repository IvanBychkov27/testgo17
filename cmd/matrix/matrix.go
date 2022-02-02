package main

import (
	"fmt"
	"gonum.org/v1/gonum/mat"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	a := creatRandomMatrix(2, 2, 10)
	fmt.Println("a:")
	fmt.Println(mat.Formatted(a, mat.Prefix(""), mat.Squeeze()))
	fmt.Println("sum =", sumElMatrix(a))
	fmt.Println()

	b := creatRandomMatrix(2, 2, 10)
	fmt.Println("b:")
	fmt.Println(mat.Formatted(b, mat.Prefix(""), mat.Squeeze()))
	fmt.Println("sum =", sumElMatrix(b))
	fmt.Println()

	//c := sumMatrix(a, b)
	//fmt.Println("c = a + b :")
	//fmt.Println(mat.Formatted(c, mat.Prefix(""), mat.Squeeze()))
	//fmt.Println("sum =", sumElMatrix(c))
	//fmt.Println()

	d := mulMatrix(a, b)
	fmt.Println("d = a * b :")
	fmt.Println(mat.Formatted(d, mat.Prefix(""), mat.Squeeze()))
	fmt.Println("sum =", sumElMatrix(d))

}

// mulMatrix умножение двух матриц
func mulMatrix(a, b *mat.Dense) *mat.Dense {
	var c mat.Dense
	c.Mul(a, b)
	return &c
}

// sumElMatrix сумма всех элементом матриц
func sumElMatrix(a *mat.Dense) float64 {
	return mat.Sum(a)
}

// sumMatrix сумма двух матриц
func sumMatrix(a, b *mat.Dense) *mat.Dense {
	var c mat.Dense
	c.Add(a, b)
	return &c
}

// creatRandomMatrix создает матрицу размера n на m с произвольными числами от 0 до random
func creatRandomMatrix(n, m, random int) *mat.Dense {
	size := n * m
	data := make([]float64, 0, size)

	for i := 0; i < size; i++ {
		d := rand.Intn(random)
		data = append(data, float64(d))
	}

	a := mat.NewDense(n, m, data)
	return a
}
