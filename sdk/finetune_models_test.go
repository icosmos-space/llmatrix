package sdk_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/icosmos-space/llmatrix/sdk"
	"github.com/icosmos-space/llmatrix/sdk/finetuning"
)

func TestFineTuningAndListModelsSDK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/v1/fine_tuning/jobs":
			var body map[string]any
			_ = json.NewDecoder(r.Body).Decode(&body)
			if body["model"] != "gpt-4o-mini" {
				t.Fatalf("model=%v", body["model"])
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"id": "ftjob_1", "status": "queued"})
		case r.Method == http.MethodGet && r.URL.Path == "/v1/models":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"object": "list",
				"data":   []map[string]any{{"id": "gpt-4o-mini"}},
			})
		default:
			t.Fatalf("%s %s", r.Method, r.URL.Path)
		}
	}))
	defer srv.Close()

	client := sdk.NewClient(
		sdk.WithAPIKey("k"),
		sdk.WithBaseURL(srv.URL+"/v1"),
		sdk.WithProvider("openai"),
	)
	job, err := client.CreateFineTuningJob(context.Background(), &finetuning.CreateJobRequest{
		Model: "openai/gpt-4o-mini", TrainingFile: "file_1",
	})
	if err != nil || job.ID != "ftjob_1" {
		t.Fatalf("%v %+v", err, job)
	}
	models, err := client.ListModels(context.Background(), "openai/gpt-4o-mini")
	if err != nil || len(models.Data) != 1 {
		t.Fatalf("%v %+v", err, models)
	}
}
