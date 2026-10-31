package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type diskRecord struct {
	Expires int64           `json:"expires"`
	Payload json.RawMessage `json:"payload"`
}

// Disk is a file-backed Cache (LiteLLM Cache(type="disk")).
// Values are stored as JSON; Get returns []byte JSON payloads.
type Disk struct {
	mu         sync.Mutex
	dir        string
	defaultTTL time.Duration
}

// NewDisk creates a disk cache under dir (created if missing).
func NewDisk(dir string, defaultTTL time.Duration) (*Disk, error) {
	if dir == "" {
		return nil, os.ErrInvalid
	}
	if defaultTTL <= 0 {
		defaultTTL = DefaultTTL
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	return &Disk{dir: dir, defaultTTL: defaultTTL}, nil
}

func (d *Disk) path(key string) string {
	// Keep filenames filesystem-safe (keys are hex hashes).
	return filepath.Join(d.dir, key+".json")
}

// Get returns a JSON payload if present and not expired.
func (d *Disk) Get(key string) (any, bool) {
	d.mu.Lock()
	defer d.mu.Unlock()
	b, err := os.ReadFile(d.path(key))
	if err != nil {
		return nil, false
	}
	var rec diskRecord
	if err := json.Unmarshal(b, &rec); err != nil {
		return nil, false
	}
	if rec.Expires > 0 && time.Now().UnixNano() > rec.Expires {
		_ = os.Remove(d.path(key))
		return nil, false
	}
	return []byte(rec.Payload), true
}

// Set stores a value as JSON with TTL.
func (d *Disk) Set(key string, value any, ttl time.Duration) {
	if ttl <= 0 {
		ttl = d.defaultTTL
	}
	payload, err := json.Marshal(value)
	if err != nil {
		return
	}
	rec := diskRecord{
		Expires: time.Now().Add(ttl).UnixNano(),
		Payload: payload,
	}
	b, err := json.Marshal(rec)
	if err != nil {
		return
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	_ = os.WriteFile(d.path(key), b, 0o644)
}

// Delete removes a key.
func (d *Disk) Delete(key string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	_ = os.Remove(d.path(key))
}
