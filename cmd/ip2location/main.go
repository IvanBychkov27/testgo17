package main

import (
	"bytes"
	"fmt"
	"github.com/ip2location/ip2location-go/v9"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	err := OpenIP2location("/home/ivan/data/IP2LOCATION.BIN_V4")
	if err != nil {
		fmt.Println("error open file ip2location", err.Error())
		return
	}
	defer Close()

	fNamePosfData := "cmd/ip2location/data/posfdata21ok.txt"
	dataPosf := getDataPosf(fNamePosfData)

	fmt.Println("len data posf", len(dataPosf))

	//res := compareIP_MSS(dataPosf)
	//fmt.Println("len connection type", len(res))

	//res := comparePlName_MSS(dataPosf)
	res := comparePlName_IPTTL(dataPosf)
	//res := comparePlName_TCPFlags(dataPosf)

	fmt.Println("len platform name", len(res))

	printPosfData(res)
}

func compareIP_MSS(dataPosf []posfData) map[string]map[string]int {
	var err error
	data := New()
	res := make(map[string]map[string]int) // map[ConnectionType]map[MSS]count
	for _, d := range dataPosf {
		err = data.Parse(d.ip)
		if err != nil {
			fmt.Println("error data parse ip", err.Error())
			return nil
		}

		m, ok := res[data.ConnectionType]
		if !ok {
			m = make(map[string]int)
			res[data.ConnectionType] = m
		}
		m[d.mss]++
	}
	return res
}

func comparePlName_MSS(dataPosf []posfData) map[string]map[string]int {
	res := make(map[string]map[string]int) // map[plName]map[MSS]count
	for _, d := range dataPosf {
		m, ok := res[d.plName]
		if !ok {
			m = make(map[string]int)
			res[d.plName] = m
		}
		m[d.mss]++
	}
	return res
}

func comparePlName_IPTTL(dataPosf []posfData) map[string]map[string]int {
	res := make(map[string]map[string]int) // map[plName]map[IPTTL]count
	for _, d := range dataPosf {
		m, ok := res[d.plName]
		if !ok {
			m = make(map[string]int)
			res[d.plName] = m
		}
		m[d.ipTTL]++
	}
	return res
}

func comparePlName_TCPFlags(dataPosf []posfData) map[string]map[string]int {
	res := make(map[string]map[string]int) // map[plName]map[TCPFlags]count
	for _, d := range dataPosf {
		m, ok := res[d.plName]
		if !ok {
			m = make(map[string]int)
			res[d.plName] = m
		}
		m[d.tcpFlags]++
	}
	return res
}

func printPosfData(data map[string]map[string]int) {
	//fmt.Println("conn type: --- mss: count")
	for connType, mss := range data {
		fmt.Println(connType)
		printMap(mss)
	}
}

func printMap(data map[string]int) {
	type dT struct {
		m string // mss
		c int    // count
	}
	ds := make([]dT, 0, len(data))
	for m, c := range data {
		d := dT{m, c}
		ds = append(ds, d)
	}
	sort.SliceStable(ds, func(i, j int) bool {
		return ds[i].c > ds[j].c // сортировка по убыванию
	})

	count := 0
	for _, v := range ds {
		count += v.c
	}

	for _, val := range ds {
		fmt.Printf("--- %4s: %d  %0.6f \n", val.m, val.c, float64(val.c)/float64(count))
	}
}

func printIP2locationAll(data *IP2LocationData) {
	fmt.Println("Country: ", data.Country)
	fmt.Println("Region : ", data.Region)
	fmt.Println("City   :", data.City)
	fmt.Println("Isp    :", data.Isp)
	fmt.Println("ConnectionType:", data.ConnectionType)
	fmt.Println("MobileBrand :", data.MobileBrand)
	fmt.Println("Timezone    :", data.Timezone)
	fmt.Println("TZHours     :", data.TZHours)
	fmt.Println("TZMinute    :", data.TZMinute)
	fmt.Println("TimezoneText:", data.TimezoneText)
}

//---------------------------------------------------

var (
	defaultDB *ip2location.DB
)

func OpenIP2location(dbpath string) error {
	fileData, err := os.ReadFile(dbpath)
	if err != nil {
		return err
	}

	r := &dbReader{buf: bytes.NewReader(fileData)}
	defaultDB, err = ip2location.OpenDBWithReader(r)
	if err != nil {
		return err
	}

	return nil
}

func Close() {
	defaultDB.Close()
}

type dbReader struct {
	buf *bytes.Reader
}

func (b *dbReader) ReadAt(p []byte, off int64) (n int, err error) {
	return b.buf.ReadAt(p, off)
}

func (b *dbReader) Read(p []byte) (n int, err error) {
	return b.buf.Read(p)
}

func (b *dbReader) Close() error {
	return nil
}

