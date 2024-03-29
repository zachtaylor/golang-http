package session

import "context"

// key is unexported context key type
type key struct{}

// ctxKey is unexported context keying literal
var ctxKey key

// NewContext returns a new Context that carries *session.T
//
// set session=nil to prevent repeated lookups for anonymous connections
func NewContext(ctx context.Context, session *T) context.Context {
	return context.WithValue(ctx, ctxKey, session)
}

// NewContext exposes package-level NewContext
func (t *T) NewContext(ctx context.Context) context.Context { return NewContext(ctx, t) }

// FromContext returns the Session value stored in ctx, if any
//
// returns (nil, true) to indicate Session lookup has already failed
func FromContext(ctx context.Context) (session *T, ok bool) {
	session, ok = ctx.Value(ctxKey).(*T)
	return
}
