package handler

import (
	"net/http"

	"qonaqzhai-backend/internal/adapter/http/httpx"
	"qonaqzhai-backend/internal/usecase/chat"
)

// Chat HTTP handler.
type Chat struct{ svc *chat.Service }

// NewChat constructs a Chat handler.
func NewChat(svc *chat.Service) *Chat { return &Chat{svc: svc} }

type chatReq struct {
	Message string `json:"message"`
}

type chatBlockDTO struct {
	Type string         `json:"type"`
	Data map[string]any `json:"data"`
}

type chatReplyDTO struct {
	Reply  string         `json:"reply"`
	Blocks []chatBlockDTO `json:"blocks"`
}

// Generate handles a single chat message and returns AI reply + structured blocks.
func (h *Chat) Generate(w http.ResponseWriter, r *http.Request) {
	var req chatReq
	if err := httpx.ReadJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	reply, err := h.svc.Generate(r.Context(), req.Message)
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	out := chatReplyDTO{Reply: reply.Reply, Blocks: make([]chatBlockDTO, 0, len(reply.Blocks))}
	for _, b := range reply.Blocks {
		out.Blocks = append(out.Blocks, chatBlockDTO{Type: b.Type, Data: b.Data})
	}
	httpx.WriteJSON(w, http.StatusOK, out)
}
