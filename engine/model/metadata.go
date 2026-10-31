package model

// Metadata describes model capabilities and limits.
type Metadata struct {
	DisplayName       string
	ContextWindow     int
	MaxOutputTokens   int
	SupportsTools     bool
	SupportsVision    bool
	SupportsStreaming bool
}
