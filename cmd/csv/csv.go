package main

import (
	"encoding/csv"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
)

type posfData struct {
	isMobile        bool
	remateAddr      string
	userAgent       string
	platformName    string
	platformVersion string
	deviceVendor    string
	deviceModel     string
	stamp           *stampData
}

type stampData struct {
	TCPWindowSize   string
	IPTTL           string
	IPFlags         string
	TCPFlags        string
	TCPHeaderLength string
	TCPOptions      string
	MSS             string
}

func main() {
	fileName := "cmd/csv/posfdata20.txt" // posfdata20.txt - 9 mill
	//file, errOpen := os.Open(fileName)
	//if errOpen != nil {
	//	fmt.Println("error open file", errOpen.Error())
	//	return
	//}
	//defer file.Close()
	//
	//reader := csv.NewReader(file)
	//reader.FieldsPerRecord = 9
	//
	//data, errReaderRead := parseReaderDataCSV(reader)
	//if errReaderRead != nil {
	//	fmt.Println("error reader read", errReaderRead.Error())
	//}

	file, errOpen := ioutil.ReadFile(fileName)
	if errOpen != nil {
		fmt.Println("error open file", errOpen.Error())
		return
	}

	data, errParse := parseDataFileTXT(file)
	if errParse != nil {
		fmt.Println("error parse data file", errOpen.Error())
		return
	}

	fmt.Println("all data", len(data))
	fmt.Println("------------------------")

	statData := stat(data)
	printStat(statData)

	statStampData(statData)

	color.Red("Done...")
}

func statStampData(statData map[string][]*stampData) {
	windowSize := make(map[string]map[string]float64)   // map[keyWindowSize]map[plName]вероятность -  - вес 1 будем умножать на float64
	headerLength := make(map[string]map[string]float64) // map[keyHeaderLength]map[plName]вероятность  - вес 2 будем умножать на float64
	options := make(map[string]map[string]float64)      // map[keyOptions]map[plName]вероятность - вес 1 будем умножать на float64 (и далее все складываем - получим итоговую вероятность)

	ipttlData := make(map[string]map[string]float64) // map[keyIPTTL]map[plName]вероятность -  - вес 1 будем умножать на float64
	mssData := make(map[string]map[string]float64)   // map[keyMSS]map[plName]вероятность -  - вес 1 будем умножать на float64

	for plName, _ := range statData {
		stamps, ok := statData[plName]
		if !ok {
			fmt.Println("no data", plName)
			return
		}
		lenStamp := len(stamps)
		if lenStamp < 90 {
			continue
		}

		tcpWindowSize := make(map[string]int)
		ipTTL := make(map[string]int)
		ipFlags := make(map[string]int)
		tcpFlags := make(map[string]int)
		tcpHeaderLength := make(map[string]int)
		tcpOptions := make(map[string]int)
		mss := make(map[string]int)

		for _, s := range stamps {
			tcpWindowSize[s.TCPWindowSize]++
			ipTTL[s.IPTTL]++
			ipFlags[s.IPFlags]++
			tcpFlags[s.TCPFlags]++
			tcpHeaderLength[s.TCPHeaderLength]++
			tcpOptions[s.TCPOptions]++
			mss[s.MSS]++
		}
		// windowSize
		interestCalculation(plName, lenStamp, tcpWindowSize, windowSize)
		// headerLength
		interestCalculation(plName, lenStamp, tcpHeaderLength, headerLength)
		// options
		interestCalculation(plName, lenStamp, tcpOptions, options)
		// ipTTL
		interestCalculation(plName, lenStamp, ipTTL, ipttlData)
		// mss
		interestCalculation(plName, lenStamp, mss, mssData)

		fmt.Print("-")
	}
	fmt.Println()
	buildReferenceFile("cmd/csv/datafiles/windowSize.txt", windowSize)
	buildReferenceFile("cmd/csv/datafiles/headerLength.txt", headerLength)
	buildReferenceFile("cmd/csv/datafiles/options.txt", options)
	buildReferenceFile("cmd/csv/datafiles/ipttl.txt", ipttlData)
	buildReferenceFile("cmd/csv/datafiles/mss.txt", mssData)
}

// расчет процентов
func interestCalculation(plName string, lenStamp int, data map[string]int, resultData map[string]map[string]float64) {
	for key, val := range data {
		d, okWindowSize := resultData[key]
		if !okWindowSize {
			d = make(map[string]float64)
			resultData[key] = d
		}
		d[plName] = float64(val) / float64(lenStamp)
	}
}

// собрать и записать файл опорные данные
func buildReferenceFile(fileName string, data map[string]map[string]float64) {
	dataStr := ""
	for keyData, plScore := range data {
		if keyData == "Unknown" {
			continue
		}
		dataStr += fmt.Sprintf("%s;", keyData)
		for pl, score := range plScore {
			if pl == "Unknown" || score > 0.99 {
				continue
			}
			dataStr += fmt.Sprintf("%s,%0.12f;", pl, score)
		}
		dataStr += "\n"
	}
	saveDataFiles(fileName, dataStr)
}

func saveDataFiles(fileName, data string) {
	df, errCreateFile := os.Create(fileName)
	if errCreateFile != nil {
		fmt.Errorf("error create file, %w", errCreateFile)
		return
	}
	defer df.Close()

	_, errWrite := df.Write([]byte(data))
	if errWrite != nil {
		fmt.Errorf("error write data %w", errWrite)
		return
	}
}

