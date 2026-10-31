package sdk

import (
	"context"
	"fmt"

	"github.com/icosmos-space/llmatrix/engine/provider"
	"github.com/icosmos-space/llmatrix/engine/types"
	"github.com/icosmos-space/llmatrix/sdk/models"
)

func (c *Client) modelsProvider(modelHint string) (provider.ModelsProvider, error) {
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
	mp, ok := p.(provider.ModelsProvider)
	if !ok {
		return nil, fmt.Errorf("%w: provider %q does not support list models", ErrUnsupported, refs[0].Provider)
	}
	return mp, nil
}

// ListModels lists available models (LiteLLM get_model_list / OpenAI GET /models).
func (c *Client) ListModels(ctx context.Context, modelHint ...string) (*models.ListResponse, error) {
	hint := ""
	if len(modelHint) > 0 {
		hint = modelHint[0]
	}
	mp, err := c.modelsProvider(hint)
	if err != nil {
		return nil, err
	}
	out, err := mp.ListModels(ctx)
	if err != nil {
		return nil, MapError(err)
	}
	return fromEngineModelList(out), nil
}

func fromEngineModelList(r *types.ModelList) *models.ListResponse {
	if r == nil {
		return nil
	}
	out := &models.ListResponse{Object: r.Object, Provider: r.Provider}
	for _, m := range r.Data {
		out.Data = append(out.Data, models.Model{
			ID:       m.ID,
			Object:   m.Object,
			Created:  m.Created,
			OwnedBy:  m.OwnedBy,
			Provider: m.Provider,
		})
	}
	return out
}
