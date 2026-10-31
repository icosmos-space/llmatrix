package types

// ImageRequest is an image generation request (OpenAI /images/generations).
type ImageRequest struct {
	Model          string `json:"model,omitempty"`
	Prompt         string `json:"prompt"`
	N              *int   `json:"n,omitempty"`
	Size           string `json:"size,omitempty"`
	Quality        string `json:"quality,omitempty"`
	Style          string `json:"style,omitempty"`
	ResponseFormat string `json:"response_format,omitempty"` // url | b64_json
	User           string `json:"user,omitempty"`
}

// ImageData is one generated image.
type ImageData struct {
	URL     string `json:"url,omitempty"`
	B64JSON string `json:"b64_json,omitempty"`
	Revised string `json:"revised_prompt,omitempty"`
}

// ImageResponse is an image generation response.
type ImageResponse struct {
	Created  int64       `json:"created,omitempty"`
	Data     []ImageData `json:"data"`
	Provider string      `json:"provider,omitempty"`
	Model    string      `json:"model,omitempty"`
}
