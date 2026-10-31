package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/icosmos-space/llmatrix/sdk"
	"github.com/icosmos-space/llmatrix/sdk/chat"
	"github.com/icosmos-space/llmatrix/sdk/streaming"
)

func main() {
	client := sdk.NewClient(
		sdk.WithAPIKey(os.Getenv("OPENAI_API_KEY")),
		sdk.WithModel(envOr("LLMATRIX_MODEL", "openai/gpt-4o-mini")),
	)

	stream, err := client.ChatCompletionStream(context.Background(), &chat.Request{
		Messages: []chat.Message{
			{Role: "user", Content: "Write a 2-line haiku about Go."},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer stream.Close()

	ctx := context.Background()
	for stream.Next(ctx) {
		ev := stream.Event()
		switch ev.Type {
		case streaming.EventTypeDelta:
			fmt.Print(ev.Delta)
		case streaming.EventTypeDone:
			fmt.Println()
		}
	}
	if err := stream.Err(); err != nil {
		log.Fatal(err)
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
