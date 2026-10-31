package openaicompatible

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/icosmos-space/llmatrix/engine/types"
)

func TestResponses(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/responses" {
			t.Fatalf("path=%s", r.URL.Path)
		}
		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		if body["input"] != "hello" {
			t.Fatalf("input=%v", body["input"])
		}
		_ = json.NewEncoder(w).Encode(types.ResponsesResponse{
			ID:     "resp_1",
			Object: "response",
			Status: "completed",
			Model:  "gpt-4o-mini",
			Output: []types.ResponsesOutputItem{{
				Type: "message", Role: "assistant", Status: "completed",
				Content: []types.ResponsesContentPart{{Type: "output_text", Text: "hi there"}},
			}},
			Usage: types.ResponsesUsage{InputTokens: 2, OutputTokens: 2, TotalTokens: 4},
		})
	}))
	defer srv.Close()

	p := New("openai", "k", srv.URL+"/v1")
	resp, err := p.Responses(context.Background(), &types.ResponsesRequest{
		Model: "gpt-4o-mini",
		Input: "hello",
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.OutputText() != "hi there" {
		t.Fatalf("text=%q", resp.OutputText())
	}
}
