package sdk

import (
	"context"
	"sync"

	"github.com/icosmos-space/llmatrix/sdk/audio"
	"github.com/icosmos-space/llmatrix/sdk/batches"
	"github.com/icosmos-space/llmatrix/sdk/chat"
	"github.com/icosmos-space/llmatrix/sdk/embedding"
	"github.com/icosmos-space/llmatrix/sdk/files"
	"github.com/icosmos-space/llmatrix/sdk/image"
	"github.com/icosmos-space/llmatrix/sdk/moderation"
	"github.com/icosmos-space/llmatrix/sdk/rerank"
	"github.com/icosmos-space/llmatrix/sdk/responses"
	"github.com/icosmos-space/llmatrix/sdk/streaming"
	"github.com/icosmos-space/llmatrix/sdk/text"
)

// Package-level helpers mirror LiteLLM's module functions:
// completion / text_completion / embedding / image_generation / image_edit / speech /
// transcription / rerank / responses / moderation / create_file / create_batch /
// token_counter / completion_cost.

var (
	defaultMu     sync.RWMutex
	defaultClient = NewClient()
	defaultOpts   []Option
)

// Configure sets options used by package-level Completion / Embedding / … helpers.
func Configure(opts ...Option) {
	defaultMu.Lock()
	defer defaultMu.Unlock()
	defaultOpts = append([]Option{}, opts...)
	defaultClient = NewClient(defaultOpts...)
}

func packageClient(extra ...Option) *Client {
	defaultMu.RLock()
	base := append([]Option{}, defaultOpts...)
	defaultMu.RUnlock()
	if len(extra) == 0 {
		defaultMu.RLock()
		c := defaultClient
		defaultMu.RUnlock()
		return c
	}
	return NewClient(append(base, extra...)...)
}

// Completion is LiteLLM completion().
func Completion(ctx context.Context, req *chat.Request, opts ...Option) (*chat.Response, error) {
	return packageClient(opts...).ChatCompletion(ctx, req)
}

// TextCompletion is LiteLLM text_completion().
func TextCompletion(ctx context.Context, req *text.Request, opts ...Option) (*text.Response, error) {
	return packageClient(opts...).TextCompletion(ctx, req)
}

// CompletionStream is LiteLLM completion(..., stream=True).
func CompletionStream(ctx context.Context, req *chat.Request, opts ...Option) (streaming.Stream, error) {
	return packageClient(opts...).ChatCompletionStream(ctx, req)
}

// Embedding is LiteLLM embedding().
func Embedding(ctx context.Context, req *embedding.Request, opts ...Option) (*embedding.Response, error) {
	return packageClient(opts...).Embeddings(ctx, req)
}

// ImageGeneration is LiteLLM image_generation().
func ImageGeneration(ctx context.Context, req *image.Request, opts ...Option) (*image.Response, error) {
	return packageClient(opts...).ImageGeneration(ctx, req)
}

// ImageEdit is LiteLLM image_edit().
func ImageEdit(ctx context.Context, req *image.EditRequest, opts ...Option) (*image.Response, error) {
	return packageClient(opts...).ImageEdit(ctx, req)
}

// Moderation is LiteLLM moderation().
func Moderation(ctx context.Context, req *moderation.Request, opts ...Option) (*moderation.Response, error) {
	return packageClient(opts...).Moderation(ctx, req)
}

// Speech is LiteLLM speech().
func Speech(ctx context.Context, req *audio.SpeechRequest, opts ...Option) (*audio.SpeechResponse, error) {
	return packageClient(opts...).Speech(ctx, req)
}

// Transcription is LiteLLM transcription().
func Transcription(ctx context.Context, req *audio.TranscriptionRequest, opts ...Option) (*audio.TranscriptionResponse, error) {
	return packageClient(opts...).Transcription(ctx, req)
}

// Rerank is LiteLLM rerank().
func Rerank(ctx context.Context, req *rerank.Request, opts ...Option) (*rerank.Response, error) {
	return packageClient(opts...).Rerank(ctx, req)
}

// Responses is LiteLLM responses().
func Responses(ctx context.Context, req *responses.Request, opts ...Option) (*responses.Response, error) {
	return packageClient(opts...).Responses(ctx, req)
}

// ResponsesStream is LiteLLM responses(..., stream=True).
func ResponsesStream(ctx context.Context, req *responses.Request, opts ...Option) (streaming.Stream, error) {
	return packageClient(opts...).ResponsesStream(ctx, req)
}

// GetResponse fetches a stored Responses object.
func GetResponse(ctx context.Context, responseID string, opts ...Option) (*responses.Response, error) {
	return packageClient(opts...).GetResponse(ctx, responseID)
}

// DeleteResponse deletes a stored Responses object.
func DeleteResponse(ctx context.Context, responseID string, opts ...Option) error {
	return packageClient(opts...).DeleteResponse(ctx, responseID)
}

// CancelResponse cancels an in-progress Responses object.
func CancelResponse(ctx context.Context, responseID string, opts ...Option) (*responses.Response, error) {
	return packageClient(opts...).CancelResponse(ctx, responseID)
}

// CreateFile is LiteLLM create_file().
func CreateFile(ctx context.Context, req *files.CreateRequest, opts ...Option) (*files.Object, error) {
	return packageClient(opts...).CreateFile(ctx, req)
}

// GetFile fetches file metadata.
func GetFile(ctx context.Context, fileID string, opts ...Option) (*files.Object, error) {
	return packageClient(opts...).GetFile(ctx, fileID)
}

// DeleteFile deletes a file.
func DeleteFile(ctx context.Context, fileID string, opts ...Option) error {
	return packageClient(opts...).DeleteFile(ctx, fileID)
}

// ListFiles lists files.
func ListFiles(ctx context.Context, purpose string, opts ...Option) (*files.ListResponse, error) {
	return packageClient(opts...).ListFiles(ctx, purpose)
}

// GetFileContent downloads file content.
func GetFileContent(ctx context.Context, fileID string, opts ...Option) ([]byte, error) {
	return packageClient(opts...).GetFileContent(ctx, fileID)
}

// CreateBatch is LiteLLM create_batch().
func CreateBatch(ctx context.Context, req *batches.CreateRequest, opts ...Option) (*batches.Batch, error) {
	return packageClient(opts...).CreateBatch(ctx, req)
}

// RetrieveBatch is LiteLLM retrieve_batch().
func RetrieveBatch(ctx context.Context, batchID string, opts ...Option) (*batches.Batch, error) {
	return packageClient(opts...).RetrieveBatch(ctx, batchID)
}

// CancelBatch is LiteLLM cancel_batch().
func CancelBatch(ctx context.Context, batchID string, opts ...Option) (*batches.Batch, error) {
	return packageClient(opts...).CancelBatch(ctx, batchID)
}

// ListBatches is LiteLLM list_batches().
func ListBatches(ctx context.Context, limit int, opts ...Option) ([]batches.Batch, error) {
	return packageClient(opts...).ListBatches(ctx, limit)
}
