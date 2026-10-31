package sdk

import (
	"net/http"
	"time"

	"github.com/icosmos-space/llmatrix/engine/retry"
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

// WithRetry sets the retry policy for non-streaming calls.
func WithRetry(policy retry.Policy) Option {
	return func(c *Config) { c.Retry = policy }
}

// WithMaxRetries sets MaxAttempts on the retry policy (0 disables retries beyond first try).
func WithMaxRetries(n int) Option {
	return func(c *Config) {
		if c.Retry.MaxAttempts == 0 {
			c.Retry = retry.DefaultPolicy()
		}
		c.Retry.MaxAttempts = n
	}
}
