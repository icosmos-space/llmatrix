package sdk

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/icosmos-space/llmatrix/sdk/responses"
)

func TestResponsesNative(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id": "resp_1", "object": "response", "status": "completed", "model": "gpt-4o-mini",
			"output": []map[string]any{{
				"type": "message", "role": "assistant", "status": "completed",
				"content": []map[string]any{{"type": "output_text", "text": "story"}},
			}},
			"usage": map[string]int{"input_tokens": 1, "output_tokens": 1, "total_tokens": 2},
		})
	}))
	defer srv.Close()

	client := NewClient(WithAPIKey("k"), WithBaseURL(srv.URL+"/v1"), WithModel("openai/gpt-4o-mini"), WithMaxRetries(1))
	resp, err := client.Responses(context.Background(), &responses.Request{Input: "tell a story"})
	if err != nil {
		t.Fatal(err)
	}
	if resp.OutputText() != "story" {
		t.Fatalf("%q", resp.OutputText())
	}
}

func TestResponsesStreamAndStore(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/v1/responses" && !strings.Contains(r.URL.Path, "cancel"):
			w.Header().Set("Content-Type", "text/event-stream")
			_, _ = w.Write([]byte("data: {\"type\":\"response.output_text.delta\",\"delta\":\"ab\"}\n\n"))
			_, _ = w.Write([]byte("data: {\"type\":\"response.completed\"}\n\n"))
		case r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/responses/resp_x"):
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id": "resp_x", "status": "completed",
				"output": []map[string]any{{
					"type": "message", "role": "assistant",
					"content": []map[string]any{{"type": "output_text", "text": "got"}},
				}},
			})
		case r.Method == http.MethodDelete && strings.HasSuffix(r.URL.Path, "/responses/resp_x"):
			w.WriteHeader(http.StatusOK)
		default:
			t.Fatalf("%s %s", r.Method, r.URL.Path)
		}
	}))
	defer srv.Close()

	client := NewClient(WithAPIKey("k"), WithBaseURL(srv.URL+"/v1"), WithModel("openai/gpt-4o-mini"), WithMaxRetries(1))
	stream, err := client.ResponsesStream(context.Background(), &responses.Request{Input: "x"})
	if err != nil {
		t.Fatal(err)
	}
	var text string
	for stream.Next(context.Background()) {
		ev := stream.Event()
		if ev.Type == "done" {
			break
		}
		text += ev.Delta
	}
	_ = stream.Close()
	if text != "ab" {
		t.Fatalf("stream=%q", text)
	}

	got, err := client.GetResponse(context.Background(), "resp_x")
	if err != nil {
		t.Fatal(err)
	}
	if got.OutputText() != "got" {
		t.Fatalf("%q", got.OutputText())
	}
	if err := client.DeleteResponse(context.Background(), "resp_x"); err != nil {
		t.Fatal(err)
	}
}
