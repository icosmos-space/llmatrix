package streaming

// EventType classifies stream events.
type EventType string

const (
	EventTypeDelta EventType = "delta"
	EventTypeError EventType = "error"
	EventTypeDone  EventType = "done"
)

// Event is a single streaming payload.
type Event struct {
	Type    EventType `json:"type"`
	Delta   string    `json:"delta,omitempty"`
	Error   string    `json:"error,omitempty"`
	RawJSON []byte    `json:"-"`
}
