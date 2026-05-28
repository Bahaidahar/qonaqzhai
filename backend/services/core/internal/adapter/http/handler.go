// Package http exposes core service operations to clients.
package http

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	pkgauth "qonaqzhai-backend/pkg/auth"
	"qonaqzhai-backend/pkg/errs"
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
	Services      ports.ServiceRepo
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
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": bs})
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
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": rv})
}

// --- photos ---

func (h *Handler) UploadPhoto(w http.ResponseWriter, r *http.Request) {
	uid, _ := pkgauth.UserIDFrom(r.Context())
	data, err := readPhotoBody(w, r)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	p, err := h.Photos.Upload(r.Context(), uid, data)
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, p)
}

// readPhotoBody returns the photo bytes regardless of whether the request is
// multipart/form-data (web clients) or a raw image body (mobile / curl). The
// multipart branch reads the "photo" form field; the raw branch reads r.Body.
func readPhotoBody(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	const cap = domain.MaxPhotoSize + 1024
	ct := r.Header.Get("Content-Type")
	if strings.HasPrefix(ct, "multipart/form-data") {
		if err := r.ParseMultipartForm(cap); err != nil {
			return nil, err
		}
		f, _, err := r.FormFile("photo")
		if err != nil {
			return nil, err
		}
		defer f.Close()
		return io.ReadAll(io.LimitReader(f, cap))
	}
	return io.ReadAll(http.MaxBytesReader(w, r.Body, cap))
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
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": ns})
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

// Chart endpoints. Implementations are stubbed empty arrays — admin dashboard
// renders a "no data yet" placeholder when items are empty. Real aggregates
// can be wired through h.Admin once analytics queries are designed.
func (h *Handler) AdminStatsBookings(w http.ResponseWriter, _ *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": []any{}})
}

func (h *Handler) AdminStatsCategories(w http.ResponseWriter, _ *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": []any{}})
}

func (h *Handler) AdminStatsFunnel(w http.ResponseWriter, _ *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": []any{}})
}

// --- vendor services ---

// servicesForVendor resolves the calling user's vendor id (errors → 0 items).
func (h *Handler) servicesForCallingVendor(r *http.Request) (string, error) {
	uid, _ := pkgauth.UserIDFrom(r.Context())
	v, err := h.Vendors.FindByUserID(r.Context(), uid)
	if err != nil {
		return "", err
	}
	return v.ID, nil
}

func (h *Handler) ListMyServices(w http.ResponseWriter, r *http.Request) {
	vendorID, err := h.servicesForCallingVendor(r)
	if err != nil {
		// vendor has no profile yet — treat as empty list rather than 404
		if errors.Is(err, errs.ErrNotFound) {
			httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": []any{}})
			return
		}
		httpx.HandleError(w, err)
		return
	}
	items, err := h.Services.ListByVendor(r.Context(), vendorID)
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": items})
}

type serviceReq struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       int64   `json:"price"`
	Unit        string  `json:"unit"`
	IsActive    *bool   `json:"isActive"`
}

func (h *Handler) AddMyService(w http.ResponseWriter, r *http.Request) {
	vendorID, err := h.servicesForCallingVendor(r)
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	var req serviceReq
	if err := httpx.ReadJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	active := true
	if req.IsActive != nil {
		active = *req.IsActive
	}
	created, err := h.Services.Create(r.Context(), &domain.Service{
		VendorID:    vendorID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Unit:        domain.ServiceUnit(req.Unit),
		IsActive:    active,
	})
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, created)
}

func (h *Handler) UpdateMyService(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req serviceReq
	if err := httpx.ReadJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	updated, err := h.Services.Update(r.Context(), id, domain.ServiceInput{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Unit:        domain.ServiceUnit(req.Unit),
		IsActive:    req.IsActive,
	})
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, updated)
}

func (h *Handler) DeleteMyService(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.Services.Delete(r.Context(), id); err != nil {
		httpx.HandleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListVendorServices(w http.ResponseWriter, r *http.Request) {
	vendorID := r.PathValue("vendorId")
	items, err := h.Services.ListByVendor(r.Context(), vendorID)
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": items})
}

// --- AI chat stub ---

// ChatStub returns a canned event-planning skeleton. The UI is built against
// `{ chatId, message: { id, role, text, blocks } }` and renders block cards;
// returning realistic-shape data here keeps the mobile + web flows alive
// while the real Gemini integration is still pending on the backend.
func (h *Handler) ChatStub(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Message string  `json:"message"`
		ChatID  *string `json:"chatId"`
	}
	if err := httpx.ReadJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	chatID := ""
	if req.ChatID != nil {
		chatID = *req.ChatID
	}
	if chatID == "" {
		chatID = "stub-" + strconv.FormatInt(int64(len(req.Message)), 10)
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"chatId": chatID,
		"message": map[string]any{
			"id":   "stub-reply",
			"role": "ai",
			"text": "Got it — here's a draft plan. Plug in a real LLM in /api/chat to replace this stub.",
			"blocks": []map[string]any{
				{
					"type": "plan",
					"data": map[string]any{
						"title":     "Draft event plan",
						"eventType": "event",
						"city":      "Almaty",
						"guests":    100,
						"budget":    3000000,
					},
				},
			},
		},
	})
}

func (h *Handler) ChatsList(w http.ResponseWriter, _ *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": []any{}})
}

func (h *Handler) ChatGet(w http.ResponseWriter, r *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"id":       r.PathValue("id"),
		"title":    "Stub chat",
		"messages": []any{},
	})
}

func (h *Handler) ChatDelete(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ChatRename(w http.ResponseWriter, r *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"id":    r.PathValue("id"),
		"title": "Renamed",
	})
}
