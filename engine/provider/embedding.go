package provider

import (
	"context"

	"github.com/icosmos-space/llmatrix/engine/types"
)

// EmbeddingProvider creates embeddings.
type EmbeddingProvider interface {
	Provider
	Embeddings(ctx context.Context, req *types.EmbeddingRequest) (*types.EmbeddingResponse, error)
}
