package openaicompatible

import (
	"net/http"

	"github.com/icosmos-space/llmatrix/engine/provider"
)

// Provider talks to any OpenAI-compatible HTTP API
// (OpenAI, DeepSeek, Ollama, Azure-compatible gateways, etc.).
type Provider struct {
	Name       string
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
}

// New constructs an OpenAI-compatible provider with defaults applied.
func New(name, apiKey, baseURL string) *Provider {
	if name == "" {
		name = "openai-compatible"
	}
	return &Provider{
		Name:    name,
		APIKey:  apiKey,
		BaseURL: baseURL,
	}
}

// Info implements provider.Provider.
func (p *Provider) Info() provider.Info {
	name := p.Name
	if name == "" {
		name = "openai-compatible"
	}
	display := name
	switch name {
	case "openai":
		display = "OpenAI"
	case "deepseek":
		display = "DeepSeek"
	case "ollama":
		display = "Ollama"
	case "groq":
		display = "Groq"
	case "together":
		display = "Together AI"
	case "mistral":
		display = "Mistral"
	}
	caps := provider.Capabilities{
		Chat: true, Streaming: true, Tools: true, Embeddings: true,
	}
	// Full multimodal surface for OpenAI and generic compatible gateways.
	switch name {
	case "openai", "openai-compatible":
		caps.Images = true
		caps.Speech = true
		caps.Transcription = true
		caps.Rerank = true
		caps.Responses = true
	case "deepseek", "groq", "together", "mistral":
		// chat-focused OpenAI-compatible APIs
	case "ollama":
		// chat + embeddings typically
	}
	return provider.Info{
		Name:         name,
		DisplayName:  display,
		Capabilities: caps,
	}
}
