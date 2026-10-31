package sdk

import (
	"net/http"
	"time"
)

// Config holds client-level settings.
type Config struct {
	Provider   string
	APIKey     string
	Model      string
	BaseURL    string
	HTTPClient *http.Client
	Timeout    time.Duration
}

func defaultConfig() Config {
	return Config{
		Timeout: 120 * time.Second,
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
