package main

import (
	"crypto/md5"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	//ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	//<-ctx.Done()

	var f float64

	f = 1.999

	resCeil := math.Ceil(f)
	resFloor := math.Floor(f)

	fmt.Println(resCeil)
	fmt.Println(resFloor)

	//fileName := "cmd/test/data/f01.txt"
	//renameFile(fileName)

	d := []int{1, 2, 3, 4}
	fmt.Println("d =", d)
	//d = nil

	for i, v := range d {
		fmt.Println(i, v)
	}

}

func renameFile(fileName string) {
	newFileName := fileName + "_" + time.Now().Format("2006-01-02_15-04-05")
	err := os.Rename(fileName, newFileName)
	if err != nil {
		fmt.Println("error:", err.Error())
	}
}

type TempStruct []struct {
	ServiceAddress string `json:"ServiceAddress"`
	ServicePort    int    `json:"ServicePort"`
}

type AutoStruct []struct {
	Address     string `json:"Address"`
	CreateIndex int    `json:"CreateIndex"`
	Datacenter  string `json:"Datacenter"`
	ID          string `json:"ID"`
	ModifyIndex int    `json:"ModifyIndex"`
	Node        string `json:"Node"`
	NodeMeta    struct {
		ConsulNetworkSegment string `json:"consul-network-segment"`
	} `json:"NodeMeta"`
	ServiceAddress string `json:"ServiceAddress"`
	ServiceConnect struct {
	} `json:"ServiceConnect"`
	ServiceEnableTagOverride bool   `json:"ServiceEnableTagOverride"`
	ServiceID                string `json:"ServiceID"`
	ServiceKind              string `json:"ServiceKind"`
	ServiceMeta              struct {
		ExternalSource string `json:"external-source"`
	} `json:"ServiceMeta"`
	ServiceName  string `json:"ServiceName"`
	ServicePort  int    `json:"ServicePort"`
	ServiceProxy struct {
		Expose struct {
		} `json:"Expose"`
		MeshGateway struct {
		} `json:"MeshGateway"`
	} `json:"ServiceProxy"`
	ServiceTaggedAddresses struct {
		LanIpv4 struct {
			Address string `json:"Address"`
			Port    int    `json:"Port"`
		} `json:"lan_ipv4"`
		WanIpv4 struct {
			Address string `json:"Address"`
			Port    int    `json:"Port"`
		} `json:"wan_ipv4"`
	} `json:"ServiceTaggedAddresses"`
	ServiceTags    []string `json:"ServiceTags"`
	ServiceWeights struct {
		Passing int `json:"Passing"`
		Warning int `json:"Warning"`
	} `json:"ServiceWeights"`
	TaggedAddresses struct {
		Lan     string `json:"lan"`
		LanIpv4 string `json:"lan_ipv4"`
		Wan     string `json:"wan"`
		WanIpv4 string `json:"wan_ipv4"`
	} `json:"TaggedAddresses"`
}

func convertOptionsInSliceInt(options string) string {
	ds := strings.Split(options, ",")
	res := make([]string, 20, 20)
	j := 0
	for i, d := range ds {
		data := []byte(d)
		sum := 0
		for _, v := range data {
			sum += int(v)
		}
		res[i] = fmt.Sprintf("0.%d", sum)
		j = i
	}
	for i := j + 1; i < len(res); i++ {
		res[i] = "0"
	}

	return strings.Join(res, ",")
}

// setCheckbox включает флажки на форме в соответствии с выбранными картами
func setCheckbox(form string) string {
	cardForm := "hand2p"
	if !strings.Contains(form, cardForm) {
		return ""
	}
	idx := strings.Index(form, cardForm) + len(cardForm) + 1
	form = form[:idx] + " checked" + form[idx:]

	return form
}

// timeNow функция выводит текущую дату и время, addHour позволяет прибавить/отнять часы (для корректировки)
func timeNow(addHour int) string {
	y := time.Now().Year()
	mec := time.Now().Month()
	d := time.Now().Day()
	h := time.Now().Hour() + addHour
	m := time.Now().Minute()
	s := time.Now().Second()
	return time.Date(y, mec, d, h, m, s, 0, time.UTC).Format("02.01.2006  15:04:05")
}

func md5Data(data string) string {
	h := md5.Sum([]byte(data))
	return fmt.Sprintf("%x", h)
}

func convertHexInDec(hex string) string {
	res := 0

	for i := 0; i < 22; i += 11 {
		d, err := ConvertInt(hex[i:i+11], 16, 10)
		if err != nil {
			fmt.Println(err.Error())
			return ""
		}
		r, _ := strconv.Atoi(d)
		res += r
	}

	return strconv.Itoa(res)
}

// ConvertInt конвертирует значение из одной системы счисления в другую, которая указана в toBase
// https://golangify.com/binary-to-decimal-octal-and-hexadecimal
func ConvertInt(val string, base, toBase int) (string, error) {
	i, err := strconv.ParseInt(val, base, 64)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(i, toBase), nil
}
