package proxy

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/icosmos-space/llmatrix/sdk/chat"
)

func TestProxyClientChatCompletion(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/chat/completions" {
			t.Fatalf("path=%s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer proxy-key" {
			t.Fatalf("auth=%q", r.Header.Get("Authorization"))
		}
		_ = json.NewEncoder(w).Encode(chat.Response{
			ID:       "1",
			Model:    "gpt-4o-mini",
			Provider: "openai",
			Choices: []chat.Choice{{
				Message: chat.Message{Role: "assistant", Content: "from-proxy"},
			}},
		})
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "proxy-key", nil)
	resp, err := c.ChatCompletion(context.Background(), &chat.Request{
		Model:    "openai/gpt-4o-mini",
		Messages: []chat.Message{{Role: "user", Content: "hi"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Choices[0].Message.Content != "from-proxy" {
		t.Fatalf("content=%q", resp.Choices[0].Message.Content)
	}
}
