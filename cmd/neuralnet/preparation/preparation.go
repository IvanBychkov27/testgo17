package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
)

type stampDataFromFile struct {
	plName string
	stamp  string
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
	fileNameStampDate := "cmd/neuralnet/preparation/data/posfdata25ok.txt"

	stamps := openStampDataFromFile(fileNameStampDate)
	fmt.Println("stamps", len(stamps))

	allPlName, countPlName := getAllPlName(stamps)
	fmt.Println(allPlName)
	_ = countPlName

	//fmt.Println("MaxLenOptions = ", getMaxLenOptions(stamps))

	data, lines := buildDataFile(stamps)

	dataName := fmt.Sprintf("%d_%d", countPlName, lines)
	fileNameSave := "cmd/neuralnet/preparation/data/train_7_" + dataName + ".csv"
	saveDataFile(fileNameSave, data)

	fmt.Println("Done...")
}

func getMaxLenOptions(stamps []stampDataFromFile) int {
	max := 0
	for _, st := range stamps {
		ds := strings.Split(st.stamp, ",")
		if len(ds) < max {
			continue
		}
		max = len(ds)
	}

	return max
}

func buildDataFile(stamps []stampDataFromFile) ([]byte, int) {
	var count, countAnd, countIOS, countWin, countMacOS, countIPadOS,
		countLinux, countLinuxChromeOS, countPlayStation4, countTizen,
		countDarwin, countNetCast, countKAIOS, countWindowsPhone,
		countSmartTV, countFreeBSD, countBlackBerry, countTrident int

	// 7 входных и 17 выходных нейронов для нейросети
	res := "TCPWindowSize,TCPHeaderLength,IPFlags,TCPFlags,IPTTL,TCPOptions,MSS," +
		"Android,iOS,Windows,macOS,iPadOS,Linux,LinuxChromeOS,PlayStation4,Tizen,Darwin,NetCast,KAIOS," +
		"WindowsPhone,SmartTV,FreeBSD,BlackBerry,Trident\n"

	//// 3 входных и 17 выходных нейронов для нейросети
	//res := "TCPWindowSize,TCPHeaderLength,Options,"+
	//	"Android,iOS,Windows,macOS,iPadOS,Linux,LinuxChromeOS,PlayStation4,Tizen,Darwin,NetCast,KAIOS,"+
	//	"WindowsPhone,SmartTV,FreeBSD,BlackBerry,Trident\n"

	n := 78 // количество каждой plName (78 - 996line)
	for _, st := range stamps {
		platform := ""
		switch st.plName {
		case "Android":
			platform = "1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0" // 1
			if countAnd > n {
				continue
			}
			countAnd++
		case "iOS":
			platform = "0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0" // 2
			if countIOS > n {
				continue
			}
			countIOS++
		case "Windows":
			platform = "0,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0" // 3
			if countWin > n {
				continue
			}
			countWin++
		case "macOS":
			platform = "0,0,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0" // 4
			if countMacOS > n {
				continue
			}
			countMacOS++
		case "iPadOS":
			platform = "0,0,0,0,1,0,0,0,0,0,0,0,0,0,0,0,0" // 5
			if countIPadOS > n {
				continue
			}
			countIPadOS++
		case "Linux":
			platform = "0,0,0,0,0,1,0,0,0,0,0,0,0,0,0,0,0" // 6
			if countLinux > n {
				continue
			}
			countLinux++

		case "LinuxChrome OS":
			platform = "0,0,0,0,0,0,1,0,0,0,0,0,0,0,0,0,0" // 7
			if countLinuxChromeOS > n {
				continue
			}
			countLinuxChromeOS++
		case "PlayStation 4":
			platform = "0,0,0,0,0,0,0,1,0,0,0,0,0,0,0,0,0" // 8
			if countPlayStation4 > n {
				continue
			}
			countPlayStation4++
		case "Tizen":
			platform = "0,0,0,0,0,0,0,0,1,0,0,0,0,0,0,0,0" // 9
			if countTizen > n {
				continue
			}
			countTizen++
		case "Darwin":
			platform = "0,0,0,0,0,0,0,0,0,1,0,0,0,0,0,0,0" // 10
			if countDarwin > n {
				continue
			}
			countDarwin++
		case "NetCast":
			platform = "0,0,0,0,0,0,0,0,0,0,1,0,0,0,0,0,0" // 11
			if countNetCast > n {
				continue
			}
			countNetCast++
		case "KAIOS":
			platform = "0,0,0,0,0,0,0,0,0,0,0,1,0,0,0,0,0" // 12
			if countKAIOS > n {
				continue
			}
			countKAIOS++
		case "Windows Phone":
			platform = "0,0,0,0,0,0,0,0,0,0,0,0,1,0,0,0,0" // 13
			if countWindowsPhone > n {
				continue
			}
			countWindowsPhone++
		case "SmartTV":
			platform = "0,0,0,0,0,0,0,0,0,0,0,0,0,1,0,0,0" // 14
			if countSmartTV > n {
				continue
			}
			countSmartTV++
		case "FreeBSD":
			platform = "0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,0,0" // 15
			if countFreeBSD > n {
				continue
			}
			countFreeBSD++
		case "BlackBerry":
			platform = "0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,0" // 16
			if countBlackBerry > n {
				continue
			}
			countBlackBerry++
		case "Trident":
			platform = "0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1" // 17
			if countTrident > n {
				continue
			}
			countTrident++

		default:
			continue
		}

		if countAnd > n && countIOS > n && countWin > n && countMacOS > n &&
			countIPadOS > n && countLinux > n && countLinuxChromeOS > n && countPlayStation4 > n &&
			countTizen > n && countDarwin > n && countNetCast > n && countKAIOS > n &&
			countWindowsPhone > n && countSmartTV > n && countFreeBSD > n && countBlackBerry > n &&
			countTrident > n {
			break
		}

		stamp := parseStamp(st.stamp)

		// 7 входных и 17 выходных данных
		res += fmt.Sprintf("0.%s,0.%s,0.%s,0.%s,0.%s,0.%s,0.%s,%s\n",
			stamp.TCPWindowSize,   // 1
			stamp.TCPHeaderLength, // 2
			stamp.IPFlags,         // 3
			stamp.TCPFlags,        // 4
			stamp.IPTTL,           // 5
			stamp.TCPOptions,      // 6
			stamp.MSS,             // 7
			platform,
		)

		// 3 входных и 17 выходных данных
		//res += fmt.Sprintf("0.%s,0.%s,0.%s,%s\n",
		//	stamp.TCPWindowSize,   // 1
		//	stamp.TCPHeaderLength, // 2
		//	//stamp.IPFlags,         // 3
		//	//stamp.TCPFlags,        // 4
		//	//stamp.IPTTL,           // 5
		//	stamp.TCPOptions,      // 6
		//	//stamp.MSS,             // 7
		//	platform,
		//)

		count++
	}
	return []byte(res), count
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

func parseStamp(stamp string) stampData {
	d := stampData{}
	s := strings.Split(stamp, ";")
	if len(s) != 6 {
		return d
	}

	d.TCPWindowSize = s[0]
	d.IPTTL = s[1]
	d.IPFlags = getIPFlags(s[2])
	d.TCPFlags = s[3]
	d.TCPHeaderLength = s[4]
	d.TCPOptions = convertHexInDec(md5Data(s[5]))
	d.MSS = getMSS(s[5])

	return d
}

func getIPFlags(f string) string {
	switch f {
	case "DF":
		return "9"
	case "MF":
		return "5"
	}
	return "0"
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
// https://golangify.com/binary-to-decimal-octal-and-hexadecimal
func convertInt(val string, base, toBase int) (string, error) {
	i, err := strconv.ParseInt(val, base, 64)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(i, toBase), nil
}

func getAllPlName(stamps []stampDataFromFile) (string, int) {
	plName := make(map[string]int, 0)
	for _, st := range stamps {
		plName[st.plName]++
	}

	res := fmt.Sprintf("\nall plName = %d\n\n", len(plName))

	type dataPlName struct {
		name  string
		count int
	}

	ds := make([]dataPlName, 0, len(plName))
	for n, c := range plName {
		d := dataPlName{n, c}
		ds = append(ds, d)
	}
	sort.SliceStable(ds, func(i, j int) bool {
		return ds[i].count > ds[j].count // сортировка по убыванию
	})

	for _, v := range ds {
		res += fmt.Sprintf("%16s = %d\n", v.name, v.count)
	}
	return res, len(plName)
}
