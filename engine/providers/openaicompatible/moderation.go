package openaicompatible

import (
	"context"

	"github.com/icosmos-space/llmatrix/engine/errors"
	"github.com/icosmos-space/llmatrix/engine/types"
)

// Moderation implements provider.ModerationProvider via POST /moderations.
func (p *Provider) Moderation(ctx context.Context, req *types.ModerationRequest) (*types.ModerationResponse, error) {
	if req == nil {
		return nil, errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "request is nil", nil)
	}
	if req.Input == nil {
		return nil, errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "input is required", nil)
	}
	body := map[string]any{"input": req.Input}
	if req.Model != "" {
		body["model"] = req.Model
	}
	var raw types.ModerationResponse
	if err := p.doJSON(ctx, "POST", "/moderations", body, &raw); err != nil {
		return nil, err
	}
	raw.Provider = p.Info().Name
	if raw.Model == "" {
		raw.Model = req.Model
	}
	return &raw, nil
}
