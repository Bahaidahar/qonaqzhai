// Package http exposes the payment service to HTTP clients.
package http

import (
	"log/slog"
	"net/http"

	"golang.org/x/time/rate"

	pkgauth "qonaqzhai-backend/pkg/auth"
	"qonaqzhai-backend/pkg/httpx"

	"qonaqzhai-backend/services/payment/internal/domain"
	"qonaqzhai-backend/services/payment/internal/usecase/card"
	"qonaqzhai-backend/services/payment/internal/usecase/payment"
)

// Handler bundles payment HTTP handlers.
type Handler struct {
	Cards    *card.Service
	Payments *payment.Service
}

// Health is a trivial liveness probe.
func (h *Handler) Health(w http.ResponseWriter, _ *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok", "svc": "payment"})
}

func (h *Handler) AddCard(w http.ResponseWriter, r *http.Request) {
	uid, _ := pkgauth.UserIDFrom(r.Context())
	var req struct {
		Number, Holder    string
		ExpMonth, ExpYear int
	}
	if err := httpx.ReadJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	c, err := h.Cards.Add(r.Context(), uid, domain.CardInput{
		Number: req.Number, ExpMonth: req.ExpMonth, ExpYear: req.ExpYear, Holder: req.Holder,
	})
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, c)
}

func (h *Handler) ListCards(w http.ResponseWriter, r *http.Request) {
	uid, _ := pkgauth.UserIDFrom(r.Context())
	cs, err := h.Cards.List(r.Context(), uid)
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, cs)
}

func (h *Handler) DeleteCard(w http.ResponseWriter, r *http.Request) {
	uid, _ := pkgauth.UserIDFrom(r.Context())
	if err := h.Cards.Delete(r.Context(), uid, r.PathValue("id")); err != nil {
		httpx.HandleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) SetDefaultCard(w http.ResponseWriter, r *http.Request) {
	uid, _ := pkgauth.UserIDFrom(r.Context())
	if err := h.Cards.SetDefault(r.Context(), uid, r.PathValue("id")); err != nil {
		httpx.HandleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListPayments(w http.ResponseWriter, r *http.Request) {
	uid, _ := pkgauth.UserIDFrom(r.Context())
	ps, err := h.Payments.ListForUser(r.Context(), uid)
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, ps)
}

// Mux wires payment HTTP routes.
func Mux(h *Handler, mw *pkgauth.Middleware, corsOrigin string, log *slog.Logger) http.Handler {
	rl := httpx.NewRateLimiter(rate.Limit(20), 40)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", h.Health)
	mux.Handle("GET /api/cards", mw.Required(http.HandlerFunc(h.ListCards)))
	mux.Handle("POST /api/cards", mw.Required(http.HandlerFunc(h.AddCard)))
	mux.Handle("DELETE /api/cards/{id}", mw.Required(http.HandlerFunc(h.DeleteCard)))
	mux.Handle("POST /api/cards/{id}/default", mw.Required(http.HandlerFunc(h.SetDefaultCard)))
	mux.Handle("GET /api/payments", mw.Required(http.HandlerFunc(h.ListPayments)))

	withLimit := rl.PerIP()(mux)
	withRecover := httpx.Recover(log)(withLimit)
	withLog := httpx.AccessLog(log)(withRecover)
	return httpx.CORS(corsOrigin, withLog)
}
