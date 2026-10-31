package openaicompatible

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/icosmos-space/llmatrix/engine/types"
)

func TestChatCompletion(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/chat/completions" {
			t.Fatalf("path=%s", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
			t.Fatalf("auth=%q", got)
		}
		var body wireChatRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body.Model != "gpt-4o-mini" || len(body.Messages) != 1 {
			t.Fatalf("bad body: %+v", body)
		}
		_ = json.NewEncoder(w).Encode(wireChatResponse{
			ID:    "chatcmpl-1",
			Model: body.Model,
			Choices: []struct {
				Index        int         `json:"index"`
				Message      wireMessage `json:"message"`
				FinishReason string      `json:"finish_reason"`
			}{
				{Index: 0, Message: wireMessage{Role: "assistant", Content: "hello"}, FinishReason: "stop"},
			},
			Usage: types.Usage{PromptTokens: 3, CompletionTokens: 1, TotalTokens: 4},
		})
	}))
	defer srv.Close()

	p := New("openai", "test-key", srv.URL+"/v1")
	resp, err := p.ChatCompletion(context.Background(), &types.ChatRequest{
		Model: "gpt-4o-mini",
		Messages: []types.Message{
			{Role: types.RoleUser, Content: "hi"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Choices[0].Message.Content != "hello" {
		t.Fatalf("content=%q", resp.Choices[0].Message.Content)
	}
	if resp.Provider != "openai" {
		t.Fatalf("provider=%q", resp.Provider)
	}
}

func TestChatCompletionHTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":{"message":"bad key","type":"auth"}}`))
	}))
	defer srv.Close()

	p := New("openai", "bad", srv.URL+"/v1")
	_, err := p.ChatCompletion(context.Background(), &types.ChatRequest{
		Model:    "gpt-4o-mini",
		Messages: []types.Message{{Role: types.RoleUser, Content: "hi"}},
	})
	if err == nil {
		t.Fatal("expected error")
	}
}
