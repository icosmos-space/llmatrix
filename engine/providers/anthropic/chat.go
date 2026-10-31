package anthropic

import (
	"context"

	"github.com/icosmos-space/llmatrix/engine/errors"
	"github.com/icosmos-space/llmatrix/engine/types"
)

// ChatCompletion implements provider.ChatProvider.
func (p *Provider) ChatCompletion(ctx context.Context, req *types.ChatRequest) (*types.ChatResponse, error) {
	body, err := buildRequest(req, false)
	if err != nil {
		return nil, err
	}
	var raw wireResponse
	if err := p.doJSON(ctx, "POST", "/v1/messages", body, &raw); err != nil {
		return nil, err
	}
	if raw.Error != nil && raw.Error.Message != "" {
		return nil, errors.Wrap(errors.KindProvider, "anthropic", raw.Error.Message, nil)
	}
	msg := types.Message{
		Role:      types.RoleAssistant,
		Content:   contentText(raw.Content),
		ToolCalls: contentToolCalls(raw.Content),
	}
	return &types.ChatResponse{
		ID:       raw.ID,
		Model:    firstNonEmpty(raw.Model, req.Model),
		Provider: "anthropic",
		Choices: []types.ChatChoice{{
			Index:        0,
			Message:      msg,
			FinishReason: mapStopReason(raw.StopReason),
		}},
		Usage: types.Usage{
			PromptTokens:     raw.Usage.InputTokens,
			CompletionTokens: raw.Usage.OutputTokens,
			TotalTokens:      raw.Usage.InputTokens + raw.Usage.OutputTokens,
		},
	}, nil
}

func buildRequest(req *types.ChatRequest, stream bool) (*wireRequest, error) {
	if req == nil {
		return nil, errors.Wrap(errors.KindInvalidRequest, "anthropic", "request is nil", nil)
	}
	if req.Model == "" {
		return nil, errors.Wrap(errors.KindInvalidRequest, "anthropic", "model is required", nil)
	}
	if len(req.Messages) == 0 {
		return nil, errors.Wrap(errors.KindInvalidRequest, "anthropic", "messages are required", nil)
	}
	system, rest := splitSystem(req.Messages)
	if len(rest) == 0 {
		return nil, errors.Wrap(errors.KindInvalidRequest, "anthropic", "at least one non-system message is required", nil)
	}
	maxTok := defaultMaxTok
	if req.MaxTokens != nil && *req.MaxTokens > 0 {
		maxTok = *req.MaxTokens
	}
	return &wireRequest{
		Model:       req.Model,
		Messages:    toWireMessages(rest),
		MaxTokens:   maxTok,
		System:      system,
		Temperature: req.Temperature,
		Stream:      stream,
		Tools:       toWireTools(req.Tools),
		ToolChoice:  toWireToolChoice(req.ToolChoice),
		Thinking:    req.Thinking,
	}, nil
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}
