package provider

import (
	"context"

	"github.com/icosmos-space/llmatrix/engine/types"
)

// FineTuningProvider manages OpenAI fine-tuning jobs (LiteLLM fine_tuning helpers).
type FineTuningProvider interface {
	Provider
	CreateFineTuningJob(ctx context.Context, req *types.FineTuningJobCreateRequest) (*types.FineTuningJob, error)
	RetrieveFineTuningJob(ctx context.Context, jobID string) (*types.FineTuningJob, error)
	ListFineTuningJobs(ctx context.Context, after string, limit int) (*types.FineTuningJobList, error)
	CancelFineTuningJob(ctx context.Context, jobID string) (*types.FineTuningJob, error)
}
