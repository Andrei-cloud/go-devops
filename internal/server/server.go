package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/andrei-cloud/go-devops/internal/router"

	"github.com/caarlos0/env"
	"github.com/go-chi/chi"
)

var cfg Config

type Config struct {
	Address string `env:"ADDRESS" envDefault:":8080"`
}

type server struct {
	r *chi.Mux
	s *http.Server
}

func init() {
	cfg = Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}
}

func NewServer() *server {
	srv := server{}
	srv.r = router.SetupRouter()
	srv.s = &http.Server{
		Addr:           cfg.Address,
		Handler:        srv.r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	return &srv
}

func (srv *server) Run() {
	go srv.s.ListenAndServe()
}

func (srv *server) Shutdown() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	<-sig

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.s.Shutdown(ctx); err != nil {
		fmt.Println(err)
	}
}
