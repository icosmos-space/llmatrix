package types

// ModerationRequest is an OpenAI /moderations request (LiteLLM moderation()).
type ModerationRequest struct {
	Model string `json:"model,omitempty"`
	// Input is a string, []string, or multimodal objects.
	Input any `json:"input"`
}

// ModerationCategoryScores holds per-category confidence scores.
type ModerationCategoryScores struct {
	Hate                  float64 `json:"hate"`
	HateThreatening       float64 `json:"hate/threatening"`
	Harassment            float64 `json:"harassment"`
	HarassmentThreatening float64 `json:"harassment/threatening"`
	SelfHarm              float64 `json:"self-harm"`
	SelfHarmIntent        float64 `json:"self-harm/intent"`
	SelfHarmInstructions  float64 `json:"self-harm/instructions"`
	Sexual                float64 `json:"sexual"`
	SexualMinors          float64 `json:"sexual/minors"`
	Violence              float64 `json:"violence"`
	ViolenceGraphic       float64 `json:"violence/graphic"`
}

// ModerationCategories holds per-category boolean flags.
type ModerationCategories struct {
	Hate                  bool `json:"hate"`
	HateThreatening       bool `json:"hate/threatening"`
	Harassment            bool `json:"harassment"`
	HarassmentThreatening bool `json:"harassment/threatening"`
	SelfHarm              bool `json:"self-harm"`
	SelfHarmIntent        bool `json:"self-harm/intent"`
	SelfHarmInstructions  bool `json:"self-harm/instructions"`
	Sexual                bool `json:"sexual"`
	SexualMinors          bool `json:"sexual/minors"`
	Violence              bool `json:"violence"`
	ViolenceGraphic       bool `json:"violence/graphic"`
}

// ModerationResult is one classification result.
type ModerationResult struct {
	Flagged        bool                     `json:"flagged"`
	Categories     ModerationCategories     `json:"categories"`
	CategoryScores ModerationCategoryScores `json:"category_scores"`
}

// ModerationResponse is an OpenAI moderation response.
type ModerationResponse struct {
	ID       string             `json:"id,omitempty"`
	Model    string             `json:"model,omitempty"`
	Provider string             `json:"provider,omitempty"`
	Results  []ModerationResult `json:"results"`
}
