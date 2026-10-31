package retry

import (
	"context"
	"testing"
	"time"

	"github.com/icosmos-space/llmatrix/engine/errors"
)

func TestDoRetriesRateLimit(t *testing.T) {
	attempts := 0
	err := Do(context.Background(), Policy{
		MaxAttempts: 3,
		Initial:     time.Millisecond,
		MaxBackoff:  5 * time.Millisecond,
		Multiplier:  2,
		Jitter:      false,
	}, func(context.Context) error {
		attempts++
		if attempts < 3 {
			return errors.WrapStatus(errors.KindRateLimited, "test", "slow down", 429, nil)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if attempts != 3 {
		t.Fatalf("attempts=%d", attempts)
	}
}

func TestDoDoesNotRetryAuth(t *testing.T) {
	attempts := 0
	err := Do(context.Background(), DefaultPolicy(), func(context.Context) error {
		attempts++
		return errors.WrapStatus(errors.KindUnauthorized, "test", "bad key", 401, nil)
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if attempts != 1 {
		t.Fatalf("attempts=%d", attempts)
	}
}
