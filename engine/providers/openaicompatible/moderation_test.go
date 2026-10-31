package openaicompatible

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/icosmos-space/llmatrix/engine/types"
)

func TestModeration(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/moderations" {
			t.Fatalf("path=%s", r.URL.Path)
		}
		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		if body["input"] != "hello" {
			t.Fatalf("body=%v", body)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id": "mod_1", "model": "omni-moderation-latest",
			"results": []map[string]any{{
				"flagged":         false,
				"categories":      map[string]bool{"hate": false, "violence": false},
				"category_scores": map[string]float64{"hate": 0.01, "violence": 0.02},
			}},
		})
	}))
	defer srv.Close()

	p := New("openai", "k", srv.URL+"/v1")
	resp, err := p.Moderation(context.Background(), &types.ModerationRequest{
		Model: "omni-moderation-latest",
		Input: "hello",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Results) != 1 || resp.Results[0].Flagged {
		t.Fatalf("%+v", resp)
	}
}
