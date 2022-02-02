//https://github.com/wcharczuk/go-chart
package main

import (
	"bytes"
	"fmt"
	"github.com/wcharczuk/go-chart/v2"
	"github.com/wcharczuk/go-chart/v2/drawing"
	"log"
	"math"
	"os"
)

func main() {
	// кол-во точек по x
	n := 200

	x := make([]float64, 0, n)
	y := make([]float64, 0, n)

	for dx := 0.0; dx < float64(n)/10; dx += 0.1 {
		x = append(x, dx)
		y = append(y, math.Sin(dx))
	}

	err := createGraph("cmd/graph/output.png", x, y)
	if err != nil {
		log.Printf("error create graph, %w", err)
	}

	err = createBarGraph("cmd/graph/graph_bar.png")
	if err != nil {
		log.Printf("error create bar graph, %w", err)
	}

	err = createText("cmd/graph/text.png")
	if err != nil {
		log.Printf("error create text, %w", err)
	}

	err = createTwoPointGraph("cmd/graph/two_point.png")
	if err != nil {
		log.Printf("error create text, %w", err)
	}

	log.Print("Done...")
}

// createGraph - создание графика функции с сохранением в файл
func createGraph(fileName string, x, y []float64) error {
	graph := chart.Chart{
		Series: []chart.Series{
			chart.ContinuousSeries{
				XValues: x,
				YValues: y,
			},
		},
	}

	pngFile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer pngFile.Close()

	err = graph.Render(chart.PNG, pngFile)
	if err != nil {
		return err
	}

	log.Printf("saved file: %s", fileName)
	return nil
}

// createBarGraph - создание столбчатого графика с сохранением в файл
func createBarGraph(fileName string) error {
	graph := chart.BarChart{
		Title: "Test Bar Chart",
		Background: chart.Style{
			Padding: chart.Box{
				Top: 40,
			},
		},
		Height:   512,
		BarWidth: 60,
		Bars: []chart.Value{
			{Value: 4.25, Label: "Blue"},
			{Value: 4.88, Label: "Green"},
			{Value: 5.74, Label: "Gray"},
			{Value: 4.22, Label: "Orange"},
			{Value: 3.3, Label: "Test"},
			{Value: 2.27, Label: "Iv"},
			{Value: 1.7, Label: "Bryansk"},
		},
	}

	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	err = graph.Render(chart.PNG, f)
	if err != nil {
		return err
	}
	log.Printf("saved file: %s", fileName)
	return nil
}

// createText - создание рисунка с текстом с сохранением в файл
func createText(fileName string) error {
	f, _ := chart.GetDefaultFont()
	r, _ := chart.PNG(1024, 1024)

	chart.Draw.Text(r, "Test", 64, 64, chart.Style{
		FontColor: drawing.ColorBlack,
		FontSize:  18,
		Font:      f,
	})

	chart.Draw.Text(r, "Ivan", 64, 64, chart.Style{
		FontColor:           drawing.ColorBlack,
		FontSize:            18,
		Font:                f,
		TextRotationDegrees: 45.0,
	})

	tb := chart.Draw.MeasureText(r, "Test", chart.Style{
		FontColor: drawing.ColorGreen,
		FontSize:  18,
		Font:      f,
	}).Shift(64, 64)

	tbc := tb.Corners().Rotate(45)

	chart.Draw.BoxCorners(r, tbc, chart.Style{
		StrokeColor: drawing.ColorRed,
		StrokeWidth: 2,
	})

	tbcb := tbc.Box()
	chart.Draw.Box(r, tbcb, chart.Style{
		StrokeColor: drawing.ColorBlue,
		StrokeWidth: 2,
	})

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	err = r.Save(file)
	if err != nil {
		return err
	}

	log.Printf("saved file: %s", fileName)
	return nil
}

// createTwoGraph - создание двух графиков функции с сохранением в файл
func createTwoGraph(fileName string) error {
	graph := chart.Chart{
		XAxis: chart.XAxis{
			TickPosition: chart.TickPositionBetweenTicks,
			ValueFormatter: func(v interface{}) string {
				typed := v.(float64)
				typedDate := chart.TimeFromFloat64(typed)
				return fmt.Sprintf("%d-%d\n%d", typedDate.Month(), typedDate.Day(), typedDate.Year())
			},
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
			},
			chart.ContinuousSeries{
				YAxis:   chart.YAxisSecondary,
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{50.0, 40.0, 30.0, 20.0, 10.0},
			},
		},
	}

	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	err = graph.Render(chart.PNG, f)
	if err != nil {
		return err
	}

	log.Printf("saved file: %s", fileName)
	return nil
}

// createTwoPointGraph - создание двух графиков функции с сохранением в файл
func createTwoPointGraph(fileName string) error {
	var b float64
	b = 1000

	ts1 := chart.ContinuousSeries{ //TimeSeries{
		Name:    "Time Series",
		XValues: []float64{10 * b, 20 * b, 30 * b, 40 * b, 50 * b, 60 * b, 70 * b, 80 * b},
		YValues: []float64{1.0, 2.0, 30.0, 4.0, 50.0, 6.0, 7.0, 88.0},
	}

	ts2 := chart.ContinuousSeries{ //TimeSeries{
		Style: chart.Style{
			StrokeColor: chart.GetDefaultColor(1),
		},

		XValues: []float64{10 * b, 20 * b, 30 * b, 40 * b, 50 * b, 60 * b, 70 * b, 80 * b},
		YValues: []float64{15.0, 52.0, 30.0, 42.0, 50.0, 26.0, 77.0, 38.0},
	}

	graph := chart.Chart{

		XAxis: chart.XAxis{
			Name:           "The XAxis",
			ValueFormatter: chart.TimeMinuteValueFormatter, //TimeHourValueFormatter,
		},

		YAxis: chart.YAxis{
			Name: "The YAxis",
		},

		Series: []chart.Series{
			ts1,
			ts2,
		},
	}

	buffer := bytes.NewBuffer([]byte{})
	err := graph.Render(chart.PNG, buffer)
	if err != nil {
		return err
	}

	fo, err := os.Create(fileName)
	if err != nil {
		return err
	}

	if _, err := fo.Write(buffer.Bytes()); err != nil {
		return err
	}

	log.Printf("saved file: %s", fileName)
	return nil
}
