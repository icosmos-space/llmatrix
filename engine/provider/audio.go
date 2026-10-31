package provider

import (
	"context"

	"github.com/icosmos-space/llmatrix/engine/types"
)

// SpeechProvider synthesizes speech audio from text.
type SpeechProvider interface {
	Provider
	Speech(ctx context.Context, req *types.SpeechRequest) (*types.SpeechResponse, error)
}

// TranscriptionProvider transcribes audio to text.
type TranscriptionProvider interface {
	Provider
	Transcription(ctx context.Context, req *types.TranscriptionRequest) (*types.TranscriptionResponse, error)
}
