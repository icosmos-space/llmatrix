package service

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/icosmos-space/llmatrix/engine/retry"
	"github.com/icosmos-space/llmatrix/engine/types"
)

func TestChatServiceComplete(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id":    "1",
			"model": "gpt-4o-mini",
			"choices": []map[string]any{{
				"index":         0,
				"message":       map[string]any{"role": "assistant", "content": "proxy-ok"},
				"finish_reason": "stop",
			}},
			"usage": map[string]int{"prompt_tokens": 1, "completion_tokens": 1, "total_tokens": 2},
		})
	}))
	defer srv.Close()

	svc := NewChatService()
	svc.Retry = retry.Policy{MaxAttempts: 1}
	svc.APIKeys = map[string]string{"openai": "k"}
	svc.BaseURLs = map[string]string{"openai": srv.URL + "/v1"}

	resp, err := svc.Complete(context.Background(), &types.ChatRequest{
		Model: "openai/gpt-4o-mini",
		Messages: []types.Message{
			{Role: types.RoleUser, Content: "hi"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Choices[0].Message.Content != "proxy-ok" {
		t.Fatalf("content=%q", resp.Choices[0].Message.Content)
	}
	if resp.Provider != "openai" {
		t.Fatalf("provider=%q", resp.Provider)
	}
}
