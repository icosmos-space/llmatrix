package azure

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/icosmos-space/llmatrix/engine/types"
)

func TestChatCompletion(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/openai/deployments/gpt-4o-mini/chat/completions") {
			t.Fatalf("path=%s", r.URL.Path)
		}
		if r.URL.Query().Get("api-version") == "" {
			t.Fatal("missing api-version")
		}
		if r.Header.Get("api-key") != "azure-key" {
			t.Fatalf("api-key=%q", r.Header.Get("api-key"))
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id": "1", "model": "gpt-4o-mini",
			"choices": []map[string]any{{
				"index": 0, "message": map[string]any{"role": "assistant", "content": "azure-ok"}, "finish_reason": "stop",
			}},
			"usage": map[string]int{"prompt_tokens": 1, "completion_tokens": 1, "total_tokens": 2},
		})
	}))
	defer srv.Close()

	p := New("azure-key", WithBaseURL(srv.URL), WithAPIVersion("2024-08-01-preview"))
	resp, err := p.ChatCompletion(context.Background(), &types.ChatRequest{
		Model: "gpt-4o-mini",
		Messages: []types.Message{
			{Role: types.RoleUser, Content: "hi"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Choices[0].Message.Content != "azure-ok" {
		t.Fatalf("%+v", resp)
	}
}
