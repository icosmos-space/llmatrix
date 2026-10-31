package repository

import "context"

// Channel is an upstream provider channel.
type Channel struct {
	ID       string
	Provider string
	BaseURL  string
	APIKey   string
	Weight   int
	Enabled  bool
}

// ChannelRepository persists channels.
type ChannelRepository interface {
	ListEnabled(ctx context.Context) ([]Channel, error)
	Get(ctx context.Context, id string) (*Channel, error)
}
