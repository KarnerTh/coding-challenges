package main

import (
	"log"
	"log/slog"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/api"
)

const (
	ListenAddress = ":8080"
	LogLevel      = slog.LevelDebug
)

func main() {
	slog.SetLogLoggerLevel(LogLevel)
	server := api.NewServer(ListenAddress)

	if err := server.Run(); err != nil {
		log.Fatal("Could not start server on ", ListenAddress)
	}
}
