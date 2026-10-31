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

func TestFilesAndBatches(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/v1/files":
			if !strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
				t.Fatalf("ct=%s", r.Header.Get("Content-Type"))
			}
			_, _ = io.ReadAll(r.Body)
			_ = json.NewEncoder(w).Encode(types.FileObject{ID: "file_1", Purpose: "batch", Filename: "a.jsonl"})
		case r.Method == http.MethodGet && r.URL.Path == "/v1/files/file_1":
			_ = json.NewEncoder(w).Encode(types.FileObject{ID: "file_1", Purpose: "batch"})
		case r.Method == http.MethodPost && r.URL.Path == "/v1/batches":
			_ = json.NewEncoder(w).Encode(types.Batch{ID: "batch_1", Status: "validating", InputFileID: "file_1"})
		case r.Method == http.MethodGet && r.URL.Path == "/v1/batches/batch_1":
			_ = json.NewEncoder(w).Encode(types.Batch{ID: "batch_1", Status: "completed"})
		case r.Method == http.MethodPost && r.URL.Path == "/v1/batches/batch_1/cancel":
			_ = json.NewEncoder(w).Encode(types.Batch{ID: "batch_1", Status: "cancelling"})
		default:
			t.Fatalf("%s %s", r.Method, r.URL.Path)
		}
	}))
	defer srv.Close()

	p := New("openai", "k", srv.URL+"/v1")
	f, err := p.CreateFile(context.Background(), &types.FileCreateRequest{
		File: []byte(`{"a":1}`), FileName: "a.jsonl", Purpose: "batch",
	})
	if err != nil || f.ID != "file_1" {
		t.Fatalf("%v %+v", err, f)
	}
	got, err := p.GetFile(context.Background(), "file_1")
	if err != nil || got.ID != "file_1" {
		t.Fatalf("%v %+v", err, got)
	}
	b, err := p.CreateBatch(context.Background(), &types.BatchCreateRequest{
		CompletionWindow: "24h", Endpoint: "/v1/chat/completions", InputFileID: "file_1",
	})
	if err != nil || b.ID != "batch_1" {
		t.Fatalf("%v %+v", err, b)
	}
	rb, err := p.RetrieveBatch(context.Background(), "batch_1")
	if err != nil || rb.Status != "completed" {
		t.Fatalf("%v %+v", err, rb)
	}
	cb, err := p.CancelBatch(context.Background(), "batch_1")
	if err != nil || cb.Status != "cancelling" {
		t.Fatalf("%v %+v", err, cb)
	}
}
