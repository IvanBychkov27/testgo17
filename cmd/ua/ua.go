// программа для получения списка UserAgent из файла
/*
     Для проекта UserAgent
	 - читает файл и выбирает все UA
	 - сортирует все UA
     - сохраняет в файл без повторов
*/

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"time"
)

func main() {
	dataFile := getDataFromFile("cmd/ua/raw/ua20211129.txt")

	uas := parseData(dataFile)

	fmt.Println("len data =", len(uas))

	data := strings.Join(uas, "\n")
	fileName := "cmd/ua/data/ua_" + time.Now().Format("20060102") + ".txt"
	saveFile(fileName, data)

}

func parseData(data []byte) []string {
	uas := make([]string, 0)

	temp := bytes.Split(data, []byte("\n"))
	for _, ua := range temp {
		if len(ua) == 0 || ua[0] != '{' {
			continue
		}
		a := bytes.Index(ua, []byte(":"))
		b := bytes.Index(ua, []byte("\","))
		if a == -1 || b == -1 || b <= a {
			continue
		}

		uas = append(uas, string(ua[a+2:b]))
	}

	sort.Strings(uas)

	res := make([]string, 0, len(uas))
	d := ""
	for _, ua := range uas {
		if d == ua {
			continue
		}
		res = append(res, ua)
		d = ua
	}

	return res
}

func getDataFromFile(fileName string) []byte {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("error", err)
	}
	return data
}

func saveFile(fileName, data string) {
	df, errCreateFile := os.Create(fileName)
	if errCreateFile != nil {
		fmt.Println("error:", errCreateFile.Error())
		return
	}
	defer df.Close()

	_, errWrite := df.WriteString(data)
	if errWrite != nil {
		fmt.Println("error:", errWrite.Error())
		return
	}
	fmt.Println("saved file:", fileName)
	return
}
