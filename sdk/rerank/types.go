package rerank

// Request ranks documents against a query.
type Request struct {
	Model           string   `json:"model,omitempty"`
	Query           string   `json:"query"`
	Documents       []string `json:"documents"`
	TopN            *int     `json:"top_n,omitempty"`
	ReturnDocuments *bool    `json:"return_documents,omitempty"`
}

// Result is one ranked hit.
type Result struct {
	Index          int     `json:"index"`
	RelevanceScore float64 `json:"relevance_score"`
	Document       string  `json:"document,omitempty"`
}

// Usage reports token-like usage when provided by the upstream.
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// Response is a rerank response.
type Response struct {
	ID       string   `json:"id,omitempty"`
	Model    string   `json:"model,omitempty"`
	Provider string   `json:"provider,omitempty"`
	Results  []Result `json:"results"`
	Usage    Usage    `json:"usage,omitempty"`
}
