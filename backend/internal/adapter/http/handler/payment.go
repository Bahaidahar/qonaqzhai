package handler

import (
	"net/http"

	"qonaqzhai-backend/internal/adapter/http/httpx"
	"qonaqzhai-backend/internal/adapter/http/middleware"
	"qonaqzhai-backend/internal/usecase/payment"
)

// Payment HTTP handler — booking checkout + PayBox callback.
type Payment struct {
	svc *payment.Service
}

// NewPayment constructs the handler.
func NewPayment(svc *payment.Service) *Payment { return &Payment{svc: svc} }

// Start initiates payment for a booking and returns the redirect URL.
func (h *Payment) Start(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFrom(r.Context())
	id := r.PathValue("id")
	url, err := h.svc.StartIntent(r.Context(), uid, id)
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]string{"redirectUrl": url})
}

// Callback handles the PayBox result webhook (form-encoded body).
func (h *Payment) Callback(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "form parse")
		return
	}
	form := map[string]string{}
	for k, v := range r.Form {
		if len(v) > 0 {
			form[k] = v[0]
		}
	}
	if err := h.svc.HandleCallback(r.Context(), form); err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
