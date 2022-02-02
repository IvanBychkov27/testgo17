package main

import (
	"fmt"
	"testing"
)

func Test_getIP(t *testing.T) {
	ipPort := "84.54.94.168:18614"
	ip := getIP(ipPort)
	fmt.Println("ip =", ip)
}

func Test_getMSS(t *testing.T) {
	stamp := "65535;49;DF;2;60;M1436,S,N,W8"
	mss := getMSS(stamp)
	fmt.Println("mss =", mss)
}
