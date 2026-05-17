// Package token also exposes a Remote verifier that delegates Parse() to
// auth-svc via gRPC. Issue() is unsupported on the remote — core-svc and
// realtime-svc are consumers only, never minters of tokens.
package token

import (
	"context"
	"errors"
	"time"

	authv1 "qonaqzhai-backend/gen/proto/auth/v1"
	"qonaqzhai-backend/internal/domain"
	"qonaqzhai-backend/internal/usecase"
)

// Remote calls auth-svc to verify JWTs.
type Remote struct {
	client authv1.AuthServiceClient
	parent context.Context
}

// NewRemote wraps an AuthServiceClient as a TokenIssuer.
// The supplied ctx scopes RPC timeouts (parent context cancelation aborts in-flight calls).
func NewRemote(ctx context.Context, c authv1.AuthServiceClient) *Remote {
	return &Remote{client: c, parent: ctx}
}

// Issue is a no-op on the remote — only auth-svc mints tokens.
func (*Remote) Issue(*domain.User, time.Duration) (string, error) {
	return "", errors.New("token.Remote does not issue tokens")
}

// Parse calls auth-svc.VerifyToken; bubbles "invalid token" on any RPC failure.
func (r *Remote) Parse(raw string) (usecase.Claims, error) {
	ctx, cancel := context.WithTimeout(r.parent, 3*time.Second)
	defer cancel()
	res, err := r.client.VerifyToken(ctx, &authv1.VerifyTokenRequest{Token: raw})
	if err != nil {
		return usecase.Claims{}, errors.New("invalid token")
	}
	return usecase.Claims{
		UserID: res.UserId,
		Email:  res.Email,
		Role:   domain.Role(res.Role),
	}, nil
}
