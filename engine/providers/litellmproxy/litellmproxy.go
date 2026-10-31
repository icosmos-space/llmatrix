package litellmproxy

import (
	"net/http"

	compat "github.com/icosmos-space/llmatrix/engine/providers/openaicompatible"
)

// Provider routes requests through a LiteLLM Proxy / AI Gateway
// (OpenAI-compatible). Use model prefix litellm_proxy/<model>.
type Provider = compat.Provider

// New creates a LiteLLM Proxy provider.
func New(apiKey string, opts ...Option) *Provider {
	p := compat.New("litellm_proxy", apiKey, "")
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// Option configures the LiteLLM Proxy provider.
type Option func(*Provider)

// WithBaseURL overrides LITELLM_PROXY_API_BASE.
func WithBaseURL(url string) Option {
	return func(p *Provider) { p.BaseURL = url }
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(c *http.Client) Option {
	return func(p *Provider) { p.HTTPClient = c }
}
