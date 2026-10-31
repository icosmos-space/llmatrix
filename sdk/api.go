package sdk

import (
	"context"
	"sync"

	"github.com/icosmos-space/llmatrix/sdk/audio"
	"github.com/icosmos-space/llmatrix/sdk/chat"
	"github.com/icosmos-space/llmatrix/sdk/embedding"
	"github.com/icosmos-space/llmatrix/sdk/image"
	"github.com/icosmos-space/llmatrix/sdk/rerank"
	"github.com/icosmos-space/llmatrix/sdk/streaming"
)

// Package-level helpers mirror LiteLLM's module functions:
// completion / embedding / image_generation / speech / transcription / rerank.

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
