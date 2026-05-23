// Command realtime runs the realtime microservice. It exposes:
//   - HTTP API on $REALTIME_HTTP_ADDR (default :8084) for threads + ws
//   - gRPC RealtimeService on $REALTIME_GRPC_ADDR (default :9084)
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

	realtimev1 "qonaqzhai-backend/gen/proto/realtime/v1"
	pkgauth "qonaqzhai-backend/pkg/auth"
	"qonaqzhai-backend/pkg/config"
	"qonaqzhai-backend/pkg/grpcutil"
	"qonaqzhai-backend/pkg/logger"

	"qonaqzhai-backend/services/realtime/internal/adapter/clock"
	realtimegrpc "qonaqzhai-backend/services/realtime/internal/adapter/grpc"
	"qonaqzhai-backend/services/realtime/internal/adapter/grpcclient"
	realtimehttp "qonaqzhai-backend/services/realtime/internal/adapter/http"
	"qonaqzhai-backend/services/realtime/internal/adapter/idgen"
	"qonaqzhai-backend/services/realtime/internal/adapter/repo"
	"qonaqzhai-backend/services/realtime/internal/adapter/ws"
	"qonaqzhai-backend/services/realtime/internal/ports"
	"qonaqzhai-backend/services/realtime/internal/usecase/thread"
)

func main() {
	log := logger.New("realtime")
	if err := run(log); err != nil {
		log.Error("fatal", "err", err)
		os.Exit(1)
	}
}

func run(log *slog.Logger) error {
	dsn := config.EnvOr("REALTIME_DATABASE_URL",
		"postgres://qonaqzhai:qonaqzhai@localhost:5433/qonaqzhai_realtime?sslmode=disable")
	httpAddr := config.EnvOr("REALTIME_HTTP_ADDR", ":8084")
	grpcAddr := config.EnvOr("REALTIME_GRPC_ADDR", ":9084")
	cors := config.EnvOr("CORS_ORIGIN", "*")
	authAddr := config.EnvOr("AUTH_GRPC_ADDR", "localhost:9081")

	db, err := repo.Open(dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	id := idgen.New()
	_ = clock.New()
	threads := repo.NewThreadRepo(db, id)

	hub := ws.NewHub(log)

	var auth ports.AuthClient
	if authAddr != "" {
		ac, err := grpcclient.NewAuthClient(authAddr)
		if err != nil {
			return err
		}
		defer ac.Close()
		auth = ac
	}

	threadSvc := thread.New(thread.Deps{Threads: threads, Auth: auth, Publisher: hub})
	hub.OnIncoming = func(userID string, m ws.IncomingMessage) {
		if m.ThreadID == "" || m.Text == "" {
			return
		}
		_, _ = threadSvc.Send(context.Background(), userID, m.ThreadID, m.Text)
	}

	verifier, err := pkgauth.NewVerifier(authAddr)
	if err != nil {
		return err
	}
	defer verifier.Close()
	mw := pkgauth.NewMiddleware(verifier, 3*time.Second)

	handler := &realtimehttp.Handler{Threads: threadSvc, Hub: hub}

	httpSrv := &http.Server{
		Addr: httpAddr, Handler: realtimehttp.Mux(handler, mw, cors, log),
		ReadHeaderTimeout: 5 * time.Second,
	}

	grpcSrv := grpc.NewServer(
		grpc.UnaryInterceptor(grpcutil.LoggingUnaryInterceptor(log)),
		grpc.ChainUnaryInterceptor(grpcutil.RecoverUnaryInterceptor(log)),
	)
	realtimev1.RegisterRealtimeServiceServer(grpcSrv, realtimegrpc.New(threadSvc, hub))

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
			log.Error("http error", "err", err)
		}
	case err := <-grpcDone:
		if err != nil {
			log.Error("grpc error", "err", err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = httpSrv.Shutdown(ctx)
	grpcSrv.GracefulStop()
	return nil
}
