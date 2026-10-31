package sdk

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

// WithModel sets the default model identifier.
func WithModel(model string) Option {
	return func(c *Config) { c.Model = model }
}

// WithBaseURL overrides the provider or proxy base URL.
func WithBaseURL(url string) Option {
	return func(c *Config) { c.BaseURL = url }
}
