// Package auth provides a remote JWT verifier (gRPC client to auth-svc) plus an
// HTTP middleware that injects authenticated claims into request context.
//
// Every service except auth-svc itself should use this package — only auth-svc
// owns the JWT secret.
package auth

import "context"

// Claims is the verified principal returned by the auth-svc.
type Claims struct {
	UserID string
	Email  string
	Role   string
	Status string
}

// IsActive reports whether the principal is active and may act.
func (c Claims) IsActive() bool { return c.Status == "" || c.Status == "active" }

type ctxKey struct{}

// WithClaims attaches claims to ctx.
func WithClaims(ctx context.Context, c Claims) context.Context {
	return context.WithValue(ctx, ctxKey{}, c)
}

// FromContext extracts claims previously injected by middleware.
func FromContext(ctx context.Context) (Claims, bool) {
	c, ok := ctx.Value(ctxKey{}).(Claims)
	return c, ok
}

// UserIDFrom is a convenience accessor.
func UserIDFrom(ctx context.Context) (string, bool) {
	c, ok := FromContext(ctx)
	return c.UserID, ok && c.UserID != ""
}

// RoleFrom is a convenience accessor.
func RoleFrom(ctx context.Context) (string, bool) {
	c, ok := FromContext(ctx)
	return c.Role, ok
}
