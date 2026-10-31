package gemini

import (
	"context"

	"github.com/icosmos-space/llmatrix/engine/errors"
	"github.com/icosmos-space/llmatrix/engine/types"
)

// ChatCompletion implements provider.ChatProvider.
func (p *Provider) ChatCompletion(ctx context.Context, req *types.ChatRequest) (*types.ChatResponse, error) {
	body, model, err := buildRequest(req)
	if err != nil {
		return nil, err
	}
	rawURL, err := p.modelURL(model, "generateContent")
	if err != nil {
		return nil, err
	}
	var raw wireResponse
	if err := p.doJSON(ctx, "POST", rawURL, body, &raw); err != nil {
		return nil, err
	}
	if raw.Error != nil && raw.Error.Message != "" {
		return nil, errors.Wrap(errors.KindProvider, "gemini", raw.Error.Message, nil)
	}
	if len(raw.Candidates) == 0 {
		return nil, errors.Wrap(errors.KindProvider, "gemini", "empty candidates in response", nil)
	}
	cand := raw.Candidates[0]
	out := &types.ChatResponse{
		Model:    model,
		Provider: "gemini",
		Choices: []types.ChatChoice{{
			Index: 0,
			Message: types.Message{
				Role:    types.RoleAssistant,
				Content: partsText(cand.Content),
			},
			FinishReason: mapFinishReason(cand.FinishReason),
		}},
	}
	if raw.UsageMetadata != nil {
		out.Usage = types.Usage{
			PromptTokens:     raw.UsageMetadata.PromptTokenCount,
			CompletionTokens: raw.UsageMetadata.CandidatesTokenCount,
			TotalTokens:      raw.UsageMetadata.TotalTokenCount,
		}
	}
	return out, nil
}

func buildRequest(req *types.ChatRequest) (*wireRequest, string, error) {
	if req == nil {
		return nil, "", errors.Wrap(errors.KindInvalidRequest, "gemini", "request is nil", nil)
	}
	if req.Model == "" {
		return nil, "", errors.Wrap(errors.KindInvalidRequest, "gemini", "model is required", nil)
	}
	if len(req.Messages) == 0 {
		return nil, "", errors.Wrap(errors.KindInvalidRequest, "gemini", "messages are required", nil)
	}
	system, rest := splitSystem(req.Messages)
	if len(rest) == 0 {
		return nil, "", errors.Wrap(errors.KindInvalidRequest, "gemini", "at least one non-system message is required", nil)
	}
	body := &wireRequest{
		Contents: toWireContents(rest),
	}
	if system != "" {
		body.SystemInstruction = &wireContent{Parts: []wirePart{{Text: system}}}
	}
	if req.Temperature != nil || req.MaxTokens != nil {
		body.GenerationConfig = &wireGenerationConfig{
			Temperature:     req.Temperature,
			MaxOutputTokens: req.MaxTokens,
		}
	}
	return body, req.Model, nil
}
