// Package app composes all dependencies and produces the runnable http.Handler.
package app

import (
	"context"
	stdhttp "net/http"

	"qonaqzhai-backend/internal/adapter/ai"
	apphttp "qonaqzhai-backend/internal/adapter/http"
	"qonaqzhai-backend/internal/adapter/http/handler"
	"qonaqzhai-backend/internal/adapter/http/middleware"
	"qonaqzhai-backend/internal/adapter/mail"
	"qonaqzhai-backend/internal/adapter/pay"
	"qonaqzhai-backend/internal/adapter/push"
	sqliteadapter "qonaqzhai-backend/internal/adapter/repo/sqlite"
	"qonaqzhai-backend/internal/domain"
	"qonaqzhai-backend/internal/infra/clock"
	"qonaqzhai-backend/internal/infra/config"
	"qonaqzhai-backend/internal/infra/db"
	"qonaqzhai-backend/internal/infra/hasher"
	"qonaqzhai-backend/internal/infra/idgen"
	"qonaqzhai-backend/internal/infra/logger"
	"qonaqzhai-backend/internal/infra/token"
	"qonaqzhai-backend/internal/usecase"
	"qonaqzhai-backend/internal/usecase/admin"
	"qonaqzhai-backend/internal/usecase/auth"
	"qonaqzhai-backend/internal/usecase/booking"
	"qonaqzhai-backend/internal/usecase/chat"
	"qonaqzhai-backend/internal/usecase/notification"
	"qonaqzhai-backend/internal/usecase/payment"
	"qonaqzhai-backend/internal/usecase/review"
	"qonaqzhai-backend/internal/usecase/vendor"
)

// App is the fully wired application.
type App struct {
	Cfg          config.Config
	Handler      stdhttp.Handler
	Notification *notification.Service
	close        func() error
}

