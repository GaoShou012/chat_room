package utils

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

func Prometheus(addr string) {
	if addr == "" {
		return
	}
	http.Handle("/metrics", promhttp.Handler())
	log.Fatalln(http.ListenAndServe(addr, nil))
}
