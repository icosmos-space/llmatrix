package image

// Request is an image generation request.
type Request struct {
	Model          string `json:"model,omitempty"`
	Prompt         string `json:"prompt"`
	N              *int   `json:"n,omitempty"`
	Size           string `json:"size,omitempty"`
	Quality        string `json:"quality,omitempty"`
	Style          string `json:"style,omitempty"`
	ResponseFormat string `json:"response_format,omitempty"`
	User           string `json:"user,omitempty"`
}

// EditRequest is an image edit request (LiteLLM image_edit()).
type EditRequest struct {
	Model          string `json:"model,omitempty"`
	Prompt         string `json:"prompt"`
	Image          []byte `json:"-"`
	ImageFileName  string `json:"-"`
	Mask           []byte `json:"-"`
	MaskFileName   string `json:"-"`
	N              *int   `json:"n,omitempty"`
	Size           string `json:"size,omitempty"`
	ResponseFormat string `json:"response_format,omitempty"`
	User           string `json:"user,omitempty"`
}

// Data is one generated image.
type Data struct {
	URL     string `json:"url,omitempty"`
	B64JSON string `json:"b64_json,omitempty"`
	Revised string `json:"revised_prompt,omitempty"`
}

// Response is an image generation response.
type Response struct {
	Created  int64  `json:"created,omitempty"`
	Data     []Data `json:"data"`
	Provider string `json:"provider,omitempty"`
	Model    string `json:"model,omitempty"`
}
