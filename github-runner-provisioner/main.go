package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"time"
)

var start = time.Now()
var count int64

func main() {
	addr := ":8080"
	msg := "Hello World!"
	if len(os.Args) > 1 {
		addr = os.Args[1]
		words := os.Args[2:]
		if len(words) > 0 {
			msg = strings.Join(words, " ")
		}
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.String())
		requests := atomic.AddInt64(&count, 1)
		since := time.Since(start)
		w.Write([]byte(fmt.Sprintf("%s (RemoteAddr: %s, Up: %s, Requests: %d)\n", msg, r.RemoteAddr, since.String(),
			requests)))
		w.Write([]byte(fmt.Sprintf("Method: %v\n", r.Method)))
		w.Write([]byte(fmt.Sprintf("Headers: %v\n", len(r.Header))))
		for k, v := range r.Header {
			w.Write([]byte(fmt.Sprintf("  %s: %v\n", k, v)))
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		w.Write([]byte(fmt.Sprintf("Body: %s\n", string(body))))
	})

	log.Printf("Starting on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Println(err)
	}
}
