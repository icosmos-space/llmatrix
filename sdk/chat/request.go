package chat

import "github.com/icosmos-space/llmatrix/sdk/tools"

// Message is a single chat turn.
type Message struct {
	Role       string       `json:"role"`
	Content    string       `json:"content"`
	Name       string       `json:"name,omitempty"`
	ToolCallID string       `json:"tool_call_id,omitempty"`
	ToolCalls  []tools.Call `json:"tool_calls,omitempty"`
}

// Request is a chat completion request (OpenAI / LiteLLM shaped).
type Request struct {
	Model       string       `json:"model,omitempty"`
	Messages    []Message    `json:"messages"`
	Temperature *float64     `json:"temperature,omitempty"`
	MaxTokens   *int         `json:"max_tokens,omitempty"`
	Stream      bool         `json:"stream,omitempty"`
	Tools       []tools.Tool `json:"tools,omitempty"`
}
