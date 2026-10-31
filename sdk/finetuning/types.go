package finetuning

// CreateJobRequest creates a fine-tuning job (LiteLLM create_fine_tuning_job).
type CreateJobRequest struct {
	Model           string         `json:"model"`
	TrainingFile    string         `json:"training_file"`
	ValidationFile  string         `json:"validation_file,omitempty"`
	Hyperparameters map[string]any `json:"hyperparameters,omitempty"`
	Suffix          string         `json:"suffix,omitempty"`
	Integrations    []string       `json:"integrations,omitempty"`
	Seed            *int           `json:"seed,omitempty"`
}

// Job is an OpenAI fine-tuning job object.
type Job struct {
	ID             string         `json:"id"`
	Object         string         `json:"object,omitempty"`
	Model          string         `json:"model,omitempty"`
	CreatedAt      int64          `json:"created_at,omitempty"`
	FinishedAt     int64          `json:"finished_at,omitempty"`
	FineTunedModel string         `json:"fine_tuned_model,omitempty"`
	OrganizationID string         `json:"organization_id,omitempty"`
	Status         string         `json:"status,omitempty"`
	TrainingFile   string         `json:"training_file,omitempty"`
	ValidationFile string         `json:"validation_file,omitempty"`
	Hyperparameters  map[string]any `json:"hyperparameters,omitempty"`
	Error          any            `json:"error,omitempty"`
	Provider       string         `json:"provider,omitempty"`
}

// JobList lists fine-tuning jobs.
type JobList struct {
	Object   string `json:"object,omitempty"`
	Data     []Job  `json:"data"`
	HasMore  bool   `json:"has_more,omitempty"`
	Provider string `json:"provider,omitempty"`
}
