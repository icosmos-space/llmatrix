package text

// Request is a LiteLLM text_completion() request.
type Request struct {
	Model            string         `json:"model,omitempty"`
	Prompt           any            `json:"prompt"`
	MaxTokens        *int           `json:"max_tokens,omitempty"`
	Temperature      *float64       `json:"temperature,omitempty"`
	TopP             *float64       `json:"top_p,omitempty"`
	N                *int           `json:"n,omitempty"`
	Stop             []string       `json:"stop,omitempty"`
	PresencePenalty  *float64       `json:"presence_penalty,omitempty"`
	FrequencyPenalty *float64       `json:"frequency_penalty,omitempty"`
	LogitBias        map[string]int `json:"logit_bias,omitempty"`
	User             string         `json:"user,omitempty"`
	Seed             *int           `json:"seed,omitempty"`
	Suffix           string         `json:"suffix,omitempty"`
	Echo             *bool          `json:"echo,omitempty"`
	BestOf           *int           `json:"best_of,omitempty"`
	Logprobs         *int           `json:"logprobs,omitempty"`
	Caching          *bool          `json:"caching,omitempty"`
}

// Choice is one text completion candidate.
type Choice struct {
	Index        int    `json:"index"`
	Text         string `json:"text"`
	FinishReason string `json:"finish_reason,omitempty"`
}

// Usage reports token consumption.
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// Response is a text completion response.
type Response struct {
	ID       string   `json:"id,omitempty"`
	Object   string   `json:"object,omitempty"`
	Created  int64    `json:"created,omitempty"`
	Model    string   `json:"model,omitempty"`
	Provider string   `json:"provider,omitempty"`
	Choices  []Choice `json:"choices"`
	Usage    Usage    `json:"usage"`
}

// Text returns the first choice text.
func (r *Response) Text() string {
	if r == nil || len(r.Choices) == 0 {
		return ""
	}
	return r.Choices[0].Text
}
