package sdk

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/icosmos-space/llmatrix/sdk/chat"
)

func TestClientChatCompletionLiteLLMModel(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chat/completions" {
			t.Fatalf("path=%s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id":    "1",
			"model": "deepseek-chat",
			"choices": []map[string]any{
				{
					"index": 0,
					"message": map[string]any{
						"role":    "assistant",
						"content": "ok",
					},
					"finish_reason": "stop",
				},
			},
			"usage": map[string]int{"prompt_tokens": 1, "completion_tokens": 1, "total_tokens": 2},
		})
	}))
	defer srv.Close()

	client := NewClient(
		WithAPIKey("k"),
		WithBaseURL(srv.URL),
		WithModel("deepseek/deepseek-chat"),
	)
	resp, err := client.ChatCompletion(context.Background(), &chat.Request{
		Messages: []chat.Message{{Role: "user", Content: "hi"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Provider != "deepseek" {
		t.Fatalf("provider=%q", resp.Provider)
	}
	if resp.Choices[0].Message.Content != "ok" {
		t.Fatalf("content=%q", resp.Choices[0].Message.Content)
	}
}
