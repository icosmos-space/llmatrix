package chat

import "context"

// Service defines chat completion operations.
type Service interface {
	Complete(ctx context.Context, req *Request) (*Response, error)
}
