package embedding

// Request is an embedding request.
type Request struct {
	Model string   `json:"model,omitempty"`
	Input []string `json:"input"`
}

// Item is a single embedding vector.
type Item struct {
	Index     int       `json:"index"`
	Embedding []float64 `json:"embedding"`
}

// Usage reports embedding token consumption.
type Usage struct {
	PromptTokens int `json:"prompt_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

// Response is an embedding response.
type Response struct {
	Model    string `json:"model,omitempty"`
	Provider string `json:"provider,omitempty"`
	Data     []Item `json:"data"`
	Usage    Usage  `json:"usage"`
}
