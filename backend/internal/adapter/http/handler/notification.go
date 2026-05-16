package handler

import (
	"net/http"

	"qonaqzhai-backend/internal/adapter/http/httpx"
	"qonaqzhai-backend/internal/adapter/http/middleware"
	"qonaqzhai-backend/internal/usecase/notification"
)

// Notification HTTP handler — push token registration + inbox.
type Notification struct {
	svc    *notification.Service
	tokens notification.FCMTokenRepo
}

// NewNotification constructs the handler.
func NewNotification(svc *notification.Service, tokens notification.FCMTokenRepo) *Notification {
	return &Notification{svc: svc, tokens: tokens}
}

// RegisterToken accepts {"token", "platform"} and binds to the calling user.
func (h *Notification) RegisterToken(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFrom(r.Context())
	var req struct {
		Token    string `json:"token"`
		Platform string `json:"platform"`
	}
	if err := httpx.ReadJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if req.Token == "" {
		httpx.WriteError(w, http.StatusBadRequest, "token required")
		return
	}
	if err := h.tokens.Register(r.Context(), uid, req.Token, req.Platform); err != nil {
		httpx.HandleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// UnregisterToken removes a token (e.g., on logout).
func (h *Notification) UnregisterToken(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		httpx.WriteError(w, http.StatusBadRequest, "token required")
		return
	}
	if err := h.tokens.Unregister(r.Context(), token); err != nil {
		httpx.HandleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Inbox returns the calling user's notifications (limited by ?limit=).
func (h *Notification) Inbox(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFrom(r.Context())
	list, err := h.svc.ListForUser(r.Context(), uid, 50)
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": list})
}
