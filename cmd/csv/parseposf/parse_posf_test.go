package main

import (
	"fmt"
	"testing"
)

func Test_setOptions(t *testing.T) {
	stamp := "65495;64;DF;2;60;M65495,S,N,W7" // local port
	res := setOptions(stamp)
	fmt.Println(res)
}
