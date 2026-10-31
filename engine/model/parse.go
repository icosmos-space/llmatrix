package model

import (
	"strings"
)

// Ref is a resolved provider + model pair (LiteLLM-style `provider/model`).
type Ref struct {
	Provider string
	Model    string
}

// Parse splits a LiteLLM-style model identifier.
//
//	"openai/gpt-4o"      -> openai, gpt-4o
//	"deepseek/deepseek-chat" -> deepseek, deepseek-chat
//	"ollama/llama3"      -> ollama, llama3
//	"gpt-4o"             -> "", gpt-4o  (provider left empty for caller default)
func Parse(raw string) Ref {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return Ref{}
	}
	provider, model, ok := strings.Cut(raw, "/")
	if !ok {
		return Ref{Model: raw}
	}
	provider = strings.TrimSpace(provider)
	model = strings.TrimSpace(model)
	if provider == "" || model == "" {
		return Ref{Model: raw}
	}
	return Ref{Provider: strings.ToLower(provider), Model: model}
}

// Resolve applies an explicit provider hint when the model string has none.
func Resolve(raw, providerHint string) Ref {
	ref := Parse(raw)
	if ref.Provider == "" {
		ref.Provider = strings.ToLower(strings.TrimSpace(providerHint))
	}
	if ref.Provider == "" {
		ref.Provider = "openai"
	}
	return ref
}
