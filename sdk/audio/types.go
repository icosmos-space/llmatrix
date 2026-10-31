package audio

// SpeechRequest is text-to-speech.
type SpeechRequest struct {
	Model          string   `json:"model,omitempty"`
	Input          string   `json:"input"`
	Voice          string   `json:"voice"`
	ResponseFormat string   `json:"response_format,omitempty"`
	Speed          *float64 `json:"speed,omitempty"`
}

// SpeechResponse holds audio bytes.
type SpeechResponse struct {
	ContentType string `json:"content_type,omitempty"`
	Audio       []byte `json:"-"`
	Provider    string `json:"provider,omitempty"`
	Model       string `json:"model,omitempty"`
}

// TranscriptionRequest is speech-to-text.
type TranscriptionRequest struct {
	Model          string
	FileName       string
	Audio          []byte
	Language       string
	Prompt         string
	ResponseFormat string
	Temperature    *float64
}

// TranscriptionResponse is speech-to-text output.
type TranscriptionResponse struct {
	Text     string  `json:"text"`
	Language string  `json:"language,omitempty"`
	Duration float64 `json:"duration,omitempty"`
	Provider string  `json:"provider,omitempty"`
	Model    string  `json:"model,omitempty"`
}
