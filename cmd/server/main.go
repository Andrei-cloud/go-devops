package main

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/andrei-cloud/go-devops/internal/server"
)

func main() {
	s := server.NewServer()

	ctx, cancel := context.WithCancel(context.Background())
	s.Run(ctx)

	s.Shutdown() // blocking function
	cancel()
	log.Info().Msg("server quit")
}
