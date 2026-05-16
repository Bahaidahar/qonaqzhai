package handler

import (
	"net/http"

	"qonaqzhai-backend/internal/adapter/http/httpx"
	"qonaqzhai-backend/internal/adapter/http/middleware"
	"qonaqzhai-backend/internal/domain"
	"qonaqzhai-backend/internal/usecase/booking"
)

// Booking HTTP handler.
type Booking struct {
	svc *booking.Service
}

// NewBooking constructs a Booking handler.
func NewBooking(svc *booking.Service) *Booking { return &Booking{svc: svc} }

type bookingReq struct {
	VendorID   string `json:"vendorId"`
	EventDate  string `json:"eventDate"`
	GuestCount int    `json:"guestCount"`
	Note       string `json:"note"`
	Amount     int64  `json:"amount"`
}

// Create creates a booking for the calling customer.
func (h *Booking) Create(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFrom(r.Context())
	var req bookingReq
	if err := httpx.ReadJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	b, err := h.svc.Create(r.Context(), uid, booking.CreateInput{
		VendorID:   req.VendorID,
		EventDate:  req.EventDate,
		GuestCount: req.GuestCount,
		Note:       req.Note,
		Amount:     req.Amount,
	})
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, b)
}

// List returns bookings scoped by role.
func (h *Booking) List(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFrom(r.Context())
	role, _ := middleware.RoleFrom(r.Context())
	var (
		list []*domain.Booking
		err  error
	)
	switch role {
	case domain.RoleCustomer:
		list, err = h.svc.ListForCustomer(r.Context(), uid)
	case domain.RoleVendor:
		list, err = h.svc.ListForVendor(r.Context(), uid)
	case domain.RoleAdmin:
		list, err = h.svc.ListAll(r.Context())
	}
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": list})
}

// Update applies a role-appropriate status transition.
func (h *Booking) Update(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFrom(r.Context())
	role, _ := middleware.RoleFrom(r.Context())
	id := r.PathValue("id")
	var req struct {
		Status string `json:"status"`
	}
	if err := httpx.ReadJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	next := domain.BookingStatus(req.Status)
	var (
		b   *domain.Booking
		err error
	)
	switch role {
	case domain.RoleAdmin:
		b, err = h.svc.AdminTransition(r.Context(), id, next)
	case domain.RoleVendor:
		b, err = h.svc.VendorTransition(r.Context(), uid, id, next)
	case domain.RoleCustomer:
		b, err = h.svc.CustomerTransition(r.Context(), uid, id, next)
	default:
		httpx.WriteError(w, http.StatusForbidden, "forbidden")
		return
	}
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, b)
}
