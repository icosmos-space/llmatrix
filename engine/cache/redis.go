package cache

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Redis is a LiteLLM-style Redis/Valkey response cache (stdlib RESP client).
type Redis struct {
	mu         sync.Mutex
	addr       string
	password   string
	db         int
	prefix     string
	defaultTTL time.Duration
	conn       net.Conn
	reader     *bufio.Reader
}

// RedisOption configures Redis.
type RedisOption func(*Redis)

// WithRedisPassword sets AUTH password.
func WithRedisPassword(password string) RedisOption {
	return func(r *Redis) { r.password = password }
}

// WithRedisDB selects a logical database index.
func WithRedisDB(db int) RedisOption {
	return func(r *Redis) { r.db = db }
}

// WithRedisPrefix namespaces keys.
func WithRedisPrefix(prefix string) RedisOption {
	return func(r *Redis) { r.prefix = prefix }
}

// NewRedis connects lazily to addr (host:port).
func NewRedis(addr string, defaultTTL time.Duration, opts ...RedisOption) *Redis {
	if defaultTTL <= 0 {
		defaultTTL = DefaultTTL
	}
	r := &Redis{addr: addr, defaultTTL: defaultTTL, prefix: "llmatrix:"}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func (r *Redis) key(k string) string { return r.prefix + k }

func (r *Redis) dial() error {
	if r.conn != nil {
		return nil
	}
	c, err := net.DialTimeout("tcp", r.addr, 5*time.Second)
	if err != nil {
		return err
	}
	r.conn = c
	r.reader = bufio.NewReader(c)
	if r.password != "" {
		if err := r.cmd("AUTH", r.password); err != nil {
			_ = c.Close()
			r.conn = nil
			return err
		}
	}
	if r.db > 0 {
		if err := r.cmd("SELECT", strconv.Itoa(r.db)); err != nil {
			_ = c.Close()
			r.conn = nil
			return err
		}
	}
	return nil
}

func (r *Redis) reset() {
	if r.conn != nil {
		_ = r.conn.Close()
	}
	r.conn = nil
	r.reader = nil
}

func (r *Redis) cmd(args ...string) error {
	_, err := r.cmdValue(args...)
	return err
}

func (r *Redis) cmdValue(args ...string) (any, error) {
	if err := r.dial(); err != nil {
		return nil, err
	}
	var b strings.Builder
	b.WriteString("*")
	b.WriteString(strconv.Itoa(len(args)))
	b.WriteString("\r\n")
	for _, a := range args {
		b.WriteString("$")
		b.WriteString(strconv.Itoa(len(a)))
		b.WriteString("\r\n")
		b.WriteString(a)
		b.WriteString("\r\n")
	}
	_ = r.conn.SetDeadline(time.Now().Add(5 * time.Second))
	if _, err := r.conn.Write([]byte(b.String())); err != nil {
		r.reset()
		return nil, err
	}
	v, err := r.readRESP()
	if err != nil {
		r.reset()
		return nil, err
	}
	return v, nil
}

func (r *Redis) readRESP() (any, error) {
	line, err := r.reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	line = strings.TrimSuffix(strings.TrimSuffix(line, "\n"), "\r")
	if line == "" {
		return nil, fmt.Errorf("redis: empty response")
	}
	switch line[0] {
	case '+':
		return line[1:], nil
	case '-':
		return nil, fmt.Errorf("redis: %s", line[1:])
	case ':':
		n, err := strconv.ParseInt(line[1:], 10, 64)
		return n, err
	case '$':
		n, err := strconv.Atoi(line[1:])
		if err != nil {
			return nil, err
		}
		if n < 0 {
			return nil, nil // nil bulk
		}
		buf := make([]byte, n+2)
		if _, err := ioReadFull(r.reader, buf); err != nil {
			return nil, err
		}
		return buf[:n], nil
	case '*':
		n, err := strconv.Atoi(line[1:])
		if err != nil {
			return nil, err
		}
		if n < 0 {
			return nil, nil
		}
		out := make([]any, 0, n)
		for i := 0; i < n; i++ {
			v, err := r.readRESP()
			if err != nil {
				return nil, err
			}
			out = append(out, v)
		}
		return out, nil
	default:
		return nil, fmt.Errorf("redis: unknown prefix %q", line)
	}
}

func ioReadFull(r *bufio.Reader, buf []byte) (int, error) {
	n := 0
	for n < len(buf) {
		nn, err := r.Read(buf[n:])
		n += nn
		if err != nil {
			return n, err
		}
	}
	return n, nil
}

// Get returns a JSON payload ([]byte) if present.
func (r *Redis) Get(key string) (any, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	v, err := r.cmdValue("GET", r.key(key))
	if err != nil || v == nil {
		return nil, false
	}
	b, ok := v.([]byte)
	if !ok || len(b) == 0 {
		return nil, false
	}
	return b, true
}

// Set stores value as JSON with TTL via SETEX.
func (r *Redis) Set(key string, value any, ttl time.Duration) {
	if ttl <= 0 {
		ttl = r.defaultTTL
	}
	payload, err := json.Marshal(value)
	if err != nil {
		return
	}
	secs := int(ttl.Seconds())
	if secs < 1 {
		secs = 1
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	_ = r.cmd("SETEX", r.key(key), strconv.Itoa(secs), string(payload))
}

// Delete removes a key.
func (r *Redis) Delete(key string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	_ = r.cmd("DEL", r.key(key))
}

// Close closes the underlying connection.
func (r *Redis) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.conn == nil {
		return nil
	}
	err := r.conn.Close()
	r.conn = nil
	r.reader = nil
	return err
}
