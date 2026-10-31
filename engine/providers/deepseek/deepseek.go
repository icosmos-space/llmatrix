package deepseek

import (
	"net/http"

	compat "github.com/icosmos-space/llmatrix/engine/providers/openaicompatible"
)

// Provider implements the DeepSeek API (OpenAI-compatible).
type Provider = compat.Provider

// New creates a DeepSeek provider.
func New(apiKey string, opts ...Option) *Provider {
	p := compat.New("deepseek", apiKey, "")
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// Option configures a DeepSeek provider.
type Option func(*Provider)

// WithBaseURL overrides the API base URL.
func WithBaseURL(url string) Option {
	return func(p *Provider) { p.BaseURL = url }
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(c *http.Client) Option {
	return func(p *Provider) { p.HTTPClient = c }
}
