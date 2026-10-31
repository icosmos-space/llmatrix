package openaicompatible

import (
	"context"
	"io"

	"github.com/icosmos-space/llmatrix/engine/errors"
	"github.com/icosmos-space/llmatrix/engine/types"
)

// Speech implements provider.SpeechProvider.
func (p *Provider) Speech(ctx context.Context, req *types.SpeechRequest) (*types.SpeechResponse, error) {
	if req == nil {
		return nil, errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "request is nil", nil)
	}
	if req.Model == "" || req.Input == "" || req.Voice == "" {
		return nil, errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "model, input, and voice are required", nil)
	}
	body := map[string]any{
		"model": req.Model,
		"input": req.Input,
		"voice": req.Voice,
	}
	if req.ResponseFormat != "" {
		body["response_format"] = req.ResponseFormat
	}
	if req.Speed != nil {
		body["speed"] = *req.Speed
	}

	httpReq, err := p.newRequest(ctx, "POST", "/audio/speech", body)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Accept", "application/octet-stream")

	resp, err := p.httpClient().Do(httpReq)
	if err != nil {
		return nil, errors.Wrap(errors.KindProvider, p.Info().Name, "http request failed", err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(errors.KindProvider, p.Info().Name, "read response", err)
	}
	if resp.StatusCode >= 400 {
		return nil, mapHTTPError(p.Info().Name, resp.StatusCode, data)
	}
	ct := resp.Header.Get("Content-Type")
	if ct == "" {
		ct = "audio/mpeg"
	}
	return &types.SpeechResponse{
		ContentType: ct,
		Audio:       data,
		Provider:    p.Info().Name,
		Model:       req.Model,
	}, nil
}
