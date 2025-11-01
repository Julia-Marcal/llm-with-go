package server

import (
	"context"
	"log/slog"

	cli "github.com/julia-marcal/llm-with-go/internal/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func Start(serverTransport, clientTransport mcp.Transport) {
	ctx := context.Background()
	srv := mcp.NewServer(&mcp.Implementation{Name: "brain", Version: "v0.0.1"}, nil)

	RegisterTools(srv)
	SetActiveServer(srv)

	serverSession, err := srv.Connect(ctx, serverTransport, nil)
	if err != nil {
		slog.Error("server connect failed", slog.Any("error", err))
		return
	}
	cli.Main(ctx, clientTransport)

	serverSession.Wait()
}
