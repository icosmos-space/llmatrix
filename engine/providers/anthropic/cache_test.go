package anthropic

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/icosmos-space/llmatrix/engine/types"
)

func TestPromptCacheControl(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body wireRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		sys, ok := body.System.([]any)
		if !ok || len(sys) != 1 {
			t.Fatalf("system=%T %+v", body.System, body.System)
		}
		block, _ := sys[0].(map[string]any)
		cc, _ := block["cache_control"].(map[string]any)
		if block["text"] != "long system" || cc["type"] != "ephemeral" {
			t.Fatalf("system block=%+v", block)
		}
		if len(body.Messages) != 1 {
			t.Fatalf("messages=%+v", body.Messages)
		}
		_ = json.NewEncoder(w).Encode(wireResponse{
			ID: "msg_1", Role: "assistant", Model: body.Model,
			Content:    []wireContent{{Type: "text", Text: "ok"}},
			StopReason: "end_turn",
		})
	}))
	defer srv.Close()

	p := New("k", WithBaseURL(srv.URL))
	resp, err := p.ChatCompletion(context.Background(), &types.ChatRequest{
		Model: "claude-sonnet-4-5",
		Messages: []types.Message{
			{
				Role:         types.RoleSystem,
				Content:      "long system",
				CacheControl: &types.CacheControl{Type: "ephemeral"},
			},
			{Role: types.RoleUser, Content: "hi"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Choices[0].Message.Content != "ok" {
		t.Fatalf("%q", resp.Choices[0].Message.Content)
	}
}
