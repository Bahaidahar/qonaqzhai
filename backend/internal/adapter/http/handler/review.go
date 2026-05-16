package handler

import (
	"net/http"

	"qonaqzhai-backend/internal/adapter/http/httpx"
	"qonaqzhai-backend/internal/adapter/http/middleware"
	"qonaqzhai-backend/internal/usecase/review"
)

// Review HTTP handler.
type Review struct {
	svc *review.Service
}

// NewReview constructs a Review handler.
func NewReview(svc *review.Service) *Review { return &Review{svc: svc} }

type reviewReq struct {
	BookingID string `json:"bookingId"`
	Rating    int    `json:"rating"`
	Text      string `json:"text"`
}

// Submit creates a review for a completed booking.
func (h *Review) Submit(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFrom(r.Context())
	var req reviewReq
	if err := httpx.ReadJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	rv, err := h.svc.Submit(r.Context(), uid, review.SubmitInput{
		BookingID: req.BookingID, Rating: req.Rating, Text: req.Text,
	})
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, rv)
}

// ListByVendor returns reviews of a vendor.
func (h *Review) ListByVendor(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	list, err := h.svc.ListByVendor(r.Context(), id)
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": list})
}

// AdminDelete removes a review.
func (h *Review) AdminDelete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.AdminDelete(r.Context(), id); err != nil {
		httpx.HandleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
