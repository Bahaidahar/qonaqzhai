package handler

import (
	"net/http"
	"strconv"
	"time"

	"qonaqzhai-backend/internal/adapter/http/httpx"
	"qonaqzhai-backend/internal/adapter/http/middleware"
	"qonaqzhai-backend/internal/domain"
	"qonaqzhai-backend/internal/usecase/admin"
)

// Admin HTTP handler.
type Admin struct {
	svc *admin.Service
}

// NewAdmin constructs an Admin handler.
func NewAdmin(svc *admin.Service) *Admin { return &Admin{svc: svc} }

// Users lists every user.
func (h *Admin) Users(w http.ResponseWriter, r *http.Request) {
	users, err := h.svc.ListUsers(r.Context())
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": users})
}

// UserStatus suspends or restores a user.
func (h *Admin) UserStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	actor, _ := middleware.UserIDFrom(r.Context())
	email, _ := middleware.EmailFrom(r.Context())
	var req struct {
		Status string `json:"status"`
	}
	if err := httpx.ReadJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	u, err := h.svc.SetUserStatus(r.Context(), actor, email, id, domain.UserStatus(req.Status))
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, u)
}

// VendorStatus moderates a vendor.
func (h *Admin) VendorStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	actor, _ := middleware.UserIDFrom(r.Context())
	email, _ := middleware.EmailFrom(r.Context())
	var req struct {
		Status string `json:"status"`
	}
	if err := httpx.ReadJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	v, err := h.svc.SetVendorStatus(r.Context(), actor, email, id, domain.VendorStatus(req.Status))
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, v)
}

// AuditLog returns the admin audit log.
func (h *Admin) AuditLog(w http.ResponseWriter, r *http.Request) {
	limit := parseInt(r.URL.Query().Get("limit"), 100)
	entries, err := h.svc.AuditLog(r.Context(), limit)
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": entries})
}

// Stats returns aggregate platform KPIs.
func (h *Admin) Stats(w http.ResponseWriter, r *http.Request) {
	st, err := h.svc.Stats(r.Context())
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, st)
}

// BookingsTimeseries returns booking counts per day, filterable by from/to (RFC3339 dates).
func (h *Admin) BookingsTimeseries(w http.ResponseWriter, r *http.Request) {
	from := parseDate(r.URL.Query().Get("from"))
	to := parseDate(r.URL.Query().Get("to"))
	pts, err := h.svc.BookingsTimeseries(r.Context(), from, to)
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": pts})
}

// TopCategories returns the most popular vendor categories.
func (h *Admin) TopCategories(w http.ResponseWriter, r *http.Request) {
	limit := parseInt(r.URL.Query().Get("limit"), 10)
	cats, err := h.svc.TopCategories(r.Context(), limit)
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": cats})
}

// ApprovalFunnel returns the vendor moderation funnel.
func (h *Admin) ApprovalFunnel(w http.ResponseWriter, r *http.Request) {
	stages, err := h.svc.ApprovalFunnel(r.Context())
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": stages})
}

func parseDate(s string) time.Time {
	if s == "" {
		return time.Time{}
	}
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t
	}
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return t
	}
	return time.Time{}
}

func parseInt(s string, def int) int {
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return n
}
