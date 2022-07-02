package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"

	"github.com/andrei-cloud/go-devops/internal/agent"
	"github.com/andrei-cloud/go-devops/internal/collector"
)

func main() {
	collector := collector.NewCollector()
	agent := agent.NewAgent(collector, nil)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
		<-sig

		cancel()
	}()

	agent.Run(ctx)
	log.Info().Msg("agent quit")
}
