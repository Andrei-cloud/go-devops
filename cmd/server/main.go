package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"

	"github.com/andrei-cloud/go-devops/internal/server"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {
	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)
	s := server.NewServer()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	s.Run(ctx)

	s.Shutdown(ctx) // blocking function
	log.Info().Msg("server quit")
}
