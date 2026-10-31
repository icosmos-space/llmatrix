package sdk

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/icosmos-space/llmatrix/sdk/chat"
)

func TestPackageCompletion(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id": "1", "model": "gpt-4o-mini",
			"choices": []map[string]any{{
				"index": 0, "message": map[string]any{"role": "assistant", "content": "hi"}, "finish_reason": "stop",
			}},
			"usage": map[string]int{"prompt_tokens": 1, "completion_tokens": 1, "total_tokens": 2},
		})
	}))
	defer srv.Close()

	Configure(WithAPIKey("k"), WithBaseURL(srv.URL+"/v1"), WithModel("gpt-4o-mini"), WithMaxRetries(1))
	t.Cleanup(func() { Configure() })

	resp, err := Completion(context.Background(), &chat.Request{
		Messages: []chat.Message{{Role: "user", Content: "hi"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Choices[0].Message.Content != "hi" {
		t.Fatalf("%+v", resp)
	}
}
