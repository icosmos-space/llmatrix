package model

// Model identifies a logical model across providers.
type Model struct {
	ID       string
	Provider string
	Metadata Metadata
}
