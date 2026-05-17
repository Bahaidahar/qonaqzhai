// Binary core-svc serves the core REST API (vendors, bookings, reviews, chat,
// cards, threads, admin, notifications). JWT verification is delegated to
// auth-svc via gRPC; core-svc never reads the JWT secret.
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authv1 "qonaqzhai-backend/gen/proto/auth/v1"
	"qonaqzhai-backend/internal/adapter/ai"
	apphttp "qonaqzhai-backend/internal/adapter/http"
	"qonaqzhai-backend/internal/adapter/http/handler"
	"qonaqzhai-backend/internal/adapter/http/httpx"
	"qonaqzhai-backend/internal/adapter/http/middleware"
	"qonaqzhai-backend/internal/adapter/mail"
	"qonaqzhai-backend/internal/adapter/pay"
	"qonaqzhai-backend/internal/adapter/push"
	sqliteadapter "qonaqzhai-backend/internal/adapter/repo/sqlite"
	"qonaqzhai-backend/internal/infra/clock"
	"qonaqzhai-backend/internal/infra/config"
	"qonaqzhai-backend/internal/infra/db"
	"qonaqzhai-backend/internal/infra/idgen"
	"qonaqzhai-backend/internal/infra/logger"
	"qonaqzhai-backend/internal/infra/token"
	"qonaqzhai-backend/internal/usecase"
	"qonaqzhai-backend/internal/usecase/admin"
	"qonaqzhai-backend/internal/usecase/booking"
	"qonaqzhai-backend/internal/usecase/card"
	"qonaqzhai-backend/internal/usecase/chat"
	"qonaqzhai-backend/internal/usecase/notification"
	"qonaqzhai-backend/internal/usecase/payment"
	"qonaqzhai-backend/internal/usecase/review"
	"qonaqzhai-backend/internal/usecase/thread"
	"qonaqzhai-backend/internal/usecase/vendor"
)

