package deepinfra

import (
	"net/http"

	compat "github.com/icosmos-space/llmatrix/engine/providers/openaicompatible"
)

// Provider implements the DeepInfra OpenAI-compatible API (LiteLLM deepinfra/).
type Provider = compat.Provider

// New creates a DeepInfra provider.
func New(apiKey string, opts ...Option) *Provider {
	p := compat.New("deepinfra", apiKey, "")
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// Option configures DeepInfra.
type Option func(*Provider)

// WithBaseURL overrides the API base URL.
func WithBaseURL(url string) Option {
	return func(p *Provider) { p.BaseURL = url }
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(c *http.Client) Option {
	return func(p *Provider) { p.HTTPClient = c }
}
