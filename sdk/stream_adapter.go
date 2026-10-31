package sdk

import (
	"context"

	enginestream "github.com/icosmos-space/llmatrix/engine/stream"
	"github.com/icosmos-space/llmatrix/sdk/streaming"
)

type engineStreamAdapter struct {
	inner enginestream.Stream
	event streaming.Event
}

func wrapEngineStream(s enginestream.Stream) streaming.Stream {
	return &engineStreamAdapter{inner: s}
}

func (a *engineStreamAdapter) Next(ctx context.Context) bool {
	if !a.inner.Next(ctx) {
		return false
	}
	ch := a.inner.Chunk()
	if ch.Done {
		a.event = streaming.Event{Type: streaming.EventTypeDone}
		return true
	}
	a.event = streaming.Event{
		Type:         streaming.EventTypeDelta,
		Delta:        ch.Delta,
		FinishReason: ch.FinishReason,
	}
	return true
}

func (a *engineStreamAdapter) Event() streaming.Event { return a.event }

func (a *engineStreamAdapter) Err() error {
	if err := a.inner.Err(); err != nil {
		a.event = streaming.Event{Type: streaming.EventTypeError, Error: err.Error()}
		return err
	}
	return nil
}

func (a *engineStreamAdapter) Close() error { return a.inner.Close() }
