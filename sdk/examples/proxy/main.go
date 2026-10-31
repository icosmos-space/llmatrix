package main

import (
	"context"
	"fmt"
	"log"

	"github.com/icosmos-space/llmatrix/sdk/chat"
	proxyclient "github.com/icosmos-space/llmatrix/sdk/proxy"
)

func main() {
	client := proxyclient.NewClient("http://localhost:8080", "YOUR_PROXY_API_KEY", nil)

	resp, err := client.ChatCompletion(context.Background(), &chat.Request{
		Messages: []chat.Message{
			{Role: "user", Content: "Hello from LLMatrix proxy example."},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)
}
