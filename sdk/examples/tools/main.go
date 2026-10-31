package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/icosmos-space/llmatrix/sdk"
	"github.com/icosmos-space/llmatrix/sdk/chat"
	"github.com/icosmos-space/llmatrix/sdk/tools"
)

func main() {
	client := sdk.NewClient(
		sdk.WithAPIKey(os.Getenv("OPENAI_API_KEY")),
		sdk.WithModel(envOr("LLMATRIX_MODEL", "openai/gpt-4o-mini")),
	)

	resp, err := client.ChatCompletion(context.Background(), &chat.Request{
		Messages: []chat.Message{
			{Role: "user", Content: "What's the weather in San Francisco?"},
		},
		Tools: []tools.Tool{{
			Type: "function",
			Function: tools.Function{
				Name:        "get_weather",
				Description: "Get weather by city",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"city": map[string]any{"type": "string"},
					},
					"required": []string{"city"},
				},
			},
		}},
		ToolChoice: "auto",
	})
	if err != nil {
		log.Fatal(err)
	}
	b, _ := json.MarshalIndent(resp.Choices[0].Message, "", "  ")
	fmt.Println(string(b))
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
