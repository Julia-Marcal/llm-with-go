package models

import (
	"flag"
	"log"
	"os"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

func ConfigureOpenRouter() (*float64, llms.Model) {
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	orModel := os.Getenv("OPENROUTER_MODEL")

	model := flag.String("model", orModel, "Model to use")
	temperature := flag.Float64("temp", 0.8, "Temperature for response")
	flag.Parse()

	llm, err := openai.New(
		openai.WithModel(*model),
		openai.WithBaseURL("https://openrouter.ai/api/v1"),
		openai.WithToken(apiKey),
	)
	if err != nil {
		log.Fatal(err)
	}

	return temperature, llm
}
