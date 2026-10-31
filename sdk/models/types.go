package models

// Model is an OpenAI /models list entry.
type Model struct {
	ID       string `json:"id"`
	Object   string `json:"object,omitempty"`
	Created  int64  `json:"created,omitempty"`
	OwnedBy  string `json:"owned_by,omitempty"`
	Provider string `json:"provider,omitempty"`
}

// ListResponse is a GET /models response.
type ListResponse struct {
	Object   string  `json:"object,omitempty"`
	Data     []Model `json:"data"`
	Provider string  `json:"provider,omitempty"`
}
