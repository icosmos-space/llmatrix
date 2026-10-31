package openaicompatible

import "github.com/icosmos-space/llmatrix/engine/types"

type wireMessage struct {
	Role       string         `json:"role"`
	Content    string         `json:"content"`
	Name       string         `json:"name,omitempty"`
	ToolCallID string         `json:"tool_call_id,omitempty"`
	ToolCalls  []wireToolCall `json:"tool_calls,omitempty"`
}

type wireToolCall struct {
	ID       string `json:"id,omitempty"`
	Type     string `json:"type,omitempty"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

type wireChatRequest struct {
	Model       string        `json:"model"`
	Messages    []wireMessage `json:"messages"`
	Temperature *float64      `json:"temperature,omitempty"`
	MaxTokens   *int          `json:"max_tokens,omitempty"`
	Stream      bool          `json:"stream,omitempty"`
	Tools       []wireTool    `json:"tools,omitempty"`
}

type wireTool struct {
	Type     string `json:"type"`
	Function struct {
		Name        string         `json:"name"`
		Description string         `json:"description,omitempty"`
		Parameters  map[string]any `json:"parameters,omitempty"`
	} `json:"function"`
}

type wireChatResponse struct {
	ID      string `json:"id"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int         `json:"index"`
		Message      wireMessage `json:"message"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage types.Usage `json:"usage"`
	Error *wireError  `json:"error,omitempty"`
}

type wireError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    any    `json:"code"`
}

type wireStreamChunk struct {
	ID      string `json:"id"`
	Model   string `json:"model"`
	Choices []struct {
		Index int `json:"index"`
		Delta struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"delta"`
		FinishReason *string `json:"finish_reason"`
	} `json:"choices"`
	Usage *types.Usage `json:"usage,omitempty"`
	Error *wireError   `json:"error,omitempty"`
}

type wireEmbeddingRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

type wireEmbeddingResponse struct {
	Model string `json:"model"`
	Data  []struct {
		Index     int       `json:"index"`
		Embedding []float64 `json:"embedding"`
	} `json:"data"`
	Usage types.Usage `json:"usage"`
	Error *wireError  `json:"error,omitempty"`
}

func toWireMessages(msgs []types.Message) []wireMessage {
	out := make([]wireMessage, 0, len(msgs))
	for _, m := range msgs {
		wm := wireMessage{
			Role:       string(m.Role),
			Content:    m.Content,
			Name:       m.Name,
			ToolCallID: m.ToolCallID,
		}
		for _, tc := range m.ToolCalls {
			wtc := wireToolCall{ID: tc.ID, Type: tc.Type}
			if wtc.Type == "" {
				wtc.Type = "function"
			}
			wtc.Function.Name = tc.Name
			wtc.Function.Arguments = tc.Arguments
			wm.ToolCalls = append(wm.ToolCalls, wtc)
		}
		out = append(out, wm)
	}
	return out
}

func toWireTools(tools []types.Tool) []wireTool {
	if len(tools) == 0 {
		return nil
	}
	out := make([]wireTool, 0, len(tools))
	for _, t := range tools {
		wt := wireTool{Type: t.Type}
		if wt.Type == "" {
			wt.Type = "function"
		}
		wt.Function.Name = t.Function.Name
		wt.Function.Description = t.Function.Description
		wt.Function.Parameters = t.Function.Parameters
		out = append(out, wt)
	}
	return out
}

func fromWireMessage(m wireMessage) types.Message {
	msg := types.Message{
		Role:       types.Role(m.Role),
		Content:    m.Content,
		Name:       m.Name,
		ToolCallID: m.ToolCallID,
	}
	for _, tc := range m.ToolCalls {
		msg.ToolCalls = append(msg.ToolCalls, types.ToolCall{
			ID:        tc.ID,
			Type:      tc.Type,
			Name:      tc.Function.Name,
			Arguments: tc.Function.Arguments,
		})
	}
	return msg
}
