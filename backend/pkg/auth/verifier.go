package auth

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authv1 "qonaqzhai-backend/gen/proto/auth/v1"
)

// Verifier turns a raw access token into Claims by calling the AuthService
// over gRPC. Cheap to share across goroutines.
type Verifier struct {
	client authv1.AuthServiceClient
	conn   *grpc.ClientConn
}

// NewVerifier dials the auth-svc gRPC endpoint and returns a Verifier.
// Call Close when the service shuts down.
//
// addr should be host:port (e.g. "auth-svc:9001"). The connection is
// insecure — TLS is the deployment platform's responsibility (service mesh /
// Kubernetes NetworkPolicy / VPC).
func NewVerifier(addr string) (*Verifier, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("dial auth %s: %w", addr, err)
	}
	return &Verifier{client: authv1.NewAuthServiceClient(conn), conn: conn}, nil
}

// Close releases the underlying gRPC connection.
func (v *Verifier) Close() error { return v.conn.Close() }

// Verify decodes a raw access token. Returns Claims when valid.
func (v *Verifier) Verify(ctx context.Context, token string) (Claims, error) {
	resp, err := v.client.VerifyToken(ctx, &authv1.VerifyTokenRequest{Token: token})
	if err != nil {
		return Claims{}, err
	}
	return Claims{
		UserID: resp.GetUserId(),
		Email:  resp.GetEmail(),
		Role:   resp.GetRole(),
		Status: resp.GetStatus(),
	}, nil
}
