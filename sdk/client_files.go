package sdk

import (
	"context"
	"fmt"

	"github.com/icosmos-space/llmatrix/engine/provider"
	"github.com/icosmos-space/llmatrix/sdk/files"
)

func (c *Client) filesProvider(modelHint string) (provider.FilesProvider, string, error) {
	raw := modelHint
	if raw == "" {
		raw = c.cfg.Model
	}
	if raw == "" {
		prov := c.cfg.Provider
		if prov == "" {
			prov = "openai"
		}
		raw = prov + "/default"
	}
	ref := c.candidates(raw)
	if len(ref) == 0 {
		return nil, "", fmt.Errorf("%w: model is required", ErrInvalidConfig)
	}
	p, err := c.newProvider(ref[0], true)
	if err != nil {
		return nil, "", err
	}
	fp, ok := p.(provider.FilesProvider)
	if !ok {
		return nil, "", fmt.Errorf("%w: provider %q does not support files", ErrUnsupported, ref[0].Provider)
	}
	return fp, ref[0].Provider, nil
}

// CreateFile uploads a file (LiteLLM create_file()).
func (c *Client) CreateFile(ctx context.Context, req *files.CreateRequest) (*files.Object, error) {
	if req == nil || len(req.File) == 0 || req.Purpose == "" {
		return nil, ErrInvalidRequest
	}
	fp, _, err := c.filesProvider(req.Model)
	if err != nil {
		return nil, err
	}
	out, err := fp.CreateFile(ctx, toEngineFileCreate(req))
	if err != nil {
		return nil, MapError(err)
	}
	return fromEngineFile(out), nil
}

// GetFile fetches file metadata.
func (c *Client) GetFile(ctx context.Context, fileID string, modelHint ...string) (*files.Object, error) {
	if fileID == "" {
		return nil, ErrInvalidRequest
	}
	hint := ""
	if len(modelHint) > 0 {
		hint = modelHint[0]
	}
	fp, _, err := c.filesProvider(hint)
	if err != nil {
		return nil, err
	}
	out, err := fp.GetFile(ctx, fileID)
	if err != nil {
		return nil, MapError(err)
	}
	return fromEngineFile(out), nil
}

// DeleteFile deletes a file.
func (c *Client) DeleteFile(ctx context.Context, fileID string, modelHint ...string) error {
	if fileID == "" {
		return ErrInvalidRequest
	}
	hint := ""
	if len(modelHint) > 0 {
		hint = modelHint[0]
	}
	fp, _, err := c.filesProvider(hint)
	if err != nil {
		return err
	}
	if err := fp.DeleteFile(ctx, fileID); err != nil {
		return MapError(err)
	}
	return nil
}

// ListFiles lists uploaded files.
func (c *Client) ListFiles(ctx context.Context, purpose string, modelHint ...string) (*files.ListResponse, error) {
	hint := ""
	if len(modelHint) > 0 {
		hint = modelHint[0]
	}
	fp, _, err := c.filesProvider(hint)
	if err != nil {
		return nil, err
	}
	out, err := fp.ListFiles(ctx, purpose)
	if err != nil {
		return nil, MapError(err)
	}
	return fromEngineFileList(out), nil
}

// GetFileContent downloads file bytes.
func (c *Client) GetFileContent(ctx context.Context, fileID string, modelHint ...string) ([]byte, error) {
	if fileID == "" {
		return nil, ErrInvalidRequest
	}
	hint := ""
	if len(modelHint) > 0 {
		hint = modelHint[0]
	}
	fp, _, err := c.filesProvider(hint)
	if err != nil {
		return nil, err
	}
	data, err := fp.GetFileContent(ctx, fileID)
	if err != nil {
		return nil, MapError(err)
	}
	return data, nil
}
