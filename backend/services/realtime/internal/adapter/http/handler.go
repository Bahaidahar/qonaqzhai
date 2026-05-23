// Package http exposes realtime operations to HTTP + WebSocket clients.
package http

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
	"golang.org/x/time/rate"

	pkgauth "qonaqzhai-backend/pkg/auth"
	"qonaqzhai-backend/pkg/httpx"

	"qonaqzhai-backend/services/realtime/internal/adapter/ws"
	"qonaqzhai-backend/services/realtime/internal/usecase/thread"
)

// Handler bundles realtime HTTP + WebSocket endpoints.
type Handler struct {
	Threads *thread.Service
	Hub     *ws.Hub
}

// Health is a trivial liveness probe.
func (h *Handler) Health(w http.ResponseWriter, _ *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok", "svc": "realtime"})
}

func (h *Handler) ListThreads(w http.ResponseWriter, r *http.Request) {
	uid, _ := pkgauth.UserIDFrom(r.Context())
	summaries, err := h.Threads.ListSummaries(r.Context(), uid)
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, summaries)
}

func (h *Handler) GetThread(w http.ResponseWriter, r *http.Request) {
	uid, _ := pkgauth.UserIDFrom(r.Context())
	t, msgs, err := h.Threads.Get(r.Context(), uid, r.PathValue("id"))
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"thread": t, "messages": msgs})
}

func (h *Handler) SendMessage(w http.ResponseWriter, r *http.Request) {
	uid, _ := pkgauth.UserIDFrom(r.Context())
	var req struct{ Text string }
	if err := httpx.ReadJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	m, err := h.Threads.Send(r.Context(), uid, r.PathValue("id"), req.Text)
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, m)
}

// upgrader is the gorilla/websocket configuration used for /api/ws.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(_ *http.Request) bool { return true },
}

// WSConnect upgrades the connection and attaches it to the hub. Auth comes via
// the standard middleware which already injected claims into ctx.
func (h *Handler) WSConnect(w http.ResponseWriter, r *http.Request) {
	uid, ok := pkgauth.UserIDFrom(r.Context())
	if !ok || uid == "" {
		httpx.WriteError(w, http.StatusUnauthorized, "auth required")
		return
	}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return // upgrader already wrote the response
	}
	h.Hub.Attach(c, uid)
}

// Mux wires every HTTP route.
func Mux(h *Handler, mw *pkgauth.Middleware, corsOrigin string, log *slog.Logger) http.Handler {
	rl := httpx.NewRateLimiter(rate.Limit(40), 80)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", h.Health)
	mux.Handle("GET /api/threads", mw.Required(http.HandlerFunc(h.ListThreads)))
	mux.Handle("GET /api/threads/{id}", mw.Required(http.HandlerFunc(h.GetThread)))
	mux.Handle("POST /api/threads/{id}/messages", mw.Required(http.HandlerFunc(h.SendMessage)))
	mux.Handle("GET /api/ws", mw.Required(http.HandlerFunc(h.WSConnect)))

	withLimit := rl.PerIP()(mux)
	withRecover := httpx.Recover(log)(withLimit)
	withLog := httpx.AccessLog(log)(withRecover)
	return httpx.CORS(corsOrigin, withLog)
}
