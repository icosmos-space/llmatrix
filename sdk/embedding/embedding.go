package embedding

import "context"

// Service defines embedding operations.
type Service interface {
	Embed(ctx context.Context, req *Request) (*Response, error)
}
