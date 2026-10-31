package anthropic

import (
	"net/http"
	"os"
	"time"

	"github.com/icosmos-space/llmatrix/engine/provider"
)

const (
	defaultBaseURL = "https://api.anthropic.com"
	apiVersion     = "2023-06-01"
	defaultMaxTok  = 4096
)

// Provider implements the Anthropic Messages API.
type Provider struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
	Version    string
}

// New creates an Anthropic provider.
func New(apiKey string, opts ...Option) *Provider {
	p := &Provider{APIKey: apiKey, BaseURL: defaultBaseURL, Version: apiVersion}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// Option configures Anthropic.
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

// WithVersion sets the anthropic-version header.
func WithVersion(v string) Option {
	return func(p *Provider) {
		if v != "" {
			p.Version = v
		}
	}
}

// Info implements provider.Provider.
func (p *Provider) Info() provider.Info {
	return provider.Info{
		Name:        "anthropic",
		DisplayName: "Anthropic",
		Capabilities: provider.Capabilities{
			Chat: true, Streaming: true, Tools: true, Embeddings: false,
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
	return os.Getenv("ANTHROPIC_API_KEY")
}

func (p *Provider) resolvedBaseURL() string {
	if p.BaseURL != "" {
		return p.BaseURL
	}
	return defaultBaseURL
}

func (p *Provider) resolvedVersion() string {
	if p.Version != "" {
		return p.Version
	}
	return apiVersion
}
