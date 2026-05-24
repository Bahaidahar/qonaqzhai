// Command core runs the core microservice. It exposes:
//   - HTTP API on $CORE_HTTP_ADDR (default :8082) for end-user vendor /
//     booking / review / photo / notification flows.
//   - gRPC CoreService on $CORE_GRPC_ADDR (default :9082) for sibling services.
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

	corev1 "qonaqzhai-backend/gen/proto/core/v1"
	pkgauth "qonaqzhai-backend/pkg/auth"
	"qonaqzhai-backend/pkg/config"
	"qonaqzhai-backend/pkg/grpcutil"
	"qonaqzhai-backend/pkg/logger"

	coregrpc "qonaqzhai-backend/services/core/internal/adapter/grpc"
	"qonaqzhai-backend/services/core/internal/adapter/grpcclient"
	corehttp "qonaqzhai-backend/services/core/internal/adapter/http"
	"qonaqzhai-backend/pkg/idgen"
	"qonaqzhai-backend/services/core/internal/adapter/repo"
	"qonaqzhai-backend/services/core/internal/usecase/admin"
	"qonaqzhai-backend/services/core/internal/usecase/booking"
	"qonaqzhai-backend/services/core/internal/usecase/notification"
	"qonaqzhai-backend/services/core/internal/usecase/photo"
	"qonaqzhai-backend/services/core/internal/usecase/review"
	"qonaqzhai-backend/services/core/internal/usecase/vendor"
)

func main() {
	log := logger.New("core")
	if err := run(log); err != nil {
		log.Error("fatal", "err", err)
		os.Exit(1)
	}
}

func run(log *slog.Logger) error {
	dsn := config.EnvOr("CORE_DATABASE_URL",
		"postgres://qonaqzhai:qonaqzhai@localhost:5433/qonaqzhai_core?sslmode=disable")
	httpAddr := config.EnvOr("CORE_HTTP_ADDR", ":8082")
	grpcAddr := config.EnvOr("CORE_GRPC_ADDR", ":9082")
	authAddr := config.EnvOr("AUTH_GRPC_ADDR", "localhost:9081")
	paymentAddr := os.Getenv("PAYMENT_GRPC_ADDR")
	realtimeAddr := os.Getenv("REALTIME_GRPC_ADDR")

	db, err := repo.Open(dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	id := idgen.New()

	vendors := repo.NewVendorRepo(db, id)
	bookings := repo.NewBookingRepo(db, id)
	photos := repo.NewPhotoRepo(db, id)
	reviews := repo.NewReviewRepo(db, id)
	notifications := repo.NewNotificationRepo(db, id)
	fcmTokens := repo.NewFCMTokenRepo(db, id)

	vendorSvc := vendor.New(vendor.Deps{Vendors: vendors, Reviews: reviews})
	reviewSvc := review.New(review.Deps{Reviews: reviews, Bookings: bookings, Vendors: vendors, Logger: log})
	photoSvc := photo.New(photo.Deps{Photos: photos, Vendors: vendors})
	notifSvc := notification.New(notification.Deps{Notifications: notifications, FCMTokens: fcmTokens})
	adminSvc := admin.New(admin.Deps{Vendors: vendors, Bookings: bookings})

	bookingDeps := booking.Deps{
		Bookings:      bookings,
		Vendors:       vendors,
		Notifications: notifications,
	}
	if paymentAddr != "" {
		pc, err := grpcclient.NewPaymentClient(paymentAddr)
		if err != nil {
			return err
		}
		defer pc.Close()
		bookingDeps.Payments = pc
	}
	if realtimeAddr != "" {
		rc, err := grpcclient.NewRealtimeClient(realtimeAddr)
		if err != nil {
			return err
		}
		defer rc.Close()
		bookingDeps.Realtime = rc
	}
	bookingSvc := booking.New(bookingDeps)

	// Auth verifier (remote). All HTTP middleware delegates to auth-svc.
	verifier, err := pkgauth.NewVerifier(authAddr)
	if err != nil {
		return err
	}
	defer verifier.Close()
	mw := pkgauth.NewMiddleware(verifier, 3*time.Second)

	handler := &corehttp.Handler{
		Vendors: vendorSvc, Bookings: bookingSvc, Reviews: reviewSvc,
		Photos: photoSvc, Notifications: notifSvc, Admin: adminSvc,
	}

	httpSrv := &http.Server{
		Addr: httpAddr, Handler: corehttp.Mux(handler, mw, log),
		ReadHeaderTimeout: 5 * time.Second,
	}

	grpcSrv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcutil.LoggingUnaryInterceptor(log),
			grpcutil.RecoverUnaryInterceptor(log),
		),
	)
	corev1.RegisterCoreServiceServer(grpcSrv, coregrpc.New(vendorSvc, bookingSvc, adminSvc))

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
