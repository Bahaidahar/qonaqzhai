// Package http exposes core service operations to clients.
package http

import (
	"io"
	"net/http"
	"strconv"
	"strings"

	pkgauth "qonaqzhai-backend/pkg/auth"
	"qonaqzhai-backend/pkg/httpx"

	"qonaqzhai-backend/services/core/internal/domain"
	"qonaqzhai-backend/services/core/internal/ports"
	"qonaqzhai-backend/services/core/internal/usecase/admin"
	"qonaqzhai-backend/services/core/internal/usecase/booking"
	"qonaqzhai-backend/services/core/internal/usecase/notification"
	"qonaqzhai-backend/services/core/internal/usecase/photo"
	"qonaqzhai-backend/services/core/internal/usecase/review"
	"qonaqzhai-backend/services/core/internal/usecase/vendor"
)

// Handler bundles all core HTTP handlers.
type Handler struct {
	Vendors       *vendor.Service
	Bookings      *booking.Service
	Reviews       *review.Service
	Photos        *photo.Service
	Notifications *notification.Service
	Admin         *admin.Service
}

func (h *Handler) Health(w http.ResponseWriter, _ *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok", "svc": "core"})
}

// --- vendors ---

func (h *Handler) UpsertVendor(w http.ResponseWriter, r *http.Request) {
	uid, _ := pkgauth.UserIDFrom(r.Context())
	var req struct {
		Name, Category, City, Description string
		PriceFrom                         int64
	}
	if err := httpx.ReadJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	v, err := h.Vendors.Submit(r.Context(), uid, domain.VendorInput{
		Name: req.Name, Category: req.Category, City: req.City,
		Description: req.Description, PriceFrom: req.PriceFrom,
	})
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, v)
}

func (h *Handler) GetVendor(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	v, err := h.Vendors.FindPublic(r.Context(), id)
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, v)
}

func (h *Handler) MyVendor(w http.ResponseWriter, r *http.Request) {
	uid, _ := pkgauth.UserIDFrom(r.Context())
	v, err := h.Vendors.FindByUserID(r.Context(), uid)
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, v)
}

func (h *Handler) SearchVendors(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	atoi64 := func(s string) int64 { n, _ := strconv.ParseInt(s, 10, 64); return n }
	atoi := func(s string) int { n, _ := strconv.Atoi(s); return n }
	atof := func(s string) float64 { n, _ := strconv.ParseFloat(s, 64); return n }
	vs, total, err := h.Vendors.Search(r.Context(), ports.VendorQuery{
		Q:         q.Get("q"),
		Category:  q.Get("category"),
		City:      q.Get("city"),
		MinPrice:  atoi64(q.Get("min_price")),
		MaxPrice:  atoi64(q.Get("max_price")),
		MinRating: atof(q.Get("min_rating")),
		Sort:      q.Get("sort"),
		Page:      atoi(q.Get("page")),
		Limit:     atoi(q.Get("limit")),
	})
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": vs, "total": total})
}

func (h *Handler) AdminSetVendorStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req struct{ Status string }
	if err := httpx.ReadJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	v, err := h.Vendors.SetStatus(r.Context(), id, domain.VendorStatus(req.Status))
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, v)
}

// --- bookings ---

func (h *Handler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	uid, _ := pkgauth.UserIDFrom(r.Context())
	var req struct {
		VendorID, ServiceID, EventDate, Note string
		GuestCount                           int
		Amount                               int64
	}
	if err := httpx.ReadJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	b, err := h.Bookings.Create(r.Context(), uid, booking.CreateInput{
		VendorID: req.VendorID, ServiceID: req.ServiceID, EventDate: req.EventDate,
		GuestCount: req.GuestCount, Note: req.Note, Amount: req.Amount,
	})
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, b)
}

func (h *Handler) ListMyBookings(w http.ResponseWriter, r *http.Request) {
	uid, _ := pkgauth.UserIDFrom(r.Context())
	role, _ := pkgauth.RoleFrom(r.Context())
	page := readPage(r)
	var (
		bs  []*domain.Booking
		err error
	)
	if role == "vendor" {
		bs, err = h.Bookings.ListForVendor(r.Context(), uid, page)
	} else {
		bs, err = h.Bookings.ListForCustomer(r.Context(), uid, page)
	}
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, bs)
}

