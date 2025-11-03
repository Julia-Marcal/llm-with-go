package server

import (
	"context"
	"log/slog"

	"github.com/julia-marcal/llm-with-go/internal/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func Start(serverTransport, clientTransport mcp.Transport) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv := mcp.NewServer(&mcp.Implementation{Name: "brain", Version: "v0.0.1"}, nil)

	RegisterTools(srv)
	SetActiveServer(srv)

	// If a clientTransport is provided (e.g., in-memory demo), start the in-process client
	if clientTransport != nil {
		go client.Main(ctx, clientTransport)
	}

	serverSession, err := srv.Connect(ctx, serverTransport, nil)
	if err != nil {
		slog.Error("server connect failed", slog.Any("error", err))
		return
	}

	slog.Info("server connected; waiting for session end")
	serverSession.Wait()
}
