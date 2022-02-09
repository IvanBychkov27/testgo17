package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/fxsjy/gonn/gonn"
	"io/ioutil"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type dataScore struct {
	plName string
	score  float64
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

const nPlatform = 1100

func main() {
	//saveResultDataInFile()

	fileNameNN := "cmd/csv/parseposf/nn/gonn20000_119"
	//fileNameNN := "cmd/csv/parseposf/nn/gonn49505_300"

	//manualCalculation(fileNameNN)
	calculatingDataFromFile(fileNameNN)
	//sumMyAndNN(fileNameNN)
}

func manualCalculation(fileNameNN string) {
	fileNameHL := "cmd/csv/parseposf/referencefiles/headerLength.txt"
	fileNameWin := "cmd/csv/parseposf/referencefiles/windowSize.txt"
	fileNameOpt := "cmd/csv/parseposf/referencefiles/options.txt"
	fileNameIpTtl := "cmd/csv/parseposf/referencefiles/ipttl.txt"
	fileNameMss := "cmd/csv/parseposf/referencefiles/mss.txt"

	scoreHL := parseDataFile(fileNameHL, 1)
	scoreWin := parseDataFile(fileNameWin, 1)
	scoreOpt := parseDataFile(fileNameOpt, 1)
	scoreIpTtl := parseDataFile(fileNameIpTtl, 1)
	scoreMss := parseDataFile(fileNameMss, 1)

	stamp := "65535;49;DF;2;60;M1436,S,N,W8" // Android
	//stamp := "8192;121;DF;194;52;M1460,N,W8,N,N,S" // Android ??? - BOT!!! - это в реальности Windows
	//stamp := "65535;56;DF;2;60;M1380,S,N,N,N,N,N,N,N,N,N,N,N,W8" // Android
	//stamp := "64240;119;DF;2;52;M1460,N,W8,N,N,S" // Windows
	//stamp := "65535;53;DF;2;64;M1410,N,W5,N,N,S,E" // iOS
	//stamp := "65535;49;DF;2;64;M1460,N,W6,N,N,S,E" // macOS
	//stamp := "65535;52;DF;2;64;M1452,N,W7,N,N,S,E" // iPadOS
	//stamp := "29200;57;DF;2;60;M1460,S,N,W7" // iPadOS
	//stamp := "14600;46;DF;2;60;M1460,S,N,W6" // NetCast
	//stamp := "14600;50;DF;2;60;M1460,S,N,W7" // Tizen
	//stamp := "29200;53;DF;2;60;M1460,S,N,W8" // Darwin
	//stamp := "65535;244;DF;2;64;M1460,N,N,S,N,W3,N,N" // Windows
	//stamp := "14600;51;DF;2;60;M1460,S,N,W7" // LinuxChrome OS

	//stamp := "14600;55;DF;2;60;M1460,S,N,W7" // LinuxChrome OS
	//stamp := "29200;55;DF;2;60;M1412,S,N,W7" // Tizen
	//stamp := "65535;45;DF;2;60;M1414,N,W7,S" // PlayStation 4
	//stamp := "29200;53;DF;194;60;M1412,S,N,W7" // NetCast ??? (Tizen/NetCast)

	//stamp := "64240;51;DF;2;60;M1460,S,N,W7" // Linux ??? (LinuxChrome OS/Tizen/Linux - 33/17/14)
	//stamp := "65535;241;DF;2;60;M1448,S,N,W8" // Linux ??? (Android/Linux - 61/9)
	//stamp := "14200;56;DF;2;52;M1400,N,N,S,N,W2" // Linux ??? (Windows/macOS/Linux - 50/43/2)
	//stamp := "65535;54;DF;2;60;M1440,S,N,W8" // Linux ??? (Android/Linux - 47/20)

	// -------------------------- combined ----------------------------------
	//stamp := "65535;55;DF;2;176;M1460,N,W5,N,N,T,S,E,E" // iOS - ok 1
	//stamp := "65535;55;DF;2;176;M1460,N,W6,N,N,T,S,E,E" // macOS - ok 2
	//stamp := "64240;55;DF;2;128;M1460,N,W8,N,N,S" // windows - ok 1
	//stamp := "64800;55;DF;2;160;M1350,S,T,N,W7" // Linux - ok 1

	//stamp := "65535;50;DF;2;60;M1452,S,N,W6" // Windows
	//stamp := "65535;138;DF;2;60;M1460,S,N,W6" // Windows
	//stamp := "29200;56;DF;2;52;M1452,N,N,S,N,W7" // iOS

	//stamp := "29200;54;DF;2;52;M1460,N,N,S,N,W7" // Android
	//stamp := "29200;50;DF;2;52;M1380,N,W6,N,N,S" // Android ???
	//stamp := "65535;51;DF;2;60;M1400,S,N,W9" // Windows ???
	//stamp := "64240;50;DF;2;60;M1460,S,N,W7" // iOS ???

	// ----------------------------------------------------------------------

	//stamp := "65495;64;DF;2;60;M65495,S,N,W7" // local port

	d := parseStamp(stamp)
	fmt.Println("dataStamp=", d)

	scorePlatforms := sumScorePlatforms(d.TCPOptions, scoreHL[d.TCPHeaderLength], scoreWin[d.TCPWindowSize], scoreOpt[d.TCPOptions], scoreIpTtl[d.IPTTL], scoreMss[d.MSS])
	printScorePlatforms(scorePlatforms)

	//fmt.Println()
	//fmt.Println("Opt", scoreOpt[d.TCPOptions])
	//fmt.Println("HL ", scoreHL[d.TCPHeaderLength])
	//fmt.Println("Win", scoreWin[d.TCPWindowSize])
	//fmt.Println("IPTTL", scoreIpTtl[d.IPTTL])
	//fmt.Println("MSS", scoreMss[d.MSS])

	fmt.Println()
	fmt.Println("===== Нейросеть ===============================================")

	nn := openFileNeuralNet(fileNameNN)

	scPlNameNN := dataPlNameFromNeuralNet(nn, stamp)
	printScPlNameNN(scPlNameNN)

}

//----------------------------------------

type dif struct {
	count   int
	dat     map[int]int
	nameDat map[string]int
	stamp   string
}

func calculatingDataFromFile(fileNameNN string) {
	fileNameHL := "cmd/csv/parseposf/referencefiles/headerLength.txt"
	fileNameWin := "cmd/csv/parseposf/referencefiles/windowSize.txt"
	fileNameOpt := "cmd/csv/parseposf/referencefiles/options.txt"
	fileNameIpTtl := "cmd/csv/parseposf/referencefiles/ipttl.txt"
	fileNameMss := "cmd/csv/parseposf/referencefiles/mss.txt"

	scoreHL := parseDataFile(fileNameHL, 1)   // 1
	scoreWin := parseDataFile(fileNameWin, 1) // 1
	scoreOpt := parseDataFile(fileNameOpt, 1) // 1
	scoreIpTtl := parseDataFile(fileNameIpTtl, 1)
	scoreMss := parseDataFile(fileNameMss, 1)

	fileNameStampDate := "cmd/csv/posfdata25ok.txt" // 289938 (87.217, 5.769)
	//fileNameStampDate := "cmd/csv/posfdata24ok.txt"   // 209086 (89.766, 4.440)
	//fileNameStampDate := "cmd/csv/posfdata23ok.txt" // 112355 (91.664, 3.563)
	//fileNameStampDate := "cmd/csv/posfdata17.txt"
	stamps := openStampDataFromFile(fileNameStampDate)
	fmt.Println("stamps", len(stamps))

	res := make(map[int]dif)
	for _, st := range stamps {
		resDat := make(map[int]int)
		resNameDat := make(map[string]int)

		dS := parseStamp(st.stamp)
		scorePlatforms := sumScorePlatforms(dS.TCPOptions, scoreHL[dS.TCPHeaderLength], scoreWin[dS.TCPWindowSize], scoreOpt[dS.TCPOptions], scoreIpTtl[dS.IPTTL], scoreMss[dS.MSS])

		i, dat, nameD := compareResult(scorePlatforms, st.plName)
		d, ok := res[i]
		if !ok {
			d = dif{dat: resDat, nameDat: resNameDat}
		}
		d.count++
		d.dat[dat]++
		d.nameDat[nameD]++
		if i == nPlatform {
			if len(d.stamp) < 3000 {
				d.stamp += st.plName + " - " + st.stamp + ";  "
			}
		}
		res[i] = d
	}
	printMap(res, len(stamps))

	fmt.Println()
	fmt.Println("===== Нейросеть 7 входных нейрона ==========")

	nn := openFileNeuralNet(fileNameNN)
	resNN := make(map[int]dif)
	for _, st := range stamps {
		resDat := make(map[int]int)
		resNameDat := make(map[string]int)

		scorePlatforms := dataPlNameFromNeuralNet(nn, st.stamp)

		i, dat, nameD := compareResult(scorePlatforms, st.plName)
		d, ok := resNN[i]
		if !ok {
			d = dif{dat: resDat, nameDat: resNameDat}
		}
		d.count++
		d.dat[dat]++
		d.nameDat[nameD]++
		if i == nPlatform {
			if len(d.stamp) < 3000 {
				d.stamp += st.plName + " - " + st.stamp + ";  "
			}
		}
		resNN[i] = d
	}
	printMap(resNN, len(stamps))

	//
	//fmt.Println()
	//fmt.Println("===== Нейросеть 3 входных нейрона ==========")
	//
	//fileNameNN_3 :="cmd/csv/parseposf/nn/gonn996_3_17_54"
	//nn3 := openFileNeuralNet(fileNameNN_3)
	//resNN3 := make(map[int]dif)
	//for _, st := range stamps {
	//	resDat := make(map[int]int)
	//	resNameDat := make(map[string]int)
	//
	//	scorePlatforms := dataPlNameFromNeuralNet3(nn3, st.stamp)
	//
	//	i, dat, nameD := compareResult(scorePlatforms, st.plName)
	//	d, ok := resNN3[i]
	//	if !ok {
	//		d = dif{dat: resDat, nameDat: resNameDat}
	//	}
	//	d.count++
	//	d.dat[dat]++
	//	d.nameDat[nameD]++
	//	if i == nPlatform {
	//		if len(d.stamp) < 3000 {
	//			d.stamp += st.plName + " - " + st.stamp + ";  "
	//		}
	//	}
	//	resNN3[i] = d
	//}
	//printMap(resNN3, len(stamps))

}

func printMap(data map[int]dif, rec int) {
	type dCount struct {
		p int // номер платформы
		c int // count
	}
	ds := make([]dCount, 0, len(data))
	for p, c := range data {
		d := dCount{p, c.count}
		ds = append(ds, d)
	}
	sort.SliceStable(ds, func(i, j int) bool {
		return ds[i].p < ds[j].p // сортировка по возрастанию
	})
	for _, v := range ds {
		fmt.Printf("%3d: %5d  %.3f \n", v.p, v.c, float64(v.c)/float64(rec)*100)
		switch v.p {
		case 1:
		//fmt.Println(data[v.p].dat)
		//fmt.Println(data[v.p].nameDat)
		case nPlatform:
			fmt.Println(data[v.p].nameDat)
			fmt.Println(data[v.p].stamp)
		}
	}
}

type stampDataFromFile struct {
	plName string
	stamp  string
}

func openStampDataFromFile(fileName string) []stampDataFromFile {
	file, errOpen := ioutil.ReadFile(fileName)
	if errOpen != nil {
		fmt.Println("error open file stamp data", errOpen.Error())
		return nil
	}

	dataFile := strings.Split(string(file), "\n")
	fmt.Println("dataFile", len(dataFile))

	dStamp := make([]stampDataFromFile, 0, 150000)
	for _, lineFile := range dataFile {
		lf := strings.Split(lineFile, `";"`)
		if len(lf) != 9 {
			continue
		}
		platformName := strings.Trim(lf[3], " \"\n\t")
		stamp := strings.Trim(lf[8], " \"\n\t")

		if platformName == "" || stamp == "" {
			continue
		}

		d := stampDataFromFile{platformName, stamp}
		dStamp = append(dStamp, d)
	}
	return dStamp
}

func compareResult(scorePlatforms map[string]float64, plName string) (int, int, string) {
	ds := make([]dataScore, 0, len(scorePlatforms))
	for pl, sc := range scorePlatforms {
		d := dataScore{pl, sc}
		ds = append(ds, d)
	}
	sort.SliceStable(ds, func(i, j int) bool {
		return ds[i].score > ds[j].score // сортировка по убыванию рейтинга
	})
	dat := 100
	nameD := ""
	for i, v := range ds {
		if v.plName == plName {
			if i != 0 {
				dat = int(ds[0].score - v.score)
				nameD = ds[0].plName + " - " + v.plName
			}
			return i, dat, nameD
		}
	}
	return 100, 100, plName
}

//-----------------------------

func printScorePlatforms(scorePlatforms map[string]float64) {
	ds := make([]dataScore, 0, len(scorePlatforms))
	for pl, sc := range scorePlatforms {
		d := dataScore{pl, sc}
		ds = append(ds, d)
	}
	sort.SliceStable(ds, func(i, j int) bool {
		return ds[i].score > ds[j].score // сортировка по убыванию рейтинга
	})

	for _, v := range ds {
		res := fmt.Sprintf("%14s = %f", v.plName, v.score)
		fmt.Println(res)
	}

	//fmt.Println()
	//dataJSON := convertInJSON(ds)
	//fmt.Println(comperestring(dataJSON))
}

func printScPlNameNN(scorePlatforms map[string]float64) {
	ds := make([]dataScore, 0, len(scorePlatforms))
	for pl, sc := range scorePlatforms {
		d := dataScore{pl, sc}
		ds = append(ds, d)
	}
	sort.SliceStable(ds, func(i, j int) bool {
		return ds[i].score > ds[j].score // сортировка по убыванию рейтинга
	})

	for _, v := range ds {
		res := fmt.Sprintf("%14s = %f", v.plName, v.score)
		fmt.Println(res)
	}

}

func convertInJSON(ds []dataScore) []byte {
	type mapPlNameScore map[string]float64

	res := make([]mapPlNameScore, 0, 16)

	sort.SliceStable(ds, func(i, j int) bool {
		return ds[i].score > ds[j].score // сортировка по убыванию рейтинга
	})

	for _, v := range ds {
		ns := mapPlNameScore{v.plName: v.score}
		res = append(res, ns)
	}

	dataJSON, errJSON := json.Marshal(res)
	if errJSON != nil {
		fmt.Println("error json marshal", errJSON.Error())
		return nil
	}

	return dataJSON
}

// суммируем вероятности по платформам
func sumScorePlatforms(stampOpt string, hl, win, opt, ipttl, mss []dataScore) map[string]float64 {
	res := make(map[string]float64)
	sumResult(hl, res)
	sumResult(win, res)
	//sumResult(opt, res)
	//sumResult(ipttl, res)
	//sumResult(mss, res)

	scOpt := setOptions(stampOpt)
	sumResult(scOpt, res)

	return res
}

func sumResult(data []dataScore, res map[string]float64) {
	for _, v := range data {
		//res[v.plName] += v.score  * 100/float64(len(data))   //* float64(10) // * 100/float64(len(data))
		res[v.plName] += v.score * float64(10)
	}
}

// Для Options важна последовательность (например, stamp := "65495;64;DF;2;60;M65495,S,N,W7")
//     Options = "M65495,S,N,W7"
// так для Android'a и Linux на позиции idx=1 чаще всего стоит "S" и последний idx чаще "W"
// для Windows последний idx чаще "S"
// для iOS и macOS последний idx чаще "E"

func setOptions(stampOpt string) []dataScore {
	d := strings.Split(stampOpt, ";")
	if len(d) == 0 {
		return nil
	}
	data := d[len(d)-1]

	el := strings.Split(data, ",")
	l := len(el)
	if l < 2 {
		return nil
	}

	n := float64(0.895) // 0.90-0.89  0.895
	res := make([]dataScore, 0)
	if el[1] == "S" {
		dsAnd := dataScore{"Android", n}
		dsLin := dataScore{"Linux", n}
		dsWin := dataScore{"Windows", 0.17}
		dsIOS := dataScore{"iOS", 0.5}
		dsMacOS := dataScore{"macOS", 0.5}
		if strings.Contains(el[l-1], "W") {
			dsAnd = dataScore{"Android", n + 0.8} //n+0.8
			dsLin = dataScore{"Linux", n + 0.41}  //n+0.5 0.404
			//dsWin = dataScore{"Windows", 0.00001}
			dsIOS = dataScore{"iOS", n}
			dsMacOS = dataScore{"macOS", n}
		}
		res = append(res, dsAnd, dsLin, dsWin, dsIOS, dsMacOS)
	}
	if el[l-1] == "S" {
		dsWin := dataScore{"Windows", 0.664} // 0.664
		res = append(res, dsWin)
	}
	if el[l-1] == "E" {
		dsIOS := dataScore{"iOS", n}
		dsMac := dataScore{"macOS", n}
		res = append(res, dsIOS, dsMac)
	}

	if l < 4 {
		return res
	}

	//   0   1 2 3 4 5
	// M1460,N,N,S,N,W7 - Android, iOS
	if el[1] == "N" && el[2] == "N" && el[3] == "S" {
		dsAnd := dataScore{"Android", n}
		dsIOS := dataScore{"iOS", 0.00001}
		if strings.Contains(el[l-1], "W") {
			dsAnd = dataScore{"Android", n + 0.4} //n+0.4
			dsIOS = dataScore{"iOS", n}
		}
		res = append(res, dsAnd, dsIOS)
	}
	//   0   1  2 3
	// M1452,N,W7,S - PlayStation 4, macOS, iOS
	if el[1] == "N" && el[2] == "W7" && el[3] == "S" {
		dsPlayStation4 := dataScore{"PlayStation 4", n + 0.1}
		res = append(res, dsPlayStation4)
	}

	//   0   1  2 3
	// M1460,N,W9,S - macOS, iOS
	if el[1] == "N" && strings.Contains(el[2], "W") && el[3] == "S" {
		dsMacOS := dataScore{"macOS", 0.76}
		dsIOS := dataScore{"iOS", 0.79}
		res = append(res, dsMacOS, dsIOS)
	}

	//   0   1 2  3
	// M1460,S,N,W7 - iPadOS, NetCast, Tizen, Darwin, LinuxChrome OS, KAIOS, Windows, iOS, macOS, Linux
	if el[1] == "S" && el[2] == "N" && strings.Contains(el[3], "W") {
		dsIPadOS := dataScore{"iPadOS", n + 0.3}
		dsNetCast := dataScore{"NetCast", n + 0.2}
		dsTizen := dataScore{"Tizen", n + 0.2}
		dsDarwin := dataScore{"Darwin", n + 0.2}
		dsLinuxChromeOS := dataScore{"LinuxChrome OS", n + 0.2}
		dsKAIOS := dataScore{"KAIOS", n + 0.2}
		dsWin := dataScore{"Windows", n}
		dsIOS := dataScore{"iOS", 0.6}
		dsMacOS := dataScore{"macOS", 0.4}
		dsLin := dataScore{"Linux", 0.2}
		res = append(res, dsIPadOS, dsNetCast, dsTizen, dsDarwin, dsLinuxChromeOS, dsKAIOS, dsWin, dsIOS, dsMacOS, dsLin)
	}

	if l < 5 {
		return res
	}
	//   0   1  2 3 4
	// M1460,N,W5,S,E  - Android
	if el[1] == "N" && strings.Contains(el[2], "W") && el[3] == "S" && el[4] == "E" {
		dsAnd := dataScore{"Android", n}
		res = append(res, dsAnd)
	}

	if l < 6 {
		return res
	}
	//   0   1  2 3 4 5
	// M1460,N,W8,N,N,S  - Android, iOS
	if el[1] == "N" && strings.Contains(el[2], "W") && el[3] == "N" && el[4] == "N" && el[5] == "S" {
		dsAnd := dataScore{"Android", n}
		dsIOS := dataScore{"iOS", 0.04}
		res = append(res, dsAnd, dsIOS)
	}

	if l < 8 {
		return res
	}
	//   0   1  2 3 4 5 6 7
	// M1380,N,W6,N,N,S,N,N  - Android
	if el[1] == "N" && strings.Contains(el[2], "W") &&
		el[3] == "N" && el[4] == "N" && el[5] == "S" && el[6] == "N" && el[7] == "N" {
		dsAnd := dataScore{"Android", n}
		res = append(res, dsAnd)
	}

	//   0   1 2 3 4  5 6 7
	// M1460,N,N,S,N,W3,N,N - Windows
	if el[1] == "N" && el[2] == "N" && el[3] == "S" && el[4] == "N" &&
		strings.Contains(el[5], "W") && el[6] == "N" && el[7] == "N" {
		dsWin := dataScore{"Windows", n + 0.2}
		res = append(res, dsWin)
	}

	return res
}

func parseStamp(stamp string) stampData {
	d := stampData{}
	s := strings.Split(stamp, ";")
	if len(s) != 6 {
		return d
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

func parseDataFile(fileName string, weight float64) map[string][]dataScore {
	file, errOpen := ioutil.ReadFile(fileName)
	if errOpen != nil {
		fmt.Println("error open file", errOpen.Error())
		return nil
	}

	dataFile := strings.Split(string(file), "\n")

	dataSc := make(map[string][]dataScore, 0)
	for _, lineFile := range dataFile {
		lf := strings.Split(lineFile, ";")
		if len(lf) == 0 {
			continue
		}
		for i, v := range lf {
			if len(v) == 0 || i == 0 {
				continue
			}
			el := strings.Split(v, ",")
			if len(el) == 0 {
				continue
			}
			key := lf[0]
			d, ok := dataSc[key]
			if !ok {
				d = make([]dataScore, 0, len(lf))
				dataSc[key] = d
			}

			sc, errParse := strconv.ParseFloat(el[1], 64)
			if errParse != nil {
				fmt.Println("error open parse float", errParse.Error())
			}
			dScore := dataScore{el[0], sc}
			dataSc[key] = append(dataSc[key], dScore)
		}
	}
	return setScoreForPlatforms(dataSc, weight)
}

func setScoreForPlatforms(data map[string][]dataScore, weight float64) map[string][]dataScore {
	res := make(map[string][]dataScore)
	for key, dScore := range data {
		var div float64
		newScore := make([]dataScore, 0, len(dScore))
		for _, d := range dScore {
			div += d.score
		}
		for _, d := range dScore {
			proc := (d.score / div)
			//if proc == float64(1) {
			//	proc *=10
			//}
			dS := dataScore{d.plName, proc * weight}
			newScore = append(newScore, dS)
		}
		res[key] = newScore
	}
	return res
}

// ===== Нейросеть 7 нейронов ===============================================

// openFileNeuralNet загружает обученную нейросеть
func openFileNeuralNet(fileNameNN string) *gonn.NeuralNetwork {
	return gonn.LoadNN(fileNameNN)
}

func dataPlNameFromNeuralNet(nn *gonn.NeuralNetwork, stamp string) map[string]float64 {

	d := parseStampForNN(stamp)

	windowSize, _ := strconv.ParseFloat("0."+d.TCPWindowSize, 64)
	ipTTL, _ := strconv.ParseFloat("0."+d.IPTTL, 64)
	ipFlags, _ := strconv.ParseFloat("0."+d.IPFlags, 64)
	tcpFlags, _ := strconv.ParseFloat("0."+d.TCPFlags, 64)
	tcpHeaderLength, _ := strconv.ParseFloat("0."+d.TCPHeaderLength, 64)
	options, _ := strconv.ParseFloat("0."+d.TCPOptions, 64)
	mss, _ := strconv.ParseFloat("0."+d.MSS, 64)

	// Получаем ответ от НС
	// (массив весов: TCPWindowSize,TCPHeaderLength,IPFlags,TCPFlags,IPTTL,TCPOptions,MSS)
	dataNN := []float64{
		windowSize,
		tcpHeaderLength,
		ipFlags,
		tcpFlags,
		ipTTL,
		options,
		mss,
	}
	out := nn.Forward(dataNN)

	return getResult(out)
}

// --- так записано в обучающем файле - важна последовательность данных
//res += fmt.Sprintf("0.%s,0.%s,0.%s,0.%s,0.%s,0.%s,0.%s,%s\n",
//              stamp.TCPWindowSize,   // 1
//              stamp.TCPHeaderLength, // 2
//              stamp.IPFlags,         // 3
//              stamp.TCPFlags,        // 4
//              stamp.IPTTL,           // 5
//              stamp.TCPOptions,      // 6
//              stamp.MSS,             // 7
//              platform,
//              )

// ===== Нейросеть 3 нейронов ===============================================

func dataPlNameFromNeuralNet3(nn *gonn.NeuralNetwork, stamp string) map[string]float64 {

	d := parseStampForNN(stamp)

	windowSize, _ := strconv.ParseFloat("0."+d.TCPWindowSize, 64)
	//ipTTL, _ := strconv.ParseFloat("0."+d.IPTTL, 64)
	//ipFlags, _ := strconv.ParseFloat("0."+d.IPFlags, 64)
	//tcpFlags, _ := strconv.ParseFloat("0."+d.TCPFlags, 64)
	tcpHeaderLength, _ := strconv.ParseFloat("0."+d.TCPHeaderLength, 64)
	options, _ := strconv.ParseFloat("0."+d.TCPOptions, 64)
	//mss, _ := strconv.ParseFloat("0."+d.MSS, 64)

	// Получаем ответ от НС
	// (массив весов: TCPWindowSize,TCPHeaderLength,IPFlags,TCPFlags,IPTTL,TCPOptions,MSS)
	dataNN := []float64{
		windowSize,
		tcpHeaderLength,
		//ipFlags,
		//tcpFlags,
		//ipTTL,
		options,
		//mss,
	}
	out := nn.Forward(dataNN)

	return getResult(out)
}

// getResult формирует ответ нейросети
func getResult(output []float64) map[string]float64 {
	if len(output) != 17 {
		return nil
	}
	res := make(map[string]float64)

	res["Android"] = output[0]
	res["iOS"] = output[1]
	res["Windows"] = output[2]
	res["macOS"] = output[3]
	res["iPadOS"] = output[4]
	res["Linux"] = output[5]
	res["LinuxChrome OS"] = output[6]
	res["PlayStation 4"] = output[7]
	res["Tizen"] = output[8]
	res["Darwin"] = output[9]
	res["NetCast"] = output[10]
	res["KAIOS"] = output[11]
	res["Windows Phone"] = output[12]
	res["SmartTV"] = output[13]
	res["FreeBSD"] = output[14]
	res["BlackBerry"] = output[15]
	res["Trident"] = output[16]

	return res
}

func parseStampForNN(stamp string) stampData {
	d := stampData{}
	s := strings.Split(stamp, ";")
	if len(s) != 6 {
		return d
	}

	d.TCPWindowSize = s[0]
	d.IPTTL = s[1]
	d.IPFlags = getIPFlagsForNN(s[2])
	d.TCPFlags = s[3]
	d.TCPHeaderLength = s[4]
	d.TCPOptions = convertHexInDec(md5Data(s[5]))
	d.MSS = getMSS(s[5])

	return d
}

func getIPFlagsForNN(f string) string {
	switch f {
	case "DF":
		return "9"
	case "MF":
		return "5"
	}
	return "0"
}

func md5Data(data string) string {
	h := md5.Sum([]byte(data))
	return fmt.Sprintf("%x", h)
}

func convertHexInDec(hex string) string {
	res := 0

	for i := 0; i < 22; i += 11 {
		d, err := convertInt(hex[i:i+11], 16, 10)
		if err != nil {
			fmt.Println(err.Error())
			return ""
		}
		r, _ := strconv.Atoi(d)
		res += r
	}

	return strconv.Itoa(res)
}

// convertInt конвертирует значение из одной системы счисления в другую, которая указана в toBase
func convertInt(val string, base, toBase int) (string, error) {
	i, err := strconv.ParseInt(val, base, 64)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(i, toBase), nil
}

func sumMyAndNN(fileNameNN string) {
	fileNameHL := "cmd/csv/parseposf/referencefiles/headerLength.txt"
	fileNameWin := "cmd/csv/parseposf/referencefiles/windowSize.txt"
	fileNameOpt := "cmd/csv/parseposf/referencefiles/options.txt"
	fileNameIpTtl := "cmd/csv/parseposf/referencefiles/ipttl.txt"
	fileNameMss := "cmd/csv/parseposf/referencefiles/mss.txt"

	scoreHL := parseDataFile(fileNameHL, 1)   // 1
	scoreWin := parseDataFile(fileNameWin, 1) // 1
	scoreOpt := parseDataFile(fileNameOpt, 1) // 1
	scoreIpTtl := parseDataFile(fileNameIpTtl, 1)
	scoreMss := parseDataFile(fileNameMss, 1)

	fileNameStampDate := "cmd/csv/posfdata25ok.txt" // 289938 (87.217, 5.769)
	//fileNameStampDate := "cmd/csv/posfdata24ok.txt"   // 209086 (89.766, 4.440)
	//fileNameStampDate := "cmd/csv/posfdata23ok.txt" // 112355 (91.664, 3.563)
	//fileNameStampDate := "cmd/csv/posfdata17.txt"
	stamps := openStampDataFromFile(fileNameStampDate)
	fmt.Println("stamps", len(stamps))

	nn := openFileNeuralNet(fileNameNN)

	res := make(map[int]dif)
	for _, st := range stamps {
		resDat := make(map[int]int)
		resNameDat := make(map[string]int)

		dS := parseStamp(st.stamp)
		scMy := sumScorePlatforms(dS.TCPOptions, scoreHL[dS.TCPHeaderLength], scoreWin[dS.TCPWindowSize], scoreOpt[dS.TCPOptions], scoreIpTtl[dS.IPTTL], scoreMss[dS.MSS])
		scNN := dataPlNameFromNeuralNet(nn, st.stamp)

		scorePlatforms := sumDataMyNN(scMy, scNN)

		i, dat, nameD := compareResult(scorePlatforms, st.plName)
		d, ok := res[i]
		if !ok {
			d = dif{dat: resDat, nameDat: resNameDat}
		}
		d.count++
		d.dat[dat]++
		d.nameDat[nameD]++
		if i == nPlatform {
			if len(d.stamp) < 3000 {
				d.stamp += st.plName + " - " + st.stamp + ";  "
			}
		}
		res[i] = d
	}
	printMap(res, len(stamps))
}

func sumDataMyNN(dataMy, dataNN map[string]float64) map[string]float64 {
	res := make(map[string]float64)

	for plName, scMy := range dataMy {
		scNN, ok := dataNN[plName]
		if !ok {
			res[plName] = scMy
			continue
		}

		res[plName] = math.Sqrt(scMy*scMy + scNN*scNN*10000)
	}

	return res
}

func saveResultDataInFile() {
	fileNameStampDate := "cmd/csv/posfdata25ok.txt"
	//fileNameStampDate := "cmd/csv/pd.txt"
	file, errOpen := ioutil.ReadFile(fileNameStampDate)
	if errOpen != nil {
		fmt.Println("error open file stamp data", errOpen.Error())
		return
	}

	dataFile := strings.Split(string(file), "\n")
	fmt.Println("dataFile", len(dataFile))

	fileNameHL := "cmd/csv/parseposf/referencefiles/headerLength.txt"
	fileNameWin := "cmd/csv/parseposf/referencefiles/windowSize.txt"
	fileNameOpt := "cmd/csv/parseposf/referencefiles/options.txt"
	fileNameIpTtl := "cmd/csv/parseposf/referencefiles/ipttl.txt"
	fileNameMss := "cmd/csv/parseposf/referencefiles/mss.txt"

	scoreHL := parseDataFile(fileNameHL, 1)   // 1
	scoreWin := parseDataFile(fileNameWin, 1) // 1
	scoreOpt := parseDataFile(fileNameOpt, 1) // 1
	scoreIpTtl := parseDataFile(fileNameIpTtl, 1)
	scoreMss := parseDataFile(fileNameMss, 1)

	dStamp := make([]string, 0, len(dataFile))

	nn := openFileNeuralNet("cmd/csv/parseposf/nn/gonn20000_119")

	for _, lineFile := range dataFile {
		lf := strings.Split(lineFile, `";"`)
		if len(lf) != 9 {
			continue
		}
		platformName := strings.Trim(lf[3], " \"\n\t")
		stamp := strings.Trim(lf[8], " \"\n\t")

		if platformName == "" || stamp == "" {
			continue
		}

		dS := parseStamp(stamp)
		scorePlatforms := sumScorePlatforms(dS.TCPOptions, scoreHL[dS.TCPHeaderLength], scoreWin[dS.TCPWindowSize], scoreOpt[dS.TCPOptions], scoreIpTtl[dS.IPTTL], scoreMss[dS.MSS])

		resPosf, _ := json.Marshal(scorePlatforms)
		d := lineFile + `;"` + string(resPosf) + `"`

		scPlNameNN := dataPlNameFromNeuralNet(nn, stamp)
		resNN, _ := json.Marshal(scPlNameNN)
		d = d + `;"` + string(resNN) + `"`

		dStamp = append(dStamp, d)
	}
	saveData := strings.Join(dStamp, "\n")
	saveDataFile("cmd/csv/pd_res.txt", []byte(saveData))

}

func saveDataFile(fileName string, data []byte) {
	df, errCreateFile := os.Create(fileName)
	if errCreateFile != nil {
		fmt.Errorf("error create file, %w", errCreateFile)
		return
	}
	defer df.Close()

	_, errWrite := df.Write(data)
	if errWrite != nil {
		fmt.Errorf("error write data %w", errWrite)
		return
	}
	fmt.Printf("file saved: %s\n", fileName)
}
