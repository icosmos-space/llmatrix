package streaming

import "context"

// Stream yields incremental events until closed or cancelled.
type Stream interface {
	// Next advances to the next event. Returns false when the stream ends.
	Next(ctx context.Context) bool
	// Event returns the current event after a successful Next call.
	Event() Event
	// Err returns the first error encountered, if any.
	Err() error
	// Close releases stream resources.
	Close() error
}
