package cid

import "context"

type cidKey struct{}

// NewContext returns a new context with given correlationId value.
func NewContext(ctx context.Context, cid string) context.Context {
	return context.WithValue(ctx, cidKey{}, cid)
}

// FromContext returns the correlationId from the given context,
// or false as the second return value if no correlation id was found in context.
func FromContext(ctx context.Context) (string, bool) {
	cid, ok := ctx.Value(cidKey{}).(string)
	return cid, ok
}
