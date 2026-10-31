package gemini

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
		if !strings.Contains(r.URL.Path, "/models/gemini-2.0-flash:generateContent") {
			t.Fatalf("path=%s", r.URL.Path)
		}
		if r.URL.Query().Get("key") != "test-key" {
			t.Fatalf("key=%q", r.URL.Query().Get("key"))
		}
		var body wireRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body.SystemInstruction == nil || len(body.Contents) != 1 {
			t.Fatalf("bad body %+v", body)
		}
		_ = json.NewEncoder(w).Encode(wireResponse{
			Candidates: []wireCandidate{{
				Content:      wireContent{Role: "model", Parts: []wirePart{{Text: "hi gemini"}}},
				FinishReason: "STOP",
			}},
			UsageMetadata: &wireUsage{PromptTokenCount: 4, CandidatesTokenCount: 2, TotalTokenCount: 6},
		})
	}))
	defer srv.Close()

	p := New("test-key", WithBaseURL(srv.URL))
	resp, err := p.ChatCompletion(context.Background(), &types.ChatRequest{
		Model: "gemini-2.0-flash",
		Messages: []types.Message{
			{Role: types.RoleSystem, Content: "brief"},
			{Role: types.RoleUser, Content: "hello"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Choices[0].Message.Content != "hi gemini" {
		t.Fatalf("content=%q", resp.Choices[0].Message.Content)
	}
	if resp.Choices[0].FinishReason != "stop" {
		t.Fatalf("finish=%q", resp.Choices[0].FinishReason)
	}
}

func TestChatCompletionStream(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "streamGenerateContent") {
			t.Fatalf("path=%s", r.URL.Path)
		}
		if r.URL.Query().Get("alt") != "sse" {
			t.Fatalf("alt=%q", r.URL.Query().Get("alt"))
		}
		w.Header().Set("Content-Type", "text/event-stream")
		flusher := w.(http.Flusher)
		events := []string{
			`data: {"candidates":[{"content":{"parts":[{"text":"Hel"}]}}]}`,
			`data: {"candidates":[{"content":{"parts":[{"text":"lo"}]},"finishReason":"STOP"}]}`,
		}
		for _, e := range events {
			_, _ = w.Write([]byte(e + "\n\n"))
			flusher.Flush()
		}
	}))
	defer srv.Close()

	p := New("k", WithBaseURL(srv.URL))
	s, err := p.ChatCompletionStream(context.Background(), &types.ChatRequest{
		Model:    "gemini-2.0-flash",
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
