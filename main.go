package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const ENDPOINT_TAG string = "endpoint"
const STATUS_TAG string = "status"

var HOSTNAME string = getEnv("HOSTNAME", "kubmetheusDefaultHost")

func main() {
	http.ListenAndServe("0.0.0.0:8000", kubHandler())
}

var hits = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "ticks_counter",
	Help: "the number of tick calls",
	ConstLabels: prometheus.Labels{
		"HOSTNAME": HOSTNAME,
	},
}, []string{ENDPOINT_TAG, STATUS_TAG})

type LoggingMetricsServeMux struct {
	*http.ServeMux
}

type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

func (sr *StatusRecorder) WriteHeader(statusCode int) {
	sr.Status = statusCode
	sr.ResponseWriter.WriteHeader(statusCode)
}

func (mux *LoggingMetricsServeMux) Handle(pattern string, handler http.Handler) {
	wrapped := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		statusRecordedRw := &StatusRecorder{rw, 200}
		handler.ServeHTTP(statusRecordedRw, req)
		log.Printf("host: %s: request completed on pathioso %s", HOSTNAME, pattern)
		hits.With(prometheus.Labels{
			ENDPOINT_TAG: pattern,
			STATUS_TAG:   fmt.Sprint(statusRecordedRw.Status),
		}).Inc()
	})
	mux.ServeMux.Handle(pattern, wrapped)
}

func (mux *LoggingMetricsServeMux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	if handler == nil {
		panic("http: nil handler")
	}
	mux.Handle(pattern, http.HandlerFunc(handler))
}

func kubHandler() http.Handler {
	sm := LoggingMetricsServeMux{http.NewServeMux()}
	sm.HandleFunc("/tick", tickHandler)
	sm.Handle("/metrics", promhttp.Handler())
	return sm
}

func tickHandler(r http.ResponseWriter, req *http.Request) {
	r.WriteHeader(200)
	fmt.Fprintf(r, "we ticked %s", "things")
}

func getEnv(key string, fallback string) string {
	e := os.Getenv(key)
	if len(e) == 0 {
		return fallback
	}
	return e
}
