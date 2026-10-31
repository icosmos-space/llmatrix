package sdk

import (
	"context"
	"fmt"

	"github.com/icosmos-space/llmatrix/engine/model"
	"github.com/icosmos-space/llmatrix/engine/provider"
	"github.com/icosmos-space/llmatrix/engine/registry"
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
// "provider/model" strings are supported (e.g. "deepseek/deepseek-chat").
func (c *Client) ChatCompletion(ctx context.Context, req *chat.Request) (*chat.Response, error) {
	if req == nil {
		return nil, ErrInvalidRequest
	}
	ref, p, err := c.resolveChatProvider(req.Model)
	if err != nil {
		return nil, err
	}
	chatProv, ok := p.(provider.ChatProvider)
	if !ok {
		return nil, fmt.Errorf("%w: provider %q does not support chat", ErrUnsupported, ref.Provider)
	}
	engineReq := toEngineChatRequest(req, ref.Model)
	engineReq.Stream = false
	resp, err := chatProv.ChatCompletion(ctx, engineReq)
	if err != nil {
		return nil, err
	}
	return fromEngineChatResponse(resp), nil
}

// ChatCompletionStream performs a streaming chat completion.
func (c *Client) ChatCompletionStream(ctx context.Context, req *chat.Request) (streaming.Stream, error) {
	if req == nil {
		return nil, ErrInvalidRequest
	}
	ref, p, err := c.resolveChatProvider(req.Model)
	if err != nil {
		return nil, err
	}
	streamProv, ok := p.(provider.ChatStreamProvider)
	if !ok {
		return nil, fmt.Errorf("%w: provider %q does not support streaming", ErrUnsupported, ref.Provider)
	}
	engineReq := toEngineChatRequest(req, ref.Model)
	engineReq.Stream = true
	s, err := streamProv.ChatCompletionStream(ctx, engineReq)
	if err != nil {
		return nil, err
	}
	return wrapEngineStream(s), nil
}

// Embeddings creates embeddings for the given input.
func (c *Client) Embeddings(ctx context.Context, req *embedding.Request) (*embedding.Response, error) {
	if req == nil {
		return nil, ErrInvalidRequest
	}
	ref, p, err := c.resolveChatProvider(req.Model)
	if err != nil {
		return nil, err
	}
	embProv, ok := p.(provider.EmbeddingProvider)
	if !ok {
		return nil, fmt.Errorf("%w: provider %q does not support embeddings", ErrUnsupported, ref.Provider)
	}
	resp, err := embProv.Embeddings(ctx, toEngineEmbeddingRequest(req, ref.Model))
	if err != nil {
		return nil, err
	}
	return fromEngineEmbeddingResponse(resp), nil
}

func (c *Client) resolveChatProvider(requestModel string) (model.Ref, provider.Provider, error) {
	rawModel := requestModel
	if rawModel == "" {
		rawModel = c.cfg.Model
	}
	if rawModel == "" {
		return model.Ref{}, nil, fmt.Errorf("%w: model is required (set request.Model or WithModel)", ErrInvalidConfig)
	}
	ref := model.Resolve(rawModel, c.cfg.Provider)
	p, err := registry.NewProvider(registry.ProviderConfig{
		Name:       ref.Provider,
		APIKey:     c.cfg.APIKey,
		BaseURL:    c.cfg.BaseURL,
		HTTPClient: c.cfg.httpClient(),
	})
	if err != nil {
		return ref, nil, err
	}
	return ref, p, nil
}
