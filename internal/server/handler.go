package server

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/julia-marcal/llm-with-go/internal/infrastructures"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type QuestionParams struct {
	Question string `json:"question"`
}

func RegisterTools(s *mcp.Server) {
	mcp.AddTool[any, any](s, &mcp.Tool{
		Name:        "ping",
		Description: "Health check. Returns 'pong'",
	}, func(ctx context.Context, req *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, any, error) {
		_ = ctx
		_ = req
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: "pong"}},
		}, nil, nil
	})

	mcp.AddTool[QuestionParams, any](s, &mcp.Tool{
		Name:        "llm-caller",
		Description: "Ask a question to the LLM and return the answer",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args QuestionParams) (*mcp.CallToolResult, any, error) {
		_ = req
		if args.Question == "" {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: "question is required"}},
				IsError: true,
			}, nil, nil
		}

		answer, err := executeLLM(ctx, args.Question)
		if err != nil {
			return nil, nil, fmt.Errorf("llm call failed: %w", err)
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: answer}},
		}, nil, nil
	})
}

func executeLLM(ctx context.Context, promptQuestion string) (string, error) {
	response, err := infrastructures.ExecuteCall(ctx, promptQuestion)
	if err != nil {
		slog.Error("Failed to execute LLM call", slog.Any("error", err))
		return "", err
	}
	return response, nil
}
