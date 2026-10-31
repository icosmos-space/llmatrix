package model

import "sync"

// AliasMap maps short model names to LiteLLM-style provider/model refs.
type AliasMap struct {
	mu   sync.RWMutex
	data map[string]string
}

// DefaultAliases returns common short-name aliases (LiteLLM-inspired).
func DefaultAliases() *AliasMap {
	return &AliasMap{data: map[string]string{
		"gpt-4o":                 "openai/gpt-4o",
		"gpt-4o-mini":            "openai/gpt-4o-mini",
		"gpt-4.1":                "openai/gpt-4.1",
		"gpt-4.1-mini":           "openai/gpt-4.1-mini",
		"o1":                     "openai/o1",
		"o1-mini":                "openai/o1-mini",
		"o3-mini":                "openai/o3-mini",
		"text-embedding-3-small": "openai/text-embedding-3-small",
		"text-embedding-3-large": "openai/text-embedding-3-large",
		"dall-e-3":               "openai/dall-e-3",
		"tts-1":                  "openai/tts-1",
		"whisper-1":              "openai/whisper-1",
		"omni-moderation-latest": "openai/omni-moderation-latest",
		"text-moderation-stable": "openai/text-moderation-stable",
		"gpt-image-1":            "openai/gpt-image-1",
		"claude-3-5-sonnet":      "anthropic/claude-3-5-sonnet-latest",
		"claude-sonnet-4-5":      "anthropic/claude-sonnet-4-5",
		"claude-opus-4":          "anthropic/claude-opus-4",
		"gemini-2.0-flash":       "gemini/gemini-2.0-flash",
		"gemini-1.5-pro":         "gemini/gemini-1.5-pro",
		"deepseek-chat":          "deepseek/deepseek-chat",
		"deepseek-reasoner":      "deepseek/deepseek-reasoner",
		"llama3":                 "ollama/llama3",
		"llama-3.3-70b":          "groq/llama-3.3-70b-versatile",
		"mixtral-8x7b":           "groq/mixtral-8x7b-32768",
		"mistral-large":          "mistral/mistral-large-latest",
		"mistral-small":          "mistral/mistral-small-latest",
		"grok-2":                 "xai/grok-2",
		"grok-2-latest":          "xai/grok-2-latest",
		"gpt-3.5-turbo-instruct": "openai/gpt-3.5-turbo-instruct",
		// Azure: use azure/<deployment>; no short aliases by default.
		// Together / OpenRouter / Fireworks: use provider/<model-id>.
	}}
}

// Set registers or overrides an alias.
func (a *AliasMap) Set(alias, target string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.data == nil {
		a.data = make(map[string]string)
	}
	a.data[alias] = target
}

// Expand replaces a known alias with its provider/model target.
// Unknown names are returned unchanged.
func (a *AliasMap) Expand(name string) string {
	if a == nil {
		return name
	}
	a.mu.RLock()
	defer a.mu.RUnlock()
	if target, ok := a.data[name]; ok {
		return target
	}
	return name
}

// ResolveAlias expands aliases then parses provider/model.
func ResolveAlias(aliases *AliasMap, raw, providerHint string) Ref {
	if aliases != nil {
		raw = aliases.Expand(raw)
	}
	return Resolve(raw, providerHint)
}
