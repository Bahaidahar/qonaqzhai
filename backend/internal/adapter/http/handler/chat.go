package handler

import (
	"net/http"

	"qonaqzhai-backend/internal/adapter/http/httpx"
	"qonaqzhai-backend/internal/adapter/http/middleware"
	"qonaqzhai-backend/internal/domain"
	"qonaqzhai-backend/internal/usecase/chat"
)

// Chat HTTP handler.
type Chat struct{ svc *chat.Service }

// NewChat constructs a Chat handler.
func NewChat(svc *chat.Service) *Chat { return &Chat{svc: svc} }

type chatReq struct {
	Message string `json:"message"`
	ChatID  string `json:"chatId,omitempty"`
}

type chatBlockDTO struct {
	Type string         `json:"type"`
	Data map[string]any `json:"data"`
}

type chatReplyDTO struct {
	ChatID string         `json:"chatId,omitempty"`
	Reply  string         `json:"reply"`
	Blocks []chatBlockDTO `json:"blocks"`
}

// Generate persists user message + AI reply for the calling user.
// When req.ChatID is empty, a new chat is created and returned in the response.
func (h *Chat) Generate(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFrom(r.Context())
	var req chatReq
	if err := httpx.ReadJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	out, err := h.svc.Send(r.Context(), uid, req.ChatID, req.Message)
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	dto := chatReplyDTO{
		ChatID: out.ChatID,
		Reply:  out.Reply,
		Blocks: make([]chatBlockDTO, 0, len(out.Blocks)),
	}
	for _, b := range out.Blocks {
		dto.Blocks = append(dto.Blocks, chatBlockDTO{Type: b.Type, Data: b.Data})
	}
	httpx.WriteJSON(w, http.StatusOK, dto)
}

// ListChats returns the user's chat history.
func (h *Chat) ListChats(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFrom(r.Context())
	chats, err := h.svc.ListChats(r.Context(), uid)
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": chats})
}

type chatMessageDTO struct {
	ID        string         `json:"id"`
	Role      string         `json:"role"`
	Text      string         `json:"text"`
	Blocks    []chatBlockDTO `json:"blocks,omitempty"`
	CreatedAt string         `json:"createdAt"`
}

type chatDetailDTO struct {
	*domain.Chat
	Messages []chatMessageDTO `json:"messages"`
}

// GetChat returns chat details + messages, enforcing ownership.
func (h *Chat) GetChat(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFrom(r.Context())
	id := r.PathValue("id")
	c, msgs, err := h.svc.GetChat(r.Context(), uid, id)
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	out := chatDetailDTO{Chat: c, Messages: make([]chatMessageDTO, 0, len(msgs))}
	for _, m := range msgs {
		dto := chatMessageDTO{
			ID:        m.ID,
			Role:      string(m.Role),
			Text:      m.Text,
			CreatedAt: m.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		for _, b := range m.Blocks {
			if mp, ok := b.(map[string]any); ok {
				t, _ := mp["type"].(string)
				data, _ := mp["data"].(map[string]any)
				dto.Blocks = append(dto.Blocks, chatBlockDTO{Type: t, Data: data})
			}
		}
		out.Messages = append(out.Messages, dto)
	}
	httpx.WriteJSON(w, http.StatusOK, out)
}

// DeleteChat removes a chat owned by the caller.
func (h *Chat) DeleteChat(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFrom(r.Context())
	id := r.PathValue("id")
	if err := h.svc.DeleteChat(r.Context(), uid, id); err != nil {
		httpx.HandleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// RenameChat updates the chat title.
func (h *Chat) RenameChat(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFrom(r.Context())
	id := r.PathValue("id")
	var req struct {
		Title string `json:"title"`
	}
	if err := httpx.ReadJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if err := h.svc.RenameChat(r.Context(), uid, id, req.Title); err != nil {
		httpx.HandleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
