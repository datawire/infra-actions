package main

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"log"
	"net/http"
)

func init() {
	config = &Config{}
	if err := envconfig.Process("", config); err != nil {
		log.Fatalf("Error loading environment configuration: %v", err)
	}
}

func main() {
	makeHandler := func(name string) http.Handler {
		mux := http.NewServeMux()
		mux.HandleFunc("/", handleRequest)
		return mux
	}

	addr := ":8080"
	fmt.Println("Started GiHub provisioner")
	if err := http.ListenAndServe(addr, makeHandler("main")); err != nil {
		log.Println(err)
	}
}
