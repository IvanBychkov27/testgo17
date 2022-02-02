package main

import (
	"encoding/base64"
	"fmt"
)

func main() {

	//m := make(map[string]string)

	res := base64.StdEncoding.EncodeToString(nil)

	fmt.Println(res)
}
