//  https://tproger.ru/translations/neural-net-from-scratch-in-go/

package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
)

// neuralNet содержит всю информацию, которая определяет обученную сеть
type neuralNet struct {
	config  neuralNetConfig
	wHidden *mat.Dense
	bHidden *mat.Dense
	wOut    *mat.Dense
	bOut    *mat.Dense
}

// neuralNetConfig определяет архитектуру и параметры обучения нашей сети
type neuralNetConfig struct {
	inputNeurons  int
	outputNeurons int
	hiddenNeurons int
	numEpochs     int
	learningRate  float64
}

func main() {
	path := "cmd/neuralnet/data/"
	// Открываем файл с обучающими данными
	inputs, labels := makeInputsAndLabels(path + "train.csv")

	// Определяем нашу сетевую архитектуру и параметры обучения
	config := neuralNetConfig{
		inputNeurons:  4,
		outputNeurons: 3,
		hiddenNeurons: 3,
		numEpochs:     5000,
		learningRate:  0.3,
	}

	// Обучите нейронную сеть
	network := newNetwork(config)
	if err := network.train(inputs, labels); err != nil {
		log.Fatal(err)
	}

	// Открываем файл с тестовыми данными
	testInputs, testLabels := makeInputsAndLabels(path + "test.csv")

	// Делает прогнозы, используя обученную модель
	predictions, err := network.predict(testInputs)
	if err != nil {
		log.Fatal(err)
	}

	// Рассчитывает точность нашей модели
	var truePosNeg int
	numPreds, _ := predictions.Dims()
	for i := 0; i < numPreds; i++ {

		// Получаем вид
		labelRow := mat.Row(nil, i, testLabels)
		//log.Printf("len labelRow = %d, [0] = %f, [1] = %f, [2] = %f", len(labelRow), labelRow[0],labelRow[1],labelRow[2])
		var prediction int
		for idx, label := range labelRow {
			if label == 1.0 {
				prediction = idx
				break
			}
		}

		resAt := predictions.At(i, prediction)
		resMax := floats.Max(mat.Row(nil, i, predictions))
		log.Printf("resMax = %f, resAt = %f", resMax, resAt)

		// Считаем количество верных предсказаний
		if predictions.At(i, prediction) == floats.Max(mat.Row(nil, i, predictions)) {
			truePosNeg++
		}

	}
	log.Printf("truePosNeg = %d, numPreds = %d", truePosNeg, numPreds)
	// Подсчитываем точность предсказаний
	accuracy := float64(truePosNeg) / float64(numPreds)

	fmt.Printf("\nAccuracy = %0.2f\n\n", accuracy)
}

// newNetwork инициализирует новую нейронную сеть
func newNetwork(config neuralNetConfig) *neuralNet {
	return &neuralNet{config: config}
}

// train обучает нейронную сеть, используя обратное распространение
func (nn *neuralNet) train(x, y *mat.Dense) error {

	// Инициализируем смещения/веса
	randSource := rand.NewSource(time.Now().UnixNano())
	randGen := rand.New(randSource)

	wHidden := mat.NewDense(nn.config.inputNeurons, nn.config.hiddenNeurons, nil)
	bHidden := mat.NewDense(1, nn.config.hiddenNeurons, nil)
	wOut := mat.NewDense(nn.config.hiddenNeurons, nn.config.outputNeurons, nil)
	bOut := mat.NewDense(1, nn.config.outputNeurons, nil)

	wHiddenRaw := wHidden.RawMatrix().Data
	bHiddenRaw := bHidden.RawMatrix().Data
	wOutRaw := wOut.RawMatrix().Data
	bOutRaw := bOut.RawMatrix().Data

	for _, param := range [][]float64{wHiddenRaw, bHiddenRaw, wOutRaw, bOutRaw} {
		for i := range param {
			param[i] = randGen.Float64()
		}
	}
	// Определяем выход сети
	output := new(mat.Dense)

	// Используем обратное распространение для регулировки весов и смещений
	if err := nn.backpropagate(x, y, wHidden, bHidden, wOut, bOut, output); err != nil {
		return err
	}

	// Определяем обученную сеть
	nn.wHidden = wHidden
	nn.bHidden = bHidden
	nn.wOut = wOut
	nn.bOut = bOut

	return nil
}

