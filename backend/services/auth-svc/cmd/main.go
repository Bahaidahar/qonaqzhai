// Binary auth-svc serves auth REST endpoints + AuthService gRPC contract.
package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	authv1 "qonaqzhai-backend/gen/proto/auth/v1"
	"qonaqzhai-backend/internal/adapter/http/handler"
	"qonaqzhai-backend/internal/adapter/http/httpx"
	"qonaqzhai-backend/internal/adapter/http/middleware"
	"qonaqzhai-backend/internal/adapter/mail"
	sqliteadapter "qonaqzhai-backend/internal/adapter/repo/sqlite"
	"qonaqzhai-backend/internal/infra/clock"
	"qonaqzhai-backend/internal/infra/config"
	"qonaqzhai-backend/internal/infra/db"
	"qonaqzhai-backend/internal/infra/hasher"
	"qonaqzhai-backend/internal/infra/idgen"
	"qonaqzhai-backend/internal/infra/token"
	"qonaqzhai-backend/internal/usecase/auth"

	grpcserver "qonaqzhai-backend/services/auth-svc/internal/grpcserver"
)

func main() {
	_ = godotenv.Load()
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	httpAddr := envOr("AUTH_HTTP_ADDR", ":8081")
	grpcAddr := envOr("AUTH_GRPC_ADDR", ":9091")

	conn, err := db.Open(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer conn.Close()

	ids := idgen.New()
	clk := clock.New()
	cost := cfg.BcryptCost
	if cost == 0 {
		cost = 12
	}
	pwHasher := hasher.New(cost)
	jwt := token.New([]byte(cfg.JWTSecret))

	users := sqliteadapter.NewUserRepo(conn, ids)
	refreshTokens := sqliteadapter.NewRefreshTokenRepo(conn, ids)
	resetTokens := sqliteadapter.NewPasswordResetRepo(conn, ids)

	authSvc := auth.New(auth.Deps{
		Users:          users,
		Refresh:        refreshTokens,
		PasswordResets: resetTokens,
		Hasher:         pwHasher,
		Tokens:         jwt,
		Clock:          clk,
		IDs:            ids,
		AccessTTL:      cfg.AccessTTL,
		RefreshTTL:     cfg.RefreshTTL,
		ResetTTL:       cfg.ResetTTL,
	})
	if m := mail.New(mail.Config{
		Host: cfg.SMTPHost, Port: cfg.SMTPPort,
		Username: cfg.SMTPUser, Password: cfg.SMTPPassword,
		From: cfg.SMTPFrom,
	}); m != nil {
		authSvc.SetMailer(m)
	}

	authHandler := handler.NewAuth(authSvc)
	authMW := middleware.NewAuth(jwt)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, _ *http.Request) {
		httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "auth"})
	})
	mux.HandleFunc("POST /api/signup", authHandler.Signup)
	mux.HandleFunc("POST /api/login", authHandler.Login)
	mux.HandleFunc("POST /api/auth/refresh", authHandler.Refresh)
	mux.HandleFunc("POST /api/auth/logout", authHandler.Logout)
	mux.HandleFunc("POST /api/auth/forgot-password", authHandler.ForgotPassword)
	mux.HandleFunc("POST /api/auth/reset-password", authHandler.ResetPassword)
	// /api/me is auth-owned (user identity)
	meHandler := handler.NewMe(users)
	mux.Handle("GET /api/me", authMW.Required(http.HandlerFunc(meHandler.Get)))

	httpSrv := &http.Server{
		Addr:              httpAddr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// gRPC server — internal-only AuthService.
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("grpc listen: %v", err)
	}
	gs := grpc.NewServer()
	authv1.RegisterAuthServiceServer(gs, grpcserver.New(jwt))

	go func() {
		log.Printf("auth-svc HTTP listening on %s", httpAddr)
		if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("http: %v", err)
		}
	}()
	go func() {
		log.Printf("auth-svc gRPC listening on %s", grpcAddr)
		if err := gs.Serve(lis); err != nil {
			log.Fatalf("grpc: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("auth-svc shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = httpSrv.Shutdown(ctx)
	gs.GracefulStop()
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
