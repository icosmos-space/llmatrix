package openaicompatible

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/icosmos-space/llmatrix/engine/types"
)

func TestImageEdit(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/images/edits" {
			t.Fatalf("path=%s", r.URL.Path)
		}
		if !strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
			t.Fatalf("ct=%s", r.Header.Get("Content-Type"))
		}
		_, _ = io.ReadAll(r.Body)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"created": 1,
			"data":    []map[string]any{{"b64_json": "aaa"}},
		})
	}))
	defer srv.Close()

	p := New("openai", "k", srv.URL+"/v1")
	n := 1
	resp, err := p.ImageEdit(context.Background(), &types.ImageEditRequest{
		Model:  "dall-e-2",
		Prompt: "add a hat",
		Image:  []byte("fakepng"),
		N:      &n,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Data) != 1 || resp.Data[0].B64JSON != "aaa" {
		t.Fatalf("%+v", resp)
	}
}
