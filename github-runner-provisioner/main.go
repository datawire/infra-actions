package main

import (
	"net/http"
	"os"

	"github.com/datawire/infra-actions/github-runner-provisioner/internal/monitoring"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{
		FieldMap: log.FieldMap{
			log.FieldKeyLevel: "severity",
			log.FieldKeyTime:  "timestamp",
			log.FieldKeyMsg:   "message",
		},
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

var cfg *Config

func main() {
	cfg = NewConfig()

	makeHandler := func(name string) http.Handler {
		mux := http.NewServeMux()
		mux.HandleFunc("/", handleProvisioningRequest)
		mux.HandleFunc("/github-runner-provisioner/healthz", handleHealthCheckRequest)
		mux.Handle("/metrics", promhttp.Handler())
		return mux
	}

	go monitoring.UpdateActionRunnersRuntimeMetric()

	addr := ":8080"
	log.Infof("Started GitHub provisioner. Listening on %s", addr)
	if err := http.ListenAndServe(addr, makeHandler("main")); err != nil {
		log.Fatal(err)
	}
}
