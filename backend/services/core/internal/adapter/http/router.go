package http

import (
	"log/slog"
	"net/http"

	"golang.org/x/time/rate"

	pkgauth "qonaqzhai-backend/pkg/auth"
	"qonaqzhai-backend/pkg/httpx"
)

// Mux wires every core HTTP route. mw verifies tokens via auth-svc gRPC.
func Mux(h *Handler, mw *pkgauth.Middleware, corsOrigin string, log *slog.Logger) http.Handler {
	rl := httpx.NewRateLimiter(rate.Limit(30), 60)
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/health", h.Health)

	// Public vendor catalog + photo serve.
	mux.HandleFunc("GET /api/vendors", h.SearchVendors)
	mux.HandleFunc("GET /api/vendors/{id}", h.GetVendor)
	mux.HandleFunc("GET /api/photos/{id}", h.ServePhoto)
	mux.HandleFunc("GET /api/vendors/{vendorId}/reviews", h.ListReviews)

	// Authenticated user endpoints.
	mux.Handle("GET /api/me/vendor", mw.Required(http.HandlerFunc(h.MyVendor)))
	mux.Handle("PUT /api/me/vendor", mw.Required(http.HandlerFunc(h.UpsertVendor)))
	mux.Handle("POST /api/me/vendor", mw.Required(http.HandlerFunc(h.UpsertVendor)))
	mux.Handle("POST /api/me/vendor/photos", mw.Required(http.HandlerFunc(h.UploadPhoto)))
	mux.Handle("DELETE /api/me/vendor/photos/{id}", mw.Required(http.HandlerFunc(h.DeletePhoto)))

	mux.Handle("POST /api/bookings", mw.Required(http.HandlerFunc(h.CreateBooking)))
	mux.Handle("GET /api/bookings", mw.Required(http.HandlerFunc(h.ListMyBookings)))
	mux.Handle("GET /api/bookings/{id}", mw.Required(http.HandlerFunc(h.GetBooking)))
	mux.Handle("PATCH /api/bookings/{id}", mw.Required(http.HandlerFunc(h.BookingTransition)))
	mux.Handle("POST /api/bookings/{id}/pay", mw.Required(http.HandlerFunc(h.PayBooking)))

	mux.Handle("POST /api/reviews", mw.Required(http.HandlerFunc(h.SubmitReview)))

	mux.Handle("GET /api/notifications", mw.Required(http.HandlerFunc(h.ListNotifications)))
	mux.Handle("POST /api/notifications/fcm", mw.Required(http.HandlerFunc(h.RegisterFCM)))

	// Admin endpoints.
	mux.Handle("PATCH /api/admin/vendors/{id}/status", mw.RequireRole("admin")(http.HandlerFunc(h.AdminSetVendorStatus)))
	mux.Handle("GET /api/admin/stats", mw.RequireRole("admin")(http.HandlerFunc(h.AdminStats)))

	withLimit := rl.PerIP()(mux)
	withRecover := httpx.Recover(log)(withLimit)
	withLog := httpx.AccessLog(log)(withRecover)
	return httpx.CORS(corsOrigin, withLog)
}
