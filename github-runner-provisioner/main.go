package main

import (
	"fmt"
	"github.com/datawire/infra-actions/github-runner-provisioner/internal/aws"
	"github.com/datawire/infra-actions/github-runner-provisioner/internal/monitoring"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

var ec2Client *aws.Ec2Client
var cfg *Config

func main() {
	cfg = NewConfig()
	ec2Client = aws.NewEc2Client()

	makeHandler := func(name string) http.Handler {
		mux := http.NewServeMux()
		mux.HandleFunc("/", handleProvisioningRequest)
		mux.HandleFunc("/github-runner-provisioner/healthz", handleHealthCheckRequest)
		mux.Handle("/metrics", promhttp.Handler())
		return mux
	}

	go monitoring.UpdateActionRunnersRuntimeMetric()

	addr := ":8080"
	fmt.Println("Started GitHub provisioner")
	if err := http.ListenAndServe(addr, makeHandler("main")); err != nil {
		log.Println(err)
	}
}
