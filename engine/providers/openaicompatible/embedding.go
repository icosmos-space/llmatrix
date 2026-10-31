package openaicompatible

import (
	"context"

	"github.com/icosmos-space/llmatrix/engine/errors"
	"github.com/icosmos-space/llmatrix/engine/types"
)

// Embeddings implements provider.EmbeddingProvider.
func (p *Provider) Embeddings(ctx context.Context, req *types.EmbeddingRequest) (*types.EmbeddingResponse, error) {
	if req == nil {
		return nil, errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "request is nil", nil)
	}
	if req.Model == "" {
		return nil, errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "model is required", nil)
	}
	if len(req.Input) == 0 {
		return nil, errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "input is required", nil)
	}

	body := wireEmbeddingRequest{
		Model:          req.Model,
		Input:          req.Input,
		Dimensions:     req.Dimensions,
		EncodingFormat: req.EncodingFormat,
		User:           req.User,
	}
	var raw wireEmbeddingResponse
	if err := p.doJSON(ctx, "POST", "/embeddings", body, &raw); err != nil {
		return nil, err
	}
	if raw.Error != nil && raw.Error.Message != "" {
		return nil, errors.Wrap(errors.KindProvider, p.Info().Name, raw.Error.Message, nil)
	}

	out := &types.EmbeddingResponse{
		Model:    raw.Model,
		Provider: p.Info().Name,
		Usage:    raw.Usage,
		Data:     make([]types.EmbeddingItem, 0, len(raw.Data)),
	}
	for _, d := range raw.Data {
		out.Data = append(out.Data, types.EmbeddingItem{
			Index:     d.Index,
			Embedding: d.Embedding,
		})
	}
	if out.Model == "" {
		out.Model = req.Model
	}
	return out, nil
}
