package http

import (
	"log/slog"
	"net/http"

	"golang.org/x/time/rate"

	pkgauth "qonaqzhai-backend/pkg/auth"
	"qonaqzhai-backend/pkg/httpx"
)

// Mux wires routes for auth HTTP endpoints. The middleware is constructed
// in-process — auth verifies its own JWTs locally rather than calling itself.
func Mux(h *Handler, mw *pkgauth.Middleware, corsOrigin string, log *slog.Logger) http.Handler {
	rl := httpx.NewRateLimiter(rate.Limit(20), 40)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", h.Health)
	mux.HandleFunc("POST /api/signup", h.Signup)
	mux.HandleFunc("POST /api/login", h.Login)
	mux.HandleFunc("POST /api/refresh", h.Refresh)
	mux.HandleFunc("POST /api/logout", h.Logout)
	mux.HandleFunc("POST /api/forgot-password", h.ForgotPassword)
	mux.HandleFunc("POST /api/reset-password", h.ResetPassword)
	mux.Handle("GET /api/me", mw.Required(http.HandlerFunc(h.Me)))

	withLimits := rl.PerIP()(mux)
	withRecover := httpx.Recover(log)(withLimits)
	withLogger := httpx.AccessLog(log)(withRecover)
	return httpx.CORS(corsOrigin, withLogger)
}
