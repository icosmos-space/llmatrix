package types

// BatchCreateRequest creates a batch job (LiteLLM create_batch()).
type BatchCreateRequest struct {
	CompletionWindow string            `json:"completion_window"` // typically "24h"
	Endpoint         string            `json:"endpoint"`          // e.g. /v1/chat/completions
	InputFileID      string            `json:"input_file_id"`
	Metadata         map[string]string `json:"metadata,omitempty"`
}

// BatchRequestCounts summarizes batch progress.
type BatchRequestCounts struct {
	Total     int `json:"total"`
	Completed int `json:"completed"`
	Failed    int `json:"failed"`
}

// Batch is an OpenAI batch object.
type Batch struct {
	ID               string             `json:"id"`
	Object           string             `json:"object,omitempty"`
	Endpoint         string             `json:"endpoint,omitempty"`
	Errors           any                `json:"errors,omitempty"`
	InputFileID      string             `json:"input_file_id,omitempty"`
	CompletionWindow string             `json:"completion_window,omitempty"`
	Status           string             `json:"status,omitempty"`
	OutputFileID     string             `json:"output_file_id,omitempty"`
	ErrorFileID      string             `json:"error_file_id,omitempty"`
	CreatedAt        int64              `json:"created_at,omitempty"`
	InProgressAt     int64              `json:"in_progress_at,omitempty"`
	ExpiresAt        int64              `json:"expires_at,omitempty"`
	CompletedAt      int64              `json:"completed_at,omitempty"`
	FailedAt         int64              `json:"failed_at,omitempty"`
	CancelledAt      int64              `json:"cancelled_at,omitempty"`
	RequestCounts    BatchRequestCounts `json:"request_counts"`
	Metadata         map[string]string  `json:"metadata,omitempty"`
	Provider         string             `json:"provider,omitempty"`
}
