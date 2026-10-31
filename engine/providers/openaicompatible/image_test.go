package openaicompatible

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/icosmos-space/llmatrix/engine/types"
)

func TestImageGeneration(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/images/generations" {
			t.Fatalf("path=%s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"created": 1,
			"data":    []map[string]any{{"url": "https://example.com/a.png"}},
		})
	}))
	defer srv.Close()

	p := New("openai", "k", srv.URL+"/v1")
	resp, err := p.ImageGeneration(context.Background(), &types.ImageRequest{
		Model:  "dall-e-3",
		Prompt: "a cat",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Data) != 1 || resp.Data[0].URL == "" {
		t.Fatalf("%+v", resp)
	}
}
