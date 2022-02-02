// https://github.com/valyala/fasthttp/blob/master/reuseport/reuseport.go

// Пакет reuseport обеспечивает сеть TCP.Прослушиватель с поддержкой SO_REUSEPORT.
// SO_REUSEPORT позволяет линейно масштабировать производительность сервера на многопроцессорных серверах.
// См. https://www.nginx.com/blog/socket-sharding-nginx-release-1-9-1/ для получения более подробной информации :)
// Пакет основан на https://github.com/kavu/go_reuseport .

package main

import (
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/reuseport"
	"github.com/valyala/tcplisten"
	"log"
	"net"
	"os"
	"strings"
)

// Возвращенный прослушиватель пытается включить следующие параметры TCP, которые обычно оказывают положительное влияние на производительность:
// - TCP_DEFER_ACCEPT. Этот параметр предполагает, что сервер считывает данные из принятого соединения, прежде чем писать им.
// - TCP_FASTOPEN. Видишь https://lwn.net/Articles/508865/ для подробностей.
// Использовать https://github.com/valyala/tcplisten если вы хотите настроить эти варианты.
// Поддерживаются только сети tcp4 и tcp6.
// ErrNoReusePort ошибка возвращается, если система не поддерживает SO_REUSEPORT.

var cfg = &tcplisten.Config{
	ReusePort:   true,
	DeferAccept: true,
	FastOpen:    true,
}

func Listen(network, addr string) (net.Listener, error) {
	ln, err := cfg.NewListener(network, addr)
	if err != nil && strings.Contains(err.Error(), "SO_REUSEPORT") {
		return nil, err
	}
	return ln, err
}

// Программу можно запустить в нескольких окнах на одном порту!!! (запрос будет выполняться произвольно в разных программах)
// для запуска go build fasthttp.go
// chmod +x fasthttp - дать права на запуск
// ./fasthttp - запустить
// запустить программу в двух-трех окнах и дергаем в браузере "127.0.0.1:8080" или в другом окне $ http 127.0.0.1:8080

func main() {
	app := New()
	//	ln, err := Listen("tcp4", "127.0.0.1:8080")
	ln, err := reuseport.Listen("tcp4", "127.0.0.1:8080")
	if err != nil {
		log.Println("error listen", err.Error())
		os.Exit(1)
	}
	defer ln.Close()
	log.Println("listen 127.0.0.1:8080")

	errSer := fasthttp.Serve(ln, app.ServeHTTP)
	if errSer != nil {
		log.Println("error", errSer.Error())
	}
}

type Application struct {
}

func New() *Application {
	return &Application{}
}

func (app *Application) ServeHTTP(ctx *fasthttp.RequestCtx) {
	res := "127.0.0.1:8080 - Hello!"
	_, err := ctx.Write([]byte(res))
	if err != nil {
		log.Println("error write ", err.Error())
	}
	log.Println("call ", ctx.RemoteAddr().String())
}
