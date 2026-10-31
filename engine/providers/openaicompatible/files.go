package openaicompatible

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/icosmos-space/llmatrix/engine/errors"
	"github.com/icosmos-space/llmatrix/engine/types"
)

// CreateFile implements provider.FilesProvider via POST /files.
func (p *Provider) CreateFile(ctx context.Context, req *types.FileCreateRequest) (*types.FileObject, error) {
	if req == nil || len(req.File) == 0 || req.Purpose == "" {
		return nil, errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "file bytes and purpose are required", nil)
	}
	name := req.FileName
	if name == "" {
		name = "upload.bin"
	}
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	_ = w.WriteField("purpose", req.Purpose)
	part, err := w.CreateFormFile("file", name)
	if err != nil {
		return nil, errors.Wrap(errors.KindInternal, p.Info().Name, "create form file", err)
	}
	if _, err := part.Write(req.File); err != nil {
		return nil, errors.Wrap(errors.KindInternal, p.Info().Name, "write file", err)
	}
	if err := w.Close(); err != nil {
		return nil, errors.Wrap(errors.KindInternal, p.Info().Name, "close multipart", err)
	}

	endpoint, err := p.endpoint("/files")
	if err != nil {
		return nil, err
	}
	httpReq, err := http.NewRequestWithContext(ctx, "POST", endpoint, &buf)
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
	var out types.FileObject
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, errors.Wrap(errors.KindProvider, p.Info().Name, "decode response", err)
	}
	out.Provider = p.Info().Name
	return &out, nil
}

// GetFile implements provider.FilesProvider via GET /files/{id}.
func (p *Provider) GetFile(ctx context.Context, fileID string) (*types.FileObject, error) {
	if fileID == "" {
		return nil, errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "file id is required", nil)
	}
	var out types.FileObject
	if err := p.doJSON(ctx, "GET", "/files/"+url.PathEscape(fileID), nil, &out); err != nil {
		return nil, err
	}
	out.Provider = p.Info().Name
	return &out, nil
}

// DeleteFile implements provider.FilesProvider via DELETE /files/{id}.
func (p *Provider) DeleteFile(ctx context.Context, fileID string) error {
	if fileID == "" {
		return errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "file id is required", nil)
	}
	return p.doJSON(ctx, "DELETE", "/files/"+url.PathEscape(fileID), nil, nil)
}

// ListFiles implements provider.FilesProvider via GET /files.
func (p *Provider) ListFiles(ctx context.Context, purpose string) (*types.FileListResponse, error) {
	path := "/files"
	if purpose != "" {
		path += "?purpose=" + url.QueryEscape(purpose)
	}
	var out types.FileListResponse
	if err := p.doJSON(ctx, "GET", path, nil, &out); err != nil {
		return nil, err
	}
	out.Provider = p.Info().Name
	return &out, nil
}

// GetFileContent implements provider.FilesProvider via GET /files/{id}/content.
func (p *Provider) GetFileContent(ctx context.Context, fileID string) ([]byte, error) {
	if fileID == "" {
		return nil, errors.Wrap(errors.KindInvalidRequest, p.Info().Name, "file id is required", nil)
	}
	endpoint, err := p.endpoint("/files/" + url.PathEscape(fileID) + "/content")
	if err != nil {
		return nil, err
	}
	httpReq, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, errors.Wrap(errors.KindInternal, p.Info().Name, "create request", err)
	}
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
	return data, nil
}
