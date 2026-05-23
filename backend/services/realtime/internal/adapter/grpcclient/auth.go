// Package grpcclient holds gRPC clients realtime-svc uses to call other services.
package grpcclient

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authv1 "qonaqzhai-backend/gen/proto/auth/v1"
	"qonaqzhai-backend/services/realtime/internal/ports"
)

// AuthClient adapts auth-svc gRPC to ports.AuthClient.
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

// Close releases the connection.
func (a *AuthClient) Close() error { return a.conn.Close() }

// GetUsersBatch returns the user records matching userIDs.
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
		out = append(out, &ports.ExternalUser{ID: u.GetId(), Email: u.GetEmail(), Name: u.GetName()})
	}
	return out, nil
}

var _ ports.AuthClient = (*AuthClient)(nil)
