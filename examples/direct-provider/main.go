package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/icosmos-space/llmatrix/engine/engine"
	"github.com/icosmos-space/llmatrix/engine/providers/openai"
	"github.com/icosmos-space/llmatrix/engine/registry"
	"github.com/icosmos-space/llmatrix/engine/types"
)

// Cross-module example: call a provider through Engine directly.
func main() {
	reg := registry.New()
	if err := reg.Register(openai.New(os.Getenv("OPENAI_API_KEY"))); err != nil {
		log.Fatal(err)
	}

	eng := engine.New(reg)
	resp, err := eng.Executor().ChatCompletion(context.Background(), "openai", &types.ChatRequest{
		Model: "gpt-4o-mini",
		Messages: []types.Message{
			{Role: types.RoleUser, Content: "Hello from examples/direct-provider"},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.Choices[0].Message.Content)
}
