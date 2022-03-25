package main

import (
	"fmt"
)

func main() {
	//rateValue()
	a()
}

func a() {
	d := 5

	for d > 0 {
		d--
		fmt.Println(d)
	}

}

type data struct {
	ctr  float64
	rate int
}

func rateValue() {
	ctr := 0.25
	price := 100
	ds := []data{
		{ctr: 0.1, rate: 90},
		{ctr: 0.2, rate: 80},
		{ctr: 0.3, rate: 70},
	}

	rate := 100
	for _, d := range ds {
		if d.ctr > ctr {
			break
		}
		rate = d.rate
	}

	newPrice := int(float64(price) * float64(rate) / float64(100))
	deltaPrice := price - newPrice

	fmt.Println("rate =", rate)
	fmt.Println("price      =", price)
	fmt.Println("newPrice   =", newPrice)
	fmt.Println("deltaPrice =", deltaPrice)

}
