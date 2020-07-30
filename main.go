package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	version = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "version",
		Help: "Version information about this binary",
		ConstLabels: map[string]string{
			"version": "v0.1.0",
		},
	})
	httpRequestsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Count of all HTTP requests",
	}, []string{"code", "method"})
)

func main() {
	bind := ""
	num := 0
	flagset := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flagset.StringVar(&bind, "bind", ":8080", "The socket to bind to.")
	flagset.IntVar(&num, "num", 100, "The number of additional metrics to export.")
	flagset.Parse(os.Args[1:])

	r := prometheus.NewRegistry()
	r.MustRegister(httpRequestsTotal)
	r.MustRegister(version)

	if num > 0 {
		for i := 0; i < num; i++ {
			metric := prometheus.NewCounter(prometheus.CounterOpts{
				Name: fmt.Sprintf("prom_example_counter_%d", i),
				Help: "additional counter",
			})
			metric.Add(rand.Float64()*100)
			r.MustRegister(metric)
		}
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello from example application."))
	})
	notfound := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	mux := http.NewServeMux()
	mux.HandleFunc("/", promhttp.InstrumentHandlerCounter(httpRequestsTotal, handler))
	mux.HandleFunc("/err", promhttp.InstrumentHandlerCounter(httpRequestsTotal, notfound))
	mux.HandleFunc("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{}).ServeHTTP)

	// serve kube-state-metrics
	mux.HandleFunc("/ksm", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/bin/kube_state_metrics")
	})

	s := &http.Server{
		Addr:           ":8443",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1048576
	}
	log.Fatal(s.ListenAndServeTLS("/ssl/prom-example.pem", "/ssl/prom-example.key"))
}
