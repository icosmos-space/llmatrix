package openaicompatible

import (
	"context"

	"github.com/icosmos-space/llmatrix/engine/types"
)

// ListModels implements provider.ModelsProvider via GET /models.
func (p *Provider) ListModels(ctx context.Context) (*types.ModelList, error) {
	var out types.ModelList
	if err := p.doJSON(ctx, "GET", "/models", nil, &out); err != nil {
		return nil, err
	}
	out.Provider = p.Info().Name
	for i := range out.Data {
		if out.Data[i].Provider == "" {
			out.Data[i].Provider = p.Info().Name
		}
		if out.Data[i].Object == "" {
			out.Data[i].Object = "model"
		}
	}
	if out.Object == "" {
		out.Object = "list"
	}
	return &out, nil
}
