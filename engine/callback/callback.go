package callback

import (
	"context"
	"time"

	"github.com/icosmos-space/llmatrix/engine/types"
)

// Event is emitted around provider calls (LiteLLM success_callback / failure_callback).
type Event struct {
	Method    string // chat, embedding, image, speech, transcription, rerank
	Provider  string
	Model     string
	Start     time.Time
	End       time.Time
	Err       error
	Usage     types.Usage
	RequestID string
}

// Handler receives completion events.
type Handler func(ctx context.Context, ev Event)

// Bus fans out events to registered handlers.
type Bus struct {
	success []Handler
	failure []Handler
}

// OnSuccess registers a success handler.
func (b *Bus) OnSuccess(h Handler) { b.success = append(b.success, h) }

// OnFailure registers a failure handler.
func (b *Bus) OnFailure(h Handler) { b.failure = append(b.failure, h) }

// EmitSuccess notifies success handlers.
func (b *Bus) EmitSuccess(ctx context.Context, ev Event) {
	for _, h := range b.success {
		h(ctx, ev)
	}
}

// EmitFailure notifies failure handlers.
func (b *Bus) EmitFailure(ctx context.Context, ev Event) {
	for _, h := range b.failure {
		h(ctx, ev)
	}
}
