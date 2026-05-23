// Package grpc exposes the auth service over gRPC. Other microservices call
// these methods instead of decoding JWTs locally.
package grpc

import (
	"context"

	authv1 "qonaqzhai-backend/gen/proto/auth/v1"
	"qonaqzhai-backend/pkg/grpcutil"

	"qonaqzhai-backend/services/auth/internal/domain"
	"qonaqzhai-backend/services/auth/internal/ports"
	"qonaqzhai-backend/services/auth/internal/usecase"
)

// Server implements authv1.AuthServiceServer.
type Server struct {
	authv1.UnimplementedAuthServiceServer
	svc *usecase.Service
}

// New constructs the gRPC server.
func New(svc *usecase.Service) *Server { return &Server{svc: svc} }

// VerifyToken decodes a Bearer access token and returns the principal.
func (s *Server) VerifyToken(ctx context.Context, req *authv1.VerifyTokenRequest) (*authv1.VerifyTokenResponse, error) {
	c, exp, err := s.svc.VerifyAccessToken(ctx, req.GetToken())
	if err != nil {
		return nil, grpcutil.ToStatus(err)
	}
	return &authv1.VerifyTokenResponse{
		UserId: c.UserID,
		Email:  c.Email,
		Role:   c.Role,
		Status: c.Status,
		Exp:    exp.Unix(),
	}, nil
}

// GetUser returns a single user by id.
func (s *Server) GetUser(ctx context.Context, req *authv1.GetUserRequest) (*authv1.User, error) {
	u, err := s.svc.FindUser(ctx, req.GetUserId())
	if err != nil {
		return nil, grpcutil.ToStatus(err)
	}
	return userProto(u), nil
}

// GetUsersBatch returns multiple users by id.
func (s *Server) GetUsersBatch(ctx context.Context, req *authv1.GetUsersBatchRequest) (*authv1.GetUsersBatchResponse, error) {
	users, err := s.svc.FindUsers(ctx, req.GetUserIds())
	if err != nil {
		return nil, grpcutil.ToStatus(err)
	}
	out := make([]*authv1.User, len(users))
	for i, u := range users {
		out[i] = userProto(u)
	}
	return &authv1.GetUsersBatchResponse{Users: out}, nil
}

// AdminListUsers returns a paginated list for admin UIs.
func (s *Server) AdminListUsers(ctx context.Context, req *authv1.AdminListUsersRequest) (*authv1.AdminListUsersResponse, error) {
	users, err := s.svc.ListUsersForAdmin(ctx, ports.ListUsersOpts{
		Limit: int(req.GetLimit()), Offset: int(req.GetOffset()), Role: req.GetRole(),
	})
	if err != nil {
		return nil, grpcutil.ToStatus(err)
	}
	out := make([]*authv1.User, len(users))
	for i, u := range users {
		out[i] = userProto(u)
	}
	return &authv1.AdminListUsersResponse{Users: out}, nil
}

// AdminSetStatus changes a user's lifecycle state.
func (s *Server) AdminSetStatus(ctx context.Context, req *authv1.AdminSetStatusRequest) (*authv1.User, error) {
	u, err := s.svc.SetUserStatus(ctx, req.GetUserId(), domain.UserStatus(req.GetStatus()))
	if err != nil {
		return nil, grpcutil.ToStatus(err)
	}
	return userProto(u), nil
}

func userProto(u *domain.User) *authv1.User {
	if u == nil {
		return nil
	}
	return &authv1.User{
		Id:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		Role:      string(u.Role),
		Status:    string(u.Status),
		CreatedAt: u.CreatedAt.Unix(),
	}
}
