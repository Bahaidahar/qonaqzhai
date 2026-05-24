// Command payment runs the payment microservice. It exposes:
//   - HTTP API on $PAYMENT_HTTP_ADDR (default :8083) for cards + payment list
//   - gRPC PaymentService on $PAYMENT_GRPC_ADDR (default :9083) for core-svc
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

	paymentv1 "qonaqzhai-backend/gen/proto/payment/v1"
	pkgauth "qonaqzhai-backend/pkg/auth"
	"qonaqzhai-backend/pkg/config"
	"qonaqzhai-backend/pkg/grpcutil"
	"qonaqzhai-backend/pkg/logger"

	"qonaqzhai-backend/pkg/clock"
	"qonaqzhai-backend/services/payment/internal/adapter/gateway"
	paymentgrpc "qonaqzhai-backend/services/payment/internal/adapter/grpc"
	"qonaqzhai-backend/services/payment/internal/adapter/grpcclient"
	paymenthttp "qonaqzhai-backend/services/payment/internal/adapter/http"
	"qonaqzhai-backend/pkg/idgen"
	"qonaqzhai-backend/services/payment/internal/adapter/repo"
	"qonaqzhai-backend/services/payment/internal/ports"
	"qonaqzhai-backend/services/payment/internal/usecase/card"
	"qonaqzhai-backend/services/payment/internal/usecase/payment"
)

func main() {
	log := logger.New("payment")
	if err := run(log); err != nil {
		log.Error("fatal", "err", err)
		os.Exit(1)
	}
}

func run(log *slog.Logger) error {
	dsn := config.EnvOr("PAYMENT_DATABASE_URL",
		"postgres://qonaqzhai:qonaqzhai@localhost:5433/qonaqzhai_payment?sslmode=disable")
	httpAddr := config.EnvOr("PAYMENT_HTTP_ADDR", ":8083")
	grpcAddr := config.EnvOr("PAYMENT_GRPC_ADDR", ":9083")
	authAddr := config.EnvOr("AUTH_GRPC_ADDR", "localhost:9081")
	coreAddr := os.Getenv("CORE_GRPC_ADDR")

	db, err := repo.Open(dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	id := idgen.New()
	clk := clock.New()
	cards := repo.NewCardRepo(db, id)
	payments := repo.NewPaymentRepo(db, id)

	var gw ports.Gateway = gateway.NewMock()
	if pb := gateway.NewPayBox(gateway.PayBoxConfig{
		MerchantID: os.Getenv("PAYBOX_MERCHANT_ID"),
		SecretKey:  os.Getenv("PAYBOX_SECRET_KEY"),
		Sandbox:    config.BoolEnv("PAYBOX_SANDBOX", true),
	}); pb != nil {
		gw = pb
	}

	var core ports.CoreClient
	if coreAddr != "" {
		cc, err := grpcclient.NewCoreClient(coreAddr)
		if err != nil {
			return err
		}
		defer cc.Close()
		core = cc
	}

	cardSvc := card.New(card.Deps{Cards: cards, Clock: clk})
	paymentSvc := payment.New(payment.Deps{
		Payments: payments, Cards: cards, Gateway: gw, Core: core, Clock: clk,
	})

	verifier, err := pkgauth.NewVerifier(authAddr)
	if err != nil {
		return err
	}
	defer verifier.Close()
	mw := pkgauth.NewMiddleware(verifier, 3*time.Second)

	handler := &paymenthttp.Handler{Cards: cardSvc, Payments: paymentSvc}

	httpSrv := &http.Server{
		Addr: httpAddr, Handler: paymenthttp.Mux(handler, mw, log),
		ReadHeaderTimeout: 5 * time.Second,
	}

	grpcSrv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcutil.LoggingUnaryInterceptor(log),
			grpcutil.RecoverUnaryInterceptor(log),
		),
	)
	paymentv1.RegisterPaymentServiceServer(grpcSrv, paymentgrpc.New(paymentSvc, cardSvc))

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
