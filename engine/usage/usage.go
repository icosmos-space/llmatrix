package usage

import "github.com/icosmos-space/llmatrix/engine/types"

// Recorder records usage for observability / billing hooks.
type Recorder interface {
	Record(provider, model string, u types.Usage)
}

// NopRecorder discards usage events.
type NopRecorder struct{}

// Record implements Recorder.
func (NopRecorder) Record(provider, model string, u types.Usage) {}

// Sum merges two Usage values.
func Sum(a, b types.Usage) types.Usage {
	return types.Usage{
		PromptTokens:     a.PromptTokens + b.PromptTokens,
		CompletionTokens: a.CompletionTokens + b.CompletionTokens,
		TotalTokens:      a.TotalTokens + b.TotalTokens,
	}
}
