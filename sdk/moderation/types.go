package moderation

// Request is a LiteLLM moderation() request.
type Request struct {
	Model string `json:"model,omitempty"`
	Input any    `json:"input"`
}

// CategoryScores holds per-category confidence scores.
type CategoryScores struct {
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

// Categories holds per-category boolean flags.
type Categories struct {
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

// Result is one classification result.
type Result struct {
	Flagged        bool           `json:"flagged"`
	Categories     Categories     `json:"categories"`
	CategoryScores CategoryScores `json:"category_scores"`
}

// Response is a moderation response.
type Response struct {
	ID       string   `json:"id,omitempty"`
	Model    string   `json:"model,omitempty"`
	Provider string   `json:"provider,omitempty"`
	Results  []Result `json:"results"`
}

// Flagged reports whether any result is flagged.
func (r *Response) Flagged() bool {
	if r == nil {
		return false
	}
	for _, res := range r.Results {
		if res.Flagged {
			return true
		}
	}
	return false
}
