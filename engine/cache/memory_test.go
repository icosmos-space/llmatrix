package cache

import (
	"testing"
	"time"
)

func TestMemoryCache(t *testing.T) {
	c := NewMemory(50 * time.Millisecond)
	c.Set("a", "v1", 0)
	if v, ok := c.Get("a"); !ok || v != "v1" {
		t.Fatalf("get=%v %v", v, ok)
	}
	time.Sleep(60 * time.Millisecond)
	if _, ok := c.Get("a"); ok {
		t.Fatal("expected expire")
	}
}

func TestKeyStable(t *testing.T) {
	k1 := Key("chat", "m", map[string]any{"x": 1})
	k2 := Key("chat", "m", map[string]any{"x": 1})
	if k1 != k2 || k1 == "" {
		t.Fatalf("%s %s", k1, k2)
	}
}
