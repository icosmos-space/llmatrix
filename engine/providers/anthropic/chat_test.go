package anthropic

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
		if r.URL.Path != "/v1/messages" {
			t.Fatalf("path=%s", r.URL.Path)
		}
		if r.Header.Get("x-api-key") != "test-key" {
			t.Fatalf("missing api key")
		}
		if r.Header.Get("anthropic-version") == "" {
			t.Fatal("missing anthropic-version")
		}
		var body wireRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body.System != "be brief" || len(body.Messages) != 1 {
			t.Fatalf("bad body %+v", body)
		}
		_ = json.NewEncoder(w).Encode(wireResponse{
			ID:    "msg_1",
			Model: body.Model,
			Role:  "assistant",
			Content: []wireContent{
				{Type: "text", Text: "hello from claude"},
			},
			StopReason: "end_turn",
			Usage: struct {
				InputTokens  int `json:"input_tokens"`
				OutputTokens int `json:"output_tokens"`
			}{InputTokens: 5, OutputTokens: 3},
		})
	}))
	defer srv.Close()

	p := New("test-key", WithBaseURL(srv.URL))
	resp, err := p.ChatCompletion(context.Background(), &types.ChatRequest{
		Model: "claude-sonnet-4-5",
		Messages: []types.Message{
			{Role: types.RoleSystem, Content: "be brief"},
			{Role: types.RoleUser, Content: "hi"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Choices[0].Message.Content != "hello from claude" {
		t.Fatalf("content=%q", resp.Choices[0].Message.Content)
	}
	if resp.Choices[0].FinishReason != "stop" {
		t.Fatalf("finish=%q", resp.Choices[0].FinishReason)
	}
	if resp.Usage.TotalTokens != 8 {
		t.Fatalf("usage=%+v", resp.Usage)
	}
}

func TestChatCompletionStream(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		flusher := w.(http.Flusher)
		events := []string{
			`data: {"type":"content_block_delta","delta":{"type":"text_delta","text":"Hel"}}`,
			`data: {"type":"content_block_delta","delta":{"type":"text_delta","text":"lo"}}`,
			`data: {"type":"message_stop"}`,
		}
		for _, e := range events {
			_, _ = w.Write([]byte(e + "\n\n"))
			flusher.Flush()
		}
	}))
	defer srv.Close()

	p := New("k", WithBaseURL(srv.URL))
	s, err := p.ChatCompletionStream(context.Background(), &types.ChatRequest{
		Model:    "claude-sonnet-4-5",
		Messages: []types.Message{{Role: types.RoleUser, Content: "hi"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	var b strings.Builder
	for s.Next(context.Background()) {
		ch := s.Chunk()
		if ch.Done {
			break
		}
		b.WriteString(ch.Delta)
	}
	if err := s.Err(); err != nil {
		t.Fatal(err)
	}
	if b.String() != "Hello" {
		t.Fatalf("got %q", b.String())
	}
}
