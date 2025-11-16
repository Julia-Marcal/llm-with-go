package main

import (
	"flag"
	"log/slog"

	"github.com/julia-marcal/llm-with-go/internal/server"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	transport := flag.String("transport", "tcp", "transport type: in-memory, tcp, unix")
	addr := flag.String("addr", ":35461", "address for tcp/unix transport (e.g. :12345 or /tmp/mcp.sock)")
	flag.Parse()

	tempSrv := mcp.NewServer(&mcp.Implementation{Name: "brain", Version: "v0.0.1"}, nil)
	server.SetActiveServer(tempSrv)

	clientTransport, serverTransport, err := server.NewTransport(*transport, *addr)
	if err != nil {
		slog.Error("Failed to create transport", slog.Any("error", err))
		return
	}

	server.Start(serverTransport, clientTransport)
}
