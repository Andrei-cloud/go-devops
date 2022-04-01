package main

import "github.com/andrei-cloud/go-devops/internal/server"

func main() {
	s := server.NewServer()

	s.Run()

	s.Shutdown() //blocking function
}
