package openaicompatible

import "os"

// Defaults for well-known OpenAI-compatible providers.
type Defaults struct {
	BaseURL string
	EnvKey  string
}

var known = map[string]Defaults{
	"openai": {
		BaseURL: "https://api.openai.com/v1",
		EnvKey:  "OPENAI_API_KEY",
	},
	"deepseek": {
		BaseURL: "https://api.deepseek.com/v1",
		EnvKey:  "DEEPSEEK_API_KEY",
	},
	"ollama": {
		BaseURL: "http://localhost:11434/v1",
		EnvKey:  "",
	},
	"groq": {
		BaseURL: "https://api.groq.com/openai/v1",
		EnvKey:  "GROQ_API_KEY",
	},
	"together": {
		BaseURL: "https://api.together.xyz/v1",
		EnvKey:  "TOGETHER_API_KEY",
	},
	"mistral": {
		BaseURL: "https://api.mistral.ai/v1",
		EnvKey:  "MISTRAL_API_KEY",
	},
	"openrouter": {
		BaseURL: "https://openrouter.ai/api/v1",
		EnvKey:  "OPENROUTER_API_KEY",
	},
	"xai": {
		BaseURL: "https://api.x.ai/v1",
		EnvKey:  "XAI_API_KEY",
	},
	"fireworks": {
		BaseURL: "https://api.fireworks.ai/inference/v1",
		EnvKey:  "FIREWORKS_API_KEY",
	},
	"openai-compatible": {
		BaseURL: "",
		EnvKey:  "OPENAI_API_KEY",
	},
}

// LookupDefaults returns known defaults for a provider name.
func LookupDefaults(name string) (Defaults, bool) {
	d, ok := known[name]
	return d, ok
}

// ResolveAPIKey prefers an explicit key, then the provider env var.
func ResolveAPIKey(name, explicit string) string {
	if explicit != "" {
		return explicit
	}
	if d, ok := known[name]; ok && d.EnvKey != "" {
		return os.Getenv(d.EnvKey)
	}
	return os.Getenv("OPENAI_API_KEY")
}

// ResolveBaseURL prefers an explicit URL, then the provider default.
func ResolveBaseURL(name, explicit string) string {
	if explicit != "" {
		return explicit
	}
	if d, ok := known[name]; ok {
		return d.BaseURL
	}
	return ""
}
