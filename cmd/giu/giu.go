// https://uproger.com/kak-sozdat-krutoj-graficheskij-interfejs-s-pomoshhyu-golang/
package main

import "github.com/AllenDang/giu"

func loop() {
	w1 := giu.Window("window 1")
	w2 := giu.Window("window 2")

	w1W, w1H := w1.CurrentSize()
	w1X, w1Y := w1.CurrentPosition()

	w1Layout := giu.Layout{
		giu.Labelf("Focused state %t", w1.HasFocus()),
		giu.Labelf("Position: %d, %d", int(w1W), int(w1H)),
		giu.Labelf("Size: %d, %d", int(w1X), int(w1Y)),
		giu.Button("bring to front window 2").OnClick(w2.BringToFront),
	}
	w2Layout := giu.Layout{
		giu.Labelf("Focused state %t", w2.HasFocus()),
	}

	w1.Layout(w1Layout)
	w2.Layout(w2Layout)
}

func main() {
	wnd := giu.NewMasterWindow("windows [DEMO]", 640, 480, 0)
	wnd.Run(loop)
}

//
//import (
//	"fmt"
//
//	g "github.com/AllenDang/giu"
//)
//
//func onClickMe() {
//	fmt.Println("Hello world!")
//}
//
//func onImSoCute() {
//	fmt.Println("Im sooooooo cute!!")
//}
//
//func loop() {
//	g.SingleWindow().Layout(
//		g.Label("Hello world from giu"),
//		g.Row(
//			g.Button("Click Me").OnClick(onClickMe),
//			g.Button("I'm so cute").OnClick(onImSoCute),
//		),
//	)
//}
//
//func main() {
//	wnd := g.NewMasterWindow("Hello world", 400, 200, g.MasterWindowFlagsNotResizable)
//	wnd.Run(loop)
//}
