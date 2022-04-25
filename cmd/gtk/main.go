// https://habr.com/ru/post/420035/
// https://uproger.com/kak-sozdat-krutoj-graficheskij-interfejs-s-pomoshhyu-golang/

// https://github.com/gotk3/gotk3-examples/tree/master/gtk-examples

// Пример обработки событий кнопки мыши.
// Демонстрирует, как обрабатывать щелчки левой, средней и правой кнопок мыши.

package main

import (
	"fmt"
	"log"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

func main() {
	gtk.Init(nil)

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win.SetTitle("Mouse Events")
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	win.Add(windowWidget())
	win.ShowAll()

	gtk.Main()
}

func windowWidget() *gtk.Widget {
	lbl, err := gtk.LabelNew("Use the mouse buttons")
	lbl.SetSizeRequest(200, 200)
	if err != nil {
		panic(err)
	}
	evtBox, err := gtk.EventBoxNew()
	if err != nil {
		panic(err)
	}
	evtBox.Add(lbl)
	evtBox.Connect("button-press-event", func(tree *gtk.EventBox, ev *gdk.Event) bool {
		btn := gdk.EventButtonNewFromEvent(ev)
		fmt.Println("button pressed")
		switch btn.Button() {
		case gdk.BUTTON_PRIMARY:
			lbl.SetText("left-click detected!")
			return true
		case gdk.BUTTON_MIDDLE:
			lbl.SetText("middle-click detected!")
			return true
		case gdk.BUTTON_SECONDARY:
			lbl.SetText("right-click detected!")
			return true
		default:
			return false
		}
	})

	return &evtBox.Widget
}

//
//import (
//	"log"
//
//	"github.com/gotk3/gotk3/gtk"
//)
//
//func main() {
//	// Инициализируем GTK.
//	gtk.Init(nil)
//
//	// Создаём окно верхнего уровня, устанавливаем заголовок
//	// И соединяем с сигналом "destroy" чтобы можно было закрыть
//	// приложение при закрытии окна
//	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
//	if err != nil {
//		log.Fatal("Не удалось создать окно:", err)
//	}
//	win.SetTitle("Простой пример gotk3")
//	win.Connect("destroy", func() {
//		gtk.MainQuit()
//	})
//
//	// Создаём новую метку чтобы показать её в окне
//	l, err := gtk.LabelNew("Привет, мир!")
//	if err != nil {
//		log.Fatal("Не удалось создать метку:", err)
//	}
//
//	// Добавляем метку в окно
//	win.Add(l)
//
//	// Устанавливаем размер окна по умолчанию
//	win.SetDefaultSize(800, 600)
//
//	// Отображаем все виджеты в окне
//	win.ShowAll()
//
//	// Выполняем главный цикл GTK (для отрисовки). Он остановится когда
//	// выполнится gtk.MainQuit()
//	gtk.Main()
//}
