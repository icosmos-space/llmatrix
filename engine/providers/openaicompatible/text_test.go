package openaicompatible

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/icosmos-space/llmatrix/engine/types"
)

func TestTextCompletion(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/completions" {
			t.Fatalf("path=%s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id": "cmpl_1", "object": "text_completion", "model": "gpt-3.5-turbo-instruct",
			"choices": []map[string]any{{"index": 0, "text": " this is a test", "finish_reason": "stop"}},
			"usage":   map[string]int{"prompt_tokens": 5, "completion_tokens": 4, "total_tokens": 9},
		})
	}))
	defer srv.Close()

	p := New("openai", "k", srv.URL+"/v1")
	resp, err := p.TextCompletion(context.Background(), &types.TextCompletionRequest{
		Model:  "gpt-3.5-turbo-instruct",
		Prompt: "Say this is a test",
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Choices[0].Text != " this is a test" {
		t.Fatalf("%q", resp.Choices[0].Text)
	}
}
