package sdk

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/icosmos-space/llmatrix/sdk/moderation"
)

func TestModerationPackage(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id": "mod_1", "model": "omni-moderation-latest",
			"results": []map[string]any{{"flagged": true, "categories": map[string]bool{"hate": true}, "category_scores": map[string]float64{"hate": 0.9}}},
		})
	}))
	defer srv.Close()

	Configure(WithAPIKey("k"), WithBaseURL(srv.URL+"/v1"), WithMaxRetries(1))
	t.Cleanup(func() { Configure() })

	resp, err := Moderation(context.Background(), &moderation.Request{Input: "bad"})
	if err != nil {
		t.Fatal(err)
	}
	if !resp.Flagged() {
		t.Fatalf("%+v", resp)
	}
}