// backpropagate завершает метод прямого распространения
func (nn *neuralNet) backpropagate(x, y, wHidden, bHidden, wOut, bOut, output *mat.Dense) error {
	var count int
	// Обучаем нашу модель в течение определенного количества эпох, используя обратное распространение
	for i := 0; i < nn.config.numEpochs; i++ {

		// Завершаем процесс прямого распространения
		hiddenLayerInput := new(mat.Dense)
		hiddenLayerInput.Mul(x, wHidden)
		addBHidden := func(_, col int, v float64) float64 { return v + bHidden.At(0, col) }
		hiddenLayerInput.Apply(addBHidden, hiddenLayerInput)

		hiddenLayerActivations := new(mat.Dense)
		applySigmoid := func(_, _ int, v float64) float64 { return sigmoid(v) }
		hiddenLayerActivations.Apply(applySigmoid, hiddenLayerInput)

		outputLayerInput := new(mat.Dense)
		outputLayerInput.Mul(hiddenLayerActivations, wOut)
		addBOut := func(_, col int, v float64) float64 { return v + bOut.At(0, col) }
		outputLayerInput.Apply(addBOut, outputLayerInput)
		output.Apply(applySigmoid, outputLayerInput)

		// Завершаем обратное распространение
		networkError := new(mat.Dense)
		networkError.Sub(y, output)

		slopeOutputLayer := new(mat.Dense)
		applySigmoidPrime := func(_, _ int, v float64) float64 { return sigmoidPrime(v) }
		slopeOutputLayer.Apply(applySigmoidPrime, output)
		slopeHiddenLayer := new(mat.Dense)
		slopeHiddenLayer.Apply(applySigmoidPrime, hiddenLayerActivations)

		dOutput := new(mat.Dense)
		dOutput.MulElem(networkError, slopeOutputLayer)
		errorAtHiddenLayer := new(mat.Dense)
		errorAtHiddenLayer.Mul(dOutput, wOut.T())

		dHiddenLayer := new(mat.Dense)
		dHiddenLayer.MulElem(errorAtHiddenLayer, slopeHiddenLayer)

		// Регулируем параметры
		wOutAdj := new(mat.Dense)
		wOutAdj.Mul(hiddenLayerActivations.T(), dOutput)
		wOutAdj.Scale(nn.config.learningRate, wOutAdj)
		wOut.Add(wOut, wOutAdj)

		bOutAdj, err := sumAlongAxis(0, dOutput)
		if err != nil {
			return err
		}
		bOutAdj.Scale(nn.config.learningRate, bOutAdj)
		bOut.Add(bOut, bOutAdj)

		wHiddenAdj := new(mat.Dense)
		wHiddenAdj.Mul(x.T(), dHiddenLayer)
		wHiddenAdj.Scale(nn.config.learningRate, wHiddenAdj)
		wHidden.Add(wHidden, wHiddenAdj)

		bHiddenAdj, err := sumAlongAxis(0, dHiddenLayer)
		if err != nil {
			return err
		}
		bHiddenAdj.Scale(nn.config.learningRate, bHiddenAdj)
		bHidden.Add(bHidden, bHiddenAdj)

		count++
		if count == 500 {
			fmt.Print("-")
			count = 0
		}
	}
	fmt.Println()
	return nil
}

