package types

// TextCompletionRequest is an OpenAI legacy /completions request (LiteLLM text_completion()).
type TextCompletionRequest struct {
	Model            string         `json:"model"`
	Prompt           any            `json:"prompt"` // string or []string
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
	Stream           bool           `json:"stream,omitempty"`
}

// TextCompletionChoice is one text completion candidate.
type TextCompletionChoice struct {
	Index        int    `json:"index"`
	Text         string `json:"text"`
	FinishReason string `json:"finish_reason,omitempty"`
}

// TextCompletionResponse is an OpenAI legacy completions response.
type TextCompletionResponse struct {
	ID       string                 `json:"id,omitempty"`
	Object   string                 `json:"object,omitempty"`
	Created  int64                  `json:"created,omitempty"`
	Model    string                 `json:"model,omitempty"`
	Provider string                 `json:"provider,omitempty"`
	Choices  []TextCompletionChoice `json:"choices"`
	Usage    Usage                  `json:"usage"`
}
