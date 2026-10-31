package ollama

import (
	"net/http"

	compat "github.com/icosmos-space/llmatrix/engine/providers/openaicompatible"
)

// Provider implements the Ollama local OpenAI-compatible API.
type Provider = compat.Provider

// New creates an Ollama provider. API key is optional.
func New(opts ...Option) *Provider {
	p := compat.New("ollama", "", "")
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// Option configures an Ollama provider.
type Option func(*Provider)

// WithBaseURL overrides the API base URL (default http://localhost:11434/v1).
func WithBaseURL(url string) Option {
	return func(p *Provider) { p.BaseURL = url }
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(c *http.Client) Option {
	return func(p *Provider) { p.HTTPClient = c }
}

// WithAPIKey sets an optional API key for authenticated Ollama gateways.
func WithAPIKey(key string) Option {
	return func(p *Provider) { p.APIKey = key }
}
