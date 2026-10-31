package cache

import (
	"encoding/json"
	"path/filepath"
	"testing"
	"time"
)

func TestDiskCache(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "c")
	c, err := NewDisk(dir, 50*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	c.Set("k", map[string]string{"a": "b"}, 0)
	v, ok := c.Get("k")
	if !ok {
		t.Fatal("miss")
	}
	b, ok := v.([]byte)
	if !ok {
		t.Fatalf("%T", v)
	}
	var m map[string]string
	if err := json.Unmarshal(b, &m); err != nil || m["a"] != "b" {
		t.Fatalf("%v %v", err, m)
	}
	time.Sleep(60 * time.Millisecond)
	if _, ok := c.Get("k"); ok {
		t.Fatal("expected expire")
	}
}
