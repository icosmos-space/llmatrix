package sdk

import (
	"context"
	"fmt"

	engineerrors "github.com/icosmos-space/llmatrix/engine/errors"
	"github.com/icosmos-space/llmatrix/engine/model"
	"github.com/icosmos-space/llmatrix/engine/provider"
	"github.com/icosmos-space/llmatrix/engine/types"
	"github.com/icosmos-space/llmatrix/sdk/chat"
	"github.com/icosmos-space/llmatrix/sdk/streaming"
)

// ChatCompletion performs a non-streaming chat completion (LiteLLM completion()).
func (c *Client) ChatCompletion(ctx context.Context, req *chat.Request) (*chat.Response, error) {
	if req == nil {
		return nil, ErrInvalidRequest
	}
	out, err := c.withCandidates(ctx, "chat", req.Model, func(ctx context.Context, ref model.Ref, p provider.Provider) (types.Usage, any, error) {
		chatProv, ok := p.(provider.ChatProvider)
		if !ok {
			return types.Usage{}, nil, fmt.Errorf("%w: provider %q does not support chat", ErrUnsupported, ref.Provider)
		}
		engineReq := toEngineChatRequest(req, ref.Model)
		engineReq.Stream = false
		resp, err := chatProv.ChatCompletion(ctx, engineReq)
		if err != nil {
			return types.Usage{}, nil, err
		}
		return resp.Usage, fromEngineChatResponse(resp), nil
	})
	if err != nil {
		return nil, err
	}
	return out.(*chat.Response), nil
}

// ChatCompletionStream performs a streaming chat completion (LiteLLM completion(stream=True)).
func (c *Client) ChatCompletionStream(ctx context.Context, req *chat.Request) (streaming.Stream, error) {
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
		streamProv, ok := p.(provider.ChatStreamProvider)
		if !ok {
			last = fmt.Errorf("%w: provider %q does not support streaming", ErrUnsupported, ref.Provider)
			continue
		}
		engineReq := toEngineChatRequest(req, ref.Model)
		engineReq.Stream = true
		s, err := streamProv.ChatCompletionStream(ctx, engineReq)
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
