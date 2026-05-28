// Command gateway runs the public-facing edge service. Requests come in over
// HTTP, JWT is verified once against auth-svc (gRPC), and the request is
// forwarded to the right internal service.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/time/rate"

	pkgauth "qonaqzhai-backend/pkg/auth"
	"qonaqzhai-backend/pkg/config"
	"qonaqzhai-backend/pkg/httpx"
	"qonaqzhai-backend/pkg/logger"

	"qonaqzhai-backend/services/gateway/internal/proxy"
)

func main() {
	log := logger.New("gateway")
	if err := run(log); err != nil {
		log.Error("fatal", "err", err)
		os.Exit(1)
	}
}

func run(log *slog.Logger) error {
	addr := config.EnvOr("GATEWAY_ADDR", ":8080")
	cors := config.EnvOr("CORS_ORIGIN", "*")
	authAddr := config.EnvOr("AUTH_GRPC_ADDR", "localhost:9081")

	authURL := config.EnvOr("AUTH_HTTP_URL", "http://localhost:8081")
	coreURL := config.EnvOr("CORE_HTTP_URL", "http://localhost:8082")
	paymentURL := config.EnvOr("PAYMENT_HTTP_URL", "http://localhost:8083")
	realtimeURL := config.EnvOr("REALTIME_HTTP_URL", "http://localhost:8084")

	verifier, err := pkgauth.NewVerifier(authAddr)
	if err != nil {
		return err
	}
	defer verifier.Close()
	mw := pkgauth.NewMiddleware(verifier, 3*time.Second)

	// Routes are first-match. List specific prefixes before /api catch-all.
	routes := []proxy.Route{
		// Auth endpoints (auth-svc owns /api/me, /api/signup, etc.).
		{Prefix: "/api/signup", Target: authURL},
		{Prefix: "/api/login", Target: authURL},
		{Prefix: "/api/refresh", Target: authURL},
		{Prefix: "/api/logout", Target: authURL},
		{Prefix: "/api/forgot-password", Target: authURL},
		{Prefix: "/api/reset-password", Target: authURL},

		// Realtime
		{Prefix: "/api/threads", Target: realtimeURL},
		{Prefix: "/api/ws", Target: realtimeURL},

		// Payment
		{Prefix: "/api/cards", Target: paymentURL},
		{Prefix: "/api/payments", Target: paymentURL},

		// /api/me is owned by auth-svc. /api/me/vendor lives in core; the
		// prefix /api/me/vendor must be listed BEFORE /api/me so it wins.
		{Prefix: "/api/me/vendor", Target: coreURL},
		{Prefix: "/api/me", Target: authURL},

		// Admin endpoints live across auth (users) and core (vendors/stats).
		// Match the user-management subtree to auth before the /api catch-all.
		{Prefix: "/api/admin/users", Target: authURL},

		// Catch-all: vendors, bookings, photos, reviews, notifications,
		// admin endpoints → core.
		{Prefix: "/api", Target: coreURL},
	}
	router, err := proxy.New(routes)
	if err != nil {
		return err
	}

	withAuth := mw.Optional(proxy.WithAuthForwarding(router))

	rl := httpx.NewRateLimiter(rate.Limit(100), 200)
	withLimit := rl.PerIP()(withAuth)
	withRecover := httpx.Recover(log)(withLimit)
	withLog := httpx.AccessLog(log)(withRecover)
	handler := httpx.CORS(cors, withLog)

	srv := &http.Server{Addr: addr, Handler: handler, ReadHeaderTimeout: 5 * time.Second}

	done := make(chan error, 1)
	go func() {
		log.Info("listen", "addr", addr,
			"auth", authURL, "core", coreURL, "payment", paymentURL, "realtime", realtimeURL,
		)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			done <- err
		}
		close(done)
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-sigCh:
		log.Info("shutdown signal")
	case err := <-done:
		if err != nil {
			log.Error("server error", "err", err)
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	return nil
}
