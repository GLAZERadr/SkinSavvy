package main

import (
	"context"
	"log"

	"github.com/tmc/langchaingo/llms/openai"
)

func main() {
	llm, err := openai.New()
	if err != nil {
		log.Fatal(err)
	}
	prompt := "What would be a good company name for a company that makes colorful socks?"
	completion, err := llm.Call(context.Background(), prompt)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(completion)
}