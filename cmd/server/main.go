package main

import (
	"context"

	"github.com/andrei-cloud/go-devops/internal/server"
	"github.com/rs/zerolog/log"
)

func main() {
	s := server.NewServer()

	ctx, cancel := context.WithCancel(context.Background())
	s.Run(ctx)

	s.Shutdown() //blocking function
	cancel()
	log.Info().Msg("server quit")
}
