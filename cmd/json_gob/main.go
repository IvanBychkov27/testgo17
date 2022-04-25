package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/pkg/profile"
)

type D struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	defer profile.Start(profile.MemProfile, profile.ProfilePath("./cmd/json_gob")).Stop()
	// go tool pprof mem.pprof
	// top

	data := []D{
		{1, "a"},
		{2, "b"},
		{3, "c"},
		{4, "d"},
		{5, "e"},
		{6, "f"},
	}
	resJsonMarshal, err := jsonMarshal(data)
	if err != nil {
		fmt.Println("error:", err.Error())
	}
	fmt.Println("resJsonMarshal   =", string(resJsonMarshal))

	resJsonUnmarshal, errUnmarshal := jsonUnmarshal(resJsonMarshal)
	if errUnmarshal != nil {
		fmt.Println("error:", errUnmarshal.Error())
	}
	fmt.Println("resJsonUnmarshal =", resJsonUnmarshal)

	resGobMarshal, errGobMarshal := gobMarshal(data)
	if err != nil {
		fmt.Println("error:", errGobMarshal.Error())
	}
	fmt.Println("resGobMarshal   =", resGobMarshal)
	//s := ""
	//for _, d := range resGobMarshal {
	//	s += fmt.Sprintf("%d,", d)
	//}
	//fmt.Println("resGobMarshal   =", s)
	//s = 13,255,131,2,1,2,255,132,0,1,255,130,0,0,31,255,129,3,1,1,1,68,1,255,130,0,1,2,1,2,73,68,1,4,0,1,4,78,97,109,101,1,12,0,0,0,40,255,132,0,6,1,2,1,1,97,0,1,4,1,1,98,0,1,6,1,1,99,0,1,8,1,1,100,0,1,10,1,1,101,0,1,12,1,1,102,0

	resGobUnmarshal, errGobUnmarshal := gobUnmarshal(resGobMarshal)
	if errGobUnmarshal != nil {
		fmt.Println("error:", errGobUnmarshal.Error())
	}
	fmt.Println("resGobUnmarshal  =", resGobUnmarshal)
}

func jsonMarshal(data []D) ([]byte, error) {
	return json.Marshal(data)
}

func jsonUnmarshal(data []byte) ([]D, error) {
	res := make([]D, 0)
	err := json.Unmarshal(data, &res)
	return res, err
}

func gobMarshal(data []D) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	c := gob.NewEncoder(buf)
	err := c.Encode(data)
	return buf.Bytes(), err
}

func gobUnmarshal(data []byte) ([]D, error) {
	res := make([]D, 0)
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&res)
	return res, err
}
