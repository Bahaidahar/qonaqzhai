package handler

import (
	"net/http"
	"strings"

	"github.com/gorilla/websocket"

	"qonaqzhai-backend/internal/adapter/http/httpx"
	"qonaqzhai-backend/internal/adapter/ws"
	"qonaqzhai-backend/internal/usecase"
)

// WS handles the WebSocket upgrade for realtime messaging.
type WS struct {
	hub      *ws.Hub
	tokens   usecase.TokenIssuer
	upgrader websocket.Upgrader
}

// NewWS constructs a WS handler. Origin check is permissive — pin it before prod.
func NewWS(hub *ws.Hub, tokens usecase.TokenIssuer) *WS {
	return &WS{
		hub:    hub,
		tokens: tokens,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(_ *http.Request) bool { return true },
		},
	}
}

// Connect performs the WebSocket upgrade. JWT comes via `?token=` query param
// because browsers can't set Authorization headers on a WS handshake.
func (h *WS) Connect(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		// Allow bearer header as a fallback (e.g. native clients).
		if raw := r.Header.Get("Authorization"); strings.HasPrefix(raw, "Bearer ") {
			token = strings.TrimPrefix(raw, "Bearer ")
		}
	}
	if token == "" {
		httpx.WriteError(w, http.StatusUnauthorized, "missing token")
		return
	}
	claims, err := h.tokens.Parse(token)
	if err != nil {
		httpx.WriteError(w, http.StatusUnauthorized, "invalid token")
		return
	}
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return // Upgrade already wrote the error response
	}
	h.hub.Attach(conn, claims.UserID)
}
