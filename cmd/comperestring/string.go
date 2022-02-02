package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
)

func main() {
	//str2 := []byte("a:598:{i:11677;d:0;i:7453;d:7766279631.452242;i:11694;d:0.005409632000002773;i:11785;d:0;i:12853;d:-0.75;i:11776;d:1573.4199719589997;}")

	// user
	forUser := "user" // для user
	fileNameRedis := "cmd/comperestring/data/user_redis01.txt"
	fileName2 := "cmd/comperestring/data/user_balance01.txt"

	//forUser := "" // для cash и camp
	// cash
	//fileNameRedis := "cmd/comperestring/data/cash_redis01.txt"
	//fileName2 := "cmd/comperestring/data/cash_balance01.txt"

	//camp
	//fileNameRedis := "cmd/comperestring/data/camp_redis01.txt"
	//fileName2 := "cmd/comperestring/data/camp_balance01.txt"

	strRedis := openFile(fileNameRedis)
	str2 := openFile(fileName2)

	resReis := parseStr(strRedis)
	fmt.Println("len redis  =", len(resReis))

	res2 := parseStr(str2)
	fmt.Println("len result =", len(res2))

	different, ok := compareBalanceID(resReis, res2, forUser)
	if !ok {
		fmt.Println("len different = ", len(different))
		for id, bal := range different {
			fmt.Printf("%6s: %s \n", id, bal)
		}
	}
	fmt.Println("compare = ", ok)
}

func openFile(fileName string) []byte {
	file, errOpen := ioutil.ReadFile(fileName)
	if errOpen != nil {
		fmt.Println("error open file", errOpen.Error())
		return nil
	}
	return file
}

func parseStr(s []byte) map[string]string {
	if !bytes.Contains(s, []byte("{")) {
		return nil
	}
	var id string
	res := make(map[string]string)

	idx := bytes.Index(s, []byte("{"))
	var n, k int
	for i := idx; i < len(s); i++ {
		if s[i] == 'i' || s[i] == 'd' {
			idx = i
			i += 2
			n = i
		}
		if s[i] == ';' {
			k = i
		}
		if n > 0 && k > 0 {
			d := s[n:k]
			n, k = 0, 0

			if s[idx] == 'i' {
				if id != "" {
					res[id] = string(d)
					id = ""
				} else {
					id = string(d)
				}
			}
			if s[idx] == 'd' {
				res[id] = string(d)
				id = ""
			}
		}
	}

	return res
}

// сравнение двух мап: standard - эталон с которым сравниваем и item - то что сравниваем с эталоном
// Результат мапа с отличиями different[id] = balItem + ", " + balanceStandard
// И если результат = true, то полное совпадение
func compareBalanceID(standard, item map[string]string, f string) (map[string]string, bool) {
	different := make(map[string]string)
	for id, bal := range item {
		balance, ok := standard[id]
		if !ok {
			different[id] = bal
			continue
		}
		if balance != bal {
			dStr := ""
			if f == "user" {
				dStr = difFloat(bal, balance)
			} else {
				dStr = difInt(bal, balance)
			}

			if dStr == "" {
				continue
			}
			different[id] = dStr
		}
	}
	res := false
	if len(different) == 0 {
		res = true
	}

	return different, res
}

func difFloat(bal, balance string) string {
	r, _ := strconv.ParseFloat(bal, 64)
	redis, _ := strconv.ParseFloat(balance, 64)
	dif := r - redis
	if math.Abs(dif) < float64(0.1) {
		return ""
	}
	return fmt.Sprintf("%22s, %22s, %22s", bal, balance, strconv.FormatFloat(dif, 'f', -1, 64))
}

func difInt(bal, balance string) string {
	r, _ := strconv.Atoi(bal)
	redis, _ := strconv.Atoi(balance)
	dif := r - redis
	if math.Abs(float64(dif)) < float64(100000000) {
		return ""
	}
	return fmt.Sprintf("%15s, %15s, %10s", bal, balance, strconv.Itoa(dif))
}
