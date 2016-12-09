package main

import (
	"flag"
	"log"
	"net/http"

	"time"

	"github.com/gorilla/mux"
	"github.com/martinhoefling/go-dmarc-report/handler"
	"github.com/martinhoefling/go-dmarc-report/report"
)

func main() {
	flag.Parse()
	var cmdargs = flag.Args()
	if len(cmdargs) != 1 {
		panic("Report path not specified on cmdline.")
	}
	var reportPath = cmdargs[0]

	go report.Repository(reportPath)

	r := mux.NewRouter()
	r.HandleFunc("/", handler.Index)
	r.HandleFunc("/domains/{domain}/", handler.Domain)

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8123",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
