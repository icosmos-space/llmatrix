package types

// Role identifies the speaker of a message.
type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleTool      Role = "tool"
)

// Message is a unified chat message.
type Message struct {
	Role       Role       `json:"role"`
	Content    string     `json:"content"`
	Name       string     `json:"name,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	// Parts, when non-empty, is preferred over Content for providers that
	// support multimodal / cache_control content blocks (Anthropic, etc.).
	Parts []ContentPart `json:"parts,omitempty"`
	// CacheControl attaches to a single text Content block when Parts is empty.
	CacheControl *CacheControl `json:"cache_control,omitempty"`
}
