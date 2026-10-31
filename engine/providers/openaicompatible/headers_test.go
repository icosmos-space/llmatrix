package openaicompatible

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/icosmos-space/llmatrix/engine/types"
)

func TestExtraHeadersAndOrg(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("OpenAI-Organization") != "org_1" {
			t.Fatalf("org=%q", r.Header.Get("OpenAI-Organization"))
		}
		if r.Header.Get("X-Custom") != "1" {
			t.Fatalf("custom=%q", r.Header.Get("X-Custom"))
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
	p.Organization = "org_1"
	p.ExtraHeaders = map[string]string{"X-Custom": "1"}
	_, err := p.ChatCompletion(context.Background(), &types.ChatRequest{
		Model: "gpt-4o-mini",
		Messages: []types.Message{
			{Role: types.RoleUser, Content: "hi"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
}
