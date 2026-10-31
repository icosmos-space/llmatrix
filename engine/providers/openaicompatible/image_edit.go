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

// ImageEdit implements provider.ImageEditProvider via POST /images/edits.
func (p *Provider) ImageEdit(ctx context.Context, req *types.ImageEditRequest) (*types.ImageResponse, error) {
	if req == nil {
		return nil, errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "request is nil", nil)
	}
	if req.Prompt == "" {
		return nil, errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "prompt is required", nil)
	}
	if len(req.Image) == 0 {
		return nil, errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "image is required", nil)
	}

	fileName := req.ImageFileName
	if fileName == "" {
		fileName = "image.png"
	}

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	_ = w.WriteField("prompt", req.Prompt)
	if req.Model != "" {
		_ = w.WriteField("model", req.Model)
	}
	if req.N != nil {
		_ = w.WriteField("n", fmt.Sprintf("%d", *req.N))
	}
	if req.Size != "" {
		_ = w.WriteField("size", req.Size)
	}
	if req.ResponseFormat != "" {
		_ = w.WriteField("response_format", req.ResponseFormat)
	}
	if req.User != "" {
		_ = w.WriteField("user", req.User)
	}
	part, err := w.CreateFormFile("image", fileName)
	if err != nil {
		return nil, errors.Wrap(errors.KindInternal, p.Info().Name, "create image form", err)
	}
	if _, err := part.Write(req.Image); err != nil {
		return nil, errors.Wrap(errors.KindInternal, p.Info().Name, "write image", err)
	}
	if len(req.Mask) > 0 {
		maskName := req.MaskFileName
		if maskName == "" {
			maskName = "mask.png"
		}
		mp, err := w.CreateFormFile("mask", maskName)
		if err != nil {
			return nil, errors.Wrap(errors.KindInternal, p.Info().Name, "create mask form", err)
		}
		if _, err := mp.Write(req.Mask); err != nil {
			return nil, errors.Wrap(errors.KindInternal, p.Info().Name, "write mask", err)
		}
	}
	if err := w.Close(); err != nil {
		return nil, errors.Wrap(errors.KindInternal, p.Info().Name, "close multipart", err)
	}

	url, err := p.endpoint("/images/edits")
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

	var raw struct {
		Created int64 `json:"created"`
		Data    []struct {
			URL           string `json:"url"`
			B64JSON       string `json:"b64_json"`
			RevisedPrompt string `json:"revised_prompt"`
		} `json:"data"`
		Error *wireError `json:"error"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, errors.Wrap(errors.KindProvider, p.Info().Name, "decode response", err)
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
			URL: d.URL, B64JSON: d.B64JSON, Revised: d.RevisedPrompt,
		})
	}
	return out, nil
}
