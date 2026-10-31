package cache

import (
	"sync"
	"time"
)

type entry struct {
	value   any
	expires time.Time
}

// Memory is an in-memory Cache (LiteLLM Cache() default).
type Memory struct {
	mu         sync.RWMutex
	items      map[string]entry
	defaultTTL time.Duration
}

// NewMemory creates an in-memory cache with the given default TTL.
func NewMemory(ttl time.Duration) *Memory {
	if ttl <= 0 {
		ttl = DefaultTTL
	}
	return &Memory{items: make(map[string]entry), defaultTTL: ttl}
}

// Get returns a cached value if present and not expired.
func (m *Memory) Get(key string) (any, bool) {
	m.mu.RLock()
	e, ok := m.items[key]
	m.mu.RUnlock()
	if !ok {
		return nil, false
	}
	if !e.expires.IsZero() && time.Now().After(e.expires) {
		m.Delete(key)
		return nil, false
	}
	return e.value, true
}

// Set stores a value with TTL (0 uses the cache default).
func (m *Memory) Set(key string, value any, ttl time.Duration) {
	if ttl <= 0 {
		ttl = m.defaultTTL
	}
	m.mu.Lock()
	m.items[key] = entry{value: value, expires: time.Now().Add(ttl)}
	m.mu.Unlock()
}

// Delete removes a key.
func (m *Memory) Delete(key string) {
	m.mu.Lock()
	delete(m.items, key)
	m.mu.Unlock()
}
