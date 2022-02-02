package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
)

type NeuralNetwork struct {
	InputLayer   []float64   `json:"input_layer"`
	HiddenLayer  []float64   `json:"hidden_layer"`
	OutputLayer  []float64   `json:"output_layer"`
	TargetLayer  []float64   `json:"-"`
	WeightHidden [][]float64 `json:"weight_hidden"`
	WeightOutput [][]float64 `json:"weight_output"`
	ErrOutput    []float64   `json:"err_output"`
	ErrHidden    []float64   `json:"-"`
	AllErrOutput []float64   `json:"-"`
	Rate         float64     `json:"rate"`
	MSE          float64     `json:"mse"`
}

func main() {
	input, target := dataSet()
	nn := creatNN(3, 32, 1, 0.07)
	nn.train(input, target, 2000)

	//fileNameNN := "cmd/my_neuralnet/mynet_3_32_1.nn"
	//nn := loadNN(fileNameNN)
	//fmt.Printf("NN : %s \nMSE: %0.15f\n",fileNameNN, nn.MSE)

	d := []float64{1, 0, 1}
	res := nn.resultNN(d)
	fmt.Printf("\nres = %d (%0.3f)\n", int(math.Round(res[0]*10)), res[0])
	//fmt.Println(nn.AllErrOutput)

	//dataFile := fmt.Sprintf("%d_%d_%d", len(nn.InputLayer), len(nn.HiddenLayer), len(nn.OutputLayer))
	//fileNameSave := "cmd/my_neuralnet/mynet_" + dataFile + ".nn"
	//nn.saveNN(fileNameSave)

	fmt.Println("Done...")
}

func (nn *NeuralNetwork) train(inputData, targetData [][]float64, iteration int) {
	nn.AllErrOutput = make([]float64, len(targetData))
	message := iteration / 10

	for it := 0; it < iteration; it++ {
		for i := range inputData {
			nn.InputLayer = inputData[i]
			nn.TargetLayer = targetData[i]
			nn.sumWeightHidden()
			nn.sumWeightOutput()
			nn.calcErrOut(i)
			nn.adjustmentWeightOutput()
			nn.adjustmentWeightHidden()
		}

		if it != 0 && (it%message == 0 || it == iteration-1) {
			fmt.Printf("iteration:%6d  MSE:%0.6f \n", it, nn.mse())
		}
	}
	//fmt.Println(nn.AllErrOutput)
	nn.MSE = nn.mse()
}

// регулировка веса внутренних нейронов
func (nn *NeuralNetwork) adjustmentWeightHidden() {
	for h, dH := range nn.HiddenLayer {
		//if h == len(nn.HiddenLayer)-1 {
		//	continue
		//}
		//dW := nn.ErrHidden[h] * sigmoidPrime(dH) // расчет дельта веса
		dW := nn.ErrHidden[h] * thPrime(dH) // расчет дельта веса через производную гиперболического тангенса
		for i, dI := range nn.InputLayer {
			nn.WeightHidden[h][i] -= dI * dW * nn.Rate // регулировка веса внутреннего нейрона
		}
	}
}

// регулировка веса выходных нейронов и расчет ошибки внутренних нейронов
func (nn *NeuralNetwork) adjustmentWeightOutput() {
	for o, dO := range nn.OutputLayer {
		//dW := nn.ErrOutput[o] * sigmoidPrime(dO) // расчет дельта веса
		dW := nn.ErrOutput[o] * thPrime(dO) // расчет дельта веса через производную гиперболического тангенса
		for h, dH := range nn.HiddenLayer {
			nn.WeightOutput[o][h] -= dH * dW * nn.Rate   // регулировка веса выходного нейрона
			nn.ErrHidden[h] = nn.WeightOutput[o][h] * dW // расчет ошибки внутреннего нейрона
		}
	}
}

// расчет ошибки выходных нейронов
func (nn *NeuralNetwork) calcErrOut(i int) {
	for t, dT := range nn.TargetLayer {
		nn.ErrOutput[t] = nn.OutputLayer[t] - dT
		nn.AllErrOutput[i] = nn.ErrOutput[t]
	}
}

// среднеквадратичная ошибка нейронной сети
func (nn *NeuralNetwork) mse() float64 {
	var errSum float64

	for _, dE := range nn.AllErrOutput {
		errSum += dE * dE
	}
	return math.Sqrt(errSum / float64(len(nn.AllErrOutput))) // нахождение среднеквадратичного отклонения
}

// получаем результат от нейронной сети nn по входным данным input
func (nn *NeuralNetwork) resultNN(input []float64) []float64 {
	nn.InputLayer = input
	nn.sumWeightHidden()
	nn.sumWeightOutput()
	return nn.OutputLayer
}

