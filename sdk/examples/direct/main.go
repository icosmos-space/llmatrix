package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/icosmos-space/llmatrix/sdk"
	"github.com/icosmos-space/llmatrix/sdk/chat"
)

func main() {
	// LiteLLM-style model: "provider/model"
	model := envOr("LLMATRIX_MODEL", "openai/gpt-4o-mini")
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("DEEPSEEK_API_KEY")
	}

	client := sdk.NewClient(
		sdk.WithAPIKey(apiKey),
		sdk.WithModel(model),
	)

	resp, err := client.ChatCompletion(context.Background(), &chat.Request{
		Messages: []chat.Message{
			{Role: "user", Content: "Reply with one short sentence: what is LLMatrix?"},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	if len(resp.Choices) == 0 {
		log.Fatal("empty choices")
	}
	fmt.Printf("[%s/%s] %s\n", resp.Provider, resp.Model, resp.Choices[0].Message.Content)
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
