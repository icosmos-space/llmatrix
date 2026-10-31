package openaicompatible

import (
	"context"
	"net/url"
	"strconv"

	"github.com/icosmos-space/llmatrix/engine/errors"
	"github.com/icosmos-space/llmatrix/engine/types"
)

// CreateFineTuningJob implements provider.FineTuningProvider.
func (p *Provider) CreateFineTuningJob(ctx context.Context, req *types.FineTuningJobCreateRequest) (*types.FineTuningJob, error) {
	if req == nil || req.Model == "" || req.TrainingFile == "" {
		return nil, errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "model and training_file are required", nil)
	}
	var out types.FineTuningJob
	if err := p.doJSON(ctx, "POST", "/fine_tuning/jobs", req, &out); err != nil {
		return nil, err
	}
	out.Provider = p.Info().Name
	return &out, nil
}

// RetrieveFineTuningJob implements provider.FineTuningProvider.
func (p *Provider) RetrieveFineTuningJob(ctx context.Context, jobID string) (*types.FineTuningJob, error) {
	if jobID == "" {
		return nil, errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "job id is required", nil)
	}
	var out types.FineTuningJob
	if err := p.doJSON(ctx, "GET", "/fine_tuning/jobs/"+url.PathEscape(jobID), nil, &out); err != nil {
		return nil, err
	}
	out.Provider = p.Info().Name
	return &out, nil
}

// ListFineTuningJobs implements provider.FineTuningProvider.
func (p *Provider) ListFineTuningJobs(ctx context.Context, after string, limit int) (*types.FineTuningJobList, error) {
	path := "/fine_tuning/jobs"
	q := url.Values{}
	if after != "" {
		q.Set("after", after)
	}
	if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}
	if enc := q.Encode(); enc != "" {
		path += "?" + enc
	}
	var out types.FineTuningJobList
	if err := p.doJSON(ctx, "GET", path, nil, &out); err != nil {
		return nil, err
	}
	out.Provider = p.Info().Name
	return &out, nil
}

// CancelFineTuningJob implements provider.FineTuningProvider.
func (p *Provider) CancelFineTuningJob(ctx context.Context, jobID string) (*types.FineTuningJob, error) {
	if jobID == "" {
		return nil, errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "job id is required", nil)
	}
	var out types.FineTuningJob
	if err := p.doJSON(ctx, "POST", "/fine_tuning/jobs/"+url.PathEscape(jobID)+"/cancel", map[string]any{}, &out); err != nil {
		return nil, err
	}
	out.Provider = p.Info().Name
	return &out, nil
}
