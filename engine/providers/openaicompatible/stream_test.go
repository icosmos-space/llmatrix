package openaicompatible

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/icosmos-space/llmatrix/engine/types"
)

func TestChatCompletionStream(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		flusher, ok := w.(http.Flusher)
		if !ok {
			t.Fatal("no flusher")
		}
		chunks := []string{
			`data: {"id":"1","choices":[{"index":0,"delta":{"content":"Hel"}}]}`,
			`data: {"id":"1","choices":[{"index":0,"delta":{"content":"lo"}}]}`,
			`data: [DONE]`,
		}
		for _, c := range chunks {
			_, _ = w.Write([]byte(c + "\n\n"))
			flusher.Flush()
		}
	}))
	defer srv.Close()

	p := New("openai", "k", srv.URL+"/v1")
	s, err := p.ChatCompletionStream(context.Background(), &types.ChatRequest{
		Model:    "gpt-4o-mini",
		Messages: []types.Message{{Role: types.RoleUser, Content: "hi"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	var b strings.Builder
	ctx := context.Background()
	for s.Next(ctx) {
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
