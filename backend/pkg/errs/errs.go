// Package errs defines sentinel errors shared by all services so that HTTP and
// gRPC adapters can map them to status codes consistently.
package errs

import "errors"

var (
	ErrNotFound       = errors.New("not found")
	ErrAlreadyExists  = errors.New("already exists")
	ErrInvalidInput   = errors.New("invalid input")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrForbidden      = errors.New("forbidden")
	ErrConflict       = errors.New("conflict")
	ErrSuspended      = errors.New("account suspended")
	ErrBadCredentials = errors.New("invalid credentials")
	ErrTooLarge       = errors.New("payload too large")
	ErrRateLimited    = errors.New("rate limited")
	ErrUpstream       = errors.New("upstream failure")
)
