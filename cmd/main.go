package main

import (
	"context"
	"fmt"
	"log"

	"github.com/julia-marcal/llm-with-go/internal/infrastructures"
)

func main() {
	ctx := context.Background()

	response, err := infrastructures.ExecuteCall(ctx)
	if err != nil {
		fmt.Println("Exceeded available tokens")
		log.Fatal(err)
	}

	fmt.Println(response)
}
