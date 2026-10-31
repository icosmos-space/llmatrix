package registry

import (
	"fmt"
	"net/http"

	"github.com/icosmos-space/llmatrix/engine/provider"
	"github.com/icosmos-space/llmatrix/engine/providers/anthropic"
	"github.com/icosmos-space/llmatrix/engine/providers/azure"
	"github.com/icosmos-space/llmatrix/engine/providers/deepseek"
	"github.com/icosmos-space/llmatrix/engine/providers/fireworks"
	"github.com/icosmos-space/llmatrix/engine/providers/gemini"
	"github.com/icosmos-space/llmatrix/engine/providers/groq"
	"github.com/icosmos-space/llmatrix/engine/providers/mistral"
	"github.com/icosmos-space/llmatrix/engine/providers/ollama"
	"github.com/icosmos-space/llmatrix/engine/providers/openai"
	compat "github.com/icosmos-space/llmatrix/engine/providers/openaicompatible"
	"github.com/icosmos-space/llmatrix/engine/providers/openrouter"
	"github.com/icosmos-space/llmatrix/engine/providers/together"
	"github.com/icosmos-space/llmatrix/engine/providers/xai"
)

// ProviderConfig configures a concrete provider instance.
type ProviderConfig struct {
	Name         string
	APIKey       string
	BaseURL      string
	APIVersion   string // Azure api-version
	HTTPClient   *http.Client
	Organization string
	ExtraHeaders map[string]string
}

func applyCompatExtras(p *compat.Provider, cfg ProviderConfig) {
	if cfg.Organization != "" {
		p.Organization = cfg.Organization
	}
	if len(cfg.ExtraHeaders) > 0 {
		p.ExtraHeaders = cfg.ExtraHeaders
	}
}

// NewProvider builds a provider by name.
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
		applyCompatExtras(p, cfg)
		return p, nil
	case "deepseek":
		p := deepseek.New(cfg.APIKey, deepseek.WithBaseURL(cfg.BaseURL))
		if cfg.HTTPClient != nil {
			p.HTTPClient = cfg.HTTPClient
		}
		applyCompatExtras(p, cfg)
		return p, nil
	case "groq":
		p := groq.New(cfg.APIKey, groq.WithBaseURL(cfg.BaseURL))
		if cfg.HTTPClient != nil {
			p.HTTPClient = cfg.HTTPClient
		}
		applyCompatExtras(p, cfg)
		return p, nil
	case "together":
		p := together.New(cfg.APIKey, together.WithBaseURL(cfg.BaseURL))
		if cfg.HTTPClient != nil {
			p.HTTPClient = cfg.HTTPClient
		}
		applyCompatExtras(p, cfg)
		return p, nil
	case "mistral":
		p := mistral.New(cfg.APIKey, mistral.WithBaseURL(cfg.BaseURL))
		if cfg.HTTPClient != nil {
			p.HTTPClient = cfg.HTTPClient
		}
		applyCompatExtras(p, cfg)
		return p, nil
	case "openrouter":
		p := openrouter.New(cfg.APIKey, openrouter.WithBaseURL(cfg.BaseURL))
		if cfg.HTTPClient != nil {
			p.HTTPClient = cfg.HTTPClient
		}
		applyCompatExtras(p, cfg)
		return p, nil
	case "xai":
		p := xai.New(cfg.APIKey, xai.WithBaseURL(cfg.BaseURL))
		if cfg.HTTPClient != nil {
			p.HTTPClient = cfg.HTTPClient
		}
		applyCompatExtras(p, cfg)
		return p, nil
	case "fireworks":
		p := fireworks.New(cfg.APIKey, fireworks.WithBaseURL(cfg.BaseURL))
		if cfg.HTTPClient != nil {
			p.HTTPClient = cfg.HTTPClient
		}
		applyCompatExtras(p, cfg)
		return p, nil
	case "ollama":
		opts := []ollama.Option{ollama.WithBaseURL(cfg.BaseURL), ollama.WithAPIKey(cfg.APIKey)}
		if cfg.HTTPClient != nil {
			opts = append(opts, ollama.WithHTTPClient(cfg.HTTPClient))
		}
		p := ollama.New(opts...)
		applyCompatExtras(p, cfg)
		return p, nil
	case "anthropic":
		opts := []anthropic.Option{anthropic.WithBaseURL(cfg.BaseURL)}
		if cfg.HTTPClient != nil {
			opts = append(opts, anthropic.WithHTTPClient(cfg.HTTPClient))
		}
		return anthropic.New(cfg.APIKey, opts...), nil
	case "gemini", "google", "vertex_ai":
		opts := []gemini.Option{gemini.WithBaseURL(cfg.BaseURL)}
		if cfg.HTTPClient != nil {
			opts = append(opts, gemini.WithHTTPClient(cfg.HTTPClient))
		}
		return gemini.New(cfg.APIKey, opts...), nil
	case "azure", "azure_ai", "azure_openai":
		opts := []azure.Option{azure.WithBaseURL(cfg.BaseURL), azure.WithAPIVersion(cfg.APIVersion)}
		if cfg.HTTPClient != nil {
			opts = append(opts, azure.WithHTTPClient(cfg.HTTPClient))
		}
		return azure.New(cfg.APIKey, opts...), nil
	case "openai-compatible", "openaicompatible":
		p := compat.New("openai-compatible", cfg.APIKey, cfg.BaseURL)
		p.HTTPClient = cfg.HTTPClient
		applyCompatExtras(p, cfg)
		return p, nil
	default:
		return nil, fmt.Errorf("registry: unsupported provider %q (supported: openai, deepseek, groq, together, mistral, openrouter, xai, fireworks, ollama, anthropic, gemini, azure, openai-compatible)", name)
	}
}
