package bedrock

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/icosmos-space/llmatrix/engine/errors"
	"github.com/icosmos-space/llmatrix/engine/types"
)

type converseRequest struct {
	Messages        []converseMessage `json:"messages"`
	System          []systemBlock     `json:"system,omitempty"`
	InferenceConfig *inferenceConfig  `json:"inferenceConfig,omitempty"`
	ToolConfig      *toolConfig       `json:"toolConfig,omitempty"`
}

type systemBlock struct {
	Text string `json:"text"`
}

type converseMessage struct {
	Role    string         `json:"role"`
	Content []contentBlock `json:"content"`
}

type contentBlock struct {
	Text       string      `json:"text,omitempty"`
	ToolUse    *toolUse    `json:"toolUse,omitempty"`
	ToolResult *toolResult `json:"toolResult,omitempty"`
}

type toolUse struct {
	ToolUseID string          `json:"toolUseId"`
	Name      string          `json:"name"`
	Input     json.RawMessage `json:"input"`
}

type toolResult struct {
	ToolUseID string              `json:"toolUseId"`
	Content   []toolResultContent `json:"content"`
	Status    string              `json:"status,omitempty"`
}

type toolResultContent struct {
	Text string `json:"text,omitempty"`
}

type inferenceConfig struct {
	MaxTokens     *int     `json:"maxTokens,omitempty"`
	Temperature   *float64 `json:"temperature,omitempty"`
	TopP          *float64 `json:"topP,omitempty"`
	StopSequences []string `json:"stopSequences,omitempty"`
}

type toolConfig struct {
	Tools []toolSpec `json:"tools"`
}

type toolSpec struct {
	ToolSpec toolSpecBody `json:"toolSpec"`
}

type toolSpecBody struct {
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	InputSchema json.RawMessage `json:"inputSchema"`
}

type converseResponse struct {
	Output     converseOutput `json:"output"`
	StopReason string         `json:"stopReason"`
	Usage      usage          `json:"usage"`
}

type converseOutput struct {
	Message converseMessage `json:"message"`
}

type usage struct {
	InputTokens  int `json:"inputTokens"`
	OutputTokens int `json:"outputTokens"`
	TotalTokens  int `json:"totalTokens"`
}

func toConverseRequest(req *types.ChatRequest) (*converseRequest, error) {
	out := &converseRequest{}
	for _, m := range req.Messages {
		role := strings.ToLower(string(m.Role))
		switch role {
		case "system", "developer":
			text := messageText(m)
			if text != "" {
				out.System = append(out.System, systemBlock{Text: text})
			}
		default:
			cm, err := toConverseMessage(m)
			if err != nil {
				return nil, err
			}
			out.Messages = append(out.Messages, *cm)
		}
	}
	if len(out.Messages) == 0 {
		return nil, errors.Wrap(errors.KindInvalidRequest, "bedrock", "no non-system messages", nil)
	}
	cfg := &inferenceConfig{}
	has := false
	if req.MaxTokens != nil {
		cfg.MaxTokens = req.MaxTokens
		has = true
	} else if req.MaxCompletionTokens != nil {
		cfg.MaxTokens = req.MaxCompletionTokens
		has = true
	}
	if req.Temperature != nil {
		cfg.Temperature = req.Temperature
		has = true
	}
	if req.TopP != nil {
		cfg.TopP = req.TopP
		has = true
	}
	if len(req.Stop) > 0 {
		cfg.StopSequences = req.Stop
		has = true
	}
	if has {
		out.InferenceConfig = cfg
	}
	if len(req.Tools) > 0 {
		tools := make([]toolSpec, 0, len(req.Tools))
		for _, t := range req.Tools {
			params := any(t.Function.Parameters)
			if params == nil {
				params = map[string]any{"type": "object", "properties": map[string]any{}}
			}
			schema, err := json.Marshal(map[string]any{"json": params})
			if err != nil {
				return nil, errors.Wrap(errors.KindInvalidRequest, "bedrock", "marshal tool schema", err)
			}
			tools = append(tools, toolSpec{ToolSpec: toolSpecBody{
				Name:        t.Function.Name,
				Description: t.Function.Description,
				InputSchema: schema,
			}})
		}
		out.ToolConfig = &toolConfig{Tools: tools}
	}
	return out, nil
}

