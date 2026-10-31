package xai

import (
	"net/http"

	compat "github.com/icosmos-space/llmatrix/engine/providers/openaicompatible"
)

// Provider implements the xAI (Grok) OpenAI-compatible API.
type Provider = compat.Provider

// New creates an xAI provider.
func New(apiKey string, opts ...Option) *Provider {
	p := compat.New("xai", apiKey, "")
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// Option configures xAI.
type Option func(*Provider)

// WithBaseURL overrides the API base URL.
func WithBaseURL(url string) Option {
	return func(p *Provider) { p.BaseURL = url }
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(c *http.Client) Option {
	return func(p *Provider) { p.HTTPClient = c }
}