// New builds and wires every component.
func New(ctx context.Context, cfg config.Config) (*App, error) {
	conn, err := db.Open(cfg.DBPath)
	if err != nil {
		return nil, err
	}
	log := logger.New("info")
	ids := idgen.New()
	clk := clock.New()
	cost := cfg.BcryptCost
	if cost == 0 {
		cost = 12
	}
	pwHasher := hasher.New(cost)
	jwt := token.New([]byte(cfg.JWTSecret))

	// Repos
	users := sqliteadapter.NewUserRepo(conn, ids)
	vendors := sqliteadapter.NewVendorRepo(conn, ids)
	photos := sqliteadapter.NewPhotoRepo(conn, ids)
	bookings := sqliteadapter.NewBookingRepo(conn, ids)
	reviews := sqliteadapter.NewReviewRepo(conn, ids)
	refreshTokens := sqliteadapter.NewRefreshTokenRepo(conn, ids)
	resetTokens := sqliteadapter.NewPasswordResetRepo(conn, ids)
	notifs := sqliteadapter.NewNotificationRepo(conn, ids)
	fcmTokens := sqliteadapter.NewFCMTokenRepo(conn, ids)
	auditRepo := sqliteadapter.NewAuditRepo(conn, ids)
	services := sqliteadapter.NewServiceRepo(conn, ids)

	// AI (optional)
	var aiClient usecase.AIClient
	if g, err := ai.New(ctx, cfg.GeminiAPIKey, cfg.GeminiModel); err == nil && g != nil {
		aiClient = g
		log.Info("gemini enabled", "model", cfg.GeminiModel)
	} else {
		if err != nil {
			log.Warn("ai init failed, fallback enabled", "err", err.Error())
		} else {
			log.Info("gemini disabled — chat uses keyword fallback")
		}
	}

	// Mailer (optional)
	var mailer usecase.Mailer
	if m := mail.New(mail.Config{
		Host: cfg.SMTPHost, Port: cfg.SMTPPort,
		Username: cfg.SMTPUser, Password: cfg.SMTPPassword,
		From: cfg.SMTPFrom,
	}); m != nil {
		mailer = m
		log.Info("smtp enabled", "host", cfg.SMTPHost)
	}

	// Pusher (optional)
	var pusher usecase.Pusher
	if p, err := push.New(push.Config{
		ProjectID:         cfg.FCMProjectID,
		ServiceAccountKey: []byte(cfg.FCMServiceAccountKey),
	}); err == nil && p != nil {
		pusher = p
		log.Info("fcm enabled", "project", cfg.FCMProjectID)
	} else if err != nil {
		log.Warn("fcm init failed", "err", err.Error())
	}

	// Payment gateway (optional)
	var gateway usecase.PaymentGateway
	if pb := pay.New(pay.Config{
		MerchantID: cfg.PayBoxMerchantID,
		SecretKey:  cfg.PayBoxSecretKey,
		Sandbox:    cfg.PayBoxSandbox,
	}); pb != nil {
		gateway = pb
		log.Info("paybox enabled", "sandbox", cfg.PayBoxSandbox)
	}

	// Notification service
	notifSvc := notification.New(notification.Deps{
		Notifications: notifs,
		Users:         users,
		FCMTokens:     fcmTokens,
		Mailer:        mailer,
		Pusher:        pusher,
		Logger:        log,
	})

	// Use cases
	authSvc := auth.New(auth.Deps{
		Users: users, Refresh: refreshTokens, PasswordResets: resetTokens,
		Hasher: pwHasher, Tokens: jwt, Clock: clk, IDs: ids,
		AccessTTL: cfg.AccessTTL, RefreshTTL: cfg.RefreshTTL, ResetTTL: cfg.ResetTTL,
	})
	if mailer != nil {
		authSvc.SetMailer(mailer)
	}
	vendorSvc := vendor.New(vendor.Deps{Vendors: vendors, Photos: photos, Services: services, Clock: clk})
	bookingSvc := booking.New(booking.Deps{Bookings: bookings, Vendors: vendors, Services: services, Clock: clk, Notifier: notifSvc})
	reviewSvc := review.New(review.Deps{Reviews: reviews, Bookings: bookings, Vendors: vendors, Clock: clk})
	chatSvc := chat.New(chat.Deps{Vendors: vendors, AI: aiClient, Logger: log})
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

	if err := seedAdmin(ctx, users, pwHasher); err != nil {
		return nil, err
	}

	// HTTP
	authMW := middleware.NewAuth(jwt)
	handlers := apphttp.Handlers{
		Auth:         handler.NewAuth(authSvc),
		Me:           handler.NewMe(users),
		Vendor:       handler.NewVendor(vendorSvc),
		Service:      handler.NewService(vendorSvc),
		Booking:      handler.NewBooking(bookingSvc),
		Review:       handler.NewReview(reviewSvc),
		Chat:         handler.NewChat(chatSvc),
		Admin:        handler.NewAdmin(adminSvc),
		Notification: handler.NewNotification(notifSvc, fcmTokens),
	}
	if paymentSvc != nil {
		handlers.Payment = handler.NewPayment(paymentSvc)
	}

	routerCfg := apphttp.RouterConfig{Auth: authMW}
	if !cfg.RateLimitDisabled {
		routerCfg.AuthRate = apphttp.AuthRate()
		routerCfg.ChatRate = apphttp.ChatRate()
	}
	router := apphttp.NewRouter(handlers, routerCfg)

	rec := middleware.Recover(log)
	lg := middleware.Logger(log)
	root := middleware.CORS(cfg.CORSOrigin, rec(lg(router)))

	return &App{
		Cfg:          cfg,
		Handler:      root,
		Notification: notifSvc,
		close: func() error {
			notifSvc.Stop()
			return conn.Close()
		},
	}, nil
}

// Close releases held resources (DB connection + notification worker).
func (a *App) Close() error {
	if a == nil || a.close == nil {
		return nil
	}
	return a.close()
}

// seedAdmin ensures the well-known admin@qonaqzhai.kz account exists.
func seedAdmin(ctx context.Context, users usecase.UserRepo, h usecase.PasswordHasher) error {
	const email = "admin@qonaqzhai.kz"
	if u, err := users.FindByEmail(ctx, email); err == nil && u != nil {
		return nil
	}
	hash, err := h.Hash("admin12345")
	if err != nil {
		return err
	}
	_, err = users.Create(ctx, &domain.User{
		Email:        email,
		Name:         "Admin",
		PasswordHash: hash,
		Role:         domain.RoleAdmin,
		Status:       domain.UserActive,
	})
	if err == domain.ErrAlreadyExists {
		return nil
	}
	return err
}
