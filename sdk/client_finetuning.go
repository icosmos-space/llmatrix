package sdk

import (
	"context"
	"fmt"

	"github.com/icosmos-space/llmatrix/engine/provider"
	"github.com/icosmos-space/llmatrix/engine/types"
	"github.com/icosmos-space/llmatrix/sdk/finetuning"
)

func (c *Client) fineTuningProvider(modelHint string) (provider.FineTuningProvider, error) {
	raw := modelHint
	if raw == "" {
		raw = c.cfg.Model
	}
	if raw == "" {
		prov := c.cfg.Provider
		if prov == "" {
			prov = "openai"
		}
		raw = prov + "/default"
	}
	refs := c.candidates(raw)
	if len(refs) == 0 {
		return nil, fmt.Errorf("%w: model is required", ErrInvalidConfig)
	}
	p, err := c.newProvider(refs[0], true)
	if err != nil {
		return nil, err
	}
	fp, ok := p.(provider.FineTuningProvider)
	if !ok {
		return nil, fmt.Errorf("%w: provider %q does not support fine-tuning", ErrUnsupported, refs[0].Provider)
	}
	return fp, nil
}

// CreateFineTuningJob starts a fine-tuning job (LiteLLM create_fine_tuning_job).
func (c *Client) CreateFineTuningJob(ctx context.Context, req *finetuning.CreateJobRequest) (*finetuning.Job, error) {
	if req == nil || req.Model == "" || req.TrainingFile == "" {
		return nil, ErrInvalidRequest
	}
	refs := c.candidates(req.Model)
	if len(refs) == 0 {
		return nil, fmt.Errorf("%w: model is required", ErrInvalidConfig)
	}
	p, err := c.newProvider(refs[0], true)
	if err != nil {
		return nil, MapError(err)
	}
	fp, ok := p.(provider.FineTuningProvider)
	if !ok {
		return nil, fmt.Errorf("%w: provider %q does not support fine-tuning", ErrUnsupported, refs[0].Provider)
	}
	engineReq := toEngineFineTuneCreate(req)
	engineReq.Model = refs[0].Model
	out, err := fp.CreateFineTuningJob(ctx, engineReq)
	if err != nil {
		return nil, MapError(err)
	}
	return fromEngineFineTuneJob(out), nil
}

// RetrieveFineTuningJob fetches a fine-tuning job.
func (c *Client) RetrieveFineTuningJob(ctx context.Context, jobID string, modelHint ...string) (*finetuning.Job, error) {
	if jobID == "" {
		return nil, ErrInvalidRequest
	}
	hint := ""
	if len(modelHint) > 0 {
		hint = modelHint[0]
	}
	fp, err := c.fineTuningProvider(hint)
	if err != nil {
		return nil, err
	}
	out, err := fp.RetrieveFineTuningJob(ctx, jobID)
	if err != nil {
		return nil, MapError(err)
	}
	return fromEngineFineTuneJob(out), nil
}

// ListFineTuningJobs lists fine-tuning jobs.
func (c *Client) ListFineTuningJobs(ctx context.Context, after string, limit int, modelHint ...string) (*finetuning.JobList, error) {
	hint := ""
	if len(modelHint) > 0 {
		hint = modelHint[0]
	}
	fp, err := c.fineTuningProvider(hint)
	if err != nil {
		return nil, err
	}
	out, err := fp.ListFineTuningJobs(ctx, after, limit)
	if err != nil {
		return nil, MapError(err)
	}
	return fromEngineFineTuneList(out), nil
}

// CancelFineTuningJob cancels a fine-tuning job.
func (c *Client) CancelFineTuningJob(ctx context.Context, jobID string, modelHint ...string) (*finetuning.Job, error) {
	if jobID == "" {
		return nil, ErrInvalidRequest
	}
	hint := ""
	if len(modelHint) > 0 {
		hint = modelHint[0]
	}
	fp, err := c.fineTuningProvider(hint)
	if err != nil {
		return nil, err
	}
	out, err := fp.CancelFineTuningJob(ctx, jobID)
	if err != nil {
		return nil, MapError(err)
	}
	return fromEngineFineTuneJob(out), nil
}

func toEngineFineTuneCreate(req *finetuning.CreateJobRequest) *types.FineTuningJobCreateRequest {
	return &types.FineTuningJobCreateRequest{
		Model:           req.Model,
		TrainingFile:    req.TrainingFile,
		ValidationFile:  req.ValidationFile,
		Hyperparameters: req.Hyperparameters,
		Suffix:          req.Suffix,
		Integrations:    req.Integrations,
		Seed:            req.Seed,
	}
}

func fromEngineFineTuneJob(j *types.FineTuningJob) *finetuning.Job {
	if j == nil {
		return nil
	}
	return &finetuning.Job{
		ID:             j.ID,
		Object:         j.Object,
		Model:          j.Model,
		CreatedAt:      j.CreatedAt,
		FinishedAt:     j.FinishedAt,
		FineTunedModel: j.FineTunedModel,
		OrganizationID: j.OrganizationID,
		Status:         j.Status,
		TrainingFile:   j.TrainingFile,
		ValidationFile: j.ValidationFile,
		Hyperparameters:  j.Hyperparameters,
		Error:          j.Error,
		Provider:       j.Provider,
	}
}

func fromEngineFineTuneList(r *types.FineTuningJobList) *finetuning.JobList {
	if r == nil {
		return nil
	}
	out := &finetuning.JobList{Object: r.Object, HasMore: r.HasMore, Provider: r.Provider}
	for i := range r.Data {
		out.Data = append(out.Data, *fromEngineFineTuneJob(&r.Data[i]))
	}
	return out
}
