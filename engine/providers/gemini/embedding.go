package gemini

import (
	"context"
	"fmt"
	"strings"

	"github.com/icosmos-space/llmatrix/engine/errors"
	"github.com/icosmos-space/llmatrix/engine/types"
)

// Embeddings implements provider.EmbeddingProvider via embedContent (one input at a time).
func (p *Provider) Embeddings(ctx context.Context, req *types.EmbeddingRequest) (*types.EmbeddingResponse, error) {
	if req == nil {
		return nil, errors.Wrap(errors.KindInvalidRequest, "gemini", "request is nil", nil)
	}
	if req.Model == "" {
		return nil, errors.Wrap(errors.KindInvalidRequest, "gemini", "model is required", nil)
	}
	if len(req.Input) == 0 {
		return nil, errors.Wrap(errors.KindInvalidRequest, "gemini", "input is required", nil)
	}

	model := req.Model
	if !strings.Contains(model, "embedding") && !strings.HasPrefix(model, "text-") {
		// allow short names like text-embedding-004
	}

	out := &types.EmbeddingResponse{
		Model:    model,
		Provider: "gemini",
		Data:     make([]types.EmbeddingItem, 0, len(req.Input)),
	}
	for i, text := range req.Input {
		rawURL, err := p.modelURL(model, "embedContent")
		if err != nil {
			return nil, err
		}
		body := wireEmbedRequest{
			Model:   "models/" + strings.TrimPrefix(model, "models/"),
			Content: wireContent{Parts: []wirePart{{Text: text}}},
		}
		var raw wireEmbedResponse
		if err := p.doJSON(ctx, "POST", rawURL, body, &raw); err != nil {
			return nil, err
		}
		if raw.Error != nil && raw.Error.Message != "" {
			return nil, errors.Wrap(errors.KindProvider, "gemini", raw.Error.Message, nil)
		}
		if raw.Embedding == nil {
			return nil, errors.Wrap(errors.KindProvider, "gemini", fmt.Sprintf("empty embedding for input %d", i), nil)
		}
		out.Data = append(out.Data, types.EmbeddingItem{
			Index:     i,
			Embedding: raw.Embedding.Values,
		})
	}
	return out, nil
}
