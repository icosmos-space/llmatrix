package openaicompatible

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/icosmos-space/llmatrix/engine/types"
)

func TestTextCompletionStream(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/completions" {
			t.Fatalf("path=%s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = w.Write([]byte("data: {\"choices\":[{\"text\":\"hel\"}]}\n\n"))
		_, _ = w.Write([]byte("data: {\"choices\":[{\"text\":\"lo\",\"finish_reason\":\"stop\"}]}\n\n"))
		_, _ = w.Write([]byte("data: [DONE]\n\n"))
	}))
	defer srv.Close()

	p := New("openai", "k", srv.URL+"/v1")
	s, err := p.TextCompletionStream(context.Background(), &types.TextCompletionRequest{
		Model: "gpt-3.5-turbo-instruct", Prompt: "hi",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	var got string
	for s.Next(context.Background()) {
		ch := s.Chunk()
		if ch.Done {
			break
		}
		got += ch.Delta
	}
	if got != "hello" {
		t.Fatalf("%q", got)
	}
}
