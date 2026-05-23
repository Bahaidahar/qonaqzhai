package auth

import (
	"context"
	"net/http"
	"strings"
	"time"

	"qonaqzhai-backend/pkg/httpx"
)

// TokenVerifier is the minimal interface the middleware needs; satisfied by
// *Verifier and easy to fake in tests.
type TokenVerifier interface {
	Verify(ctx context.Context, token string) (Claims, error)
}

// Middleware authenticates Bearer tokens by delegating to a TokenVerifier.
type Middleware struct {
	v       TokenVerifier
	timeout time.Duration
}

// NewMiddleware constructs an auth middleware. Each verify call uses timeout
// (default 3s when zero).
func NewMiddleware(v TokenVerifier, timeout time.Duration) *Middleware {
	if timeout == 0 {
		timeout = 3 * time.Second
	}
	return &Middleware{v: v, timeout: timeout}
}

// Required rejects requests without a valid Bearer token.
func (m *Middleware) Required(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, ok := m.parse(r)
		if !ok {
			httpx.WriteError(w, http.StatusUnauthorized, "missing or invalid token")
			return
		}
		next.ServeHTTP(w, r.WithContext(WithClaims(r.Context(), c)))
	})
}

// Optional attaches claims when present, but does not reject anonymous calls.
func (m *Middleware) Optional(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if c, ok := m.parse(r); ok {
			next.ServeHTTP(w, r.WithContext(WithClaims(r.Context(), c)))
			return
		}
		next.ServeHTTP(w, r)
	})
}

// RequireRole returns a middleware allowing only the listed roles.
func (m *Middleware) RequireRole(roles ...string) func(http.Handler) http.Handler {
	allowed := make(map[string]bool, len(roles))
	for _, r := range roles {
		allowed[r] = true
	}
	return func(next http.Handler) http.Handler {
		return m.Required(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, _ := RoleFrom(r.Context())
			if !allowed[role] {
				httpx.WriteError(w, http.StatusForbidden, "forbidden")
				return
			}
			next.ServeHTTP(w, r)
		}))
	}
}

func (m *Middleware) parse(r *http.Request) (Claims, bool) {
	raw := r.Header.Get("Authorization")
	if !strings.HasPrefix(raw, "Bearer ") {
		return Claims{}, false
	}
	tok := strings.TrimPrefix(raw, "Bearer ")
	ctx, cancel := context.WithTimeout(r.Context(), m.timeout)
	defer cancel()
	c, err := m.v.Verify(ctx, tok)
	if err != nil || c.UserID == "" {
		return Claims{}, false
	}
	if !c.IsActive() {
		return Claims{}, false
	}
	return c, true
}
