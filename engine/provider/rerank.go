package provider

import (
	"context"

	"github.com/icosmos-space/llmatrix/engine/types"
)

// RerankProvider ranks documents by relevance to a query.
type RerankProvider interface {
	Provider
	Rerank(ctx context.Context, req *types.RerankRequest) (*types.RerankResponse, error)
}
