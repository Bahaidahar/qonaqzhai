// Package http wires HTTP handlers into a single http.Handler.
package http

import (
	"net/http"
	"time"

	"golang.org/x/time/rate"

	"qonaqzhai-backend/internal/adapter/http/handler"
	"qonaqzhai-backend/internal/adapter/http/httpx"
	"qonaqzhai-backend/internal/adapter/http/middleware"
	"qonaqzhai-backend/internal/domain"
)

// Handlers aggregates the per-feature HTTP handlers.
type Handlers struct {
	Auth         *handler.Auth
	Me           *handler.Me
	Vendor       *handler.Vendor
	Service      *handler.Service
	Booking      *handler.Booking
	Review       *handler.Review
	Chat         *handler.Chat
	Admin        *handler.Admin
	Payment      *handler.Payment      // optional — nil disables routes
	Notification *handler.Notification // optional — nil disables routes
	Thread       *handler.Thread
}

// RouterConfig bundles router-level dependencies.
type RouterConfig struct {
	Auth     *middleware.Auth
	AuthRate *middleware.RateLimiter // per-IP, applied to /api/login + /api/signup
	ChatRate *middleware.RateLimiter // per-user, applied to /api/chat
}

// NewRouter builds the application HTTP mux.
func NewRouter(h Handlers, cfg RouterConfig) http.Handler {
	mux := http.NewServeMux()

	authLimit := func(next http.Handler) http.Handler { return next }
	if cfg.AuthRate != nil {
		authLimit = cfg.AuthRate.PerIP()
	}
	chatLimit := func(next http.Handler) http.Handler { return next }
	if cfg.ChatRate != nil {
		chatLimit = cfg.ChatRate.PerUser()
	}

	// public
	mux.HandleFunc("GET /api/health", health)
	mux.HandleFunc("GET /api/docs", handler.SwaggerUI)
	mux.HandleFunc("GET /api/docs/openapi.yaml", handler.OpenAPI)
	mux.Handle("POST /api/signup", authLimit(http.HandlerFunc(h.Auth.Signup)))
	mux.Handle("POST /api/login", authLimit(http.HandlerFunc(h.Auth.Login)))
	mux.HandleFunc("POST /api/auth/refresh", h.Auth.Refresh)
	mux.HandleFunc("POST /api/auth/logout", h.Auth.Logout)
	mux.Handle("POST /api/auth/forgot-password", authLimit(http.HandlerFunc(h.Auth.ForgotPassword)))
	mux.HandleFunc("POST /api/auth/reset-password", h.Auth.ResetPassword)
	mux.Handle("GET /api/vendors", cfg.Auth.Optional(http.HandlerFunc(h.Vendor.List)))
	mux.Handle("GET /api/vendors/{id}", cfg.Auth.Optional(http.HandlerFunc(h.Vendor.Detail)))
	mux.HandleFunc("GET /api/vendors/{id}/reviews", h.Review.ListByVendor)
	mux.Handle("GET /api/vendors/{id}/services", cfg.Auth.Optional(http.HandlerFunc(h.Service.List)))
	mux.HandleFunc("GET /api/photos/{id}", h.Vendor.ServePhoto)

	// authenticated (any role)
	mux.Handle("GET /api/me", cfg.Auth.Required(http.HandlerFunc(h.Me.Get)))
	mux.Handle("POST /api/chat", cfg.Auth.Required(chatLimit(http.HandlerFunc(h.Chat.Generate))))
	mux.Handle("GET /api/chats", cfg.Auth.Required(http.HandlerFunc(h.Chat.ListChats)))
	mux.Handle("GET /api/chats/{id}", cfg.Auth.Required(http.HandlerFunc(h.Chat.GetChat)))
	mux.Handle("PATCH /api/chats/{id}", cfg.Auth.Required(http.HandlerFunc(h.Chat.RenameChat)))
	mux.Handle("DELETE /api/chats/{id}", cfg.Auth.Required(http.HandlerFunc(h.Chat.DeleteChat)))
	mux.Handle("GET /api/bookings", cfg.Auth.Required(http.HandlerFunc(h.Booking.List)))
	mux.Handle("PATCH /api/bookings/{id}", cfg.Auth.Required(http.HandlerFunc(h.Booking.Update)))
	mux.Handle("GET /api/threads", cfg.Auth.Required(http.HandlerFunc(h.Thread.List)))
	mux.Handle("GET /api/threads/{id}", cfg.Auth.Required(http.HandlerFunc(h.Thread.Get)))
	mux.Handle("POST /api/threads/{id}/messages", cfg.Auth.Required(http.HandlerFunc(h.Thread.Send)))

	// vendor
	vendorOnly := cfg.Auth.RequireRole(domain.RoleVendor)
	mux.Handle("POST /api/vendor", vendorOnly(http.HandlerFunc(h.Vendor.Upsert)))
	mux.Handle("GET /api/vendor", vendorOnly(http.HandlerFunc(h.Vendor.My)))
	mux.Handle("POST /api/vendor/photos", vendorOnly(http.HandlerFunc(h.Vendor.UploadPhoto)))
	mux.Handle("DELETE /api/vendor/photos/{id}", vendorOnly(http.HandlerFunc(h.Vendor.DeletePhoto)))
	mux.Handle("GET /api/vendor/services", vendorOnly(http.HandlerFunc(h.Service.MyList)))
	mux.Handle("POST /api/vendor/services", vendorOnly(http.HandlerFunc(h.Service.Create)))
	mux.Handle("PATCH /api/vendor/services/{id}", vendorOnly(http.HandlerFunc(h.Service.Update)))
	mux.Handle("DELETE /api/vendor/services/{id}", vendorOnly(http.HandlerFunc(h.Service.Delete)))

	// customer
	customerOnly := cfg.Auth.RequireRole(domain.RoleCustomer)
	mux.Handle("POST /api/bookings", customerOnly(http.HandlerFunc(h.Booking.Create)))
	mux.Handle("POST /api/reviews", customerOnly(http.HandlerFunc(h.Review.Submit)))

	// admin
	adminOnly := cfg.Auth.RequireRole(domain.RoleAdmin)
	mux.Handle("GET /api/admin/users", adminOnly(http.HandlerFunc(h.Admin.Users)))
	mux.Handle("PATCH /api/admin/users/{id}", adminOnly(http.HandlerFunc(h.Admin.UserStatus)))
	mux.Handle("PATCH /api/admin/vendors/{id}", adminOnly(http.HandlerFunc(h.Admin.VendorStatus)))
	mux.Handle("GET /api/admin/stats", adminOnly(http.HandlerFunc(h.Admin.Stats)))
	mux.Handle("GET /api/admin/stats/bookings", adminOnly(http.HandlerFunc(h.Admin.BookingsTimeseries)))
	mux.Handle("GET /api/admin/stats/categories", adminOnly(http.HandlerFunc(h.Admin.TopCategories)))
	mux.Handle("GET /api/admin/stats/funnel", adminOnly(http.HandlerFunc(h.Admin.ApprovalFunnel)))
	mux.Handle("GET /api/admin/audit", adminOnly(http.HandlerFunc(h.Admin.AuditLog)))
	mux.Handle("DELETE /api/admin/reviews/{id}", adminOnly(http.HandlerFunc(h.Review.AdminDelete)))

	// payments (optional — feature-flagged)
	if h.Payment != nil {
		mux.Handle("POST /api/bookings/{id}/pay", customerOnly(http.HandlerFunc(h.Payment.Start)))
		mux.HandleFunc("POST /api/webhooks/paybox", h.Payment.Callback)
	}

	// notifications (optional)
	if h.Notification != nil {
		mux.Handle("GET /api/notifications", cfg.Auth.Required(http.HandlerFunc(h.Notification.Inbox)))
		mux.Handle("POST /api/notifications/tokens", cfg.Auth.Required(http.HandlerFunc(h.Notification.RegisterToken)))
		mux.Handle("DELETE /api/notifications/tokens", cfg.Auth.Required(http.HandlerFunc(h.Notification.UnregisterToken)))
	}

	return mux
}

func health(w http.ResponseWriter, _ *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// AuthRate exposes constructor for auth-endpoint rate limit (10 req/min, burst 5).
func AuthRate() *middleware.RateLimiter {
	return middleware.NewRateLimiter(rate.Every(6*time.Second), 5) // 10/min ≈ every 6s
}

// ChatRate exposes constructor for chat rate limit (30 req/min, burst 5).
func ChatRate() *middleware.RateLimiter {
	return middleware.NewRateLimiter(rate.Every(2*time.Second), 5)
}