// получаем значения внутренних нейронов
func (nn *NeuralNetwork) sumWeightHidden() {
	for h, dH := range nn.WeightHidden { // перебираем каждый внутренний нейрон
		//if h == len(nn.WeightHidden)-1 {
		//	continue
		//}
		var dataNeuron float64
		for i, dI := range nn.InputLayer {
			dataNeuron += dI * dH[i] // складываем произведения весов внутренних и значений входных нейронов
		}
		//nn.HiddenLayer[h] = sigmoid(dataNeuron) // сохраняем значение внутреннего нейрона с его активацией функцией сигмойда
		nn.HiddenLayer[h] = th(dataNeuron) // сохраняем значение внутреннего нейрона с его активацией функцией гиперболического тангенса
	}
}

// получаем значения выходных нейронов
func (nn *NeuralNetwork) sumWeightOutput() {
	for o, dO := range nn.WeightOutput { // перебираем каждый выходной нейрон
		var dataNeuron float64
		for h, dh := range nn.HiddenLayer {
			dataNeuron += dh * dO[h] // складываем произведения весов выходных и значений внутренних нейронов
		}
		//nn.OutputLayer[o] = sigmoid(dataNeuron) // сохраняем значение выходного нейрона с его активацией функцией сигмойда
		nn.OutputLayer[o] = th(dataNeuron) // сохраняем значение выходного нейрона с его активацией функцией гиперболического тангенса
	}
}

func creatNN(input, hidden, output int, rate float64) *NeuralNetwork {
	//rand.Seed(time.Now().UTC().UnixNano())
	//input++  // добавляем входной нейрон смещения
	//hidden++ // добавляем внутренний нейрон смещения
	nn := &NeuralNetwork{
		InputLayer:   make([]float64, input),
		HiddenLayer:  make([]float64, hidden),
		OutputLayer:  make([]float64, output),
		TargetLayer:  make([]float64, output),
		WeightHidden: randomWeight(hidden, input),
		WeightOutput: randomWeight(output, hidden),
		ErrOutput:    make([]float64, output),
		ErrHidden:    make([]float64, hidden),
		AllErrOutput: make([]float64, output),
		Rate:         rate,
	}
	// присваиваем нейронам смещения значения по умолчанию
	//nn.InputLayer[input-1] = 1.0
	//nn.HiddenLayer[hidden-1] = 1.0
	return nn
}

func randomWeight(n, m int) [][]float64 {
	data := make([][]float64, 0, n)
	for i := 0; i < n; i++ {
		d := make([]float64, 0, m)
		for j := 0; j < m; j++ {
			d = append(d, 2.0*rand.Float64()-1)
		}
		data = append(data, d)
	}
	return data
}

//sigmoid является реализацией сигмоиды
func sigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}

// sigmoidPrime является реализацией производной
func sigmoidPrime(x float64) float64 {
	return sigmoid(x) * (1.0 - sigmoid(x))
}

// th - гиперболический тангенс
func th(x float64) float64 {
	return (math.Exp(2*x) - 1) / (math.Exp(2*x) + 1)
}

// thPrime - производная гиперболического тангенса
func thPrime(x float64) float64 {
	ch := (math.Exp(x) + math.Exp(-x)) / 2 // гиперболический косинус
	return 1.0 / ch * ch
}

// loadNN загрузка нейросети из файла
func loadNN(fileName string) *NeuralNetwork {
	nn := &NeuralNetwork{
		WeightHidden: make([][]float64, 0),
		WeightOutput: make([][]float64, 0),
	}
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("error open file", err.Error())
		return nil
	}
	err = json.Unmarshal(data, nn)
	if err != nil {
		fmt.Println("error json unmarshal", err.Error())
		return nil
	}
	nn.HiddenLayer = make([]float64, len(nn.WeightHidden))
	nn.OutputLayer = make([]float64, len(nn.WeightOutput))
	return nn
}

// saveNN сохранение нейросети в файл
func (nn *NeuralNetwork) saveNN(fileName string) {
	data, err := json.Marshal(nn)
	if err != nil {
		fmt.Println("error json marshal", err.Error())
		return
	}

	df, errCreateFile := os.Create(fileName)
	if errCreateFile != nil {
		fmt.Println("error create file", errCreateFile.Error())
		return
	}
	defer df.Close()

	_, errWrite := df.Write(data)
	if errWrite != nil {
		fmt.Println("error write data", errWrite.Error())
		return
	}
	fmt.Println("saved file", fileName)
}

func dataSet() ([][]float64, [][]float64) {
	input := [][]float64{
		{0, 0, 0},
		{0, 0, 1},
		{0, 1, 0},
		{0, 1, 1},
		{1, 0, 0},
		{1, 0, 1},
		{1, 1, 0},
		{1, 1, 1},
	}
	target := [][]float64{
		{0.0},
		{0.1},
		{0.2},
		{0.3},
		{0.4},
		{0.5},
		{0.6},
		{0.7},
	}
	return input, target
}
