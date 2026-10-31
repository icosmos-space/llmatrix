package anthropic

import (
	"encoding/json"

	"github.com/icosmos-space/llmatrix/engine/types"
)

type wireMessage struct {
	Role    string        `json:"role"`
	Content []wireContent `json:"content"`
}

type wireContent struct {
	Type         string              `json:"type"`
	Text         string              `json:"text,omitempty"`
	ID           string              `json:"id,omitempty"`
	Name         string              `json:"name,omitempty"`
	Input        json.RawMessage     `json:"input,omitempty"`
	ToolUseID    string              `json:"tool_use_id,omitempty"`
	Content      string              `json:"content,omitempty"` // tool_result payload
	CacheControl *types.CacheControl `json:"cache_control,omitempty"`
}

type wireTool struct {
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	InputSchema map[string]any `json:"input_schema"`
}

type wireRequest struct {
	Model       string        `json:"model"`
	Messages    []wireMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens"`
	System      any           `json:"system,omitempty"` // string or []wireContent (prompt cache)
	Temperature *float64      `json:"temperature,omitempty"`
	Stream      bool          `json:"stream,omitempty"`
	Tools       []wireTool    `json:"tools,omitempty"`
	ToolChoice  any           `json:"tool_choice,omitempty"`
	Thinking    any           `json:"thinking,omitempty"` // Anthropic extended thinking
}

type wireResponse struct {
	ID         string        `json:"id"`
	Type       string        `json:"type"`
	Role       string        `json:"role"`
	Model      string        `json:"model"`
	Content    []wireContent `json:"content"`
	StopReason string        `json:"stop_reason"`
	Usage      struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
	Error *wireError `json:"error,omitempty"`
}

type wireError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type wireStreamEvent struct {
	Type  string `json:"type"`
	Index int    `json:"index,omitempty"`
	Delta *struct {
		Type        string `json:"type"`
		Text        string `json:"text,omitempty"`
		StopReason  string `json:"stop_reason,omitempty"`
		PartialJSON string `json:"partial_json,omitempty"`
	} `json:"delta,omitempty"`
	Message      *wireResponse `json:"message,omitempty"`
	ContentBlock *struct {
		Type string `json:"type"`
		Text string `json:"text,omitempty"`
	} `json:"content_block,omitempty"`
	Usage *struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage,omitempty"`
	Error *wireError `json:"error,omitempty"`
}

func splitSystem(msgs []types.Message) (system any, rest []types.Message) {
	var blocks []wireContent
	var plain string
	useBlocks := false
	rest = make([]types.Message, 0, len(msgs))
	for _, m := range msgs {
		if m.Role != types.RoleSystem {
			rest = append(rest, m)
			continue
		}
		if len(m.Parts) > 0 || m.CacheControl != nil {
			useBlocks = true
		}
		blocks = append(blocks, messageTextBlocks(m)...)
		if m.Content != "" {
			if plain != "" {
				plain += "\n"
			}
			plain += m.Content
		} else {
			for _, part := range m.Parts {
				if part.Type == "text" || part.Type == "" {
					if plain != "" {
						plain += "\n"
					}
					plain += part.Text
				}
			}
		}
	}
	if useBlocks {
		return blocks, rest
	}
	if plain == "" {
		return nil, rest
	}
	return plain, rest
}

func messageTextBlocks(m types.Message) []wireContent {
	if len(m.Parts) > 0 {
		out := make([]wireContent, 0, len(m.Parts))
		for _, part := range m.Parts {
			typ := part.Type
			if typ == "" {
				typ = "text"
			}
			out = append(out, wireContent{
				Type:         typ,
				Text:         part.Text,
				CacheControl: part.CacheControl,
			})
		}
		return out
	}
	return []wireContent{{
		Type:         "text",
		Text:         m.Content,
		CacheControl: m.CacheControl,
	}}
}

func toWireMessages(msgs []types.Message) []wireMessage {
	out := make([]wireMessage, 0, len(msgs))
	for _, m := range msgs {
		switch m.Role {
		case types.RoleTool:
			out = append(out, wireMessage{
				Role: "user",
				Content: []wireContent{{
					Type:      "tool_result",
					ToolUseID: m.ToolCallID,
					Content:   m.Content,
				}},
			})
		case types.RoleAssistant:
			blocks := make([]wireContent, 0, 1+len(m.ToolCalls))
			if len(m.Parts) > 0 || m.Content != "" || m.CacheControl != nil {
				blocks = append(blocks, messageTextBlocks(m)...)
			}
			for _, tc := range m.ToolCalls {
				input := json.RawMessage(tc.Arguments)
				if !json.Valid(input) {
					b, _ := json.Marshal(tc.Arguments)
					input = b
				}
				blocks = append(blocks, wireContent{
					Type:  "tool_use",
					ID:    tc.ID,
					Name:  tc.Name,
					Input: input,
				})
			}
			if len(blocks) == 0 {
				blocks = append(blocks, wireContent{Type: "text", Text: ""})
			}
			out = append(out, wireMessage{Role: "assistant", Content: blocks})
		default:
			out = append(out, wireMessage{
				Role:    "user",
				Content: messageTextBlocks(m),
			})
		}
	}
	return out
}

func toWireTools(tools []types.Tool) []wireTool {
	if len(tools) == 0 {
		return nil
	}
	out := make([]wireTool, 0, len(tools))
	for _, t := range tools {
		schema := t.Function.Parameters
		if schema == nil {
			schema = map[string]any{"type": "object", "properties": map[string]any{}}
		}
		out = append(out, wireTool{
			Name:        t.Function.Name,
			Description: t.Function.Description,
			InputSchema: schema,
		})
	}
	return out
}

func toWireToolChoice(choice any) any {
	if choice == nil {
		return nil
	}
	if s, ok := choice.(string); ok {
		switch s {
		case "auto", "none", "any", "required":
			if s == "required" {
				return map[string]any{"type": "any"}
			}
			return map[string]any{"type": s}
		}
	}
	return choice
}

func contentText(blocks []wireContent) string {
	var s string
	for _, b := range blocks {
		if b.Type == "text" || b.Type == "" {
			s += b.Text
		}
	}
	return s
}

func contentToolCalls(blocks []wireContent) []types.ToolCall {
	var out []types.ToolCall
	for _, b := range blocks {
		if b.Type != "tool_use" {
			continue
		}
		args := string(b.Input)
		if args == "" {
			args = "{}"
		}
		out = append(out, types.ToolCall{
			ID:        b.ID,
			Type:      "function",
			Name:      b.Name,
			Arguments: args,
		})
	}
	return out
}

func mapStopReason(reason string) string {
	switch reason {
	case "end_turn":
		return "stop"
	case "max_tokens":
		return "length"
	case "tool_use":
		return "tool_calls"
	default:
		return reason
	}
}
