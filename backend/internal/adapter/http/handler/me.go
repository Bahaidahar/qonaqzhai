package handler

import (
	"net/http"

	"qonaqzhai-backend/internal/adapter/http/httpx"
	"qonaqzhai-backend/internal/adapter/http/middleware"
	"qonaqzhai-backend/internal/usecase"
)

// Me serves the /api/me endpoint.
type Me struct {
	users usecase.UserRepo
}

// NewMe constructs a Me handler.
func NewMe(users usecase.UserRepo) *Me { return &Me{users: users} }

// Get returns the authenticated user.
func (h *Me) Get(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFrom(r.Context())
	u, err := h.users.FindByID(r.Context(), uid)
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, u)
}