func main() {
	_ = godotenv.Load()
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	httpAddr := envOr("CORE_HTTP_ADDR", ":8082")
	authGRPC := envOr("AUTH_GRPC_TARGET", "localhost:9091")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, err := db.Open(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer conn.Close()

	logr := logger.New("info")
	ids := idgen.New()
	clk := clock.New()

	// gRPC client to auth-svc; remote token verifier instead of local JWT parsing.
	authConn, err := grpc.NewClient(authGRPC, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("dial auth: %v", err)
	}
	defer authConn.Close()
	tokVerifier := token.NewRemote(ctx, authv1.NewAuthServiceClient(authConn))

	// Repos
	users := sqliteadapter.NewUserRepo(conn, ids)
	vendors := sqliteadapter.NewVendorRepo(conn, ids)
	photos := sqliteadapter.NewPhotoRepo(conn, ids)
	bookings := sqliteadapter.NewBookingRepo(conn, ids)
	reviews := sqliteadapter.NewReviewRepo(conn, ids)
	notifs := sqliteadapter.NewNotificationRepo(conn, ids)
	fcmTokens := sqliteadapter.NewFCMTokenRepo(conn, ids)
	auditRepo := sqliteadapter.NewAuditRepo(conn, ids)
	services := sqliteadapter.NewServiceRepo(conn, ids)
	chats := sqliteadapter.NewChatRepo(conn, ids)
	threads := sqliteadapter.NewThreadRepo(conn, ids)
	cards := sqliteadapter.NewCardRepo(conn, ids)

	// External adapters
	var aiClient usecase.AIClient
	if g, err := ai.New(ctx, cfg.GeminiAPIKey, cfg.GeminiModel); err == nil && g != nil {
		aiClient = g
		logr.Info("gemini enabled", "model", cfg.GeminiModel)
	}
	var mailer usecase.Mailer
	if m := mail.New(mail.Config{
		Host: cfg.SMTPHost, Port: cfg.SMTPPort,
		Username: cfg.SMTPUser, Password: cfg.SMTPPassword,
		From: cfg.SMTPFrom,
	}); m != nil {
		mailer = m
	}
	var pusher usecase.Pusher
	if p, err := push.New(push.Config{
		ProjectID:         cfg.FCMProjectID,
		ServiceAccountKey: []byte(cfg.FCMServiceAccountKey),
	}); err == nil && p != nil {
		pusher = p
	}
	var gateway usecase.PaymentGateway
	if pb := pay.New(pay.Config{
		MerchantID: cfg.PayBoxMerchantID,
		SecretKey:  cfg.PayBoxSecretKey,
		Sandbox:    cfg.PayBoxSandbox,
	}); pb != nil {
		gateway = pb
	}

	// Use cases
	notifSvc := notification.New(notification.Deps{
		Notifications: notifs, Users: users, FCMTokens: fcmTokens,
		Mailer: mailer, Pusher: pusher, Logger: logr,
	})
	vendorSvc := vendor.New(vendor.Deps{Vendors: vendors, Photos: photos, Services: services, Clock: clk})
	threadSvc := thread.New(thread.Deps{
		Threads: threads, Bookings: bookings, Vendors: vendors, Users: users, Notifier: notifSvc,
		// Publisher is nil — realtime-svc fans out, not core.
	})
	bookingSvc := booking.New(booking.Deps{
		Bookings: bookings, Vendors: vendors, Services: services,
		Clock: clk, Notifier: notifSvc, Threads: threadSvc,
	})
	cardSvc := card.New(card.Deps{Cards: cards})
	reviewSvc := review.New(review.Deps{Reviews: reviews, Bookings: bookings, Vendors: vendors, Clock: clk})
	chatSvc := chat.New(chat.Deps{Vendors: vendors, Chats: chats, AI: aiClient, Logger: logr})
	adminSvc := admin.New(admin.Deps{
		Users: users, Vendors: vendors, Bookings: bookings, Reviews: reviews,
		Audit: auditRepo, Notifier: notifSvc,
	})
	var paymentSvc *payment.Service
	if gateway != nil {
		paymentSvc = payment.New(payment.Deps{
			Bookings: bookings, Vendors: vendors, Users: users,
			Gateway: gateway, Notifier: notifSvc, BaseURL: cfg.BaseURL,
		})
	}

	authMW := middleware.NewAuth(tokVerifier)
	handlers := apphttp.Handlers{
		Vendor:       handler.NewVendor(vendorSvc),
		Service:      handler.NewService(vendorSvc),
		Booking:      handler.NewBooking(bookingSvc, cardSvc),
		Review:       handler.NewReview(reviewSvc),
		Chat:         handler.NewChat(chatSvc),
		Admin:        handler.NewAdmin(adminSvc),
		Notification: handler.NewNotification(notifSvc, fcmTokens),
		Thread:       handler.NewThread(threadSvc),
		Card:         handler.NewCard(cardSvc),
	}
	if paymentSvc != nil {
		handlers.Payment = handler.NewPayment(paymentSvc)
	}

	routerCfg := apphttp.RouterConfig{Auth: authMW}
	if !cfg.RateLimitDisabled {
		routerCfg.ChatRate = apphttp.ChatRate()
	}
	router := apphttp.NewRouter(handlers, routerCfg)

	// Internal endpoint exposed only to realtime-svc (gateway never routes it).
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, _ *http.Request) {
		httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "core"})
	})
	mux.HandleFunc("POST /internal/threads/{id}/messages", internalThreadSend(threadSvc))
	mux.Handle("/", router)

	rec := middleware.Recover(logr)
	lg := middleware.Logger(logr)
	// CORS is applied at the gateway edge; core-svc never receives cross-origin
	// requests directly in production.
	root := rec(lg(mux))

	srv := &http.Server{
		Addr:              httpAddr,
		Handler:           root,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		log.Printf("core-svc HTTP listening on %s · auth=%s", httpAddr, authGRPC)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("http: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("core-svc shutting down")
	shutCtx, cancelShut := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShut()
	_ = srv.Shutdown(shutCtx)
	notifSvc.Stop()
}

// internalThreadSend persists a message coming from realtime-svc and returns
// the persisted message PLUS thread participants so realtime can fan out to
// both customer and vendor. Payload: {"senderId":"u1","text":"hi"}.
// No JWT — relies on network isolation (gateway never forwards /internal/*).
func internalThreadSend(svc *thread.Service) http.HandlerFunc {
	type req struct {
		SenderID string `json:"senderId"`
		Text     string `json:"text"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var body req
		if err := httpx.ReadJSON(r, &body); err != nil {
			httpx.WriteError(w, http.StatusBadRequest, "invalid json")
			return
		}
		id := r.PathValue("id")
		msg, err := svc.Send(r.Context(), body.SenderID, id, body.Text)
		if err != nil {
			httpx.HandleError(w, err)
			return
		}
		// Look up participants by reading the thread directly (sender already authorized via Send).
		t, _, err := svc.Get(r.Context(), body.SenderID, id)
		if err != nil {
			httpx.HandleError(w, err)
			return
		}
		httpx.WriteJSON(w, http.StatusCreated, map[string]any{
			"message":      msg,
			"participants": []string{t.CustomerID, t.VendorID},
		})
	}
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
