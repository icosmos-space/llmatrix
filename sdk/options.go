package sdk

import (
	"net/http"
	"time"
)

// Option configures a Client.
type Option func(*Config)

// WithProvider sets the default provider name (e.g. "openai").
// Model strings like "deepseek/deepseek-chat" override this.
func WithProvider(name string) Option {
	return func(c *Config) { c.Provider = name }
}

// WithAPIKey sets the API key used for authentication.
func WithAPIKey(key string) Option {
	return func(c *Config) { c.APIKey = key }
}

// WithModel sets the default model identifier (supports "provider/model").
func WithModel(model string) Option {
	return func(c *Config) { c.Model = model }
}

// WithBaseURL overrides the provider or proxy base URL.
func WithBaseURL(url string) Option {
	return func(c *Config) { c.BaseURL = url }
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(client *http.Client) Option {
	return func(c *Config) { c.HTTPClient = client }
}

// WithTimeout sets the default HTTP client timeout when no HTTPClient is provided.
func WithTimeout(d time.Duration) Option {
	return func(c *Config) { c.Timeout = d }
}
