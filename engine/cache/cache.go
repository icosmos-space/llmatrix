package cache

import "time"

// Cache is a LiteLLM-style response cache backend.
type Cache interface {
	Get(key string) (any, bool)
	Set(key string, value any, ttl time.Duration)
	Delete(key string)
}

// DefaultTTL is used when callers pass a non-positive TTL.
const DefaultTTL = 1 * time.Hour
