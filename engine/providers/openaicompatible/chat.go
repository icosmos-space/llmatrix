package openaicompatible

import (
	"context"
	"fmt"

	"github.com/icosmos-space/llmatrix/engine/errors"
	"github.com/icosmos-space/llmatrix/engine/types"
)

// ChatCompletion implements provider.ChatProvider.
func (p *Provider) ChatCompletion(ctx context.Context, req *types.ChatRequest) (*types.ChatResponse, error) {
	if req == nil {
		return nil, errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "request is nil", nil)
	}
	if req.Model == "" {
		return nil, errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "model is required", nil)
	}
	if len(req.Messages) == 0 {
		return nil, errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "messages are required", nil)
	}
	if req.Stream {
		return nil, errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "use ChatCompletionStream for stream=true", nil)
	}

	body := toWireChatRequest(req, false)

	var raw wireChatResponse
	if err := p.doJSON(ctx, "POST", "/chat/completions", body, &raw); err != nil {
		return nil, err
	}
	if raw.Error != nil && raw.Error.Message != "" {
		return nil, errors.Wrap(errors.KindProvider, p.Info().Name, raw.Error.Message, nil)
	}
	if len(raw.Choices) == 0 {
		return nil, errors.Wrap(errors.KindProvider, p.Info().Name, "empty choices in response", nil)
	}

	out := &types.ChatResponse{
		ID:       raw.ID,
		Object:   raw.Object,
		Created:  raw.Created,
		Model:    raw.Model,
		Provider: p.Info().Name,
		Usage:    raw.Usage,
		Choices:  make([]types.ChatChoice, 0, len(raw.Choices)),
	}
	for _, c := range raw.Choices {
		out.Choices = append(out.Choices, types.ChatChoice{
			Index:        c.Index,
			Message:      fromWireMessage(c.Message),
			FinishReason: c.FinishReason,
		})
	}
	if out.Model == "" {
		out.Model = req.Model
	}
	return out, nil
}

// Ensure compile-time interface satisfaction.
var _ interface {
	ChatCompletion(context.Context, *types.ChatRequest) (*types.ChatResponse, error)
} = (*Provider)(nil)

// String helps debugging.
func (p *Provider) String() string {
	return fmt.Sprintf("openaicompatible.Provider(%s)", p.Info().Name)
}
