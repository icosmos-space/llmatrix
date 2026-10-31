package openaicompatible

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/icosmos-space/llmatrix/engine/errors"
	"github.com/icosmos-space/llmatrix/engine/types"
)

// Transcription implements provider.TranscriptionProvider.
func (p *Provider) Transcription(ctx context.Context, req *types.TranscriptionRequest) (*types.TranscriptionResponse, error) {
	if req == nil {
		return nil, errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "request is nil", nil)
	}
	if req.Model == "" || len(req.Audio) == 0 {
		return nil, errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "model and audio are required", nil)
	}
	fileName := req.FileName
	if fileName == "" {
		fileName = "audio.wav"
	}

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	_ = w.WriteField("model", req.Model)
	if req.Language != "" {
		_ = w.WriteField("language", req.Language)
	}
	if req.Prompt != "" {
		_ = w.WriteField("prompt", req.Prompt)
	}
	if req.ResponseFormat != "" {
		_ = w.WriteField("response_format", req.ResponseFormat)
	}
	if req.Temperature != nil {
		_ = w.WriteField("temperature", fmt.Sprintf("%g", *req.Temperature))
	}
	part, err := w.CreateFormFile("file", fileName)
	if err != nil {
		return nil, errors.Wrap(errors.KindInternal, p.Info().Name, "create form file", err)
	}
	if _, err := part.Write(req.Audio); err != nil {
		return nil, errors.Wrap(errors.KindInternal, p.Info().Name, "write audio", err)
	}
	if err := w.Close(); err != nil {
		return nil, errors.Wrap(errors.KindInternal, p.Info().Name, "close multipart", err)
	}

	url, err := p.endpoint("/audio/transcriptions")
	if err != nil {
		return nil, err
	}
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, &buf)
	if err != nil {
		return nil, errors.Wrap(errors.KindInternal, p.Info().Name, "create request", err)
	}
	httpReq.Header.Set("Content-Type", w.FormDataContentType())
	if key := p.resolvedAPIKey(); key != "" {
		httpReq.Header.Set("Authorization", "Bearer "+key)
	}

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

	// text format returns plain text
	if req.ResponseFormat == "text" {
		return &types.TranscriptionResponse{
			Text:     string(data),
			Provider: p.Info().Name,
			Model:    req.Model,
		}, nil
	}

	var raw struct {
		Text     string     `json:"text"`
		Language string     `json:"language"`
		Duration float64    `json:"duration"`
		Error    *wireError `json:"error"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		// fall back to raw body as text
		return &types.TranscriptionResponse{
			Text:     string(data),
			Provider: p.Info().Name,
			Model:    req.Model,
		}, nil
	}
	if raw.Error != nil && raw.Error.Message != "" {
		return nil, errors.Wrap(errors.KindProvider, p.Info().Name, raw.Error.Message, nil)
	}
	return &types.TranscriptionResponse{
		Text:     raw.Text,
		Language: raw.Language,
		Duration: raw.Duration,
		Provider: p.Info().Name,
		Model:    req.Model,
	}, nil
}
