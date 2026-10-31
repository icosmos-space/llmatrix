package assistants

// CreateRequest creates an assistant (LiteLLM create_assistants).
type CreateRequest struct {
	Model        string   `json:"model,omitempty"`
	Name         string   `json:"name,omitempty"`
	Description  string   `json:"description,omitempty"`
	Instructions string   `json:"instructions,omitempty"`
	Tools        []any    `json:"tools,omitempty"`
	Metadata     any      `json:"metadata,omitempty"`
	Temperature  *float64 `json:"temperature,omitempty"`
}

// Assistant is an OpenAI Assistant object.
type Assistant struct {
	ID           string `json:"id"`
	Object       string `json:"object,omitempty"`
	CreatedAt    int64  `json:"created_at,omitempty"`
	Name         string `json:"name,omitempty"`
	Description  string `json:"description,omitempty"`
	Model        string `json:"model,omitempty"`
	Instructions string `json:"instructions,omitempty"`
	Tools        []any  `json:"tools,omitempty"`
	Metadata     any    `json:"metadata,omitempty"`
	Provider     string `json:"provider,omitempty"`
}

// ThreadCreateRequest creates a thread.
type ThreadCreateRequest struct {
	Messages []Message `json:"messages,omitempty"`
	Metadata any       `json:"metadata,omitempty"`
	Model    string    `json:"model,omitempty"` // credential routing hint
}

// Message is a thread message payload.
type Message struct {
	Role    string `json:"role"`
	Content any    `json:"content"`
}

// Thread is an OpenAI Thread object.
type Thread struct {
	ID        string `json:"id"`
	Object    string `json:"object,omitempty"`
	CreatedAt int64  `json:"created_at,omitempty"`
	Metadata  any    `json:"metadata,omitempty"`
	Provider  string `json:"provider,omitempty"`
}

// MessageObject is a stored thread message.
type MessageObject struct {
	ID          string `json:"id"`
	Object      string `json:"object,omitempty"`
	CreatedAt   int64  `json:"created_at,omitempty"`
	ThreadID    string `json:"thread_id,omitempty"`
	Role        string `json:"role,omitempty"`
	Content     any    `json:"content,omitempty"`
	AssistantID string `json:"assistant_id,omitempty"`
	RunID       string `json:"run_id,omitempty"`
	Provider    string `json:"provider,omitempty"`
}

// MessageList lists thread messages.
type MessageList struct {
	Object   string          `json:"object,omitempty"`
	Data     []MessageObject `json:"data"`
	Provider string          `json:"provider,omitempty"`
}

// RunCreateRequest starts a run.
type RunCreateRequest struct {
	AssistantID  string   `json:"assistant_id"`
	Model        string   `json:"model,omitempty"`
	Instructions string   `json:"instructions,omitempty"`
	Tools        []any    `json:"tools,omitempty"`
	Metadata     any      `json:"metadata,omitempty"`
	Temperature  *float64 `json:"temperature,omitempty"`
}

// Run is an OpenAI Run object.
type Run struct {
	ID           string `json:"id"`
	Object       string `json:"object,omitempty"`
	CreatedAt    int64  `json:"created_at,omitempty"`
	ThreadID     string `json:"thread_id,omitempty"`
	AssistantID  string `json:"assistant_id,omitempty"`
	Status       string `json:"status,omitempty"`
	Model        string `json:"model,omitempty"`
	Instructions string `json:"instructions,omitempty"`
	Provider     string `json:"provider,omitempty"`
}
