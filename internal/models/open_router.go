package models

import (
	"flag"
	"log"
	"os"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

func ConfigureOpenRouter(promptQuestion string) (*string, *float64, llms.Model) {
	model := flag.String("model", "meta-llama/llama-3.2-3b-instruct:free", "Model to use")
	prompt := flag.String("prompt", promptQuestion, "Prompt being sent")
	temperature := flag.Float64("temp", 0.8, "Temperature for response")
	flag.Parse()

	apiKey := os.Getenv("OPENROUTER_API_KEY")

	llm, err := openai.New(
		openai.WithModel(*model),
		openai.WithBaseURL("https://openrouter.ai/api/v1"),
		openai.WithToken(apiKey),
	)
	if err != nil {
		log.Fatal(err)
	}

	return prompt, temperature, llm
}
