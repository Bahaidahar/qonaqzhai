// Package ws hosts the WebSocket hub used by the realtime messaging layer.
package ws

import (
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Envelope is the on-wire frame shape.
//
//	{"op":"message", "data":{...}}
type Envelope struct {
	Op   string `json:"op"`
	Data any    `json:"data,omitempty"`
}

// IncomingMessage is what clients post over the socket.
type IncomingMessage struct {
	ThreadID string `json:"threadId"`
	Text     string `json:"text"`
}

// Conn represents one connected websocket peer.
type Conn struct {
	ws     *websocket.Conn
	hub    *Hub
	userID string
	send   chan []byte
}

// Hub multiplexes broadcasts to multiple connections per user.
type Hub struct {
	mu      sync.RWMutex
	clients map[string]map[*Conn]struct{} // userID → set of conns
	logger  *slog.Logger

	OnIncoming func(userID string, msg IncomingMessage)
}

// NewHub constructs an empty Hub.
func NewHub(logger *slog.Logger) *Hub {
	if logger == nil {
		logger = slog.Default()
	}
	return &Hub{
		clients: map[string]map[*Conn]struct{}{},
		logger:  logger,
	}
}

// Attach upgrades + binds a connection to userID. Runs read+write pumps.
func (h *Hub) Attach(ws *websocket.Conn, userID string) {
	c := &Conn{ws: ws, hub: h, userID: userID, send: make(chan []byte, 32)}
	h.register(c)
	go c.writePump()
	c.readPump()
}

func (h *Hub) register(c *Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	conns, ok := h.clients[c.userID]
	if !ok {
		conns = map[*Conn]struct{}{}
		h.clients[c.userID] = conns
	}
	conns[c] = struct{}{}
}

func (h *Hub) unregister(c *Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if conns, ok := h.clients[c.userID]; ok {
		delete(conns, c)
		if len(conns) == 0 {
			delete(h.clients, c.userID)
		}
	}
	close(c.send)
}

// SendToUsers broadcasts the envelope to every active connection of each userID.
// Called by usecase services after a domain event (e.g. new thread message).
func (h *Hub) SendToUsers(env Envelope, userIDs ...string) {
	raw, err := json.Marshal(env)
	if err != nil {
		h.logger.Warn("ws marshal failed", slog.String("err", err.Error()))
		return
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, uid := range userIDs {
		for c := range h.clients[uid] {
			select {
			case c.send <- raw:
			default:
				// Slow client — drop to keep hub healthy.
			}
		}
	}
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 4096
)

func (c *Conn) readPump() {
	defer func() {
		c.hub.unregister(c)
		_ = c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	_ = c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error {
		_ = c.ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, raw, err := c.ws.ReadMessage()
		if err != nil {
			return
		}
		var env Envelope
		if err := json.Unmarshal(raw, &env); err != nil {
			continue
		}
		if env.Op == "message" && c.hub.OnIncoming != nil {
			b, _ := json.Marshal(env.Data)
			var m IncomingMessage
			if err := json.Unmarshal(b, &m); err == nil {
				c.hub.OnIncoming(c.userID, m)
			}
		}
	}
}

func (c *Conn) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.ws.Close()
	}()
	for {
		select {
		case msg, ok := <-c.send:
			_ = c.ws.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = c.ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.ws.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			_ = c.ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
