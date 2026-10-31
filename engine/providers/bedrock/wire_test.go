package bedrock

import (
	"encoding/json"
	"testing"

	"github.com/icosmos-space/llmatrix/engine/types"
)

func TestToConverseRequestSystemAndTools(t *testing.T) {
	max := 100
	temp := 0.2
	req := &types.ChatRequest{
		Model: "anthropic.claude-3-haiku-20240307-v1:0",
		Messages: []types.Message{
			{Role: types.RoleSystem, Content: "be brief"},
			{Role: types.RoleUser, Content: "hi"},
		},
		MaxTokens:   &max,
		Temperature: &temp,
		Tools: []types.Tool{{
			Type: "function",
			Function: types.ToolFunction{
				Name:        "lookup",
				Description: "look up",
				Parameters:  map[string]any{"type": "object"},
			},
		}},
	}
	out, err := toConverseRequest(req)
	if err != nil {
		t.Fatal(err)
	}
	if len(out.System) != 1 || out.System[0].Text != "be brief" {
		t.Fatalf("system=%+v", out.System)
	}
	if len(out.Messages) != 1 || out.Messages[0].Role != "user" {
		t.Fatalf("messages=%+v", out.Messages)
	}
	if out.InferenceConfig == nil || out.InferenceConfig.MaxTokens == nil || *out.InferenceConfig.MaxTokens != 100 {
		t.Fatalf("inference=%+v", out.InferenceConfig)
	}
	if out.ToolConfig == nil || len(out.ToolConfig.Tools) != 1 {
		t.Fatalf("tools=%+v", out.ToolConfig)
	}
	var schema map[string]any
	if err := json.Unmarshal(out.ToolConfig.Tools[0].ToolSpec.InputSchema, &schema); err != nil {
		t.Fatal(err)
	}
	if _, ok := schema["json"]; !ok {
		t.Fatalf("schema=%v", schema)
	}
}

func TestFromConverseResponseToolUse(t *testing.T) {
	raw := &converseResponse{
		StopReason: "tool_use",
		Usage:      usage{InputTokens: 3, OutputTokens: 4},
		Output: converseOutput{Message: converseMessage{
			Role: "assistant",
			Content: []contentBlock{{
				ToolUse: &toolUse{
					ToolUseID: "tu1",
					Name:      "lookup",
					Input:     json.RawMessage(`{"q":"x"}`),
				},
			}},
		}},
	}
	resp := fromConverseResponse("m", raw)
	if resp.Choices[0].FinishReason != "tool_calls" {
		t.Fatalf("finish=%q", resp.Choices[0].FinishReason)
	}
	if len(resp.Choices[0].Message.ToolCalls) != 1 || resp.Choices[0].Message.ToolCalls[0].Name != "lookup" {
		t.Fatalf("toolcalls=%+v", resp.Choices[0].Message.ToolCalls)
	}
	if resp.Usage.TotalTokens != 7 {
		t.Fatalf("usage=%+v", resp.Usage)
	}
}
