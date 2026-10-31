package openaicompatible

import (
	"context"
	"net/url"
	"strconv"

	"github.com/icosmos-space/llmatrix/engine/errors"
	"github.com/icosmos-space/llmatrix/engine/types"
)

// CreateBatch implements provider.BatchesProvider via POST /batches.
func (p *Provider) CreateBatch(ctx context.Context, req *types.BatchCreateRequest) (*types.Batch, error) {
	if req == nil || req.InputFileID == "" || req.Endpoint == "" {
		return nil, errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "input_file_id and endpoint are required", nil)
	}
	window := req.CompletionWindow
	if window == "" {
		window = "24h"
	}
	body := map[string]any{
		"completion_window": window,
		"endpoint":          req.Endpoint,
		"input_file_id":     req.InputFileID,
	}
	if req.Metadata != nil {
		body["metadata"] = req.Metadata
	}
	var out types.Batch
	if err := p.doJSON(ctx, "POST", "/batches", body, &out); err != nil {
		return nil, err
	}
	out.Provider = p.Info().Name
	return &out, nil
}

// RetrieveBatch implements provider.BatchesProvider via GET /batches/{id}.
func (p *Provider) RetrieveBatch(ctx context.Context, batchID string) (*types.Batch, error) {
	if batchID == "" {
		return nil, errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "batch id is required", nil)
	}
	var out types.Batch
	if err := p.doJSON(ctx, "GET", "/batches/"+url.PathEscape(batchID), nil, &out); err != nil {
		return nil, err
	}
	out.Provider = p.Info().Name
	return &out, nil
}

// CancelBatch implements provider.BatchesProvider via POST /batches/{id}/cancel.
func (p *Provider) CancelBatch(ctx context.Context, batchID string) (*types.Batch, error) {
	if batchID == "" {
		return nil, errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "batch id is required", nil)
	}
	var out types.Batch
	if err := p.doJSON(ctx, "POST", "/batches/"+url.PathEscape(batchID)+"/cancel", map[string]any{}, &out); err != nil {
		return nil, err
	}
	out.Provider = p.Info().Name
	return &out, nil
}

// ListBatches implements provider.BatchesProvider via GET /batches.
func (p *Provider) ListBatches(ctx context.Context, limit int) ([]types.Batch, error) {
	path := "/batches"
	if limit > 0 {
		path += "?limit=" + strconv.Itoa(limit)
	}
	var raw struct {
		Data []types.Batch `json:"data"`
	}
	if err := p.doJSON(ctx, "GET", path, nil, &raw); err != nil {
		return nil, err
	}
	for i := range raw.Data {
		raw.Data[i].Provider = p.Info().Name
	}
	return raw.Data, nil
}
