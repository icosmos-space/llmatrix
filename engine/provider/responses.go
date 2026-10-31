package provider

import (
	"context"

	"github.com/icosmos-space/llmatrix/engine/stream"
	"github.com/icosmos-space/llmatrix/engine/types"
)

// ResponsesProvider implements the OpenAI Responses API create call.
type ResponsesProvider interface {
	Provider
	Responses(ctx context.Context, req *types.ResponsesRequest) (*types.ResponsesResponse, error)
}

// ResponsesStreamProvider streams Responses API SSE events as text deltas.
type ResponsesStreamProvider interface {
	Provider
	ResponsesStream(ctx context.Context, req *types.ResponsesRequest) (stream.Stream, error)
}

// ResponsesStoreProvider manages stored Responses objects (get/delete/cancel).
type ResponsesStoreProvider interface {
	Provider
	GetResponse(ctx context.Context, responseID string) (*types.ResponsesResponse, error)
	DeleteResponse(ctx context.Context, responseID string) error
	CancelResponse(ctx context.Context, responseID string) (*types.ResponsesResponse, error)
}
