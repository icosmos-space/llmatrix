package responses

import (
	"context"
	"fmt"

	"github.com/icosmos-space/llmatrix/engine/provider"
	"github.com/icosmos-space/llmatrix/engine/stream"
	"github.com/icosmos-space/llmatrix/engine/types"
)

// Create calls a native ResponsesProvider when available; otherwise bridges
// through ChatCompletion (LiteLLM-style fallback for providers without /responses).
func Create(ctx context.Context, p provider.Provider, req *types.ResponsesRequest) (*types.ResponsesResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("responses: request is nil")
	}
	if rp, ok := p.(provider.ResponsesProvider); ok {
		return rp.Responses(ctx, req)
	}
	chatProv, ok := p.(provider.ChatProvider)
	if !ok {
		return nil, fmt.Errorf("responses: provider %q supports neither responses nor chat", p.Info().Name)
	}
	chatReq, err := ToChatRequest(req)
	if err != nil {
		return nil, err
	}
	chatResp, err := chatProv.ChatCompletion(ctx, chatReq)
	if err != nil {
		return nil, err
	}
	return FromChatResponse(chatResp), nil
}

// CreateStream streams native Responses SSE when available; otherwise bridges
// through ChatCompletionStream.
func CreateStream(ctx context.Context, p provider.Provider, req *types.ResponsesRequest) (stream.Stream, error) {
	if req == nil {
		return nil, fmt.Errorf("responses: request is nil")
	}
	if rp, ok := p.(provider.ResponsesStreamProvider); ok {
		return rp.ResponsesStream(ctx, req)
	}
	streamProv, ok := p.(provider.ChatStreamProvider)
	if !ok {
		return nil, fmt.Errorf("responses: provider %q supports neither responses stream nor chat stream", p.Info().Name)
	}
	chatReq, err := ToChatRequest(req)
	if err != nil {
		return nil, err
	}
	chatReq.Stream = true
	return streamProv.ChatCompletionStream(ctx, chatReq)
}

// ToChatRequest converts a Responses request into a chat completion request.
func ToChatRequest(req *types.ResponsesRequest) (*types.ChatRequest, error) {
	msgs, err := inputToMessages(req.Input, req.Instructions)
	if err != nil {
		return nil, err
	}
	return &types.ChatRequest{
		Model:       req.Model,
		Messages:    msgs,
		Temperature: req.Temperature,
		TopP:        req.TopP,
		MaxTokens:   req.MaxOutputTokens,
		User:        req.User,
		Stream:      false,
	}, nil
}

func inputToMessages(input any, instructions string) ([]types.Message, error) {
	var msgs []types.Message
	if instructions != "" {
		msgs = append(msgs, types.Message{Role: types.RoleSystem, Content: instructions})
	}
	switch v := input.(type) {
	case string:
		msgs = append(msgs, types.Message{Role: types.RoleUser, Content: v})
	case []any:
		for _, item := range v {
			switch it := item.(type) {
			case string:
				msgs = append(msgs, types.Message{Role: types.RoleUser, Content: it})
			case map[string]any:
				role, _ := it["role"].(string)
				content, _ := it["content"].(string)
				if role == "" {
					role = "user"
				}
				msgs = append(msgs, types.Message{Role: types.Role(role), Content: content})
			default:
				return nil, fmt.Errorf("responses: unsupported input item type %T", item)
			}
		}
	default:
		return nil, fmt.Errorf("responses: unsupported input type %T (want string or []any)", input)
	}
	if len(msgs) == 0 || (len(msgs) == 1 && msgs[0].Role == types.RoleSystem) {
		return nil, fmt.Errorf("responses: input produced no user/assistant messages")
	}
	return msgs, nil
}

// FromChatResponse maps a chat completion into a Responses API object.
func FromChatResponse(resp *types.ChatResponse) *types.ResponsesResponse {
	if resp == nil {
		return nil
	}
	text := ""
	if len(resp.Choices) > 0 {
		text = resp.Choices[0].Message.Content
	}
	return &types.ResponsesResponse{
		ID:       resp.ID,
		Object:   "response",
		Status:   "completed",
		Model:    resp.Model,
		Provider: resp.Provider,
		Output: []types.ResponsesOutputItem{{
			Type:   "message",
			Status: "completed",
			Role:   "assistant",
			Content: []types.ResponsesContentPart{{
				Type: "output_text",
				Text: text,
			}},
		}},
		Usage: types.ResponsesUsage{
			InputTokens:  resp.Usage.PromptTokens,
			OutputTokens: resp.Usage.CompletionTokens,
			TotalTokens:  resp.Usage.TotalTokens,
		},
	}
}
