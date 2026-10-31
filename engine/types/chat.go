package types

// ChatRequest is the engine-level chat completion request (OpenAI / LiteLLM shaped).
type ChatRequest struct {
	Model               string         `json:"model"`
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
	Tools               []Tool         `json:"tools,omitempty"`
	ToolChoice          any            `json:"tool_choice,omitempty"`
	ParallelToolCalls   *bool          `json:"parallel_tool_calls,omitempty"`
	Stream              bool           `json:"stream,omitempty"`
	StreamOptions       map[string]any `json:"stream_options,omitempty"`
	Logprobs            *bool          `json:"logprobs,omitempty"`
	TopLogprobs         *int           `json:"top_logprobs,omitempty"`
	ReasoningEffort     string         `json:"reasoning_effort,omitempty"` // none|minimal|low|medium|high|...
	// PromptCacheKey is an OpenAI prompt caching key (when supported by the API).
	PromptCacheKey string `json:"prompt_cache_key,omitempty"`
	// Thinking is Anthropic extended thinking (LiteLLM thinking=...).
	Thinking any `json:"thinking,omitempty"`
}

// ChatChoice is one completion candidate.
type ChatChoice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason,omitempty"`
}

// ChatResponse is the engine-level chat completion response.
type ChatResponse struct {
	ID       string       `json:"id,omitempty"`
	Object   string       `json:"object,omitempty"`
	Created  int64        `json:"created,omitempty"`
	Model    string       `json:"model,omitempty"`
	Provider string       `json:"provider,omitempty"`
	Choices  []ChatChoice `json:"choices"`
	Usage    Usage        `json:"usage"`
}
