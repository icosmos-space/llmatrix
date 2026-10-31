package model

import "strings"

// Info describes LiteLLM-style model metadata (model_info).
type Info struct {
	Model                        string
	Provider                     string
	Mode                         string // chat | embedding | image_generation | ...
	MaxInputTokens               int
	MaxOutputTokens              int
	SupportsFunctionCalling      bool
	SupportsParallelFunctionCall bool
	SupportsVision               bool
	SupportsStreaming            bool
	SupportsPromptCaching        bool
	SupportsResponseSchema       bool
	SupportsReasoningEffort      bool
}

// GetModelInfo returns known metadata for a provider/model (or alias).
// Unknown models return a conservative zero-ish Info with Provider/Model filled.
func GetModelInfo(raw string, providerHint string) Info {
	ref := ResolveAlias(DefaultAliases(), raw, providerHint)
	key := strings.ToLower(ref.Provider + "/" + ref.Model)
	if info, ok := catalog[key]; ok {
		info.Provider = ref.Provider
		info.Model = ref.Model
		return info
	}
	// Provider-level defaults.
	info := Info{
		Provider:          ref.Provider,
		Model:             ref.Model,
		Mode:              "chat",
		SupportsStreaming: true,
	}
	switch ref.Provider {
	case "openai", "azure", "openrouter", "groq", "together", "mistral", "fireworks", "xai", "deepseek":
		info.SupportsFunctionCalling = true
		info.SupportsParallelFunctionCall = true
		info.SupportsResponseSchema = true
	case "anthropic":
		info.SupportsFunctionCalling = true
		info.SupportsPromptCaching = true
	case "gemini", "google", "vertex_ai":
		info.SupportsFunctionCalling = true
		info.SupportsVision = true
	case "ollama":
		// tools/vision vary by model; leave false by default
	}
	if strings.Contains(strings.ToLower(ref.Model), "vision") ||
		strings.Contains(strings.ToLower(ref.Model), "gpt-4o") ||
		strings.Contains(strings.ToLower(ref.Model), "claude-3") ||
		strings.Contains(strings.ToLower(ref.Model), "claude-sonnet") ||
		strings.Contains(strings.ToLower(ref.Model), "claude-opus") ||
		strings.Contains(strings.ToLower(ref.Model), "gemini") {
		info.SupportsVision = true
	}
	return info
}

// SupportsFunctionCalling mirrors litellm.supports_function_calling.
func SupportsFunctionCalling(model string, providerHint ...string) bool {
	hint := ""
	if len(providerHint) > 0 {
		hint = providerHint[0]
	}
	return GetModelInfo(model, hint).SupportsFunctionCalling
}

// SupportsParallelFunctionCalling mirrors litellm.supports_parallel_function_calling.
func SupportsParallelFunctionCalling(model string, providerHint ...string) bool {
	hint := ""
	if len(providerHint) > 0 {
		hint = providerHint[0]
	}
	return GetModelInfo(model, hint).SupportsParallelFunctionCall
}

// SupportsVision mirrors litellm.supports_vision.
func SupportsVision(model string, providerHint ...string) bool {
	hint := ""
	if len(providerHint) > 0 {
		hint = providerHint[0]
	}
	return GetModelInfo(model, hint).SupportsVision
}

// GetSupportedOpenAIParams mirrors litellm.get_supported_openai_params for chat.
func GetSupportedOpenAIParams(model string, providerHint ...string) []string {
	hint := ""
	if len(providerHint) > 0 {
		hint = providerHint[0]
	}
	ref := ResolveAlias(DefaultAliases(), model, hint)
	base := []string{
		"temperature", "top_p", "stream", "stop", "max_tokens", "tools", "tool_choice", "user",
	}
	switch ref.Provider {
	case "openai", "azure", "openrouter", "groq", "together", "mistral", "fireworks", "xai", "deepseek", "openai-compatible":
		return append(base,
			"n", "presence_penalty", "frequency_penalty", "logit_bias", "seed",
			"response_format", "logprobs", "top_logprobs", "parallel_tool_calls",
			"max_completion_tokens", "stream_options", "reasoning_effort", "prompt_cache_key",
		)
	case "anthropic":
		return []string{
			"temperature", "top_p", "stream", "stop", "max_tokens", "tools", "tool_choice",
			"user", "response_format", "thinking",
		}
	case "gemini", "google", "vertex_ai":
		return []string{
			"temperature", "top_p", "stream", "stop", "max_tokens", "tools", "tool_choice",
			"user", "response_format", "n", "seed",
		}
	case "ollama":
		return []string{"temperature", "top_p", "stream", "stop", "max_tokens", "seed", "response_format"}
	default:
		return base
	}
}

var catalog = map[string]Info{
	"openai/gpt-4o": {
		Mode: "chat", MaxInputTokens: 128000, MaxOutputTokens: 16384,
		SupportsFunctionCalling: true, SupportsParallelFunctionCall: true,
		SupportsVision: true, SupportsStreaming: true, SupportsResponseSchema: true,
	},
	"openai/gpt-4o-mini": {
		Mode: "chat", MaxInputTokens: 128000, MaxOutputTokens: 16384,
		SupportsFunctionCalling: true, SupportsParallelFunctionCall: true,
		SupportsVision: true, SupportsStreaming: true, SupportsResponseSchema: true,
	},
	"openai/o1-mini": {
		Mode: "chat", MaxInputTokens: 128000, MaxOutputTokens: 65536,
		SupportsFunctionCalling: true, SupportsStreaming: true, SupportsReasoningEffort: true,
	},
	"openai/o3-mini": {
		Mode: "chat", MaxInputTokens: 200000, MaxOutputTokens: 100000,
		SupportsFunctionCalling: true, SupportsStreaming: true, SupportsReasoningEffort: true,
	},
	"openai/gpt-3.5-turbo-instruct": {
		Mode: "completion", MaxInputTokens: 4096, MaxOutputTokens: 4096, SupportsStreaming: true,
	},
	"openai/text-embedding-3-small": {Mode: "embedding", MaxInputTokens: 8191},
	"openai/text-embedding-3-large": {Mode: "embedding", MaxInputTokens: 8191},
	"anthropic/claude-sonnet-4-5": {
		Mode: "chat", MaxInputTokens: 200000, MaxOutputTokens: 64000,
		SupportsFunctionCalling: true, SupportsVision: true, SupportsStreaming: true,
		SupportsPromptCaching: true,
	},
	"anthropic/claude-opus-4": {
		Mode: "chat", MaxInputTokens: 200000, MaxOutputTokens: 32000,
		SupportsFunctionCalling: true, SupportsVision: true, SupportsStreaming: true,
		SupportsPromptCaching: true,
	},
	"gemini/gemini-2.0-flash": {
		Mode: "chat", MaxInputTokens: 1048576, MaxOutputTokens: 8192,
		SupportsFunctionCalling: true, SupportsVision: true, SupportsStreaming: true,
	},
	"xai/grok-2": {
		Mode: "chat", SupportsFunctionCalling: true, SupportsParallelFunctionCall: true, SupportsStreaming: true,
	},
	"xai/grok-2-latest": {
		Mode: "chat", SupportsFunctionCalling: true, SupportsParallelFunctionCall: true, SupportsStreaming: true,
	},
	"deepseek/deepseek-chat": {
		Mode: "chat", SupportsFunctionCalling: true, SupportsParallelFunctionCall: true, SupportsStreaming: true,
	},
}
