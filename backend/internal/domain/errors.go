package domain

import "errors"

// Sentinel errors returned across layers.
// HTTP adapter maps them to status codes; usecase code uses errors.Is.
var (
	ErrNotFound         = errors.New("not found")
	ErrAlreadyExists    = errors.New("already exists")
	ErrInvalidInput     = errors.New("invalid input")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrForbidden        = errors.New("forbidden")
	ErrConflict         = errors.New("conflict")
	ErrSuspended        = errors.New("account suspended")
	ErrBadCredentials   = errors.New("invalid credentials")
	ErrTooLarge         = errors.New("payload too large")
	ErrRateLimited      = errors.New("rate limited")
)
