package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/andrei-cloud/go-devops/internal/handlers"
	"github.com/andrei-cloud/go-devops/internal/storage/inmem"
)

func main1() {

	repo := inmem.New()

	mux := http.NewServeMux()
	mux.HandleFunc("/update/gauge/", handlers.Gauges(repo))
	mux.HandleFunc("/update/counter/", handlers.Counters(repo))
	mux.HandleFunc("/update/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	})

	s := &http.Server{
		Addr:           ":8080",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go s.ListenAndServe()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	<-sig

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		fmt.Println(err)
	}
	//fmt.Println("server quit")
}
