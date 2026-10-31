package server

import (
	"net/http"
	"time"
)

// Server wraps http.Server with proxy defaults.
type Server struct {
	inner *http.Server
}

// New creates a Server.
func New(addr string, handler http.Handler) *Server {
	return &Server{
		inner: &http.Server{
			Addr:              addr,
			Handler:           handler,
			ReadHeaderTimeout: 10 * time.Second,
		},
	}
}

// ListenAndServe starts serving.
func (s *Server) ListenAndServe() error {
	return s.inner.ListenAndServe()
}
