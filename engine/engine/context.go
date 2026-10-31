package engine

import "context"

type ctxKey int

const (
	keyRequestID ctxKey = iota
	keyProvider
)

// WithRequestID attaches a request ID to the context.
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, keyRequestID, id)
}

// RequestID extracts the request ID from context, if present.
func RequestID(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(keyRequestID).(string)
	return v, ok
}

// WithProvider attaches the selected provider name to the context.
func WithProvider(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, keyProvider, name)
}

// ProviderName extracts the provider name from context, if present.
func ProviderName(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(keyProvider).(string)
	return v, ok
}
