package main

import "testing"

func Benchmark_convertIntInByte(b *testing.B) {
	for i := 0; i < b.N; i++ {
		convertIntInByte(1234567890)
	}
}

func Benchmark_code(b *testing.B) {
	for i := 0; i < b.N; i++ {
		code(1234567890)
	}
}

func Benchmark_encode(b *testing.B) {
	d := []byte{73, 150, 2, 210}
	for i := 0; i < b.N; i++ {
		encode(d)
	}
}
