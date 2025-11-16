package server

import (
	"context"
	"log/slog"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func Start(serverTransport, clientTransport mcp.Transport) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv := mcp.NewServer(&mcp.Implementation{Name: "brain", Version: "v0.0.1"}, nil)

	RegisterTools(srv)
	SetActiveServer(srv)

	if clientTransport != nil {
		slog.Info("client transport provided; start an external client to connect to the server")
	}

	serverSession, err := srv.Connect(ctx, serverTransport, nil)
	if err != nil {
		slog.Error("server connect failed", slog.Any("error", err))
		return
	}

	slog.Info("server connected; waiting for session end")
	serverSession.Wait()
}
