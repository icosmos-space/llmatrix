package sdk

import (
	"net/http"
	"time"

	"github.com/icosmos-space/llmatrix/engine/cache"
	"github.com/icosmos-space/llmatrix/engine/callback"
	"github.com/icosmos-space/llmatrix/engine/model"
	"github.com/icosmos-space/llmatrix/engine/retry"
)

// Option configures a Client.
type Option func(*Config)

// WithProvider sets the default provider name (e.g. "openai").
func WithProvider(name string) Option {
	return func(c *Config) { c.Provider = name }
}

// WithAPIKey sets the API key used for authentication.
func WithAPIKey(key string) Option {
	return func(c *Config) { c.APIKey = key }
}

// WithModel sets the default model identifier (supports "provider/model" and aliases).
func WithModel(model string) Option {
	return func(c *Config) { c.Model = model }
}

// WithBaseURL overrides the provider or proxy base URL.
func WithBaseURL(url string) Option {
	return func(c *Config) { c.BaseURL = url }
}

// WithAPIVersion sets Azure OpenAI api-version (and similar versioned APIs).
func WithAPIVersion(version string) Option {
	return func(c *Config) { c.APIVersion = version }
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

// WithMaxRetries sets MaxAttempts on the retry policy.
func WithMaxRetries(n int) Option {
	return func(c *Config) {
		if c.Retry.MaxAttempts == 0 {
			c.Retry = retry.DefaultPolicy()
		}
		c.Retry.MaxAttempts = n
	}
}

// WithFallback appends LiteLLM-style model refs tried when the primary model fails.
func WithFallback(models ...string) Option {
	return func(c *Config) {
		c.Fallbacks = append(c.Fallbacks, models...)
	}
}

// WithAliases replaces the alias map (nil disables aliases).
func WithAliases(aliases *model.AliasMap) Option {
	return func(c *Config) { c.Aliases = aliases }
}

// WithSuccessCallback registers a LiteLLM-style success callback.
func WithSuccessCallback(h callback.Handler) Option {
	return func(c *Config) {
		if c.Callbacks == nil {
			c.Callbacks = &callback.Bus{}
		}
		c.Callbacks.OnSuccess(h)
	}
}

// WithFailureCallback registers a LiteLLM-style failure callback.
func WithFailureCallback(h callback.Handler) Option {
	return func(c *Config) {
		if c.Callbacks == nil {
			c.Callbacks = &callback.Bus{}
		}
		c.Callbacks.OnFailure(h)
	}
}

// WithCache sets a LiteLLM-style response cache (nil disables).
func WithCache(c cache.Cache) Option {
	return func(cfg *Config) { cfg.Cache = c }
}

// WithMemoryCache installs an in-memory cache (LiteLLM Cache() default).
func WithMemoryCache(ttl time.Duration) Option {
	return func(cfg *Config) {
		cfg.Cache = cache.NewMemory(ttl)
		cfg.CacheTTL = ttl
	}
}

// WithCacheTTL sets the TTL used when writing cache entries.
func WithCacheTTL(ttl time.Duration) Option {
	return func(cfg *Config) { cfg.CacheTTL = ttl }
}

// WithDiskCache installs a disk-backed cache under dir (LiteLLM Cache(type="disk")).
func WithDiskCache(dir string, ttl time.Duration) Option {
	return func(cfg *Config) {
		d, err := cache.NewDisk(dir, ttl)
		if err != nil {
			return
		}
		cfg.Cache = d
		cfg.CacheTTL = ttl
	}
}

// WithOrganization sets the OpenAI-Organization header.
func WithOrganization(org string) Option {
	return func(c *Config) { c.Organization = org }
}

// WithExtraHeaders merges extra HTTP headers sent to providers.
func WithExtraHeaders(headers map[string]string) Option {
	return func(c *Config) {
		if c.ExtraHeaders == nil {
			c.ExtraHeaders = map[string]string{}
		}
		for k, v := range headers {
			c.ExtraHeaders[k] = v
		}
	}
}

// WithDropParams enables LiteLLM drop_params behavior for unsupported OpenAI fields.
func WithDropParams(drop bool) Option {
	return func(c *Config) { c.DropParams = drop }
}

// WithUseLiteLLMProxy routes all requests through a LiteLLM Proxy gateway.
func WithUseLiteLLMProxy(enabled bool) Option {
	return func(c *Config) { c.UseLiteLLMProxy = enabled }
}

// WithRedisCache installs a Redis/Valkey response cache (LiteLLM Cache(type="redis")).
func WithRedisCache(addr string, ttl time.Duration, opts ...cache.RedisOption) Option {
	return func(cfg *Config) {
		cfg.Cache = cache.NewRedis(addr, ttl, opts...)
		cfg.CacheTTL = ttl
	}
}
