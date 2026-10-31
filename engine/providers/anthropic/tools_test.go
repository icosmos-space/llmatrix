package anthropic

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/icosmos-space/llmatrix/engine/types"
)

func TestChatCompletionToolUse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body wireRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if len(body.Tools) != 1 || body.Tools[0].Name != "get_weather" {
			t.Fatalf("tools=%+v", body.Tools)
		}
		if body.Tools[0].InputSchema == nil {
			t.Fatal("missing input_schema")
		}
		_ = json.NewEncoder(w).Encode(wireResponse{
			ID:    "msg_tools",
			Model: body.Model,
			Content: []wireContent{
				{Type: "tool_use", ID: "toolu_1", Name: "get_weather", Input: json.RawMessage(`{"city":"SF"}`)},
			},
			StopReason: "tool_use",
			Usage: struct {
				InputTokens  int `json:"input_tokens"`
				OutputTokens int `json:"output_tokens"`
			}{InputTokens: 8, OutputTokens: 4},
		})
	}))
	defer srv.Close()

	p := New("k", WithBaseURL(srv.URL))
	resp, err := p.ChatCompletion(context.Background(), &types.ChatRequest{
		Model: "claude-sonnet-4-5",
		Messages: []types.Message{
			{Role: types.RoleUser, Content: "weather?"},
		},
		Tools: []types.Tool{{
			Type: "function",
			Function: types.ToolFunction{
				Name:        "get_weather",
				Description: "weather",
				Parameters:  map[string]any{"type": "object"},
			},
		}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Choices[0].FinishReason != "tool_calls" {
		t.Fatalf("finish=%q", resp.Choices[0].FinishReason)
	}
	if len(resp.Choices[0].Message.ToolCalls) != 1 || resp.Choices[0].Message.ToolCalls[0].Name != "get_weather" {
		t.Fatalf("%+v", resp.Choices[0].Message.ToolCalls)
	}
}
