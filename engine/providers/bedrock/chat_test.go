package bedrock

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/icosmos-space/llmatrix/engine/types"
)

func TestChatCompletionConverse(t *testing.T) {
	sawAuth := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.Header.Get("Authorization"), "AWS4-HMAC-SHA256") {
			t.Fatalf("auth=%q", r.Header.Get("Authorization"))
		}
		sawAuth = true
		body, _ := io.ReadAll(r.Body)
		var req map[string]any
		_ = json.Unmarshal(body, &req)
		msgs, _ := req["messages"].([]any)
		if len(msgs) != 1 {
			t.Fatalf("messages=%v", req["messages"])
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"output": map[string]any{
				"message": map[string]any{
					"role":    "assistant",
					"content": []map[string]any{{"text": "hello"}},
				},
			},
			"stopReason": "end_turn",
			"usage":      map[string]any{"inputTokens": 1, "outputTokens": 2, "totalTokens": 3},
		})
	}))
	defer srv.Close()

	p := New(
		WithCredentials("AKIATEST", "secretsecretsecretsecretsecret12", ""),
		WithRegion("us-east-1"),
		WithHTTPClient(srv.Client()),
	)
	p.RuntimeBaseURL = srv.URL

	resp, err := p.ChatCompletion(context.Background(), &types.ChatRequest{
		Model:    "anthropic.claude-3-haiku-20240307-v1:0",
		Messages: []types.Message{{Role: types.RoleUser, Content: "hi"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !sawAuth {
		t.Fatal("expected signed request")
	}
	if resp.Choices[0].Message.Content != "hello" || resp.Provider != "bedrock" {
		t.Fatalf("%+v", resp)
	}
}
