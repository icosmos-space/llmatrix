package stream

import "context"

// Chunk is a streaming content delta.
type Chunk struct {
	Delta        string
	FinishReason string
	Done         bool
}

// Stream yields incremental chat chunks.
type Stream interface {
	Next(ctx context.Context) bool
	Chunk() Chunk
	Err() error
	Close() error
}
