package types

// EmbeddingRequest is the engine-level embedding request.
type EmbeddingRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

// EmbeddingItem is a single vector result.
type EmbeddingItem struct {
	Index     int       `json:"index"`
	Embedding []float64 `json:"embedding"`
}

// EmbeddingResponse is the engine-level embedding response.
type EmbeddingResponse struct {
	Model    string          `json:"model,omitempty"`
	Provider string          `json:"provider,omitempty"`
	Data     []EmbeddingItem `json:"data"`
	Usage    Usage           `json:"usage"`
}
