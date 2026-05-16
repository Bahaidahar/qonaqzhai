// Package middleware holds cross-cutting HTTP middlewares: auth, CORS, recovery, logger, rate limit.
package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"

	"qonaqzhai-backend/internal/domain"
	"qonaqzhai-backend/internal/usecase"
)

type ctxKey string

const (
	ctxUserID ctxKey = "uid"
	ctxRole   ctxKey = "role"
	ctxEmail  ctxKey = "email"
)

// UserIDFrom extracts the authenticated user id from ctx.
func UserIDFrom(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(ctxUserID).(string)
	return v, ok
}

// RoleFrom extracts the authenticated user role from ctx.
func RoleFrom(ctx context.Context) (domain.Role, bool) {
	v, ok := ctx.Value(ctxRole).(string)
	return domain.Role(v), ok
}

// EmailFrom extracts the authenticated user email from ctx.
func EmailFrom(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(ctxEmail).(string)
	return v, ok
}

// Auth wraps handlers with JWT verification logic.
type Auth struct {
	tokens usecase.TokenIssuer
}

// NewAuth constructs an Auth middleware.
func NewAuth(t usecase.TokenIssuer) *Auth { return &Auth{tokens: t} }

// Optional attaches claims to context when a valid Bearer token is present;
// otherwise it proceeds without claims.
func (a *Auth) Optional(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if c, ok := a.parse(r); ok {
			next.ServeHTTP(w, r.WithContext(injectClaims(r.Context(), c)))
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Required rejects requests without a valid Bearer token.
func (a *Auth) Required(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, ok := a.parse(r)
		if !ok {
			writeError(w, http.StatusUnauthorized, "missing or invalid token")
			return
		}
		next.ServeHTTP(w, r.WithContext(injectClaims(r.Context(), c)))
	})
}

// RequireRole returns a middleware factory that allows only the listed roles.
func (a *Auth) RequireRole(roles ...domain.Role) func(http.Handler) http.Handler {
	allowed := map[domain.Role]bool{}
	for _, r := range roles {
		allowed[r] = true
	}
	return func(next http.Handler) http.Handler {
		return a.Required(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, _ := RoleFrom(r.Context())
			if !allowed[role] {
				writeError(w, http.StatusForbidden, "forbidden")
				return
			}
			next.ServeHTTP(w, r)
		}))
	}
}

func (a *Auth) parse(r *http.Request) (usecase.Claims, bool) {
	raw := r.Header.Get("Authorization")
	if !strings.HasPrefix(raw, "Bearer ") {
		return usecase.Claims{}, false
	}
	tok := strings.TrimPrefix(raw, "Bearer ")
	c, err := a.tokens.Parse(tok)
	if err != nil {
		return usecase.Claims{}, false
	}
	return c, true
}

func injectClaims(ctx context.Context, c usecase.Claims) context.Context {
	ctx = context.WithValue(ctx, ctxUserID, c.UserID)
	ctx = context.WithValue(ctx, ctxRole, string(c.Role))
	ctx = context.WithValue(ctx, ctxEmail, c.Email)
	return ctx
}

// --- CORS ---

// CORS returns a handler wrapping next with CORS headers for origin.
func CORS(origin string, next http.Handler) http.Handler {
	if origin == "" {
		origin = "http://localhost:3000"
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Vary", "Origin")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// --- Recover + Logger ---

// Recover converts panics into 500 responses and logs them.
func Recover(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					log.Error("panic", "err", rec, "path", r.URL.Path)
					writeError(w, http.StatusInternalServerError, "internal error")
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// Logger emits a one-line slog entry per request.
func Logger(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			lw := &loggingWriter{ResponseWriter: w, status: 200}
			next.ServeHTTP(lw, r)
			log.Info("http",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", lw.status),
				slog.Duration("dur", time.Since(start)),
			)
		})
	}
}

type loggingWriter struct {
	http.ResponseWriter
	status int
}

func (lw *loggingWriter) WriteHeader(c int) {
	lw.status = c
	lw.ResponseWriter.WriteHeader(c)
}

// --- Rate limit ---

// RateLimiter throttles requests per key (e.g. IP, user id).
type RateLimiter struct {
	mu     sync.Mutex
	visits map[string]*rate.Limiter
	rate   rate.Limit
	burst  int
}

// NewRateLimiter constructs a per-key limiter allowing `r` events/sec with burst `b`.
func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	return &RateLimiter{visits: map[string]*rate.Limiter{}, rate: r, burst: b}
}

// Allow reports whether a request from key may proceed.
func (l *RateLimiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	lim, ok := l.visits[key]
	if !ok {
		lim = rate.NewLimiter(l.rate, l.burst)
		l.visits[key] = lim
	}
	return lim.Allow()
}

// PerIP returns a middleware that limits by client IP.
func (l *RateLimiter) PerIP() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !l.Allow(clientIP(r)) {
				writeError(w, http.StatusTooManyRequests, "rate limited")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// PerUser returns a middleware that limits by authenticated user id (falls back to IP).
func (l *RateLimiter) PerUser() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key, ok := UserIDFrom(r.Context())
			if !ok || key == "" {
				key = clientIP(r)
			}
			if !l.Allow(key) {
				writeError(w, http.StatusTooManyRequests, "rate limited")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func clientIP(r *http.Request) string {
	if xf := r.Header.Get("X-Forwarded-For"); xf != "" {
		if idx := strings.Index(xf, ","); idx > 0 {
			return strings.TrimSpace(xf[:idx])
		}
		return strings.TrimSpace(xf)
	}
	if h := r.Header.Get("X-Real-IP"); h != "" {
		return h
	}
	host := r.RemoteAddr
	if idx := strings.LastIndex(host, ":"); idx > 0 {
		host = host[:idx]
	}
	return host
}

// --- helpers ---

func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

// silence unused-import warning when this pkg is built without features that use them
var _ = errors.New
