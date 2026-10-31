package types

// CacheControl marks a content block for prompt caching (Anthropic-style).
type CacheControl struct {
	Type string `json:"type"` // typically "ephemeral"
}

// ContentPart is a structured message content block.
type ContentPart struct {
	Type         string        `json:"type"` // "text"
	Text         string        `json:"text,omitempty"`
	CacheControl *CacheControl `json:"cache_control,omitempty"`
}
