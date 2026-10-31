package openaicompatible

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/icosmos-space/llmatrix/engine/types"
)

func TestPromptCacheKey(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		if body["prompt_cache_key"] != "session-1" {
			t.Fatalf("body=%v", body)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id": "1", "model": "gpt-4o-mini",
			"choices": []map[string]any{{
				"index": 0, "message": map[string]any{"role": "assistant", "content": "ok"}, "finish_reason": "stop",
			}},
			"usage": map[string]int{"prompt_tokens": 1, "completion_tokens": 1, "total_tokens": 2},
		})
	}))
	defer srv.Close()

	p := New("openai", "k", srv.URL+"/v1")
	resp, err := p.ChatCompletion(context.Background(), &types.ChatRequest{
		Model:          "gpt-4o-mini",
		PromptCacheKey: "session-1",
		Messages:       []types.Message{{Role: types.RoleUser, Content: "hi"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Choices[0].Message.Content != "ok" {
		t.Fatalf("%+v", resp)
	}
}
