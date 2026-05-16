// Package handler hosts thin HTTP adapters that delegate to usecase services.
package handler

import (
	"net/http"

	"qonaqzhai-backend/internal/adapter/http/httpx"
	"qonaqzhai-backend/internal/domain"
	"qonaqzhai-backend/internal/usecase/auth"
)

// Auth HTTP handler.
type Auth struct {
	svc *auth.Service
}

// NewAuth constructs an Auth handler.
func NewAuth(svc *auth.Service) *Auth { return &Auth{svc: svc} }

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

// Signup creates a new account.
func (h *Auth) Signup(w http.ResponseWriter, r *http.Request) {
	var req signupReq
	if err := httpx.ReadJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	out, err := h.svc.Signup(r.Context(), auth.SignupInput{
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

// Login authenticates with email + password.
func (h *Auth) Login(w http.ResponseWriter, r *http.Request) {
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

// Refresh rotates the refresh token.
func (h *Auth) Refresh(w http.ResponseWriter, r *http.Request) {
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

// Logout revokes a refresh token.
func (h *Auth) Logout(w http.ResponseWriter, r *http.Request) {
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

// ForgotPassword starts the password reset flow.
func (h *Auth) ForgotPassword(w http.ResponseWriter, r *http.Request) {
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

// ResetPassword consumes a reset token + sets a new password.
func (h *Auth) ResetPassword(w http.ResponseWriter, r *http.Request) {
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
