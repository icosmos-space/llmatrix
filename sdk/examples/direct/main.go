package main

import (
	"context"
	"fmt"
	"log"

	"github.com/icosmos-space/llmatrix/sdk"
	"github.com/icosmos-space/llmatrix/sdk/chat"
)

func main() {
	client := sdk.NewClient(
		sdk.WithProvider("openai"),
		sdk.WithAPIKey("YOUR_API_KEY"),
		sdk.WithModel("gpt-4o"),
	)

	resp, err := client.ChatCompletion(context.Background(), &chat.Request{
		Messages: []chat.Message{
			{Role: "user", Content: "Hello from LLMatrix direct provider example."},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)
}
