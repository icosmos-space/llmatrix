package openaicompatible

import (
	"context"

	"github.com/icosmos-space/llmatrix/engine/errors"
	"github.com/icosmos-space/llmatrix/engine/types"
)

// Responses implements provider.ResponsesProvider via POST /responses.
func (p *Provider) Responses(ctx context.Context, req *types.ResponsesRequest) (*types.ResponsesResponse, error) {
	if req == nil {
		return nil, errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "request is nil", nil)
	}
	if req.Model == "" {
		return nil, errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "model is required", nil)
	}
	if req.Input == nil {
		return nil, errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "input is required", nil)
	}
	if req.Stream {
		return nil, errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "use ResponsesStream for stream=true", nil)
	}

	var raw types.ResponsesResponse
	if err := p.doJSON(ctx, "POST", "/responses", responsesBody(req, false), &raw); err != nil {
		return nil, err
	}
	if raw.Error != nil && raw.Error.Message != "" {
		return nil, errors.Wrap(errors.KindProvider, p.Info().Name, raw.Error.Message, nil)
	}
	raw.Provider = p.Info().Name
	if raw.Object == "" {
		raw.Object = "response"
	}
	if raw.Status == "" {
		raw.Status = "completed"
	}
	if raw.Model == "" {
		raw.Model = req.Model
	}
	return &raw, nil
}

func responsesBody(req *types.ResponsesRequest, stream bool) map[string]any {
	body := map[string]any{
		"model":  req.Model,
		"input":  req.Input,
		"stream": stream,
	}
	if req.Instructions != "" {
		body["instructions"] = req.Instructions
	}
	if req.MaxOutputTokens != nil {
		body["max_output_tokens"] = *req.MaxOutputTokens
	}
	if req.Temperature != nil {
		body["temperature"] = *req.Temperature
	}
	if req.TopP != nil {
		body["top_p"] = *req.TopP
	}
	if len(req.Tools) > 0 {
		body["tools"] = req.Tools
	}
	if req.ToolChoice != nil {
		body["tool_choice"] = req.ToolChoice
	}
	if req.User != "" {
		body["user"] = req.User
	}
	if req.Metadata != nil {
		body["metadata"] = req.Metadata
	}
	if req.Store != nil {
		body["store"] = *req.Store
	}
	if req.PreviousResponseID != "" {
		body["previous_response_id"] = req.PreviousResponseID
	}
	return body
}
