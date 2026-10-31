package main

import (
	"encoding/json"
	"fmt"

	"github.com/icosmos-space/llmatrix/sdk/tools"
)

func main() {
	tool := tools.Tool{
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
	}
	b, _ := json.MarshalIndent(tool, "", "  ")
	fmt.Println(string(b))
}
