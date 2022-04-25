package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"sync"
)

func main() {
	c()
}

func c() {
	data := []int{1, 2, 3, 4, 5}
	//data := []int{}

	data = nil

	fmt.Println("len(data) =", len(data))

	for i, d := range data {
		fmt.Printf("i = %d, d = %d \n", i, d)
	}

	if data == nil {
		fmt.Println("data = nil")
		data = append(data, 100)
	}

	fmt.Println("data =", data)
	fmt.Println("Done")
}

func a() {
	d, errCode := code(1234567890)
	if errCode != nil {
		fmt.Println("error code:", errCode.Error())
		return
	}
	fmt.Println("d =", d)

	res, errEncode := encode(d)
	if errEncode != nil {
		fmt.Println("error encode:", errEncode.Error())
		return
	}
	fmt.Println("res =", res)
}

func convertIntInByte(d int) []byte {
	return []byte(strconv.Itoa(d))
}

func code(d int) ([]byte, error) {
	buf := new(bytes.Buffer)
	//err := binary.Write(buf, binary.LittleEndian, float64(d))
	err := binary.Write(buf, binary.BigEndian, uint32(d))
	if err != nil {
		return nil, fmt.Errorf("binary.Write failed: %w", err)
	}
	return buf.Bytes(), nil
}

func encode(d []byte) (int, error) {
	//var res float64
	var res uint32
	b := bytes.NewReader(d)
	//err := binary.Read(b, binary.LittleEndian, &res)
	err := binary.Read(b, binary.BigEndian, &res)
	if err != nil {
		return 0, fmt.Errorf("binary.Read failed: %w", err)
	}
	return int(res), nil
}

func b() {
	var pi float64
	b := []byte{0x18, 0x2d, 0x44, 0x54, 0xfb, 0x21, 0x09, 0x40}
	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.LittleEndian, &pi)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}
	fmt.Print(pi)
}

func mapa() {
	m := sync.Map{}
	m.Store("key1", 1)
	m.Store("key2", 2)

	d, ok := m.Load("key")
	if !ok {
		d = 0
	}
	fmt.Println("d =", d.(int))

	m.Range(f)

}
func f(key interface{}, value interface{}) bool {
	fmt.Printf("%s:%d\n", key, value)
	return true
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
