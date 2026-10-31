package sdk

import (
	"context"
	"fmt"

	engineerrors "github.com/icosmos-space/llmatrix/engine/errors"
	"github.com/icosmos-space/llmatrix/engine/model"
	"github.com/icosmos-space/llmatrix/engine/provider"
	"github.com/icosmos-space/llmatrix/engine/registry"
	"github.com/icosmos-space/llmatrix/engine/retry"
	"github.com/icosmos-space/llmatrix/sdk/chat"
	"github.com/icosmos-space/llmatrix/sdk/embedding"
	"github.com/icosmos-space/llmatrix/sdk/streaming"
)

// Client is the unified entry point for LLMatrix SDK calls
// (LiteLLM-style: one API, multiple providers).
type Client struct {
	cfg Config
}

// NewClient constructs a Client with the given options.
func NewClient(opts ...Option) *Client {
	c := &Client{cfg: defaultConfig()}
	for _, opt := range opts {
		opt(&c.cfg)
	}
	return c
}

// ChatCompletion performs a non-streaming chat completion.
// Model may be set on the request or via WithModel; LiteLLM-style
// "provider/model" strings are supported (e.g. "anthropic/claude-sonnet-4-5").
// Configured fallbacks are tried when the primary candidate fails with a fallbackable error.
func (c *Client) ChatCompletion(ctx context.Context, req *chat.Request) (*chat.Response, error) {
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
			last = err
			continue
		}
		chatProv, ok := p.(provider.ChatProvider)
		if !ok {
			last = fmt.Errorf("%w: provider %q does not support chat", ErrUnsupported, ref.Provider)
			continue
		}
		engineReq := toEngineChatRequest(req, ref.Model)
		engineReq.Stream = false

		var resp *chat.Response
		err = retry.Do(ctx, c.cfg.Retry, func(ctx context.Context) error {
			out, callErr := chatProv.ChatCompletion(ctx, engineReq)
			if callErr != nil {
				return callErr
			}
			resp = fromEngineChatResponse(out)
			return nil
		})
		if err == nil {
			return resp, nil
		}
		last = err
		if i < len(candidates)-1 && engineerrors.IsFallbackable(err) {
			continue
		}
		return nil, err
	}
	return nil, last
}

// ChatCompletionStream performs a streaming chat completion.
// Fallbacks apply only while establishing the stream (not mid-stream).
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
			last = err
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
		last = err
		if i < len(candidates)-1 && engineerrors.IsFallbackable(err) {
			continue
		}
		return nil, err
	}
	return nil, last
}

// Embeddings creates embeddings for the given input.
func (c *Client) Embeddings(ctx context.Context, req *embedding.Request) (*embedding.Response, error) {
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
			last = err
			continue
		}
		embProv, ok := p.(provider.EmbeddingProvider)
		if !ok {
			last = fmt.Errorf("%w: provider %q does not support embeddings", ErrUnsupported, ref.Provider)
			continue
		}
		var resp *embedding.Response
		err = retry.Do(ctx, c.cfg.Retry, func(ctx context.Context) error {
			out, callErr := embProv.Embeddings(ctx, toEngineEmbeddingRequest(req, ref.Model))
			if callErr != nil {
				return callErr
			}
			resp = fromEngineEmbeddingResponse(out)
			return nil
		})
		if err == nil {
			return resp, nil
		}
		last = err
		if i < len(candidates)-1 && engineerrors.IsFallbackable(err) {
			continue
		}
		return nil, err
	}
	return nil, last
}

func (c *Client) candidates(requestModel string) []model.Ref {
	rawModel := requestModel
	if rawModel == "" {
		rawModel = c.cfg.Model
	}
	if rawModel == "" {
		return nil
	}
	primary := model.Resolve(rawModel, c.cfg.Provider)
	out := []model.Ref{primary}
	for _, fb := range c.cfg.Fallbacks {
		ref := model.Resolve(fb, "")
		if ref.Provider == primary.Provider && ref.Model == primary.Model {
			continue
		}
		out = append(out, ref)
	}
	return out
}

// newProvider builds a provider. Shared APIKey/BaseURL are reused for candidates
// on the same provider family; cross-provider fallbacks use env-specific credentials.
func (c *Client) newProvider(ref model.Ref, shareClientCreds bool) (provider.Provider, error) {
	cfg := registry.ProviderConfig{
		Name:       ref.Provider,
		HTTPClient: c.cfg.httpClient(),
	}
	if shareClientCreds {
		cfg.APIKey = c.cfg.APIKey
		cfg.BaseURL = c.cfg.BaseURL
	}
	return registry.NewProvider(cfg)
}
