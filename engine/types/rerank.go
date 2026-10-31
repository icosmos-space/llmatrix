package types

// RerankRequest ranks documents by relevance to a query (Cohere / LiteLLM shaped).
type RerankRequest struct {
	Model           string   `json:"model"`
	Query           string   `json:"query"`
	Documents       []string `json:"documents"`
	TopN            *int     `json:"top_n,omitempty"`
	ReturnDocuments *bool    `json:"return_documents,omitempty"`
}

// RerankResult is one ranked document.
type RerankResult struct {
	Index          int     `json:"index"`
	RelevanceScore float64 `json:"relevance_score"`
	Document       string  `json:"document,omitempty"`
}

// RerankResponse is a rerank response.
type RerankResponse struct {
	ID       string         `json:"id,omitempty"`
	Model    string         `json:"model,omitempty"`
	Provider string         `json:"provider,omitempty"`
	Results  []RerankResult `json:"results"`
	Usage    Usage          `json:"usage,omitempty"`
}
