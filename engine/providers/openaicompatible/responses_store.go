package openaicompatible

import (
	"context"
	"net/url"

	"github.com/icosmos-space/llmatrix/engine/errors"
	"github.com/icosmos-space/llmatrix/engine/types"
)

// GetResponse implements provider.ResponsesStoreProvider via GET /responses/{id}.
func (p *Provider) GetResponse(ctx context.Context, responseID string) (*types.ResponsesResponse, error) {
	if responseID == "" {
		return nil, errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "response id is required", nil)
	}
	var raw types.ResponsesResponse
	if err := p.doJSON(ctx, "GET", "/responses/"+url.PathEscape(responseID), nil, &raw); err != nil {
		return nil, err
	}
	raw.Provider = p.Info().Name
	return &raw, nil
}

// DeleteResponse implements provider.ResponsesStoreProvider via DELETE /responses/{id}.
func (p *Provider) DeleteResponse(ctx context.Context, responseID string) error {
	if responseID == "" {
		return errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "response id is required", nil)
	}
	return p.doJSON(ctx, "DELETE", "/responses/"+url.PathEscape(responseID), nil, nil)
}

// CancelResponse implements provider.ResponsesStoreProvider via POST /responses/{id}/cancel.
func (p *Provider) CancelResponse(ctx context.Context, responseID string) (*types.ResponsesResponse, error) {
	if responseID == "" {
		return nil, errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "response id is required", nil)
	}
	var raw types.ResponsesResponse
	if err := p.doJSON(ctx, "POST", "/responses/"+url.PathEscape(responseID)+"/cancel", map[string]any{}, &raw); err != nil {
		return nil, err
	}
	raw.Provider = p.Info().Name
	return &raw, nil
}
