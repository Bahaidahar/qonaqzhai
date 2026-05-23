// Package httpx provides JSON helpers and centralized error-to-status mapping
// that any service can consume.
package httpx

import (
	"encoding/json"
	"errors"
	"net/http"

	"qonaqzhai-backend/pkg/errs"
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

// HandleError maps an errs sentinel to an HTTP response. Returns true when
// err was handled (including nil → false).
func HandleError(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}
	switch {
	case errors.Is(err, errs.ErrNotFound):
		WriteError(w, http.StatusNotFound, "not found")
	case errors.Is(err, errs.ErrAlreadyExists):
		WriteError(w, http.StatusConflict, "already exists")
	case errors.Is(err, errs.ErrConflict):
		WriteError(w, http.StatusConflict, err.Error())
	case errors.Is(err, errs.ErrInvalidInput):
		WriteError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, errs.ErrUnauthorized):
		WriteError(w, http.StatusUnauthorized, "unauthorized")
	case errors.Is(err, errs.ErrBadCredentials):
		WriteError(w, http.StatusUnauthorized, "invalid credentials")
	case errors.Is(err, errs.ErrSuspended):
		WriteError(w, http.StatusForbidden, "account suspended")
	case errors.Is(err, errs.ErrForbidden):
		WriteError(w, http.StatusForbidden, "forbidden")
	case errors.Is(err, errs.ErrTooLarge):
		WriteError(w, http.StatusRequestEntityTooLarge, "payload too large")
	case errors.Is(err, errs.ErrRateLimited):
		WriteError(w, http.StatusTooManyRequests, "rate limited")
	default:
		WriteError(w, http.StatusInternalServerError, "internal error")
	}
	return true
}
