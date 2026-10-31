package sdk

import (
	"net/http"
	"time"

	"github.com/icosmos-space/llmatrix/engine/callback"
	"github.com/icosmos-space/llmatrix/engine/model"
	"github.com/icosmos-space/llmatrix/engine/retry"
)

// Config holds client-level settings.
type Config struct {
	Provider   string
	APIKey     string
	Model      string
	BaseURL    string
	APIVersion string // Azure api-version
	HTTPClient *http.Client
	Timeout    time.Duration
	Retry      retry.Policy
	Fallbacks  []string
	Aliases    *model.AliasMap
	Callbacks  *callback.Bus
}

func defaultConfig() Config {
	return Config{
		Timeout:   120 * time.Second,
		Retry:     retry.DefaultPolicy(),
		Aliases:   model.DefaultAliases(),
		Callbacks: &callback.Bus{},
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
