package retry

import (
	"context"
	"math"
	"math/rand"
	"time"

	"github.com/icosmos-space/llmatrix/engine/errors"
)

// Policy controls retry behavior.
type Policy struct {
	MaxAttempts int
	Initial     time.Duration
	MaxBackoff  time.Duration
	Multiplier  float64
	Jitter      bool
}

// DefaultPolicy returns sensible defaults for LLM HTTP calls.
func DefaultPolicy() Policy {
	return Policy{
		MaxAttempts: 3,
		Initial:     200 * time.Millisecond,
		MaxBackoff:  5 * time.Second,
		Multiplier:  2,
		Jitter:      true,
	}
}

// Do runs fn with retries for retryable errors.
func Do(ctx context.Context, policy Policy, fn func(context.Context) error) error {
	if policy.MaxAttempts <= 0 {
		policy.MaxAttempts = 1
	}
	if policy.Initial <= 0 {
		policy.Initial = 100 * time.Millisecond
	}
	if policy.MaxBackoff <= 0 {
		policy.MaxBackoff = 5 * time.Second
	}
	if policy.Multiplier < 1 {
		policy.Multiplier = 2
	}

	var last error
	for attempt := 1; attempt <= policy.MaxAttempts; attempt++ {
		last = fn(ctx)
		if last == nil {
			return nil
		}
		if attempt == policy.MaxAttempts || !errors.IsRetryable(last) {
			return last
		}
		delay := backoff(policy, attempt)
		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
	}
	return last
}

func backoff(policy Policy, attempt int) time.Duration {
	// attempt is 1-based; delay after failure attempt n uses n as exponent base.
	exp := math.Pow(policy.Multiplier, float64(attempt-1))
	d := time.Duration(float64(policy.Initial) * exp)
	if d > policy.MaxBackoff {
		d = policy.MaxBackoff
	}
	if policy.Jitter && d > 0 {
		// Full jitter: [0, d]
		d = time.Duration(rand.Int63n(int64(d) + 1))
	}
	return d
}
