package types

// ChatRequest is the engine-level chat completion request.
type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature *float64  `json:"temperature,omitempty"`
	MaxTokens   *int      `json:"max_tokens,omitempty"`
	Tools       []Tool    `json:"tools,omitempty"`
	Stream      bool      `json:"stream,omitempty"`
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
	Model    string       `json:"model,omitempty"`
	Provider string       `json:"provider,omitempty"`
	Choices  []ChatChoice `json:"choices"`
	Usage    Usage        `json:"usage"`
}
