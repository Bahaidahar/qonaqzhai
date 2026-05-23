package httpx

import (
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// CORS wraps next with permissive CORS headers for origin.
func CORS(origin string, next http.Handler) http.Handler {
	if origin == "" {
		origin = "*"
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")
		w.Header().Set("Vary", "Origin")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Recover converts panics into 500 responses.
func Recover(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					log.Error("panic", "err", rec, "path", r.URL.Path)
					WriteError(w, http.StatusInternalServerError, "internal error")
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// AccessLog emits one slog entry per request with status + duration.
func AccessLog(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			lw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
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

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(c int) {
	w.status = c
	w.ResponseWriter.WriteHeader(c)
}

// RateLimiter throttles requests per key (IP or user id).
type RateLimiter struct {
	mu     sync.Mutex
	visits map[string]*rate.Limiter
	rate   rate.Limit
	burst  int
}

// NewRateLimiter constructs a per-key limiter allowing r events/sec with burst b.
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

// PerIP returns a middleware limiting by client IP.
func (l *RateLimiter) PerIP() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !l.Allow(ClientIP(r)) {
				WriteError(w, http.StatusTooManyRequests, "rate limited")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// ClientIP best-effort extracts the originating client IP. Honors
// X-Forwarded-For and X-Real-IP, falling back to RemoteAddr.
func ClientIP(r *http.Request) string {
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
