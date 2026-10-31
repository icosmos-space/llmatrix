package chat

import "github.com/icosmos-space/llmatrix/sdk/tools"

// CacheControl marks a content block for prompt caching (Anthropic-style).
type CacheControl struct {
	Type string `json:"type"` // typically "ephemeral"
}

// ContentPart is a structured message content block.
type ContentPart struct {
	Type         string        `json:"type"`
	Text         string        `json:"text,omitempty"`
	CacheControl *CacheControl `json:"cache_control,omitempty"`
}

// Message is a single chat turn.
type Message struct {
	Role         string        `json:"role"`
	Content      string        `json:"content"`
	Name         string        `json:"name,omitempty"`
	ToolCallID   string        `json:"tool_call_id,omitempty"`
	ToolCalls    []tools.Call  `json:"tool_calls,omitempty"`
	Parts        []ContentPart `json:"parts,omitempty"`
	CacheControl *CacheControl `json:"cache_control,omitempty"`
}

// Request is a chat completion request (OpenAI / LiteLLM shaped).
type Request struct {
	Model               string         `json:"model,omitempty"`
	Messages            []Message      `json:"messages"`
	Temperature         *float64       `json:"temperature,omitempty"`
	TopP                *float64       `json:"top_p,omitempty"`
	N                   *int           `json:"n,omitempty"`
	Stop                []string       `json:"stop,omitempty"`
	MaxTokens           *int           `json:"max_tokens,omitempty"`
	MaxCompletionTokens *int           `json:"max_completion_tokens,omitempty"`
	PresencePenalty     *float64       `json:"presence_penalty,omitempty"`
	FrequencyPenalty    *float64       `json:"frequency_penalty,omitempty"`
	LogitBias           map[string]int `json:"logit_bias,omitempty"`
	User                string         `json:"user,omitempty"`
	Seed                *int           `json:"seed,omitempty"`
	ResponseFormat      any            `json:"response_format,omitempty"`
	Stream              bool           `json:"stream,omitempty"`
	StreamOptions       map[string]any `json:"stream_options,omitempty"`
	Tools               []tools.Tool   `json:"tools,omitempty"`
	ToolChoice          any            `json:"tool_choice,omitempty"`
	ParallelToolCalls   *bool          `json:"parallel_tool_calls,omitempty"`
	Logprobs            *bool          `json:"logprobs,omitempty"`
	TopLogprobs         *int           `json:"top_logprobs,omitempty"`
	ReasoningEffort     string         `json:"reasoning_effort,omitempty"`
	PromptCacheKey      string         `json:"prompt_cache_key,omitempty"`
	Thinking            any            `json:"thinking,omitempty"`
}