var (
	tzText = map[string]string{
		"-12:00": "Pacific/Pago_Pago", // не нашел верного названия
		"-11:00": "Pacific/Pago_Pago",
		"-10:00": "America/Adak",
		"-09:30": "Pacific/Marquesas",
		"-09:00": "America/Anchorage",
		"-08:00": "America/Los_Angeles",
		"-07:00": "America/Denver",
		"-06:00": "America/Chicago",
		"-05:00": "America/Detroit",
		"-04:00": "America/Halifax",
		"-03:30": "America/St_Johns",
		"-03:00": "America/Araguaina",
		"-02:00": "America/Noronha",
		"-01:00": "Atlantic/Azores",
		"-00:00": "UTC",
		"00:00":  "UTC",
		"+00:00": "UTC",
		"+01:00": "Europe/Berlin",
		"+02:00": "Europe/Kaliningrad",
		"+03:00": "Europe/Moscow",
		"+04:00": "Asia/Baku",
		"+05:00": "Asia/Aqtau",
		"+05:30": "Asia/Colombo",
		"+05:45": "Asia/Kathmandu",
		"+06:00": "Asia/Bishkek",
		"+06:30": "Asia/Rangoon",
		"+07:00": "Asia/Bangkok",
		"+08:00": "Asia/Harbin",
		"+08:45": "Australia/Eucla",
		"+09:00": "Asia/Dili",
		"+09:30": "Australia/Adelaide",
		"+10:00": "Asia/Vladivostok",
		"+10:30": "Australia/Lord_Howe",
		"+11:00": "Asia/Magadan",
		"+12:00": "Asia/Kamchatka",
		"+12:45": "Pacific/Chatham",
		"+13:00": "Pacific/Apia",
		"+14:00": "Pacific/Kiritimati",
	}
)

type IP2LocationData struct {
	Country        string
	Region         string
	City           string
	Isp            string
	ConnectionType string
	MobileBrand    string
	Timezone       string
	TZHours        int
	TZMinute       int
	TimezoneText   string
}

func New() *IP2LocationData {
	return &IP2LocationData{}
}

func (data *IP2LocationData) Parse(ip string) error {
	if defaultDB == nil {
		return fmt.Errorf("not inited")
	}
	info, err := defaultDB.Get_all(ip)
	if err != nil {
		return err
	}

	data.Country = info.Country_short
	data.Region = info.Region
	data.City = info.City
	data.Isp = info.Isp
	data.ConnectionType = info.Usagetype
	data.MobileBrand = info.Mobilebrand
	data.Timezone = info.Timezone
	data.TZHours, data.TZMinute = data.timezoneToOffset(data.Timezone)
	data.TimezoneText = data.timezoneToText(data.Timezone)

	return nil
}

func (data *IP2LocationData) Reset() {
	data.Country = ""
	data.Region = ""
	data.City = ""
	data.Isp = ""
	data.ConnectionType = ""
	data.MobileBrand = ""
	data.Timezone = ""
	data.TZHours = 0
	data.TZMinute = 0
	data.TimezoneText = ""
}

func (data *IP2LocationData) timezoneToOffset(tz string) (hour, minutes int) {
	hour, minutes = 0, 0

	pair := strings.Split(tz, ":")
	if len(pair) != 2 {
		return 0, 0
	}

	var err error

	minutes, err = strconv.Atoi(pair[1])
	if err != nil {
		return 0, 0
	}

	// Час должен быть со знаком +01, -02 и тд
	if len(pair[0]) < 2 {
		return 0, 0
	}

	hour, err = strconv.Atoi(pair[0][1:])
	if err != nil {
		return 0, 0
	}

	if pair[0][0] == '-' {
		hour = -hour
		minutes = -minutes
	}

	return hour, minutes
}

func (data *IP2LocationData) timezoneToText(tz string) string {
	text, ok := tzText[tz]
	if !ok {
		return "UTC"
	}

	return text
}

//-----------------------------------------------

type posfData struct {
	ip       string
	plName   string
	mss      string
	ipTTL    string
	tcpFlags string
}

func getDataPosf(fileName string) []posfData {
	file, errOpen := ioutil.ReadFile(fileName)
	if errOpen != nil {
		fmt.Println("error open file stamp data", errOpen.Error())
		return nil
	}

	dataFile := strings.Split(string(file), "\n")
	//fmt.Println("data posf", len(dataFile))

	data := make([]posfData, 0, 100000)
	for _, lineFile := range dataFile {
		lf := strings.Split(lineFile, `";"`)
		if len(lf) != 9 {
			continue
		}
		ipPort := strings.Trim(lf[1], " \"\n\t")
		plName := strings.Trim(lf[3], " \"\n\t")
		stamp := strings.Trim(lf[8], " \"\n\t")

		if ipPort == "" || stamp == "" || plName == "" {
			continue
		}

		ip := getIP(ipPort)
		mss, ttl, tcpFl := getMSS_IPTTL_TCPFlags(stamp)

		data = append(data, posfData{ip: ip, plName: plName, mss: mss, ipTTL: ttl, tcpFlags: tcpFl})
	}

	return data
}

func getIP(ds string) string {
	d := strings.Split(ds, ":")
	if len(d) == 2 {
		return d[0]
	}
	return ""
}

func getMSS(ds string) string {
	d := strings.Split(ds, ";")
	if len(d) != 6 {
		return ""
	}
	m := d[5]
	ms := strings.Split(m, ",")
	if len(ms) < 1 {
		return ""
	}
	mss := ms[0]
	return mss[1:]
}

// stamp = 65535;49;DF;2;60;M1436,S,N,W8 - (0 - TCPWindowSize; 1 - IPTTL; 2 - IPFlags; 3 - TCPFlags; 4 - TCPHeaderLength; 5 - TCPOptions)
func getMSS_IPTTL_TCPFlags(ds string) (MSS, IPTTL, TCPFlags string) {
	d := strings.Split(ds, ";")
	if len(d) != 6 {
		return
	}
	m := d[5]
	ms := strings.Split(m, ",")
	if len(ms) < 1 {
		return "", d[1], d[3]
	}
	mss := ms[0]

	return mss[1:], d[1], d[3]
}
