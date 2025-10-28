package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"github.com/julia-marcal/llm-with-go/internal/infrastructures"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	ctx := context.Background()
	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	server := mcp.NewServer(&mcp.Implementation{Name: "brain", Version: "v0.0.1"}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "llm-caller",
		Description: "Answer a question using the LLM",
	}, askQuestion)

	serverSession, err := server.Connect(ctx, serverTransport, nil)
	if err != nil {
		log.Fatal(err)
	}
	client(ctx, clientTransport)

	serverSession.Wait()
}

func client(ctx context.Context, clientTransport mcp.Transport) {
	client := mcp.NewClient(&mcp.Implementation{Name: "client"}, nil)
	clientSession, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      "llm-caller",
		Arguments: map[string]any{"name": "What is Go language?"},
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

type QuestionParams struct {
	Question string `json:"name"`
}

func askQuestion(ctx context.Context, req *mcp.CallToolRequest, args QuestionParams) (*mcp.CallToolResult, any, error) {
	answer, err := executeLLM(ctx, args.Question)
	if err != nil {
		return nil, nil, err
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: answer},
		},
	}, nil, nil
}

func executeLLM(ctx context.Context, promptQuestion string) (string, error) {
	response, err := infrastructures.ExecuteCall(ctx, promptQuestion)
	if err != nil {
		slog.Error("Exceeded available tokens", slog.Any("token_error", err))
		return "", err
	}
	return response, nil
}
