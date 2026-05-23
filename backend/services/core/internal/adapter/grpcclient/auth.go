// Package grpcclient holds gRPC clients core-svc uses to call other services.
package grpcclient

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authv1 "qonaqzhai-backend/gen/proto/auth/v1"
	"qonaqzhai-backend/services/core/internal/ports"
)

// AuthClient adapts auth-svc gRPC to the ports.AuthClient interface.
type AuthClient struct {
	c    authv1.AuthServiceClient
	conn *grpc.ClientConn
}

// NewAuthClient dials auth-svc at addr.
func NewAuthClient(addr string) (*AuthClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("dial auth %s: %w", addr, err)
	}
	return &AuthClient{c: authv1.NewAuthServiceClient(conn), conn: conn}, nil
}

// Close releases the underlying connection.
func (a *AuthClient) Close() error { return a.conn.Close() }

// GetUser returns the user matching userID.
func (a *AuthClient) GetUser(ctx context.Context, userID string) (*ports.ExternalUser, error) {
	u, err := a.c.GetUser(ctx, &authv1.GetUserRequest{UserId: userID})
	if err != nil {
		return nil, err
	}
	return &ports.ExternalUser{ID: u.GetId(), Email: u.GetEmail(), Name: u.GetName(), Role: u.GetRole()}, nil
}

// GetUsersBatch returns users for the supplied ids in arbitrary order.
func (a *AuthClient) GetUsersBatch(ctx context.Context, userIDs []string) ([]*ports.ExternalUser, error) {
	if len(userIDs) == 0 {
		return nil, nil
	}
	resp, err := a.c.GetUsersBatch(ctx, &authv1.GetUsersBatchRequest{UserIds: userIDs})
	if err != nil {
		return nil, err
	}
	out := make([]*ports.ExternalUser, 0, len(resp.GetUsers()))
	for _, u := range resp.GetUsers() {
		out = append(out, &ports.ExternalUser{
			ID: u.GetId(), Email: u.GetEmail(), Name: u.GetName(), Role: u.GetRole(),
		})
	}
	return out, nil
}

var _ ports.AuthClient = (*AuthClient)(nil)
