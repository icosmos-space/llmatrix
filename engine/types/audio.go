package types

// SpeechRequest is a text-to-speech request (OpenAI /audio/speech).
type SpeechRequest struct {
	Model          string   `json:"model"`
	Input          string   `json:"input"`
	Voice          string   `json:"voice"`
	ResponseFormat string   `json:"response_format,omitempty"` // mp3, opus, aac, flac, wav, pcm
	Speed          *float64 `json:"speed,omitempty"`
}

// SpeechResponse holds synthesized audio bytes.
type SpeechResponse struct {
	ContentType string
	Audio       []byte
	Provider    string
	Model       string
}

// TranscriptionRequest is a speech-to-text request (OpenAI /audio/transcriptions).
type TranscriptionRequest struct {
	Model          string
	FileName       string
	Audio          []byte
	Language       string
	Prompt         string
	ResponseFormat string // json, text, verbose_json, srt, vtt
	Temperature    *float64
}

// TranscriptionResponse is a speech-to-text response.
type TranscriptionResponse struct {
	Text     string  `json:"text"`
	Language string  `json:"language,omitempty"`
	Duration float64 `json:"duration,omitempty"`
	Provider string  `json:"provider,omitempty"`
	Model    string  `json:"model,omitempty"`
}
