// Package grpcserver implements AuthService gRPC contract.
// VerifyToken validates JWTs issued by this service so that downstream services
// (core, realtime) never need the JWT secret.
package grpcserver

import (
	"context"

	authv1 "qonaqzhai-backend/gen/proto/auth/v1"
	"qonaqzhai-backend/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server adapts a TokenIssuer to the AuthService gRPC interface.
type Server struct {
	authv1.UnimplementedAuthServiceServer
	tokens usecase.TokenIssuer
}

// New constructs the gRPC server.
func New(tokens usecase.TokenIssuer) *Server { return &Server{tokens: tokens} }

// VerifyToken returns the JWT subject and metadata. It never logs the raw token.
func (s *Server) VerifyToken(_ context.Context, req *authv1.VerifyTokenRequest) (*authv1.VerifyTokenResponse, error) {
	if req == nil || req.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "token required")
	}
	claims, err := s.tokens.Parse(req.Token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}
	return &authv1.VerifyTokenResponse{
		UserId: claims.UserID,
		Email:  claims.Email,
		Role:   string(claims.Role),
	}, nil
}
