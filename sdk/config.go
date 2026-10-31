package sdk

// Config holds client-level settings.
type Config struct {
	Provider string
	APIKey   string
	Model    string
	BaseURL  string
}

func defaultConfig() Config {
	return Config{}
}
