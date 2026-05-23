// Package grpcutil bridges pkg/errs sentinels with gRPC status codes so service
// boundaries preserve error semantics.
package grpcutil

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"qonaqzhai-backend/pkg/errs"
)

// ToStatus converts an errs sentinel to a gRPC status error. Unknown errors map
// to Internal.
func ToStatus(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, errs.ErrNotFound):
		return status.Error(codes.NotFound, "not found")
	case errors.Is(err, errs.ErrAlreadyExists):
		return status.Error(codes.AlreadyExists, "already exists")
	case errors.Is(err, errs.ErrInvalidInput):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, errs.ErrUnauthorized), errors.Is(err, errs.ErrBadCredentials):
		return status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, errs.ErrForbidden), errors.Is(err, errs.ErrSuspended):
		return status.Error(codes.PermissionDenied, err.Error())
	case errors.Is(err, errs.ErrConflict):
		return status.Error(codes.FailedPrecondition, err.Error())
	case errors.Is(err, errs.ErrRateLimited):
		return status.Error(codes.ResourceExhausted, "rate limited")
	default:
		return status.Error(codes.Internal, err.Error())
	}
}

// FromStatus maps a gRPC status error back to an errs sentinel. The original
// status message is preserved by wrapping.
func FromStatus(err error) error {
	if err == nil {
		return nil
	}
	st, ok := status.FromError(err)
	if !ok {
		return err
	}
	switch st.Code() {
	case codes.OK:
		return nil
	case codes.NotFound:
		return errs.ErrNotFound
	case codes.AlreadyExists:
		return errs.ErrAlreadyExists
	case codes.InvalidArgument:
		return errs.ErrInvalidInput
	case codes.Unauthenticated:
		return errs.ErrUnauthorized
	case codes.PermissionDenied:
		return errs.ErrForbidden
	case codes.FailedPrecondition:
		return errs.ErrConflict
	case codes.ResourceExhausted:
		return errs.ErrRateLimited
	default:
		return errs.ErrUpstream
	}
}
