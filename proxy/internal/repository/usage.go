package repository

import (
	"context"
	"time"

	"github.com/icosmos-space/llmatrix/engine/types"
)

// UsageRecord is a persisted usage event.
type UsageRecord struct {
	ID        string
	UserID    string
	Provider  string
	Model     string
	Usage     types.Usage
	CreatedAt time.Time
}

// UsageRepository persists usage records.
type UsageRepository interface {
	Insert(ctx context.Context, record *UsageRecord) error
}
