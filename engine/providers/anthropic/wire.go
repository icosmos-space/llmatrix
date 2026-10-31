package anthropic

import "github.com/icosmos-space/llmatrix/engine/types"

type wireMessage struct {
	Role    string        `json:"role"`
	Content []wireContent `json:"content"`
}

type wireContent struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

type wireRequest struct {
	Model       string        `json:"model"`
	Messages    []wireMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens"`
	System      string        `json:"system,omitempty"`
	Temperature *float64      `json:"temperature,omitempty"`
	Stream      bool          `json:"stream,omitempty"`
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

func splitSystem(msgs []types.Message) (system string, rest []types.Message) {
	rest = make([]types.Message, 0, len(msgs))
	for _, m := range msgs {
		if m.Role == types.RoleSystem {
			if system != "" {
				system += "\n"
			}
			system += m.Content
			continue
		}
		rest = append(rest, m)
	}
	return system, rest
}

func toWireMessages(msgs []types.Message) []wireMessage {
	out := make([]wireMessage, 0, len(msgs))
	for _, m := range msgs {
		role := string(m.Role)
		if role != "user" && role != "assistant" {
			// Anthropic only accepts user/assistant in messages array.
			role = "user"
		}
		out = append(out, wireMessage{
			Role: role,
			Content: []wireContent{{
				Type: "text",
				Text: m.Content,
			}},
		})
	}
	return out
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
