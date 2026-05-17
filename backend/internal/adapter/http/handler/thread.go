package handler

import (
	"net/http"

	"qonaqzhai-backend/internal/adapter/http/httpx"
	"qonaqzhai-backend/internal/adapter/http/middleware"
	"qonaqzhai-backend/internal/usecase/thread"
)

// Thread HTTP handler — booking DM threads.
type Thread struct {
	svc *thread.Service
}

// NewThread constructs a Thread handler.
func NewThread(svc *thread.Service) *Thread { return &Thread{svc: svc} }

// List returns threads visible to the calling user, enriched with booking +
// counterpart details so the inbox renders useful info instead of raw IDs.
func (h *Thread) List(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFrom(r.Context())
	items, err := h.svc.ListSummariesForUser(r.Context(), uid)
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": items})
}

// Get returns a thread + messages, enforcing membership.
func (h *Thread) Get(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFrom(r.Context())
	id := r.PathValue("id")
	t, msgs, err := h.svc.Get(r.Context(), uid, id)
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"thread":   t,
		"messages": msgs,
	})
}

// Send appends a message to a thread.
func (h *Thread) Send(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFrom(r.Context())
	id := r.PathValue("id")
	var req struct {
		Text string `json:"text"`
	}
	if err := httpx.ReadJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	m, err := h.svc.Send(r.Context(), uid, id, req.Text)
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, m)
}
