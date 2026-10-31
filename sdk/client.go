package sdk

import (
	"context"

	"github.com/icosmos-space/llmatrix/sdk/chat"
	"github.com/icosmos-space/llmatrix/sdk/embedding"
)

// Client is the unified entry point for LLMatrix SDK calls.
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

// ChatCompletion performs a chat completion request.
func (c *Client) ChatCompletion(ctx context.Context, req *chat.Request) (*chat.Response, error) {
	_ = ctx
	_ = req
	return nil, ErrNotImplemented
}

// Embeddings creates embeddings for the given input.
func (c *Client) Embeddings(ctx context.Context, req *embedding.Request) (*embedding.Response, error) {
	_ = ctx
	_ = req
	return nil, ErrNotImplemented
}
