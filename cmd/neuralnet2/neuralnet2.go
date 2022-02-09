// https://habr.com/ru/post/343466/
package main

import (
	"encoding/csv"
	"fmt"
	"github.com/fxsjy/gonn/gonn"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	// создаем и записываем в файл обученную нейросеть (достаточно записать один раз - долее открываем и пользуемся)
	fmt.Println("start training: ", time.Now().Format(" 15:04:05"))
	fileName := "cmd/neuralnet2/data/train_7_17_420.csv" // нейросеть создается на основе данных данного файла (1000 строк - время около 20 минут) (train_12_131 - 12-это кол-во PlName 131-строк обучающих данных)
	createNN(fileName)
	fmt.Println("end of training: ", time.Now().Format(" 15:04:05"))

	// Загружем НС из файла
	//fileNameNN := "cmd/neuralnet2/net/gonn49505_300"
	//fileNameNN := "cmd/neuralnet2/net/gonn20000_119"
	fileNameNN := "gonn"
	nn := gonn.LoadNN(fileNameNN)

	// Записываем значения в переменные:
	// 7 входных нейронов
	windowSize, tcpHeaderLength, ipFlags, tcpFlags, ipTTL, options, mss := 0.65535, 0.64, 0.9, 0.2, 0.53, 0.7327042439008, 0.1460 // - iOS (macOS)    - 1
	//windowSize, tcpHeaderLength, ipFlags, tcpFlags, ipTTL, options, mss := 0.64240, 0.52, 0.9, 0.2, 0.120, 0.15822044659906, 0.1412 // - Windows       - 2
	//windowSize, tcpHeaderLength, ipFlags, tcpFlags, ipTTL, options, mss := 0.65535, 0.60, 0.9, 0.2, 0.54, 0.15207360277883, 0.1400  // - Android       - 3
	//windowSize, tcpHeaderLength, ipFlags, tcpFlags, ipTTL, options, mss := 0.65535, 0.64, 0.9, 0.2, 0.51, 0.7327042439008, 0.1460   // - macOS         - 4
	//windowSize, tcpHeaderLength, ipFlags, tcpFlags, ipTTL, options, mss := 0.65535, 0.64, 0.9, 0.2, 0.50, 0.21994876700036, 0.1460  // - iPadOS        - 5
	//windowSize, tcpHeaderLength, ipFlags, tcpFlags, ipTTL, options, mss := 0.65535, 0.52, 0.9, 0.2, 0.56, 0.10808706590817, 0.1350  // - Linux         - 6
	//windowSize, tcpHeaderLength, ipFlags, tcpFlags, ipTTL, options, mss := 0.29200, 0.60, 0.9, 0.2, 0.48, 0.3675683890415, 0.1460   // - LinuxChromeOS - 7
	//windowSize, tcpHeaderLength, ipFlags, tcpFlags, ipTTL, options, mss := 0.65535, 0.60, 0.9, 0.2, 0.50, 0.12750133365872, 0.1460  // - PlayStation4  - 8
	//windowSize, tcpHeaderLength, ipFlags, tcpFlags, ipTTL, options, mss := 0.29200, 0.60, 0.9, 0.2, 0.47, 0.8123913536086, 0.1440   // - Tizen         - 9
	//windowSize, tcpHeaderLength, ipFlags, tcpFlags, ipTTL, options, mss := 0.29200, 0.60, 0.9, 0.2, 0.53, 0.17649720258654, 0.1460  // - Darwin        - 10
	//windowSize, tcpHeaderLength, ipFlags, tcpFlags, ipTTL, options, mss := 0.29200, 0.60, 0.9, 0.2, 0.55, 0.3675683890415, 0.1460   // - NetCast       - 11
	//windowSize, tcpHeaderLength, ipFlags, tcpFlags, ipTTL, options, mss := 0.29200, 0.60, 0.9, 0.2, 0.46, 0.19823331574194, 0.1370  // - KAIOS         - 12
	//windowSize, tcpHeaderLength, ipFlags, tcpFlags, ipTTL, options, mss := 0.65535, 0.52, 0.9, 0.2, 0.117, 0.30724801112478, 0.1360 // - WindowsPhone  - 13
	//windowSize, tcpHeaderLength, ipFlags, tcpFlags, ipTTL, options, mss := 0.14600, 0.60, 0.9, 0.2, 0.51, 0.23303404777345, 0.1460  // - SmartTV       - 14
	//windowSize, tcpHeaderLength, ipFlags, tcpFlags, ipTTL, options, mss := 0.14600, 0.60, 0.9, 0.2, 0.45, 0.26547435639462, 0.1400  // - FreeBSD       - 15
	//windowSize, tcpHeaderLength, ipFlags, tcpFlags, ipTTL, options, mss := 0.65535, 0.64, 0.9, 0.2, 0.42, 0.15745999206570, 0.1412  // - BlackBerry    - 16
	//windowSize, tcpHeaderLength, ipFlags, tcpFlags, ipTTL, options, mss := 0.65535, 0.52, 0.9, 0.2, 0.116, 0.17536972182053, 0.1460 // - Trident       - 17

	fmt.Println()

	// 3 входных нейронов
	//windowSize, tcpHeaderLength, options := 0.65535, 0.64, 0.7327042439008  // - iOS (macOS)    - 1
	//windowSize, tcpHeaderLength, options := 0.64240, 0.52, 0.15822044659906 // - Windows       - 2
	//windowSize, tcpHeaderLength, options := 0.65535, 0.60, 0.15207360277883 // - Android       - 3
	//windowSize, tcpHeaderLength, options := 0.65535, 0.64, 0.7327042439008  // - macOS         - 4
	//windowSize, tcpHeaderLength, options := 0.65535, 0.64, 0.21994876700036 // - iPadOS        - 5
	//windowSize, tcpHeaderLength, options := 0.65535, 0.52, 0.10808706590817 // - Linux         - 6
	//windowSize, tcpHeaderLength, options := 0.29200, 0.60, 0.3675683890415  // - LinuxChromeOS - 7
	//windowSize, tcpHeaderLength, options := 0.65535, 0.60, 0.12750133365872 // - PlayStation4  - 8
	//windowSize, tcpHeaderLength, options := 0.29200, 0.60, 0.8123913536086  // - Tizen         - 9
	//windowSize, tcpHeaderLength, options := 0.29200, 0.60, 0.17649720258654 // - Darwin        - 10
	//windowSize, tcpHeaderLength, options := 0.29200, 0.60, 0.3675683890415  // - NetCast       - 11
	//windowSize, tcpHeaderLength, options := 0.29200, 0.60, 0.19823331574194 // - KAIOS         - 12
	//windowSize, tcpHeaderLength, options := 0.65535, 0.52, 0.30724801112478 // - WindowsPhone  - 13
	//windowSize, tcpHeaderLength, options := 0.14600, 0.60, 0.23303404777345 // - SmartTV       - 14
	//windowSize, tcpHeaderLength, options := 0.14600, 0.60, 0.26547435639462 // - FreeBSD       - 15
	//windowSize, tcpHeaderLength, options := 0.65535, 0.64, 0.15745999206570 // - BlackBerry    - 16
	//windowSize, tcpHeaderLength, options := 0.65535, 0.52, 0.17536972182053 // - Trident       - 17

	// Получаем ответ от НС (массив весов)
	out := nn.Forward([]float64{windowSize, ipTTL, ipFlags, tcpFlags, tcpHeaderLength, options, mss})
	//out := nn.Forward([]float64{windowSize, tcpHeaderLength, options})

	// Печатаем ответ на экран
	fmt.Println(getResult(out))
}

