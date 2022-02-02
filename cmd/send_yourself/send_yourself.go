package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type application struct {
	address string
	uin     string
	server  *http.Server
}

func main() {
	app := &application{address: "127.0.0.1:2999", uin: "send to yourself"}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go app.startServer(ctx, cancel, wg)

	uin := app.getData()
	if uin == app.uin {
		fmt.Println("uin ok!:", uin)
	} else {
		fmt.Println("uin bad:", uin)
	}

	app.stopServerMaster(cancel)

	<-ctx.Done()
	wg.Wait()
	fmt.Println("Done...")
}

func (app *application) getData() string {
	time.Sleep(time.Second)
	fmt.Println("request http://localhost:2999")
	resp, err := http.Get("http://localhost:2999/abc/")
	if err != nil {
		fmt.Println("error get data:", err.Error())
		return ""
	}
	defer resp.Body.Close()

	body, errBody := io.ReadAll(resp.Body)
	if errBody != nil {
		fmt.Println("error read body", errBody.Error())
		return ""
	}

	return string(body)
}

func (app *application) startServer(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup) {
	defer wg.Done()

	ln, errListen := net.Listen("tcp", app.address)
	if errListen != nil {
		fmt.Println("error listen master address", errListen.Error())
		cancel()
		return
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.handler)

	server := &http.Server{
		Handler: mux,
	}
	app.server = server

	wg.Add(1)
	go app.runServerMaster(cancel, wg, server, ln)

	<-ctx.Done()
	//app.stopServerMaster(cancel)
}

func (app *application) runServerMaster(cancel context.CancelFunc, wg *sync.WaitGroup, server *http.Server, ln net.Listener) {
	defer wg.Done()
	fmt.Println("run master server address", ln.Addr().String())
	errServe := server.Serve(ln)
	if errServe != nil && !errors.Is(errServe, http.ErrServerClosed) {
		fmt.Println("error serve master server", errServe.Error())
		cancel()
	}
}

func (app *application) stopServerMaster(cancel context.CancelFunc) {
	defer cancel()
	fmt.Println("stop master server")
	errShutdown := app.server.Shutdown(context.Background())
	if errShutdown != nil {
		fmt.Println("error shutdown master server", errShutdown.Error())
	}
}

func (app *application) handler(rw http.ResponseWriter, req *http.Request) {
	_, err := rw.Write([]byte(app.uin))
	if err != nil {
		fmt.Println("error handler", err.Error())
	}
}
