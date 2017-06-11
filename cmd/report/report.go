package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/martinhoefling/go-dmarc-report/handler"
	"github.com/martinhoefling/go-dmarc-report/report"
)

func runRepository(reportDir string, errorChannel chan error, wg *sync.WaitGroup) {
	err := report.Repository(reportDir)
	if err != nil {
		errorChannel <- err
	}
	wg.Done()
}

func runServer(listen string, errorChannel chan error, wg *sync.WaitGroup) {
	r := mux.NewRouter()
	r.HandleFunc("/", handler.Index)
	r.HandleFunc("/domains/{domain}/", handler.Domain)
	r.HandleFunc("/domains/{domain}/report/{id}", handler.SingleReport)

	srv := &http.Server{
		Handler:      r,
		Addr:         listen,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	err := srv.ListenAndServe()

	if err != nil {
		errorChannel <- err
	}
	wg.Done()
}

func sigHandler(errorChannel chan error, wg *sync.WaitGroup) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	for range signalChan {
		fmt.Println("\nInterrupt received, stopping service...")
		errorChannel <- nil
	}
	wg.Done()
}

func main() {
	var reportDir, listen string
	flag.StringVar(&reportDir, "reportDir", "reports", "Report directory")
	flag.StringVar(&listen, "listen", "127.0.0.1:8123", "Address and port to listen to")
	flag.Parse()

	var wg sync.WaitGroup
	waitGroupLength := 3
	errChannel := make(chan error, 1)
	wg.Add(waitGroupLength)
	finished := make(chan bool, 1)

	go runRepository(reportDir, errChannel, &wg)
	go runServer(listen, errChannel, &wg)
	go sigHandler(errChannel, &wg)

	go func() {
		wg.Wait()
		close(finished)
	}()

	log.Printf("Listening on http://%s", listen)

	select {
	case <-finished:
	case err := <-errChannel:
		if err != nil {
			log.Fatal(err)
			return
		}
	}

	fmt.Println("Terminated successfully")
}
