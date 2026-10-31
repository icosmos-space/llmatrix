package openaicompatible

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/icosmos-space/llmatrix/engine/types"
)

func TestResponsesStream(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/responses" {
			t.Fatalf("path=%s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = w.Write([]byte("data: {\"type\":\"response.created\"}\n\n"))
		_, _ = w.Write([]byte("data: {\"type\":\"response.output_text.delta\",\"delta\":\"hel\"}\n\n"))
		_, _ = w.Write([]byte("data: {\"type\":\"response.output_text.delta\",\"delta\":\"lo\"}\n\n"))
		_, _ = w.Write([]byte("data: {\"type\":\"response.completed\"}\n\n"))
	}))
	defer srv.Close()

	p := New("openai", "k", srv.URL+"/v1")
	s, err := p.ResponsesStream(context.Background(), &types.ResponsesRequest{
		Model: "gpt-4o-mini",
		Input: "hi",
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
	if err := s.Err(); err != nil {
		t.Fatal(err)
	}
	if got != "hello" {
		t.Fatalf("got=%q", got)
	}
}

func TestResponsesStore(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/responses/resp_1"):
			_ = json.NewEncoder(w).Encode(types.ResponsesResponse{
				ID: "resp_1", Status: "completed",
				Output: []types.ResponsesOutputItem{{
					Type: "message", Role: "assistant",
					Content: []types.ResponsesContentPart{{Type: "output_text", Text: "stored"}},
				}},
			})
		case r.Method == http.MethodDelete && strings.HasSuffix(r.URL.Path, "/responses/resp_1"):
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id":"resp_1","deleted":true}`))
		case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/responses/resp_2/cancel"):
			_ = json.NewEncoder(w).Encode(types.ResponsesResponse{ID: "resp_2", Status: "cancelled"})
		default:
			t.Fatalf("%s %s", r.Method, r.URL.Path)
		}
	}))
	defer srv.Close()

	p := New("openai", "k", srv.URL+"/v1")
	got, err := p.GetResponse(context.Background(), "resp_1")
	if err != nil {
		t.Fatal(err)
	}
	if got.OutputText() != "stored" {
		t.Fatalf("%q", got.OutputText())
	}
	if err := p.DeleteResponse(context.Background(), "resp_1"); err != nil {
		t.Fatal(err)
	}
	cancelled, err := p.CancelResponse(context.Background(), "resp_2")
	if err != nil {
		t.Fatal(err)
	}
	if cancelled.Status != "cancelled" {
		t.Fatalf("%+v", cancelled)
	}
}
