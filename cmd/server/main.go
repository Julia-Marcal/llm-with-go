package main

import (
	"log/slog"

	"github.com/julia-marcal/llm-with-go/internal/server"
)

func main() {
	serverTransport, clientTransport, err := server.NewTransport("in-memory", "")
	if err != nil {
		slog.Error("Failed to create transport", slog.Any("error", err))
		return
	}

	server.Start(serverTransport, clientTransport)
}
