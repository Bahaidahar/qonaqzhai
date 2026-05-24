// Command auth runs the auth microservice. It exposes:
//   - HTTP API on $ADDR_HTTP (default :8081) for end-user auth flows
//   - gRPC AuthService on $ADDR_GRPC (default :9081) for other services
//
// auth owns the JWT secret and the users / refresh_tokens / password_reset_tokens
// tables. No other service connects to its database directly.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"

	authv1 "qonaqzhai-backend/gen/proto/auth/v1"
	pkgauth "qonaqzhai-backend/pkg/auth"
	"qonaqzhai-backend/pkg/config"
	"qonaqzhai-backend/pkg/grpcutil"
	"qonaqzhai-backend/pkg/logger"

	"qonaqzhai-backend/pkg/clock"
	authgrpc "qonaqzhai-backend/services/auth/internal/adapter/grpc"
	"qonaqzhai-backend/services/auth/internal/adapter/hasher"
	authhttp "qonaqzhai-backend/services/auth/internal/adapter/http"
	"qonaqzhai-backend/pkg/idgen"
	"qonaqzhai-backend/services/auth/internal/adapter/mail"
	"qonaqzhai-backend/services/auth/internal/adapter/repo"
	"qonaqzhai-backend/services/auth/internal/usecase"
)

func main() {
	log := logger.New("auth")
	if err := run(log); err != nil {
		log.Error("fatal", "err", err)
		os.Exit(1)
	}
}

func run(log *slog.Logger) error {
	dsn := config.EnvOr("AUTH_DATABASE_URL",
		"postgres://qonaqzhai:qonaqzhai@localhost:5433/qonaqzhai_auth?sslmode=disable")
	httpAddr := config.EnvOr("AUTH_HTTP_ADDR", ":8081")
	grpcAddr := config.EnvOr("AUTH_GRPC_ADDR", ":9081")

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		s, err := config.RandomHex(32)
		if err != nil {
			return err
		}
		jwtSecret = s
		log.Warn("JWT_SECRET unset — generated ephemeral secret; tokens will not survive restart")
	}

	db, err := repo.Open(dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	id := idgen.New()
	clk := clock.New()
	users := repo.NewUserRepo(db, id)
	refresh := repo.NewRefreshTokenRepo(db, id)
	resets := repo.NewPasswordResetRepo(db, id)
	hash := hasher.New(0)
	signer := pkgauth.NewJWTSigner([]byte(jwtSecret), "qonaqzhai")

	svc := usecase.New(usecase.Deps{
		Users: users, Refresh: refresh, PasswordResets: resets,
		Hasher: hash, Signer: signer, Clock: clk, IDs: id,
		AccessTTL:  config.DurationEnv("JWT_ACCESS_TTL", 15*time.Minute),
		RefreshTTL: config.DurationEnv("JWT_REFRESH_TTL", 30*24*time.Hour),
		ResetTTL:   config.DurationEnv("PASSWORD_RESET_TTL", time.Hour),
	})

	if mailer := mail.New(mail.Config{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     config.EnvOr("SMTP_PORT", "587"),
		Username: os.Getenv("SMTP_USER"),
		Password: os.Getenv("SMTP_PASSWORD"),
		From:     config.FirstNonEmpty(os.Getenv("SMTP_FROM"), os.Getenv("SMTP_USER")),
	}); mailer != nil {
		svc.SetMailer(mailer)
	}

	seedAdmin(context.Background(), log, svc)

	// HTTP server. Auth verifies its own JWTs locally — wrap pkgauth.Middleware
	// with a local TokenVerifier instead of dialing gRPC into ourselves.
	mw := pkgauth.NewMiddleware(localVerifier{signer: signer}, 2*time.Second)
	httpSrv := &http.Server{
		Addr:              httpAddr,
		Handler:           authhttp.Mux(authhttp.NewHandler(svc), mw, log),
		ReadHeaderTimeout: 5 * time.Second,
	}

	// gRPC server.
	grpcSrv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcutil.LoggingUnaryInterceptor(log),
			grpcutil.RecoverUnaryInterceptor(log),
		),
	)
	authv1.RegisterAuthServiceServer(grpcSrv, authgrpc.New(svc))

	httpDone := make(chan error, 1)
	grpcDone := make(chan error, 1)

	go func() {
		log.Info("http listen", "addr", httpAddr)
		if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			httpDone <- err
		}
		close(httpDone)
	}()

	go func() {
		lis, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			grpcDone <- err
			close(grpcDone)
			return
		}
		log.Info("grpc listen", "addr", grpcAddr)
		if err := grpcSrv.Serve(lis); err != nil {
			grpcDone <- err
		}
		close(grpcDone)
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-sigCh:
		log.Info("shutdown signal")
	case err := <-httpDone:
		if err != nil {
			log.Error("http server error", "err", err)
		}
	case err := <-grpcDone:
		if err != nil {
			log.Error("grpc server error", "err", err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = httpSrv.Shutdown(ctx)
	grpcSrv.GracefulStop()
	return nil
}

// localVerifier lets the auth-svc verify its own JWTs without dialing itself.
type localVerifier struct{ signer *pkgauth.JWTSigner }

func (v localVerifier) Verify(_ context.Context, token string) (pkgauth.Claims, error) {
	c, _, err := v.signer.Parse(token)
	return c, err
}

func seedAdmin(ctx context.Context, log *slog.Logger, svc *usecase.Service) {
	email := os.Getenv("ADMIN_EMAIL")
	password := os.Getenv("ADMIN_PASSWORD")
	if email == "" || password == "" {
		log.Info("admin seed skipped — set ADMIN_EMAIL and ADMIN_PASSWORD to enable")
		return
	}
	u, err := svc.EnsureAdmin(ctx, email, password, "Admin")
	if err != nil {
		log.Warn("admin seed failed", "err", err)
		return
	}
	log.Info("admin ready", "id", u.ID, "email", u.Email, "role", u.Role)
}
