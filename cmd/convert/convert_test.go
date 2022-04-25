package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_intConvByte(t *testing.T) {
	res := intConvByte(259)
	assert.Equal(t, []byte{0x1, 0x3}, res)
}

func Benchmark_intConvByte(b *testing.B) {
	for i := 0; i < b.N; i++ {
		intConvByte(65535)
	}
}

func Benchmark_convertIntInByte(b *testing.B) {
	for i := 0; i < b.N; i++ {
		convertIntInByteByStr(65535)
	}
}

func Test_byteConvInt(t *testing.T) {
	res, _ := byteConvInt([]byte{0x1, 0x3})
	assert.Equal(t, 259, res)
}

func Benchmark_byteConvInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		byteConvInt([]byte{255, 255})
	}
}

func Test_intConvByte3(t *testing.T) {
	res := intConvByte3(16777215)
	assert.Equal(t, []byte{255, 255, 255}, res)
}

func Test_byteConvInt3(t *testing.T) {
	res, _ := byteConvInt3([]byte{255, 255, 255})
	assert.Equal(t, 16777215, res)
}

func Benchmark_intConvByte3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		intConvByte3(16777215)
	}
}

func Benchmark_byteConvInt3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		byteConvInt3([]byte{255, 255, 255})
	}
}
