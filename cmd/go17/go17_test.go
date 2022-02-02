package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_factorial(t *testing.T) {
	n := uint64(5)
	expected := uint64(120)

	res := factorialNaive(n)
	assert.Equal(t, expected, res)

	res = factorialTree(n)
	assert.Equal(t, expected, res)

	res = factorialRecursion(n)
	assert.Equal(t, expected, res)
}

func Benchmark_factorialNaive(b *testing.B) {
	n := uint64(65)
	for i := 0; i < b.N; i++ {
		factorialNaive(n)
	}
}

func Benchmark_factorialTree(b *testing.B) {
	n := uint64(65)
	for i := 0; i < b.N; i++ {
		factorialTree(n)
	}
}

func Benchmark_factorialRecursion(b *testing.B) {
	n := uint64(65)
	for i := 0; i < b.N; i++ {
		factorialRecursion(n)
	}
}
