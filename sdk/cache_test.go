package sdk

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/icosmos-space/llmatrix/sdk/chat"
)

func TestChatCaching(t *testing.T) {
	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id": "1", "model": "gpt-4o-mini",
			"choices": []map[string]any{{
				"index": 0, "message": map[string]any{"role": "assistant", "content": "joke"}, "finish_reason": "stop",
			}},
			"usage": map[string]int{"prompt_tokens": 1, "completion_tokens": 1, "total_tokens": 2},
		})
	}))
	defer srv.Close()

	trueVal := true
	client := NewClient(
		WithAPIKey("k"), WithBaseURL(srv.URL+"/v1"), WithModel("openai/gpt-4o-mini"),
		WithMemoryCache(time.Minute), WithMaxRetries(1),
	)
	req := &chat.Request{
		Messages: []chat.Message{{Role: "user", Content: "Tell me a joke."}},
		Caching:  &trueVal,
	}
	r1, err := client.ChatCompletion(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	r2, err := client.ChatCompletion(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if r1.Choices[0].Message.Content != "joke" || r2.Choices[0].Message.Content != "joke" {
		t.Fatalf("%v %v", r1, r2)
	}
	if hits.Load() != 1 {
		t.Fatalf("hits=%d want 1", hits.Load())
	}
}
