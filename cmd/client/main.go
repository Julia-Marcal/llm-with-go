package main

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/modelcontextprotocol/go-sdk/jsonrpc"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	ctx := context.Background()
	mode := flag.String("mode", "dial", "transport mode: dial|listen")
	addr := flag.String("addr", "127.0.0.1:35461", "transport address (e.g. unix:///tmp/mcp.sock or tcp://host:port)")
	flag.Parse()

	if *addr == "" {
		log.Fatal("missing --addr; implement transport creation from --mode and --addr (ex: --mode=dial --addr=unix:///tmp/mcp.sock)")
	}

	client := mcp.NewClient(&mcp.Implementation{Name: "client"}, nil)

	clientTransport, err := makeTransport(*mode, *addr)
	if err != nil {
		log.Fatalf("failed to create transport: %v", err)
	}

	if clientTransport == nil {
		log.Fatal("transport not created: implement transport creation from --mode and --addr before connecting")
	}

	clientSession, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      "llm-caller",
		Arguments: map[string]any{"question": "Give me 3 names for a orange cat"},
	})
	if err != nil {
		log.Fatal(err)
	}

	if res != nil && len(res.Content) > 0 {
		if textContent, ok := res.Content[0].(*mcp.TextContent); ok {
			fmt.Println("LLM answer:", textContent.Text)
		} else {
			fmt.Printf("Tool returned non-text content: %#v\n", res.Content[0])
		}
	} else {
		fmt.Println("Tool returned no content")
	}

	if err := clientSession.Close(); err != nil {
		log.Printf("error closing client session: %v", err)
	}
}

func makeTransport(mode, rawAddr string) (mcp.Transport, error) {
	network, addr, err := parseAddr(rawAddr)
	if err != nil {
		return nil, err
	}
	switch mode {
	case "dial":
		return &netConnTransport{mode: "dial", network: network, address: addr}, nil
	case "listen":
		return &netConnTransport{mode: "listen", network: network, address: addr}, nil
	default:
		return nil, fmt.Errorf("unknown mode %q: must be dial or listen", mode)
	}
}

func parseAddr(raw string) (network, addr string, err error) {
	if strings.Contains(raw, "://") {
		u, err := url.Parse(raw)
		if err != nil {
			return "", "", err
		}
		scheme := u.Scheme
		switch scheme {
		case "tcp", "unix":
			if scheme == "tcp" {
				return "tcp", u.Host, nil
			}
			return "unix", u.Path, nil
		default:
			return "", "", fmt.Errorf("unsupported scheme: %s", scheme)
		}
	}
	return "tcp", raw, nil
}

type netConnTransport struct {
	mode    string
	network string
	address string
}

func (t *netConnTransport) Connect(ctx context.Context) (mcp.Connection, error) {
	switch t.mode {
	case "dial":
		var d net.Dialer
		log.Printf("client dialing %s %s", t.network, t.address)
		conn, err := d.DialContext(ctx, t.network, t.address)
		if err != nil {
			return nil, err
		}
		return newNetConnection(conn), nil
	case "listen":
		log.Printf("client listening on %s %s", t.network, t.address)
		ln, err := net.Listen(t.network, t.address)
		if err != nil {
			return nil, err
		}
		acceptCh := make(chan net.Conn, 1)
		errCh := make(chan error, 1)
		go func() {
			c, err := ln.Accept()
			if err != nil {
				errCh <- err
				return
			}
			acceptCh <- c
		}()
		select {
		case <-ctx.Done():
			_ = ln.Close()
			return nil, ctx.Err()
		case err := <-errCh:
			_ = ln.Close()
			return nil, err
		case conn := <-acceptCh:
			_ = ln.Close()
			return newNetConnection(conn), nil
		}
	default:
		return nil, fmt.Errorf("unsupported transport mode: %s", t.mode)
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
