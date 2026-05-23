// Package grpcclient holds gRPC clients payment-svc uses to call other services.
package grpcclient

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	corev1 "qonaqzhai-backend/gen/proto/core/v1"
	"qonaqzhai-backend/services/payment/internal/ports"
)

// CoreClient adapts core-svc gRPC to ports.CoreClient.
type CoreClient struct {
	c    corev1.CoreServiceClient
	conn *grpc.ClientConn
}

// NewCoreClient dials core-svc at addr.
func NewCoreClient(addr string) (*CoreClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("dial core %s: %w", addr, err)
	}
	return &CoreClient{c: corev1.NewCoreServiceClient(conn), conn: conn}, nil
}

// Close releases the underlying connection.
func (c *CoreClient) Close() error { return c.conn.Close() }

// MarkBookingPaid tells core a payment captured successfully.
func (c *CoreClient) MarkBookingPaid(ctx context.Context, bookingID, paymentID string) error {
	_, err := c.c.MarkBookingPaid(ctx, &corev1.MarkBookingPaidRequest{
		BookingId: bookingID, PaymentId: paymentID,
	})
	return err
}

var _ ports.CoreClient = (*CoreClient)(nil)
