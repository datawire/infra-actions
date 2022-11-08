package main

import (
	"fmt"
	"github.com/datawire/infra-actions/github-runner-provisioner/internal/aws"
	"github.com/datawire/infra-actions/github-runner-provisioner/internal/monitoring"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

var ec2Client *aws.Ec2Client
var cfg = NewConfig()

func main() {
	ec2Client = aws.NewEc2Client()
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
	fmt.Println("Started GiHub provisioner")
	if err := http.ListenAndServe(addr, makeHandler("main")); err != nil {
		log.Println(err)
	}
}
