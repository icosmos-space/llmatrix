package responses

// Request is a LiteLLM / OpenAI Responses API create request.
type Request struct {
	Model              string            `json:"model,omitempty"`
	Input              any               `json:"input"`
	Instructions       string            `json:"instructions,omitempty"`
	MaxOutputTokens    *int              `json:"max_output_tokens,omitempty"`
	Temperature        *float64          `json:"temperature,omitempty"`
	TopP               *float64          `json:"top_p,omitempty"`
	Tools              []any             `json:"tools,omitempty"`
	ToolChoice         any               `json:"tool_choice,omitempty"`
	Stream             bool              `json:"stream,omitempty"`
	User               string            `json:"user,omitempty"`
	Metadata           map[string]string `json:"metadata,omitempty"`
	Store              *bool             `json:"store,omitempty"`
	PreviousResponseID string            `json:"previous_response_id,omitempty"`
}

// ContentPart is a content part inside a message output item.
type ContentPart struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

// OutputItem is one item in the response output array.
type OutputItem struct {
	Type    string        `json:"type"`
	ID      string        `json:"id,omitempty"`
	Status  string        `json:"status,omitempty"`
	Role    string        `json:"role,omitempty"`
	Content []ContentPart `json:"content,omitempty"`
	Result  string        `json:"result,omitempty"`
}

// Usage is token usage for a Responses call.
type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

// Response is an OpenAI Responses API response object.
type Response struct {
	ID       string       `json:"id,omitempty"`
	Object   string       `json:"object,omitempty"`
	Created  int64        `json:"created_at,omitempty"`
	Status   string       `json:"status,omitempty"`
	Model    string       `json:"model,omitempty"`
	Provider string       `json:"provider,omitempty"`
	Output   []OutputItem `json:"output"`
	Usage    Usage        `json:"usage"`
}

// OutputText concatenates assistant output_text parts.
func (r *Response) OutputText() string {
	if r == nil {
		return ""
	}
	var s string
	for _, item := range r.Output {
		for _, part := range item.Content {
			if part.Type == "output_text" || part.Type == "text" {
				s += part.Text
			}
		}
	}
	return s
}
