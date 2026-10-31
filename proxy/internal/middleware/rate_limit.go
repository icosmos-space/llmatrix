package middleware

import "net/http"

// RateLimit is a placeholder for request rate limiting.
func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: per-key / per-IP limiting
		next.ServeHTTP(w, r)
	})
}
