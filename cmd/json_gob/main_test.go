// https://habr.com/ru/company/kaspersky/blog/591725/

package main

import "testing"

func Benchmark_jsonMarshal(b *testing.B) {
	d := []D{
		{1, "a"},
		{2, "b"},
		{3, "c"},
		{4, "d"},
		{5, "e"},
		{6, "f"},
	}

	for i := 0; i < b.N; i++ {
		jsonMarshal(d)
	}
}

func Benchmark_jsonUnmarshal(b *testing.B) {
	j := []byte(`
[
  {"id":1,"name":"a"},
  {"id":2,"name":"b"},
  {"id":3,"name":"c"},
  {"id":4,"name":"d"},
  {"id":5,"name":"e"},
  {"id":6,"name":"f"}
]
`)

	for i := 0; i < b.N; i++ {
		jsonUnmarshal(j)
	}
}

func Benchmark_gobMarshal(b *testing.B) {
	d := []D{
		{1, "a"},
		{2, "b"},
		{3, "c"},
		{4, "d"},
		{5, "e"},
		{6, "f"},
	}

	for i := 0; i < b.N; i++ {
		gobMarshal(d)
	}
}

func Benchmark_gobUnmarshal(b *testing.B) {
	j := []byte{13, 255, 131, 2, 1, 2, 255, 132, 0, 1, 255, 130, 0, 0, 31, 255, 129, 3, 1, 1, 1, 68, 1, 255, 130, 0, 1, 2, 1, 2, 73, 68, 1, 4, 0, 1, 4, 78, 97, 109, 101, 1, 12, 0, 0, 0, 40, 255, 132, 0, 6, 1, 2, 1, 1, 97, 0, 1, 4, 1, 1, 98, 0, 1, 6, 1, 1, 99, 0, 1, 8, 1, 1, 100, 0, 1, 10, 1, 1, 101, 0, 1, 12, 1, 1, 102, 0}
	for i := 0; i < b.N; i++ {
		gobUnmarshal(j)
	}
}
