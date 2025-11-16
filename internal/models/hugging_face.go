package models

import (
	"flag"
	"log"
	"os"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/huggingface"
)

func ConfigureHuggingFace() (*float64, llms.Model) {
	apiKey := os.Getenv("HUGGING_FACE_KEY")
	hgModel := os.Getenv("HUGGING_FACE_MODEL")

	model := flag.String("model", hgModel, "Model to use")
	temperature := flag.Float64("temp", 0.8, "Temperature for response")
	flag.Parse()

	clientOptions := []huggingface.Option{
		huggingface.WithToken(apiKey),
		huggingface.WithModel(*model),
	}

	llm, err := huggingface.New(clientOptions...)
	if err != nil {
		log.Fatal(err)
	}

	return temperature, llm
}
