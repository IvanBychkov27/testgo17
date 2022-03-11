// https://docs.victoriametrics.com/Single-server-VictoriaMetrics.html
// https://pkg.go.dev/github.com/VictoriaMetrics/metrics#example-Counter-Vec
package main

import (
	"bytes"
	"context"
	"errors"
	"github.com/VictoriaMetrics/metrics"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	countMetrics = metrics.NewCounter("counter")
	gaugeMetrics = metrics.NewGauge("gauge", data)
	//setMetrics = metrics.NewSet()
	summaryMetrics = metrics.NewSummary("summary")
	//sMetrics = metrics.Gauge{}
)

func vm(w http.ResponseWriter, req *http.Request) {
	countMetrics.Inc()

	gaugeMetrics.Get()

	//summaryMetrics.Update(0.15)
	metrics.WritePrometheus(w, false)
}

func data() float64 {
	d := rand.Intn(10000)
	return float64(d)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	ln, errLn := net.Listen("tcp", "127.0.0.1:2025")
	if errLn != nil {
		log.Printf("error listen address, %s", errLn.Error())
	}
	defer ln.Close()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go run(ctx, cancel, wg, ln)

	w := &bytes.Buffer{}
	metrics.WritePrometheus(w, true)

	<-ctx.Done()
	wg.Wait()
	log.Println("done")

}

func run(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup, ln net.Listener) {
	defer wg.Done()

	mux := http.NewServeMux()
	mux.HandleFunc("/liveness", handlerLiveness)
	mux.HandleFunc("/metrics", vm)

	server := &http.Server{
		Handler: mux,
	}

	wg.Add(1)
	go runServer(cancel, wg, server, ln)

	<-ctx.Done()

	log.Println("stop server", ln.Addr().String())
	errShutdown := server.Shutdown(context.Background())
	if errShutdown != nil {
		log.Println("error shutdown control server", errShutdown.Error())
	}
}

func runServer(cancel context.CancelFunc, wg *sync.WaitGroup, server *http.Server, ln net.Listener) {
	defer wg.Done()

	log.Println("run server address", ln.Addr().String())

	errServe := server.Serve(ln)
	if errServe != nil && !errors.Is(errServe, http.ErrServerClosed) {
		log.Println("error server", errServe.Error())
		cancel()
	}
}

func handlerLiveness(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}
