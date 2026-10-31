package openaicompatible

import "github.com/icosmos-space/llmatrix/engine/types"

type wireMessage struct {
	Role       string         `json:"role"`
	Content    any            `json:"content"` // string or []content parts (prompt cache)
	Name       string         `json:"name,omitempty"`
	ToolCallID string         `json:"tool_call_id,omitempty"`
	ToolCalls  []wireToolCall `json:"tool_calls,omitempty"`
}

type wireContentPart struct {
	Type         string              `json:"type"`
	Text         string              `json:"text,omitempty"`
	CacheControl *types.CacheControl `json:"cache_control,omitempty"`
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
	Model             string         `json:"model"`
	Messages          []wireMessage  `json:"messages"`
	Temperature       *float64       `json:"temperature,omitempty"`
	TopP              *float64       `json:"top_p,omitempty"`
	N                 *int           `json:"n,omitempty"`
	Stop              []string       `json:"stop,omitempty"`
	MaxTokens         *int           `json:"max_tokens,omitempty"`
	PresencePenalty   *float64       `json:"presence_penalty,omitempty"`
	FrequencyPenalty  *float64       `json:"frequency_penalty,omitempty"`
	LogitBias         map[string]int `json:"logit_bias,omitempty"`
	User              string         `json:"user,omitempty"`
	Seed              *int           `json:"seed,omitempty"`
	ResponseFormat    any            `json:"response_format,omitempty"`
	Stream            bool           `json:"stream,omitempty"`
	StreamOptions     map[string]any `json:"stream_options,omitempty"`
	Tools             []wireTool     `json:"tools,omitempty"`
	ToolChoice        any            `json:"tool_choice,omitempty"`
	ParallelToolCalls *bool          `json:"parallel_tool_calls,omitempty"`
	PromptCacheKey    string         `json:"prompt_cache_key,omitempty"`
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
	Object  string `json:"object"`
	Created int64  `json:"created"`
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
	Model          string   `json:"model"`
	Input          []string `json:"input"`
	Dimensions     *int     `json:"dimensions,omitempty"`
	EncodingFormat string   `json:"encoding_format,omitempty"`
	User           string   `json:"user,omitempty"`
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

func toWireChatRequest(req *types.ChatRequest, stream bool) wireChatRequest {
	return wireChatRequest{
		Model:             req.Model,
		Messages:          toWireMessages(req.Messages),
		Temperature:       req.Temperature,
		TopP:              req.TopP,
		N:                 req.N,
		Stop:              req.Stop,
		MaxTokens:         req.MaxTokens,
		PresencePenalty:   req.PresencePenalty,
		FrequencyPenalty:  req.FrequencyPenalty,
		LogitBias:         req.LogitBias,
		User:              req.User,
		Seed:              req.Seed,
		ResponseFormat:    req.ResponseFormat,
		Stream:            stream,
		StreamOptions:     req.StreamOptions,
		Tools:             toWireTools(req.Tools),
		ToolChoice:        req.ToolChoice,
		ParallelToolCalls: req.ParallelToolCalls,
		PromptCacheKey:    req.PromptCacheKey,
	}
}

func toWireMessages(msgs []types.Message) []wireMessage {
	out := make([]wireMessage, 0, len(msgs))
	for _, m := range msgs {
		wm := wireMessage{
			Role:       string(m.Role),
			Content:    wireMessageContent(m),
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

func wireMessageContent(m types.Message) any {
	if len(m.Parts) > 0 {
		parts := make([]wireContentPart, 0, len(m.Parts))
		for _, part := range m.Parts {
			typ := part.Type
			if typ == "" {
				typ = "text"
			}
			parts = append(parts, wireContentPart{
				Type: typ, Text: part.Text, CacheControl: part.CacheControl,
			})
		}
		return parts
	}
	if m.CacheControl != nil {
		return []wireContentPart{{
			Type: "text", Text: m.Content, CacheControl: m.CacheControl,
		}}
	}
	return m.Content
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
		Content:    contentAsString(m.Content),
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

func contentAsString(v any) string {
	switch c := v.(type) {
	case string:
		return c
	case nil:
		return ""
	default:
		return ""
	}
}
