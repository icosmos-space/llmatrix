package sdk

import (
	"context"
	"fmt"

	"github.com/icosmos-space/llmatrix/engine/model"
	"github.com/icosmos-space/llmatrix/engine/provider"
	"github.com/icosmos-space/llmatrix/engine/types"
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
		if v, ok := c.cacheGet("text", modelKey, req); ok {
			if cached, ok := v.(*text.Response); ok {
				return cached, nil
			}
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
