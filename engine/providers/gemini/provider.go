package gemini

import (
	"net/http"
	"os"
	"time"

	"github.com/icosmos-space/llmatrix/engine/provider"
)

const defaultBaseURL = "https://generativelanguage.googleapis.com/v1beta"

// Provider implements the Google Gemini generateContent API.
type Provider struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
}

// New creates a Gemini provider.
func New(apiKey string, opts ...Option) *Provider {
	p := &Provider{APIKey: apiKey, BaseURL: defaultBaseURL}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// Option configures Gemini.
type Option func(*Provider)

// WithBaseURL overrides the API base URL.
func WithBaseURL(url string) Option {
	return func(p *Provider) {
		if url != "" {
			p.BaseURL = url
		}
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(c *http.Client) Option {
	return func(p *Provider) { p.HTTPClient = c }
}

// Info implements provider.Provider.
func (p *Provider) Info() provider.Info {
	return provider.Info{
		Name:        "gemini",
		DisplayName: "Google Gemini",
		Capabilities: provider.Capabilities{
			Chat: true, Streaming: true, Tools: true, Embeddings: true,
		},
	}
}

func (p *Provider) httpClient() *http.Client {
	if p.HTTPClient != nil {
		return p.HTTPClient
	}
	return &http.Client{Timeout: 120 * time.Second}
}

func (p *Provider) resolvedAPIKey() string {
	if p.APIKey != "" {
		return p.APIKey
	}
	if v := os.Getenv("GEMINI_API_KEY"); v != "" {
		return v
	}
	return os.Getenv("GOOGLE_API_KEY")
}

func (p *Provider) resolvedBaseURL() string {
	if p.BaseURL != "" {
		return p.BaseURL
	}
	return defaultBaseURL
}
