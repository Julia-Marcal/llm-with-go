package main

import (
	"context"
	"fmt"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main(ctx context.Context, clientTransport mcp.Transport) {
	client := mcp.NewClient(&mcp.Implementation{Name: "client"}, nil)
	clientSession, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      "llm-caller",
		Arguments: map[string]any{"name": "Give me 3 names for a golang project?"},
	})
	if err != nil {
		log.Fatal(err)
	}

	if len(res.Content) > 0 {
		if textContent, ok := res.Content[0].(*mcp.TextContent); ok {
			fmt.Println("LLM answer:", textContent.Text)
		}
	}

	clientSession.Close()
}