// predict делает предсказание с помощью обученной нейронной сети
func (nn *neuralNet) predict(x *mat.Dense) (*mat.Dense, error) {

	// Проверяем, представляет ли значение neuralNet обученную модель
	if nn.wHidden == nil || nn.wOut == nil {
		return nil, errors.New("the supplied weights are empty")
	}
	if nn.bHidden == nil || nn.bOut == nil {
		return nil, errors.New("the supplied biases are empty")
	}

	// Определяем выход сети
	output := new(mat.Dense)

	// Завершаем процесс прямого распространения
	hiddenLayerInput := new(mat.Dense)
	hiddenLayerInput.Mul(x, nn.wHidden)
	addBHidden := func(_, col int, v float64) float64 { return v + nn.bHidden.At(0, col) }
	hiddenLayerInput.Apply(addBHidden, hiddenLayerInput)

	hiddenLayerActivations := new(mat.Dense)
	applySigmoid := func(_, _ int, v float64) float64 { return sigmoid(v) }
	hiddenLayerActivations.Apply(applySigmoid, hiddenLayerInput)

	outputLayerInput := new(mat.Dense)
	outputLayerInput.Mul(hiddenLayerActivations, nn.wOut)
	addBOut := func(_, col int, v float64) float64 { return v + nn.bOut.At(0, col) }
	outputLayerInput.Apply(addBOut, outputLayerInput)
	output.Apply(applySigmoid, outputLayerInput)

	return output, nil
}

// sigmoid является реализацией сигмоиды (активационная функция), используемой для активации
func sigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}

//
//// sigmoidPrime является реализацией производной сигмоиды для обратного распространения.
//func sigmoidPrime(x float64) float64 {
//	return x * (1.0 - x)
//}

// sigmoidPrime является реализацией производной сигмоиды для обратного распространения
func sigmoidPrime(x float64) float64 {
	return sigmoid(x) * (1.0 - sigmoid(x))
}

// sumAlongAxis позволяет складывать значения только по столбцам или только по строкам матрицы
func sumAlongAxis(axis int, m *mat.Dense) (*mat.Dense, error) {

	numRows, numCols := m.Dims()

	var output *mat.Dense

	switch axis {
	case 0:
		data := make([]float64, numCols)
		for i := 0; i < numCols; i++ {
			col := mat.Col(nil, i, m)
			data[i] = floats.Sum(col)
		}
		output = mat.NewDense(1, numCols, data)
	case 1:
		data := make([]float64, numRows)
		for i := 0; i < numRows; i++ {
			row := mat.Row(nil, i, m)
			data[i] = floats.Sum(row)
		}
		output = mat.NewDense(numRows, 1, data)
	default:
		return nil, errors.New("invalid axis, must be 0 or 1")
	}

	return output, nil
}

// Открываем файл с данными
func makeInputsAndLabels(fileName string) (*mat.Dense, *mat.Dense) {
	// Open the dataset file
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Create a new CSV reader reading from the opened file
	reader := csv.NewReader(f)
	reader.FieldsPerRecord = 7

	// Read in all of the CSV records
	rawCSVData, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// inputsData and labelsData будет содержать все значения с плавающей запятой,
	// которые в конечном итоге будут использоваться для формирования матриц
	inputsData := make([]float64, 4*len(rawCSVData))
	labelsData := make([]float64, 3*len(rawCSVData))

	// Будет отслеживать текущий индекс значений матрицы
	var inputsIndex int
	var labelsIndex int

	for idx, record := range rawCSVData {

		// Пропускаем строку заголовка
		if idx == 0 {
			continue
		}
		// Пропускаем закомментированные строки
		if strings.Contains(record[0], "#") {
			continue
		}

		// Цикл по столбцам
		for i, val := range record {

			parsedVal, err := strconv.ParseFloat(val, 64)
			if err != nil {
				log.Fatal(err)
			}

			// Добавляем к меткам данные, если это необходимо
			if i == 4 || i == 5 || i == 6 {
				labelsData[labelsIndex] = parsedVal
				labelsIndex++
				continue
			}

			inputsData[inputsIndex] = parsedVal
			inputsIndex++
		}
	}
	inputs := mat.NewDense(len(rawCSVData), 4, inputsData)
	labels := mat.NewDense(len(rawCSVData), 3, labelsData)
	return inputs, labels
}
