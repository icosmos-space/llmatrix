package sdk

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/icosmos-space/llmatrix/engine/retry"
	"github.com/icosmos-space/llmatrix/sdk/chat"
)

func TestChatCompletionFallback(t *testing.T) {
	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		model, _ := body["model"].(string)
		hits.Add(1)
		if model == "primary-model" {
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"error":{"message":"rate limited"}}`))
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id":    "ok",
			"model": model,
			"choices": []map[string]any{{
				"index":         0,
				"message":       map[string]any{"role": "assistant", "content": "fallback-ok"},
				"finish_reason": "stop",
			}},
			"usage": map[string]int{"prompt_tokens": 1, "completion_tokens": 1, "total_tokens": 2},
		})
	}))
	defer srv.Close()

	client := NewClient(
		WithAPIKey("k"),
		WithBaseURL(srv.URL+"/v1"),
		WithModel("openai/primary-model"),
		WithFallback("openai/backup-model"),
		WithRetry(retry.Policy{MaxAttempts: 1, Initial: 0, MaxBackoff: 0, Multiplier: 1}),
	)
	resp, err := client.ChatCompletion(context.Background(), &chat.Request{
		Messages: []chat.Message{{Role: "user", Content: "hi"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Choices[0].Message.Content != "fallback-ok" {
		t.Fatalf("content=%q", resp.Choices[0].Message.Content)
	}
	if hits.Load() < 2 {
		t.Fatalf("hits=%d", hits.Load())
	}
}
