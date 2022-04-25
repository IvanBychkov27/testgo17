package main

import (
	"fmt"
	"strconv"
)

func main() {
	n := 16777215
	fmt.Println("n =", n)

	//b := intConvByte(n)
	//fmt.Println("b =", b)

	b3 := intConvByte3(n)
	fmt.Println("b3=", b3)

	i, err := byteConvInt3(b3)
	if err != nil {
		fmt.Println(err.Error(), " b3 =", b3)
	}
	fmt.Println("i =", i)

}

// intConvByte - конвертирует int в []byte, где n <= 65535. Результат 2 байта
func intConvByte(n int) []byte {
	const d = 256
	return []byte{uint8(n / d), uint8(n % d)}
}

// byteConvInt - конвертирует []byte в int, где len(b)=2
func byteConvInt(b []byte) (int, error) {
	if len(b) > 2 {
		return 0, fmt.Errorf("error: len(data) > 2")
	}
	const d = 256
	res := int(b[0])*d + int(b[1])
	return res, nil
}

func convertIntInByteByStr(d int) []byte {
	return []byte(strconv.Itoa(d))
}

// intConvByte - конвертирует int в []byte, где n <= 16 777 215. Результат 3 байта
func intConvByte3(n int) []byte {
	const d = 256
	const d3 = 65536

	if n <= d3 {
		return []byte{uint8(0), uint8(n / d), uint8(n % d)}
	}
	return []byte{uint8(n / d3), uint8((n % d3) / d), uint8((n % d3) % d)}
}

// byteConvInt - конвертирует []byte в int, где len(b)=3
func byteConvInt3(b []byte) (int, error) {
	if len(b) > 3 {
		return 0, fmt.Errorf("error: len(data) > 3")
	}
	const d = 256
	const d3 = 65536

	res := int(b[0])*d3 + int(b[1])*d + int(b[2])
	return res, nil
}
