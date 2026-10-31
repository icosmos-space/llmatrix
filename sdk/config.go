package sdk

import (
	"net/http"
	"time"

	"github.com/icosmos-space/llmatrix/engine/retry"
)

// Config holds client-level settings.
type Config struct {
	Provider   string
	APIKey     string
	Model      string
	BaseURL    string
	HTTPClient *http.Client
	Timeout    time.Duration
	Retry      retry.Policy
	Fallbacks  []string // LiteLLM-style "provider/model" refs tried after primary fails
}

func defaultConfig() Config {
	return Config{
		Timeout: 120 * time.Second,
		Retry:   retry.DefaultPolicy(),
	}
}

func (c Config) httpClient() *http.Client {
	if c.HTTPClient != nil {
		return c.HTTPClient
	}
	timeout := c.Timeout
	if timeout <= 0 {
		timeout = 120 * time.Second
	}
	return &http.Client{Timeout: timeout}
}
