package main

import (
	"encoding/base64"
	"fmt"
)

func main() {

	//m := make(map[string]string)

	password := []byte("qwerty")
	res := base64.StdEncoding.EncodeToString(password)

	fmt.Println(res)
}