// readPage extracts limit + offset query params with sane defaults.
func readPage(r *http.Request) ports.Page {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	return ports.Page{Limit: limit, Offset: offset}
}

func (h *Handler) BookingTransition(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	uid, _ := pkgauth.UserIDFrom(r.Context())
	var req struct{ Status string }
	if err := httpx.ReadJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	next := domain.BookingStatus(req.Status)
	var (
		b   *domain.Booking
		err error
	)
	if next == domain.BookingCancelled {
		b, err = h.Bookings.CustomerCancel(r.Context(), uid, id)
	} else {
		b, err = h.Bookings.VendorTransition(r.Context(), uid, id, next)
	}
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, b)
}

func (h *Handler) PayBooking(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	uid, _ := pkgauth.UserIDFrom(r.Context())
	var req struct{ CardID, Currency string }
	if err := httpx.ReadJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	b, err := h.Bookings.Pay(r.Context(), uid, id, req.CardID, req.Currency)
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, b)
}

func (h *Handler) GetBooking(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	uid, _ := pkgauth.UserIDFrom(r.Context())
	b, err := h.Bookings.Find(r.Context(), uid, id)
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, b)
}

// --- reviews ---

func (h *Handler) SubmitReview(w http.ResponseWriter, r *http.Request) {
	uid, _ := pkgauth.UserIDFrom(r.Context())
	var req struct {
		BookingID, Text string
		Rating          int
	}
	if err := httpx.ReadJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	rv, err := h.Reviews.Submit(r.Context(), uid, req.BookingID, req.Rating, req.Text)
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, rv)
}

func (h *Handler) ListReviews(w http.ResponseWriter, r *http.Request) {
	vendorID := r.PathValue("vendorId")
	rv, err := h.Reviews.ListForVendor(r.Context(), vendorID, readPage(r))
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, rv)
}

// --- photos ---

func (h *Handler) UploadPhoto(w http.ResponseWriter, r *http.Request) {
	uid, _ := pkgauth.UserIDFrom(r.Context())
	data, err := io.ReadAll(http.MaxBytesReader(w, r.Body, domain.MaxPhotoSize+1024))
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "read body")
		return
	}
	p, err := h.Photos.Upload(r.Context(), uid, data)
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, p)
}

func (h *Handler) ServePhoto(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	p, err := h.Photos.Get(r.Context(), id)
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	w.Header().Set("Content-Type", p.MIME)
	w.Header().Set("Content-Length", strconv.FormatInt(p.Size, 10))
	_, _ = w.Write(p.Data)
}

func (h *Handler) DeletePhoto(w http.ResponseWriter, r *http.Request) {
	uid, _ := pkgauth.UserIDFrom(r.Context())
	id := r.PathValue("id")
	if err := h.Photos.Delete(r.Context(), uid, id); err != nil {
		httpx.HandleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- notifications ---

func (h *Handler) ListNotifications(w http.ResponseWriter, r *http.Request) {
	uid, _ := pkgauth.UserIDFrom(r.Context())
	ns, err := h.Notifications.ListForUser(r.Context(), uid, readPage(r))
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, ns)
}

func (h *Handler) RegisterFCM(w http.ResponseWriter, r *http.Request) {
	uid, _ := pkgauth.UserIDFrom(r.Context())
	var req struct{ Token, Platform string }
	if err := httpx.ReadJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if strings.TrimSpace(req.Token) == "" {
		httpx.WriteError(w, http.StatusBadRequest, "token required")
		return
	}
	if err := h.Notifications.RegisterToken(r.Context(), uid, req.Token, req.Platform); err != nil {
		httpx.HandleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- admin ---

func (h *Handler) AdminStats(w http.ResponseWriter, r *http.Request) {
	st, err := h.Admin.Compute(r.Context())
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, st)
}
