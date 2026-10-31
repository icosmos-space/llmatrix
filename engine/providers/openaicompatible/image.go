package openaicompatible

import (
	"context"

	"github.com/icosmos-space/llmatrix/engine/errors"
	"github.com/icosmos-space/llmatrix/engine/types"
)

// ImageGeneration implements provider.ImageProvider.
func (p *Provider) ImageGeneration(ctx context.Context, req *types.ImageRequest) (*types.ImageResponse, error) {
	if req == nil {
		return nil, errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "request is nil", nil)
	}
	if req.Prompt == "" {
		return nil, errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "prompt is required", nil)
	}
	body := map[string]any{
		"prompt": req.Prompt,
	}
	if req.Model != "" {
		body["model"] = req.Model
	}
	if req.N != nil {
		body["n"] = *req.N
	}
	if req.Size != "" {
		body["size"] = req.Size
	}
	if req.Quality != "" {
		body["quality"] = req.Quality
	}
	if req.Style != "" {
		body["style"] = req.Style
	}
	if req.ResponseFormat != "" {
		body["response_format"] = req.ResponseFormat
	}
	if req.User != "" {
		body["user"] = req.User
	}

	var raw struct {
		Created int64 `json:"created"`
		Data    []struct {
			URL           string `json:"url"`
			B64JSON       string `json:"b64_json"`
			RevisedPrompt string `json:"revised_prompt"`
		} `json:"data"`
		Error *wireError `json:"error"`
	}
	if err := p.doJSON(ctx, "POST", "/images/generations", body, &raw); err != nil {
		return nil, err
	}
	if raw.Error != nil && raw.Error.Message != "" {
		return nil, errors.Wrap(errors.KindProvider, p.Info().Name, raw.Error.Message, nil)
	}
	out := &types.ImageResponse{
		Created:  raw.Created,
		Provider: p.Info().Name,
		Model:    req.Model,
		Data:     make([]types.ImageData, 0, len(raw.Data)),
	}
	for _, d := range raw.Data {
		out.Data = append(out.Data, types.ImageData{
			URL:     d.URL,
			B64JSON: d.B64JSON,
			Revised: d.RevisedPrompt,
		})
	}
	return out, nil
}
