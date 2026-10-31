package types

// ResponsesRequest is an OpenAI / LiteLLM Responses API create request.
// Input may be a string or a structured list of items (encoded as any).
type ResponsesRequest struct {
	Model           string            `json:"model"`
	Input           any               `json:"input"`
	Instructions    string            `json:"instructions,omitempty"`
	MaxOutputTokens *int              `json:"max_output_tokens,omitempty"`
	Temperature     *float64          `json:"temperature,omitempty"`
	TopP            *float64          `json:"top_p,omitempty"`
	Tools           []any             `json:"tools,omitempty"`
	ToolChoice      any               `json:"tool_choice,omitempty"`
	Stream          bool              `json:"stream,omitempty"`
	User            string            `json:"user,omitempty"`
	Metadata        map[string]string `json:"metadata,omitempty"`
	Store           *bool             `json:"store,omitempty"`
}

// ResponsesContentPart is a content part inside a message output item.
type ResponsesContentPart struct {
	Type        string `json:"type"`
	Text        string `json:"text,omitempty"`
	Annotations []any  `json:"annotations,omitempty"`
}

// ResponsesOutputItem is one item in the response output array.
type ResponsesOutputItem struct {
	Type    string                 `json:"type"`
	ID      string                 `json:"id,omitempty"`
	Status  string                 `json:"status,omitempty"`
	Role    string                 `json:"role,omitempty"`
	Content []ResponsesContentPart `json:"content,omitempty"`
	Result  string                 `json:"result,omitempty"` // image_generation_call etc.
}

// ResponsesUsage is token usage for a Responses call.
type ResponsesUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

// ResponsesResponse is an OpenAI Responses API response object.
type ResponsesResponse struct {
	ID       string                `json:"id,omitempty"`
	Object   string                `json:"object,omitempty"`
	Created  int64                 `json:"created_at,omitempty"`
	Status   string                `json:"status,omitempty"`
	Model    string                `json:"model,omitempty"`
	Provider string                `json:"provider,omitempty"`
	Output   []ResponsesOutputItem `json:"output"`
	Usage    ResponsesUsage        `json:"usage"`
	Error    *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// OutputText concatenates assistant output_text parts (convenience helper).
func (r *ResponsesResponse) OutputText() string {
	if r == nil {
		return ""
	}
	var s string
	for _, item := range r.Output {
		if item.Type != "message" && item.Type != "" {
			continue
		}
		for _, part := range item.Content {
			if part.Type == "output_text" || part.Type == "text" {
				s += part.Text
			}
		}
	}
	return s
}
