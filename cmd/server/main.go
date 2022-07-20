package main

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/andrei-cloud/go-devops/internal/server"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	if buildVersion == "" {
		buildVersion = "N/A"
	}
	if buildDate == "" {
		buildDate = "N/A"
	}
	if buildCommit == "" {
		buildCommit = "N/A"
	}
	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)
	s := server.NewServer()

	ctx, cancel := context.WithCancel(context.Background())
	s.Run(ctx)

	s.Shutdown() // blocking function
	cancel()
	log.Info().Msg("server quit")
}
