package provider

import (
	"context"

	"github.com/icosmos-space/llmatrix/engine/types"
)

// FilesProvider manages OpenAI Files API (LiteLLM create_file / file helpers).
type FilesProvider interface {
	Provider
	CreateFile(ctx context.Context, req *types.FileCreateRequest) (*types.FileObject, error)
	GetFile(ctx context.Context, fileID string) (*types.FileObject, error)
	DeleteFile(ctx context.Context, fileID string) error
	ListFiles(ctx context.Context, purpose string) (*types.FileListResponse, error)
	GetFileContent(ctx context.Context, fileID string) ([]byte, error)
}
