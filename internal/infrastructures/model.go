package infrastructures

import (
	"context"
	"time"

	"github.com/julia-marcal/llm-with-go/internal/dto"
	"github.com/julia-marcal/llm-with-go/internal/models"
	"github.com/tmc/langchaingo/llms"
)

func ExecuteCall(ctx context.Context, promptQuestion string) (string, error) {
	prompt, temperature, llm := models.ConfigureOpenRouter(promptQuestion)

	start := time.Now()

	response, err := executePrompt(ctx, llm, prompt, temperature)
	duration := time.Since(start)

	audit := dto.LLMAudit{
		Prompt:      *prompt,
		Temperature: *temperature,
		Response:    response,
		Err:         err,
		Duration:    duration,
		Timestamp:   time.Now(),
	}
	LogAudit(&audit)

	return response, err
}

func executePrompt(ctx context.Context, llm llms.Model, prompt *string, temp *float64) (string, error) {
	opts := []llms.CallOption{
		llms.WithTemperature(*temp),
	}

	response, err := llms.GenerateFromSinglePrompt(ctx, llm, *prompt, opts...)

	if err != nil {
		return "", err
	}

	return response, nil
}
