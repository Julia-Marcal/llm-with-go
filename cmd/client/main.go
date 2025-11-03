package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	ctx := context.Background()
	mode := flag.String("mode", "dial", "transport mode: dial|listen")
	addr := flag.String("addr", "", "transport address (e.g. unix:///tmp/mcp.sock or tcp://host:port)")
	flag.Parse()

	if *addr == "" {
		log.Fatal("missing --addr; implement transport creation from --mode and --addr (ex: --mode=dial --addr=unix:///tmp/mcp.sock)")
	}

	client := mcp.NewClient(&mcp.Implementation{Name: "client"}, nil)

	// TODO: criar o transport real a partir de *mode e *addr.
	// Exemplo de comportamento esperado:
	// - se mode == "dial" -> dial para *addr (tcp/unix)
	// - se mode == "listen" -> escutar em *addr e aceitar conexões
	// Substitua a linha abaixo pela criação concreta do transport usando a SDK.
	clientTransport := (mcp.Transport)(nil)

	if clientTransport == nil {
		log.Fatal("transport not created: implement transport creation from --mode and --addr before connecting")
	}

	clientSession, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      "llm-caller",
		Arguments: map[string]any{"question": "Give me 3 names for a golang project?"},
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
