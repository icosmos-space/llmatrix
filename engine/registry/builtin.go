package registry

import (
	"fmt"
	"net/http"

	"github.com/icosmos-space/llmatrix/engine/provider"
	"github.com/icosmos-space/llmatrix/engine/providers/anthropic"
	"github.com/icosmos-space/llmatrix/engine/providers/deepseek"
	"github.com/icosmos-space/llmatrix/engine/providers/gemini"
	"github.com/icosmos-space/llmatrix/engine/providers/ollama"
	"github.com/icosmos-space/llmatrix/engine/providers/openai"
	compat "github.com/icosmos-space/llmatrix/engine/providers/openaicompatible"
)

// ProviderConfig configures a concrete provider instance.
type ProviderConfig struct {
	Name       string
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
}

// NewProvider builds a provider by name.
// Supported: openai, deepseek, ollama, anthropic, gemini, openai-compatible.
func NewProvider(cfg ProviderConfig) (provider.Provider, error) {
	name := cfg.Name
	if name == "" {
		name = "openai"
	}
	switch name {
	case "openai":
		p := openai.New(cfg.APIKey, openai.WithBaseURL(cfg.BaseURL))
		if cfg.HTTPClient != nil {
			p.HTTPClient = cfg.HTTPClient
		}
		return p, nil
	case "deepseek":
		p := deepseek.New(cfg.APIKey, deepseek.WithBaseURL(cfg.BaseURL))
		if cfg.HTTPClient != nil {
			p.HTTPClient = cfg.HTTPClient
		}
		return p, nil
	case "ollama":
		opts := []ollama.Option{ollama.WithBaseURL(cfg.BaseURL), ollama.WithAPIKey(cfg.APIKey)}
		if cfg.HTTPClient != nil {
			opts = append(opts, ollama.WithHTTPClient(cfg.HTTPClient))
		}
		return ollama.New(opts...), nil
	case "anthropic":
		opts := []anthropic.Option{anthropic.WithBaseURL(cfg.BaseURL)}
		if cfg.HTTPClient != nil {
			opts = append(opts, anthropic.WithHTTPClient(cfg.HTTPClient))
		}
		return anthropic.New(cfg.APIKey, opts...), nil
	case "gemini", "google", "vertex_ai":
		// vertex_ai alias maps to Gemini developer API for now.
		opts := []gemini.Option{gemini.WithBaseURL(cfg.BaseURL)}
		if cfg.HTTPClient != nil {
			opts = append(opts, gemini.WithHTTPClient(cfg.HTTPClient))
		}
		return gemini.New(cfg.APIKey, opts...), nil
	case "openai-compatible", "openaicompatible":
		p := compat.New("openai-compatible", cfg.APIKey, cfg.BaseURL)
		p.HTTPClient = cfg.HTTPClient
		return p, nil
	default:
		return nil, fmt.Errorf("registry: unsupported provider %q (supported: openai, deepseek, ollama, anthropic, gemini, openai-compatible)", name)
	}
}
