package errors

import "testing"

func TestIsRetryable(t *testing.T) {
	cases := []struct {
		err  error
		want bool
	}{
		{WrapStatus(KindRateLimited, "x", "slow", 429, nil), true},
		{WrapStatus(KindUnauthorized, "x", "bad", 401, nil), false},
		{WrapStatus(KindInvalidRequest, "x", "bad", 400, nil), false},
		{WrapStatus(KindTimeout, "x", "timeout", 504, nil), true},
		{WrapStatus(KindProvider, "x", "oops", 500, nil), true},
	}
	for _, tc := range cases {
		if got := IsRetryable(tc.err); got != tc.want {
			t.Fatalf("%v: got %v want %v", tc.err, got, tc.want)
		}
	}
}