// createNN Создаем нейросеть
func createNN(fileName string) {
	// Параметры нейросети: кол-во нейронов (входных, внутренних, выходных, регрессия)
	nn := gonn.NewNetwork(7, 32, 17, false, 0.25, 0.1)

	//nn := gonn.DefaultNetwork(7, 12, 17, false)
	//nn := gonn.DefaultNetwork(3, 51, 17, false)

	input, target := openCSV(fileName)

	// Начинаем обучать нашу НС
	nn.Train(input, target, 100000)

	// Сохраняем готовую НС в файл
	gonn.DumpNN("gonn", nn)
}

// открывает файл с базой данных для обучения и выдаем входные и выходные данные
func openCSV(fileName string) (inputs, targets [][]float64) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("error:", err.Error())
		return
	}
	defer file.Close()

	inp := 7  // входных нейронов
	out := 17 // выходных нейронов

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = inp + out // входных + выходных нейронов
	reader.Comment = '#'

	rawCSVData, err := reader.ReadAll()
	if err != nil {
		fmt.Println("error ReadAll:", err.Error())
		return
	}

	for idx, record := range rawCSVData {

		inputsData := make([]float64, 0, inp) // входных нейронов
		labelsData := make([]float64, 0, out) //  выходных нейронов

		if idx == 0 || strings.Contains(record[0], "#") {
			continue
		}

		for i, val := range record {
			d, errPars := strconv.ParseFloat(val, 64)
			if errPars != nil {
				fmt.Println("error:", errPars)
				return
			}
			if i < inp {
				inputsData = append(inputsData, d)
			} else {
				labelsData = append(labelsData, d)
			}
		}
		inputs = append(inputs, inputsData)
		targets = append(targets, labelsData)
	}

	return
}

// getResult выберет ответ нейрона с самым большим весом
func getResult(output []float64) string {
	max := float64(-99999)
	pos := -1
	// Ищем позицию нейрона с самым большим весом
	for i, value := range output {
		if value > max {
			max = value
			pos = i
		}
	}

	res := fmt.Sprintf(" 0: %f Android \n 1: %f iOS \n 2: %f Windows \n 3: %f macOS \n 4: %f iPadOS \n 5: %f Linux \n 6: %f LinuxChrome OS \n 7: %f PlayStation 4 \n 8: %f Tizen \n 9: %f Darwin \n10: %f NetCast\n11: %f KAIOS\n12: %f Windows Phone\n13: %f SmartTV\n14: %f FreeBSD\n15: %f BlackBerry\n16: %f Trident\n",
		output[0], output[1], output[2], output[3], output[4], output[5], output[6],
		output[7], output[8], output[9], output[10], output[11],
		output[12], output[13], output[14], output[15], output[16],
	)
	res += "result: "

	// Теперь, в зависимости от позиции, возвращаем решение
	switch pos {
	case 0:
		res += "Android"
	case 1:
		res += "iOS"
	case 2:
		res += "Windows"
	case 3:
		res += "macOS"
	case 4:
		res += "iPadOS"
	case 5:
		res += "Linux"
	case 6:
		res += "LinuxChrome OS"
	case 7:
		res += "PlayStation 4"
	case 8:
		res += "Tizen"
	case 9:
		res += "Darwin"
	case 10:
		res += "NetCast"
	case 11:
		res += "KAIOS"
	case 12:
		res += "Windows Phone"
	case 13:
		res += "SmartTV"
	case 14:
		res += "FreeBSD"
	case 15:
		res += "BlackBerry"
	case 16:
		res += "Trident"
	default:
		res += "Unknown"
	}
	return res
}
