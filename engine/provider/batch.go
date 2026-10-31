package provider

import (
	"context"

	"github.com/icosmos-space/llmatrix/engine/types"
)

// BatchesProvider manages OpenAI Batches API (LiteLLM create_batch / …).
type BatchesProvider interface {
	Provider
	CreateBatch(ctx context.Context, req *types.BatchCreateRequest) (*types.Batch, error)
	RetrieveBatch(ctx context.Context, batchID string) (*types.Batch, error)
	CancelBatch(ctx context.Context, batchID string) (*types.Batch, error)
	ListBatches(ctx context.Context, limit int) ([]types.Batch, error)
}
