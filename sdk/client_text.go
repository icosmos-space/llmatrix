package sdk

import (
	"context"
	"fmt"

	engineerrors "github.com/icosmos-space/llmatrix/engine/errors"
	"github.com/icosmos-space/llmatrix/engine/model"
	"github.com/icosmos-space/llmatrix/engine/provider"
	"github.com/icosmos-space/llmatrix/engine/types"
	"github.com/icosmos-space/llmatrix/sdk/streaming"
	"github.com/icosmos-space/llmatrix/sdk/text"
)

// TextCompletion calls the OpenAI legacy completions API (LiteLLM text_completion()).
func (c *Client) TextCompletion(ctx context.Context, req *text.Request) (*text.Response, error) {
	if req == nil {
		return nil, ErrInvalidRequest
	}
	modelName := req.Model
	if modelName == "" {
		modelName = "openai/gpt-3.5-turbo-instruct"
	}
	modelKey := c.cacheModelKey(modelName)
	if c.cacheEnabled(req.Caching) {
		var cached text.Response
		if c.cacheGet("text", modelKey, req, &cached) {
			return &cached, nil
		}
	}
	out, err := c.withCandidates(ctx, "text_completion", modelName, func(ctx context.Context, ref model.Ref, p provider.Provider) (types.Usage, any, error) {
		tp, ok := p.(provider.TextCompletionProvider)
		if !ok {
			return types.Usage{}, nil, fmt.Errorf("%w: provider %q does not support text completion", ErrUnsupported, ref.Provider)
		}
		resp, err := tp.TextCompletion(ctx, toEngineTextRequest(req, ref.Model))
		if err != nil {
			return types.Usage{}, nil, err
		}
		return resp.Usage, fromEngineTextResponse(resp), nil
	})
	if err != nil {
		return nil, err
	}
	resp := out.(*text.Response)
	if c.cacheEnabled(req.Caching) {
		c.cacheSet("text", modelKey, req, resp)
	}
	return resp, nil
}

// TextCompletionStream streams a legacy text completion.
func (c *Client) TextCompletionStream(ctx context.Context, req *text.Request) (streaming.Stream, error) {
	if req == nil {
		return nil, ErrInvalidRequest
	}
	modelName := req.Model
	if modelName == "" {
		modelName = "openai/gpt-3.5-turbo-instruct"
	}
	candidates := c.candidates(modelName)
	if len(candidates) == 0 {
		return nil, fmt.Errorf("%w: model is required", ErrInvalidConfig)
	}
	primaryProvider := candidates[0].Provider
	var last error
	for i, ref := range candidates {
		p, err := c.newProvider(ref, ref.Provider == primaryProvider)
		if err != nil {
			last = MapError(err)
			continue
		}
		sp, ok := p.(provider.TextCompletionStreamProvider)
		if !ok {
			last = fmt.Errorf("%w: provider %q does not support text completion stream", ErrUnsupported, ref.Provider)
			continue
		}
		engineReq := toEngineTextRequest(req, ref.Model)
		engineReq.Stream = true
		s, err := sp.TextCompletionStream(ctx, engineReq)
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
