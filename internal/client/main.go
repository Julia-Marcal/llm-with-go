package client

import (
	"context"
	"fmt"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func Main(ctx context.Context, clientTransport mcp.Transport) {
	client := mcp.NewClient(&mcp.Implementation{Name: "client"}, nil)
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
