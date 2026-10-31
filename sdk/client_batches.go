package sdk

import (
	"context"
	"fmt"

	"github.com/icosmos-space/llmatrix/engine/provider"
	"github.com/icosmos-space/llmatrix/sdk/batches"
)

func (c *Client) batchesProvider(modelHint string) (provider.BatchesProvider, error) {
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
	bp, ok := p.(provider.BatchesProvider)
	if !ok {
		return nil, fmt.Errorf("%w: provider %q does not support batches", ErrUnsupported, refs[0].Provider)
	}
	return bp, nil
}

// CreateBatch creates a batch job (LiteLLM create_batch()).
func (c *Client) CreateBatch(ctx context.Context, req *batches.CreateRequest) (*batches.Batch, error) {
	if req == nil || req.InputFileID == "" || req.Endpoint == "" {
		return nil, ErrInvalidRequest
	}
	bp, err := c.batchesProvider(req.Model)
	if err != nil {
		return nil, err
	}
	out, err := bp.CreateBatch(ctx, toEngineBatchCreate(req))
	if err != nil {
		return nil, MapError(err)
	}
	return fromEngineBatch(out), nil
}

// RetrieveBatch fetches a batch by id (LiteLLM retrieve_batch()).
func (c *Client) RetrieveBatch(ctx context.Context, batchID string, modelHint ...string) (*batches.Batch, error) {
	if batchID == "" {
		return nil, ErrInvalidRequest
	}
	hint := ""
	if len(modelHint) > 0 {
		hint = modelHint[0]
	}
	bp, err := c.batchesProvider(hint)
	if err != nil {
		return nil, err
	}
	out, err := bp.RetrieveBatch(ctx, batchID)
	if err != nil {
		return nil, MapError(err)
	}
	return fromEngineBatch(out), nil
}

// CancelBatch cancels a batch (LiteLLM cancel_batch()).
func (c *Client) CancelBatch(ctx context.Context, batchID string, modelHint ...string) (*batches.Batch, error) {
	if batchID == "" {
		return nil, ErrInvalidRequest
	}
	hint := ""
	if len(modelHint) > 0 {
		hint = modelHint[0]
	}
	bp, err := c.batchesProvider(hint)
	if err != nil {
		return nil, err
	}
	out, err := bp.CancelBatch(ctx, batchID)
	if err != nil {
		return nil, MapError(err)
	}
	return fromEngineBatch(out), nil
}

// ListBatches lists batches (LiteLLM list_batches()).
func (c *Client) ListBatches(ctx context.Context, limit int, modelHint ...string) ([]batches.Batch, error) {
	hint := ""
	if len(modelHint) > 0 {
		hint = modelHint[0]
	}
	bp, err := c.batchesProvider(hint)
	if err != nil {
		return nil, err
	}
	out, err := bp.ListBatches(ctx, limit)
	if err != nil {
		return nil, MapError(err)
	}
	res := make([]batches.Batch, 0, len(out))
	for i := range out {
		res = append(res, *fromEngineBatch(&out[i]))
	}
	return res, nil
}
