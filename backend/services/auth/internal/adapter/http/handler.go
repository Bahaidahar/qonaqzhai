// Package http exposes the auth service over HTTP. Handlers are thin wrappers
// around the usecase.Service.
package http

import (
	"net/http"

	pkgauth "qonaqzhai-backend/pkg/auth"
	"qonaqzhai-backend/pkg/httpx"

	"qonaqzhai-backend/services/auth/internal/domain"
	"qonaqzhai-backend/services/auth/internal/usecase"
)

// Handler bundles handlers + the auth service.
type Handler struct{ svc *usecase.Service }

// NewHandler constructs a Handler.
func NewHandler(svc *usecase.Service) *Handler { return &Handler{svc: svc} }

type signupReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Role     string `json:"role"`
}

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authResp struct {
	Token        string       `json:"token"` // legacy alias for accessToken
	AccessToken  string       `json:"accessToken"`
	RefreshToken string       `json:"refreshToken"`
	User         *domain.User `json:"user"`
}

func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	var req signupReq
	if err := httpx.ReadJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	out, err := h.svc.Signup(r.Context(), usecase.SignupInput{
		Email: req.Email, Password: req.Password, Name: req.Name, Role: req.Role,
	})
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, authResp{
		Token: out.AccessToken, AccessToken: out.AccessToken, RefreshToken: out.RefreshToken, User: out.User,
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginReq
	if err := httpx.ReadJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	out, err := h.svc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, authResp{
		Token: out.AccessToken, AccessToken: out.AccessToken, RefreshToken: out.RefreshToken, User: out.User,
	})
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := httpx.ReadJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	out, err := h.svc.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, authResp{
		Token: out.AccessToken, AccessToken: out.AccessToken, RefreshToken: out.RefreshToken, User: out.User,
	})
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refreshToken"`
	}
	_ = httpx.ReadJSON(r, &req)
	if err := h.svc.Logout(r.Context(), req.RefreshToken); err != nil {
		httpx.HandleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	if err := httpx.ReadJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if _, err := h.svc.ForgotPassword(r.Context(), req.Email); err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token       string `json:"token"`
		NewPassword string `json:"newPassword"`
	}
	if err := httpx.ReadJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if err := h.svc.ResetPassword(r.Context(), req.Token, req.NewPassword); err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// Me returns the caller's principal. Requires Required-auth middleware to have
// injected claims into the request context.
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	c, _ := pkgauth.FromContext(r.Context())
	u, err := h.svc.FindUser(r.Context(), c.UserID)
	if err != nil {
		httpx.HandleError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, u)
}

// Health is a trivial liveness probe.
func (h *Handler) Health(w http.ResponseWriter, _ *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok", "svc": "auth"})
}
