package handler

import (
	"net/http"

	"qonaqzhai-backend/internal/adapter/http/httpx"
	"qonaqzhai-backend/internal/adapter/http/middleware"
	"qonaqzhai-backend/internal/domain"
	"qonaqzhai-backend/internal/usecase/card"
)

// Card HTTP handler — saved cards (mock payment instruments).
type Card struct {
	svc *card.Service
}

// NewCard constructs the handler.
func NewCard(svc *card.Service) *Card { return &Card{svc: svc} }

// List returns the calling user's saved cards.
func (h *Card) List(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFrom(r.Context())
	items, err := h.svc.List(r.Context(), uid)
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": items})
}

// Create adds a new mock card. The PAN is never stored — only last4 + brand.
func (h *Card) Create(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFrom(r.Context())
	var req struct {
		Number      string `json:"number"`
		ExpMonth    int    `json:"expMonth"`
		ExpYear     int    `json:"expYear"`
		Holder      string `json:"holder"`
		MakeDefault bool   `json:"makeDefault"`
	}
	if err := httpx.ReadJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	c, err := h.svc.Add(r.Context(), uid, domain.CardInput{
		Number: req.Number, ExpMonth: req.ExpMonth, ExpYear: req.ExpYear, Holder: req.Holder,
	}, req.MakeDefault)
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, c)
}

// Delete removes a card by id.
func (h *Card) Delete(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFrom(r.Context())
	id := r.PathValue("id")
	if err := h.svc.Delete(r.Context(), uid, id); err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// SetDefault marks a card as the user's default.
func (h *Card) SetDefault(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFrom(r.Context())
	id := r.PathValue("id")
	if err := h.svc.SetDefault(r.Context(), uid, id); err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
