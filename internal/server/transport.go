package server

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"sync"
	"time"

	"github.com/modelcontextprotocol/go-sdk/jsonrpc"
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

		l, err := net.Listen(transportTypeToNetwork(transportType), address)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to listen on %s: %w", address, err)
		}

		slog.Info("server listening", slog.String("network", transportTypeToNetwork(transportType)), slog.String("address", l.Addr().String()))

		if transportType == "unix" {
			_ = os.Chmod(address, 0o600)
		}

		serverTransport := &listenerTransport{ln: l}
		return nil, serverTransport, nil

	default:
		return nil, nil, fmt.Errorf("unsupported transport type: %s", transportType)
	}
}

func transportTypeToNetwork(t string) string {
	if t == "tcp" {
		return "tcp"
	}
	return "unix"
}

type listenerTransport struct {
	ln net.Listener
}

func (t *listenerTransport) Connect(ctx context.Context) (mcp.Connection, error) {
	acceptCh := make(chan net.Conn, 1)
	errCh := make(chan error, 1)
	go func() {
		c, err := t.ln.Accept()
		if err != nil {
			errCh <- err
			return
		}
		acceptCh <- c
	}()

	select {
	case <-ctx.Done():
		_ = t.ln.Close()
		return nil, ctx.Err()
	case err := <-errCh:
		_ = t.ln.Close()
		return nil, err
	case conn := <-acceptCh:
		return newNetConnection(conn), nil
	}
}

type netConnection struct {
	conn       net.Conn
	r          *bufio.Reader
	writeMu    sync.Mutex
	closedOnce sync.Once
	closedErr  error
	sessionID  string
}

func newNetConnection(c net.Conn) mcp.Connection {
	return &netConnection{
		conn:      c,
		r:         bufio.NewReader(c),
		sessionID: c.RemoteAddr().String(),
	}
}

func (c *netConnection) SessionID() string { return c.sessionID }

func (c *netConnection) Close() error {
	c.closedOnce.Do(func() { c.closedErr = c.conn.Close() })
	return c.closedErr
}

func (c *netConnection) Read(ctx context.Context) (jsonrpc.Message, error) {
	readCh := make(chan struct {
		msg jsonrpc.Message
		err error
	}, 1)
	go func() {
		line, err := c.r.ReadBytes('\n')
		if err != nil {
			readCh <- struct {
				msg jsonrpc.Message
				err error
			}{nil, err}
			return
		}
		msg, err := jsonrpc.DecodeMessage(bytes.TrimSpace(line))
		readCh <- struct {
			msg jsonrpc.Message
			err error
		}{msg, err}
	}()

	select {
	case <-ctx.Done():
		_ = c.conn.SetReadDeadline(time.Now())
		res := <-readCh
		if res.err != nil {
			return nil, ctx.Err()
		}
		return res.msg, nil
	case res := <-readCh:
		_ = c.conn.SetReadDeadline(time.Time{})
		return res.msg, res.err
	}
}

func (c *netConnection) Write(ctx context.Context, msg jsonrpc.Message) error {
	data, err := jsonrpc.EncodeMessage(msg)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	if dl, ok := ctx.Deadline(); ok {
		_ = c.conn.SetWriteDeadline(dl)
		defer func() { _ = c.conn.SetWriteDeadline(time.Time{}) }()
	}
	_, err = c.conn.Write(data)
	return err
}
