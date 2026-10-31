package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/icosmos-space/llmatrix/sdk/chat"
	proxyclient "github.com/icosmos-space/llmatrix/sdk/proxy"
)

// Cross-module example: call a running LLMatrix Proxy via SDK.
func main() {
	base := envOr("LLMATRIX_PROXY_URL", "http://localhost:8080")
	key := os.Getenv("LLMATRIX_API_KEY")
	client := proxyclient.NewClient(base, key, nil)
	resp, err := client.ChatCompletion(context.Background(), &chat.Request{
		Model: envOr("LLMATRIX_MODEL", "openai/gpt-4o-mini"),
		Messages: []chat.Message{
			{Role: "user", Content: "Hello from examples/proxy-client"},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.Choices[0].Message.Content)
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
