// https://habr.com/ru/company/southbridge/blog/455290/
// https://hamsterden.ru/prometheus-deleting-time-series-metrics/
// https://prometheus.io/docs/prometheus/latest/querying/api/

// https://eax.me/golang-prometheus-metrics/
// https://pkg.go.dev/github.com/VictoriaMetrics/metrics#example-Counter-Vec
package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/VictoriaMetrics/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

var (
	testMetrics = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "test_metrics",
		Help: "test metrics",
	}, []string{"id"})

	countMetrics = metrics.NewCounter("counter")
)

func victoriaMetrics(w http.ResponseWriter, req *http.Request) {
	metrics.WritePrometheus(w, false)
}

func Register() {
	prometheus.MustRegister(testMetrics)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	Register()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	ln, errLn := net.Listen("tcp", "127.0.0.1:2021")
	if errLn != nil {
		log.Printf("error listen address, %s", errLn.Error())
	}
	defer ln.Close()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go run(ctx, cancel, wg, ln)

	wg.Add(1)
	go countAdd(ctx, wg)

	<-ctx.Done()
	wg.Wait()
	log.Println("done")
}

func countAdd(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			id := rand.Intn(3)
			testMetrics.WithLabelValues(strconv.Itoa(id)).Inc()

			v1 := rand.Intn(3)
			v2 := rand.Intn(3)

			name := fmt.Sprintf(`metrics{v1="%d", v2="%d", id="%d"}`, v1, v2, id)
			metrics.GetOrCreateCounter(name).Inc()

		case <-ctx.Done():
			return
		}
		countMetrics.Inc()
	}
}

func run(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup, ln net.Listener) {
	defer wg.Done()

	mux := http.NewServeMux()
	mux.HandleFunc("/liveness", handlerLiveness)
	mux.HandleFunc("/vm", victoriaMetrics)
	mux.Handle("/metrics", promhttp.Handler())

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

func handlerLiveness(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("ok"))
}