func toConverseMessage(m types.Message) (*converseMessage, error) {
	role := strings.ToLower(string(m.Role))
	if role == "tool" {
		return &converseMessage{
			Role: "user",
			Content: []contentBlock{{
				ToolResult: &toolResult{
					ToolUseID: m.ToolCallID,
					Content:   []toolResultContent{{Text: messageText(m)}},
				},
			}},
		}, nil
	}
	cm := &converseMessage{Role: role}
	if role == "assistant" && len(m.ToolCalls) > 0 {
		if text := messageText(m); text != "" {
			cm.Content = append(cm.Content, contentBlock{Text: text})
		}
		for _, tc := range m.ToolCalls {
			input := json.RawMessage(tc.Arguments)
			if len(input) == 0 || !json.Valid(input) {
				input = json.RawMessage(`{}`)
			}
			cm.Content = append(cm.Content, contentBlock{
				ToolUse: &toolUse{ToolUseID: tc.ID, Name: tc.Name, Input: input},
			})
		}
		return cm, nil
	}
	cm.Content = []contentBlock{{Text: messageText(m)}}
	return cm, nil
}

func messageText(m types.Message) string {
	if len(m.Parts) > 0 {
		var b strings.Builder
		for _, p := range m.Parts {
			if p.Text != "" {
				b.WriteString(p.Text)
			}
		}
		return b.String()
	}
	return m.Content
}

func fromConverseResponse(model string, raw *converseResponse) *types.ChatResponse {
	msg := types.Message{Role: types.Role(raw.Output.Message.Role)}
	var texts []string
	for _, block := range raw.Output.Message.Content {
		if block.Text != "" {
			texts = append(texts, block.Text)
		}
		if block.ToolUse != nil {
			args := string(block.ToolUse.Input)
			if args == "" {
				args = "{}"
			}
			msg.ToolCalls = append(msg.ToolCalls, types.ToolCall{
				ID:        block.ToolUse.ToolUseID,
				Type:      "function",
				Name:      block.ToolUse.Name,
				Arguments: args,
			})
		}
	}
	if len(texts) > 0 {
		msg.Content = strings.Join(texts, "")
	}
	finish := mapStopReason(raw.StopReason)
	total := raw.Usage.TotalTokens
	if total == 0 {
		total = raw.Usage.InputTokens + raw.Usage.OutputTokens
	}
	return &types.ChatResponse{
		ID:       "bedrock-" + model,
		Object:   "chat.completion",
		Model:    model,
		Provider: "bedrock",
		Choices:  []types.ChatChoice{{Index: 0, Message: msg, FinishReason: finish}},
		Usage: types.Usage{
			PromptTokens:     raw.Usage.InputTokens,
			CompletionTokens: raw.Usage.OutputTokens,
			TotalTokens:      total,
		},
	}
}

func mapStopReason(s string) string {
	switch strings.ToLower(s) {
	case "end_turn", "stop":
		return "stop"
	case "tool_use", "tool_calls":
		return "tool_calls"
	case "max_tokens":
		return "length"
	default:
		if s == "" {
			return "stop"
		}
		return s
	}
}

func mapHTTPError(status int, body []byte) error {
	msg := strings.TrimSpace(string(body))
	if msg == "" {
		msg = fmt.Sprintf("bedrock http %d", status)
	}
	kind := errors.KindProvider
	switch {
	case status == 401 || status == 403:
		kind = errors.KindUnauthorized
	case status == 404:
		kind = errors.KindNotFound
	case status == 429:
		kind = errors.KindRateLimited
	case status >= 500:
		kind = errors.KindProvider
	case status >= 400:
		kind = errors.KindInvalidRequest
	}
	return errors.Wrap(kind, "bedrock", msg, nil)
}
