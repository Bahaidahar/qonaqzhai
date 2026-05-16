// Package httpx provides JSON helpers and centralized domain-error to HTTP-status mapping.
package httpx

import (
	"encoding/json"
	"errors"
	"net/http"

	"qonaqzhai-backend/internal/domain"
)

// WriteJSON serializes body as JSON with the given status.
func WriteJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

// WriteError emits a JSON error envelope with the given status.
func WriteError(w http.ResponseWriter, status int, msg string) {
	WriteJSON(w, status, map[string]string{"error": msg})
}

// ReadJSON decodes the request body strictly (rejects unknown fields).
func ReadJSON(r *http.Request, v any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}

// HandleError maps a domain/usecase error to an HTTP response.
// Returns true when err was handled.
func HandleError(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}
	switch {
	case errors.Is(err, domain.ErrNotFound):
		WriteError(w, http.StatusNotFound, "not found")
	case errors.Is(err, domain.ErrAlreadyExists):
		WriteError(w, http.StatusConflict, "already exists")
	case errors.Is(err, domain.ErrConflict):
		WriteError(w, http.StatusConflict, err.Error())
	case errors.Is(err, domain.ErrInvalidInput):
		WriteError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrUnauthorized):
		WriteError(w, http.StatusUnauthorized, "unauthorized")
	case errors.Is(err, domain.ErrBadCredentials):
		WriteError(w, http.StatusUnauthorized, "invalid credentials")
	case errors.Is(err, domain.ErrSuspended):
		WriteError(w, http.StatusForbidden, "account suspended")
	case errors.Is(err, domain.ErrForbidden):
		WriteError(w, http.StatusForbidden, "forbidden")
	case errors.Is(err, domain.ErrTooLarge):
		WriteError(w, http.StatusRequestEntityTooLarge, "payload too large")
	case errors.Is(err, domain.ErrRateLimited):
		WriteError(w, http.StatusTooManyRequests, "rate limited")
	default:
		WriteError(w, http.StatusInternalServerError, "internal error")
	}
	return true
}
