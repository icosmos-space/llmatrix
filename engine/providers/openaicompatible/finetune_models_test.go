package openaicompatible

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/icosmos-space/llmatrix/engine/types"
)

func TestFineTuningAndListModels(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/v1/fine_tuning/jobs":
			var body types.FineTuningJobCreateRequest
			_ = json.NewDecoder(r.Body).Decode(&body)
			if body.Model != "gpt-4o-mini" || body.TrainingFile != "file_1" {
				t.Fatalf("body=%+v", body)
			}
			_ = json.NewEncoder(w).Encode(types.FineTuningJob{ID: "ftjob_1", Status: "queued", Model: body.Model})
		case r.Method == http.MethodGet && r.URL.Path == "/v1/fine_tuning/jobs/ftjob_1":
			_ = json.NewEncoder(w).Encode(types.FineTuningJob{ID: "ftjob_1", Status: "running"})
		case r.Method == http.MethodGet && r.URL.Path == "/v1/fine_tuning/jobs":
			if r.URL.Query().Get("limit") != "2" {
				t.Fatalf("query=%v", r.URL.Query())
			}
			_ = json.NewEncoder(w).Encode(types.FineTuningJobList{
				Object: "list",
				Data:   []types.FineTuningJob{{ID: "ftjob_1"}},
			})
		case r.Method == http.MethodPost && r.URL.Path == "/v1/fine_tuning/jobs/ftjob_1/cancel":
			_ = json.NewEncoder(w).Encode(types.FineTuningJob{ID: "ftjob_1", Status: "cancelled"})
		case r.Method == http.MethodGet && r.URL.Path == "/v1/models":
			_ = json.NewEncoder(w).Encode(types.ModelList{
				Object: "list",
				Data:   []types.Model{{ID: "gpt-4o-mini", OwnedBy: "openai"}},
			})
		default:
			t.Fatalf("%s %s", r.Method, r.URL.Path)
		}
	}))
	defer srv.Close()

	p := New("openai", "k", srv.URL+"/v1")
	job, err := p.CreateFineTuningJob(context.Background(), &types.FineTuningJobCreateRequest{
		Model: "gpt-4o-mini", TrainingFile: "file_1",
	})
	if err != nil || job.ID != "ftjob_1" || job.Provider != "openai" {
		t.Fatalf("%v %+v", err, job)
	}
	got, err := p.RetrieveFineTuningJob(context.Background(), "ftjob_1")
	if err != nil || got.Status != "running" {
		t.Fatalf("%v %+v", err, got)
	}
	list, err := p.ListFineTuningJobs(context.Background(), "", 2)
	if err != nil || len(list.Data) != 1 {
		t.Fatalf("%v %+v", err, list)
	}
	cancelled, err := p.CancelFineTuningJob(context.Background(), "ftjob_1")
	if err != nil || cancelled.Status != "cancelled" {
		t.Fatalf("%v %+v", err, cancelled)
	}
	models, err := p.ListModels(context.Background())
	if err != nil || len(models.Data) != 1 || models.Data[0].ID != "gpt-4o-mini" {
		t.Fatalf("%v %+v", err, models)
	}
	if models.Data[0].Provider != "openai" {
		t.Fatalf("provider=%q", models.Data[0].Provider)
	}
}
