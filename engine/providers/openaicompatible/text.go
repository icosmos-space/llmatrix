package openaicompatible

import (
	"context"

	"github.com/icosmos-space/llmatrix/engine/errors"
	"github.com/icosmos-space/llmatrix/engine/types"
)

// TextCompletion implements provider.TextCompletionProvider via POST /completions.
func (p *Provider) TextCompletion(ctx context.Context, req *types.TextCompletionRequest) (*types.TextCompletionResponse, error) {
	if req == nil {
		return nil, errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "request is nil", nil)
	}
	if req.Model == "" {
		return nil, errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "model is required", nil)
	}
	if req.Prompt == nil {
		return nil, errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "prompt is required", nil)
	}
	if req.Stream {
		return nil, errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "streaming text completion not yet supported", nil)
	}

	body := map[string]any{
		"model":  req.Model,
		"prompt": req.Prompt,
	}
	if req.MaxTokens != nil {
		body["max_tokens"] = *req.MaxTokens
	}
	if req.Temperature != nil {
		body["temperature"] = *req.Temperature
	}
	if req.TopP != nil {
		body["top_p"] = *req.TopP
	}
	if req.N != nil {
		body["n"] = *req.N
	}
	if len(req.Stop) > 0 {
		body["stop"] = req.Stop
	}
	if req.PresencePenalty != nil {
		body["presence_penalty"] = *req.PresencePenalty
	}
	if req.FrequencyPenalty != nil {
		body["frequency_penalty"] = *req.FrequencyPenalty
	}
	if req.LogitBias != nil {
		body["logit_bias"] = req.LogitBias
	}
	if req.User != "" {
		body["user"] = req.User
	}
	if req.Seed != nil {
		body["seed"] = *req.Seed
	}
	if req.Suffix != "" {
		body["suffix"] = req.Suffix
	}
	if req.Echo != nil {
		body["echo"] = *req.Echo
	}
	if req.BestOf != nil {
		body["best_of"] = *req.BestOf
	}
	if req.Logprobs != nil {
		body["logprobs"] = *req.Logprobs
	}

	var raw types.TextCompletionResponse
	if err := p.doJSON(ctx, "POST", "/completions", body, &raw); err != nil {
		return nil, err
	}
	raw.Provider = p.Info().Name
	if raw.Model == "" {
		raw.Model = req.Model
	}
	if raw.Object == "" {
		raw.Object = "text_completion"
	}
	return &raw, nil
}
