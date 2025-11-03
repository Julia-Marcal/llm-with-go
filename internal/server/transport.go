package server

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var activeServer *mcp.Server

func SetActiveServer(server *mcp.Server) {
	activeServer = server
}

func NewTransport(transportType, address string) (mcp.Transport, mcp.Transport, error) {
	switch transportType {

	case "in-memory":
		clientTransport, serverTransport := mcp.NewInMemoryTransports()
		return clientTransport, serverTransport, nil

	case "unix", "tcp":
		if activeServer == nil {
			return nil, nil, fmt.Errorf("active MCP server is nil; call SetActiveServer first")
		}

		handler := mcp.NewStreamableHTTPHandler(
			func(r *http.Request) *mcp.Server { return activeServer },
			&mcp.StreamableHTTPOptions{},
		)

		if transportType == "unix" {
			_ = os.Remove(address)
			l, err := net.Listen("unix", address)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to listen on unix socket: %w", err)
			}
			_ = os.Chmod(address, 0o600)
			go func() { _ = http.Serve(l, handler) }()
		} else {
			srv := &http.Server{Addr: address, Handler: handler}
			go func() { _ = srv.ListenAndServe() }()
		}

		return nil, nil, nil

	default:
		return nil, nil, fmt.Errorf("unsupported transport type: %s", transportType)
	}
}
