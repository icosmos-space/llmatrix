package sdk

import (
	"context"
	"fmt"

	engineerrors "github.com/icosmos-space/llmatrix/engine/errors"
	"github.com/icosmos-space/llmatrix/engine/model"
	"github.com/icosmos-space/llmatrix/engine/provider"
	engineresponses "github.com/icosmos-space/llmatrix/engine/responses"
	"github.com/icosmos-space/llmatrix/engine/types"
	"github.com/icosmos-space/llmatrix/sdk/responses"
	"github.com/icosmos-space/llmatrix/sdk/streaming"
)

// Responses calls the OpenAI Responses API (LiteLLM responses()).
// Providers without native /responses support are bridged via chat completions.
func (c *Client) Responses(ctx context.Context, req *responses.Request) (*responses.Response, error) {
	if req == nil {
		return nil, ErrInvalidRequest
	}
	out, err := c.withCandidates(ctx, "responses", req.Model, func(ctx context.Context, ref model.Ref, p provider.Provider) (types.Usage, any, error) {
		engineReq := toEngineResponsesRequest(req, ref.Model)
		resp, err := engineresponses.Create(ctx, p, engineReq)
		if err != nil {
			return types.Usage{}, nil, err
		}
		usage := types.Usage{
			PromptTokens:     resp.Usage.InputTokens,
			CompletionTokens: resp.Usage.OutputTokens,
			TotalTokens:      resp.Usage.TotalTokens,
		}
		return usage, fromEngineResponsesResponse(resp), nil
	})
	if err != nil {
		return nil, err
	}
	return out.(*responses.Response), nil
}

// ResponsesStream streams Responses API output (native SSE or chat-bridge).
func (c *Client) ResponsesStream(ctx context.Context, req *responses.Request) (streaming.Stream, error) {
	if req == nil {
		return nil, ErrInvalidRequest
	}
	candidates := c.candidates(req.Model)
	if len(candidates) == 0 {
		return nil, fmt.Errorf("%w: model is required (set request.Model or WithModel)", ErrInvalidConfig)
	}
	primaryProvider := candidates[0].Provider
	var last error
	for i, ref := range candidates {
		p, err := c.newProvider(ref, ref.Provider == primaryProvider)
		if err != nil {
			last = MapError(err)
			continue
		}
		engineReq := toEngineResponsesRequest(req, ref.Model)
		engineReq.Stream = true
		s, err := engineresponses.CreateStream(ctx, p, engineReq)
		if err == nil {
			return wrapEngineStream(s), nil
		}
		last = MapError(err)
		if i < len(candidates)-1 && engineerrors.IsFallbackable(err) {
			continue
		}
		return nil, last
	}
	return nil, last
}

// GetResponse fetches a stored Responses object by id (OpenAI).
func (c *Client) GetResponse(ctx context.Context, responseID string) (*responses.Response, error) {
	if responseID == "" {
		return nil, ErrInvalidRequest
	}
	p, ref, err := c.storeProvider()
	if err != nil {
		return nil, err
	}
	store, ok := p.(provider.ResponsesStoreProvider)
	if !ok {
		return nil, fmt.Errorf("%w: provider %q does not support get response", ErrUnsupported, ref.Provider)
	}
	resp, err := store.GetResponse(ctx, responseID)
	if err != nil {
		return nil, MapError(err)
	}
	return fromEngineResponsesResponse(resp), nil
}

// DeleteResponse deletes a stored Responses object by id (OpenAI).
func (c *Client) DeleteResponse(ctx context.Context, responseID string) error {
	if responseID == "" {
		return ErrInvalidRequest
	}
	p, ref, err := c.storeProvider()
	if err != nil {
		return err
	}
	store, ok := p.(provider.ResponsesStoreProvider)
	if !ok {
		return fmt.Errorf("%w: provider %q does not support delete response", ErrUnsupported, ref.Provider)
	}
	if err := store.DeleteResponse(ctx, responseID); err != nil {
		return MapError(err)
	}
	return nil
}

// CancelResponse cancels an in-progress Responses object by id (OpenAI).
func (c *Client) CancelResponse(ctx context.Context, responseID string) (*responses.Response, error) {
	if responseID == "" {
		return nil, ErrInvalidRequest
	}
	p, ref, err := c.storeProvider()
	if err != nil {
		return nil, err
	}
	store, ok := p.(provider.ResponsesStoreProvider)
	if !ok {
		return nil, fmt.Errorf("%w: provider %q does not support cancel response", ErrUnsupported, ref.Provider)
	}
	resp, err := store.CancelResponse(ctx, responseID)
	if err != nil {
		return nil, MapError(err)
	}
	return fromEngineResponsesResponse(resp), nil
}

func (c *Client) storeProvider() (provider.Provider, model.Ref, error) {
	raw := c.cfg.Model
	if raw == "" {
		prov := c.cfg.Provider
		if prov == "" {
			prov = "openai"
		}
		raw = prov + "/default"
	}
	ref := model.ResolveAlias(c.cfg.Aliases, raw, c.cfg.Provider)
	p, err := c.newProvider(ref, true)
	return p, ref, err
}