func infoKey(data map[string]int, score map[string]map[string]float64) string {
	type dataKey struct {
		key string
		val int
	}

	ds := make([]*dataKey, 0, len(data))
	for key, val := range data {
		if strings.TrimSpace(key) == "" {
			continue
		}
		ds = append(ds, &dataKey{key, val})
	}

	sort.SliceStable(ds, func(i, j int) bool {
		return ds[i].val > ds[j].val // сортировка по убыванию рейтинга
	})

	res := ""
	for _, d := range ds {
		s := ""
		if len(score) > 0 {
			sc := score[d.key]
			for pl, ver := range sc {
				s += pl + ": " + strconv.FormatFloat(ver, 'f', -1, 64) + "; "
			}
		}
		res += fmt.Sprintf("----%6s  %d  %s\n", d.key, d.val, s)
	}
	return res
}

func infoKeyOptions(data map[string]int) string {
	type dataKey struct {
		key string
		val int
	}

	ds := make([]*dataKey, 0, len(data))
	for key, val := range data {
		if strings.TrimSpace(key) == "" {
			continue
		}
		ds = append(ds, &dataKey{key, val})
	}

	sort.SliceStable(ds, func(i, j int) bool {
		return ds[i].val > ds[j].val // сортировка по убыванию рейтинга
	})

	res := ""
	for _, d := range ds {
		//color.Yellow("----%20s  %d\n", d.key, d.val)
		res += fmt.Sprintf("----%20s  %d\n", d.key, d.val)
	}
	return res
}

func printStat(statData map[string][]*stampData) {
	type dataKey struct {
		key string
		val int
	}

	ds := make([]*dataKey, 0, len(statData))
	for key, val := range statData {
		if strings.TrimSpace(key) == "" {
			continue
		}
		ds = append(ds, &dataKey{key, len(val)})
	}

	sort.SliceStable(ds, func(i, j int) bool {
		return ds[i].val > ds[j].val // сортировка по убыванию рейтинга
	})

	for _, d := range ds {
		fmt.Printf("%15s = %d\n", d.key, d.val)
	}
}

func stat(data []*posfData) map[string][]*stampData {
	platform := make(map[string][]*stampData)
	for _, d := range data {
		if _, ok := platform[d.platformName]; !ok {
			ds := make([]*stampData, 0, len(data))
			platform[d.platformName] = ds
		}
		platform[d.platformName] = append(platform[d.platformName], d.stamp)
	}
	return platform
}

func parseStamp(stamp string) *stampData {
	d := &stampData{}
	s := strings.Split(stamp, ";")
	if len(s) != 6 {
		return nil
	}

	d.TCPWindowSize = s[0]
	d.IPTTL = s[1]
	d.IPFlags = s[2]
	d.TCPFlags = s[3]
	d.TCPHeaderLength = s[4]
	d.TCPOptions = s[5]
	d.MSS = getMSS(s[5])

	return d
}

func getMSS(dTCPOptions string) string {
	ms := strings.Split(dTCPOptions, ",")
	if len(ms) < 1 {
		return ""
	}
	mss := ms[0]
	if mss == "" {
		return ""
	}
	return mss[1:]
}

func parseDataFileTXT(file []byte) ([]*posfData, error) {
	dF := strings.Split(string(file), "\n")
	res := make([]*posfData, 0, len(dF))

	flag := false
	if len(dF) > 1000000 {
		flag = true
	}
	count := 0
	for _, s := range dF {
		if flag {
			count++
			if count > 300000 {
				count = 0
				fmt.Print("-")
			}
		}

		ds := strings.Split(s, `";"`)

		if len(ds) != 9 {
			continue
		}

		pd := &posfData{}

		if strings.Trim(ds[5], " \"\n\t") == "true" {
			pd.isMobile = true
		}

		pd.remateAddr = strings.Trim(ds[1], " \"\n\t")
		pd.userAgent = strings.Trim(ds[2], " \"\n\t")
		pd.platformName = strings.Trim(ds[3], " \"\n\t")
		pd.platformVersion = strings.Trim(ds[4], " \"\n\t")
		pd.deviceVendor = strings.Trim(ds[6], " \"\n\t")
		pd.deviceModel = strings.Trim(ds[7], " \"\n\t")
		pd.stamp = parseStamp(strings.Trim(ds[8], " \"\n\t"))

		res = append(res, pd)
	}
	fmt.Println()
	return res, nil
}

func parseReaderDataCSV(reader *csv.Reader) ([]*posfData, error) {
	data := make([]*posfData, 0, 1000000)
	for {
		record, errRead := reader.Read()
		if errRead != nil {
			if errRead.Error() == "EOF" {
				break
			}
			return nil, errRead
		}
		d := savePosfData(record)
		if d.stamp == nil {
			continue
		}
		data = append(data, d)
	}
	return data, nil
}

func savePosfData(rec []string) *posfData {
	d := &posfData{}
	if rec[5] == "true" {
		d.isMobile = true
	}
	d.remateAddr = rec[1]
	d.userAgent = rec[2]
	d.platformName = rec[3]
	d.platformVersion = rec[4]
	d.deviceVendor = rec[6]
	d.deviceModel = rec[7]
	d.stamp = parseStamp(rec[8])
	return d
}
