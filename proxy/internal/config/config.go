package config

import (
	"fmt"
	"os"
)

// Config is the proxy process configuration.
type Config struct {
	ListenAddr string
	APIKey     string
}

// Load reads configuration. YAML parsing will be added later; for now it
// applies sensible defaults and optional env overrides.
func Load(path string) (*Config, error) {
	cfg := &Config{
		ListenAddr: ":8080",
	}
	if v := os.Getenv("LLMATRIX_LISTEN"); v != "" {
		cfg.ListenAddr = v
	}
	if v := os.Getenv("LLMATRIX_API_KEY"); v != "" {
		cfg.APIKey = v
	}
	if path != "" {
		if _, err := os.Stat(path); err != nil && !os.IsNotExist(err) {
			return nil, fmt.Errorf("config: stat %s: %w", path, err)
		}
	}
	return cfg, nil
}
