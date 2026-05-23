package grpcclient

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	realtimev1 "qonaqzhai-backend/gen/proto/realtime/v1"
	"qonaqzhai-backend/services/core/internal/ports"
)

// RealtimeClient adapts realtime-svc gRPC to ports.RealtimeClient.
type RealtimeClient struct {
	c    realtimev1.RealtimeServiceClient
	conn *grpc.ClientConn
}

// NewRealtimeClient dials realtime-svc at addr.
func NewRealtimeClient(addr string) (*RealtimeClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("dial realtime %s: %w", addr, err)
	}
	return &RealtimeClient{c: realtimev1.NewRealtimeServiceClient(conn), conn: conn}, nil
}

// Close releases the underlying connection.
func (r *RealtimeClient) Close() error { return r.conn.Close() }

// EnsureThread idempotently creates the DM thread for a booking.
func (r *RealtimeClient) EnsureThread(ctx context.Context, bookingID, customerID, vendorUserID string) error {
	_, err := r.c.EnsureThread(ctx, &realtimev1.EnsureThreadRequest{
		BookingId: bookingID, CustomerId: customerID, VendorId: vendorUserID,
	})
	return err
}

// Publish fans out an event to the supplied user ids.
func (r *RealtimeClient) Publish(ctx context.Context, event string, payloadJSON []byte, userIDs ...string) error {
	_, err := r.c.PublishEvent(ctx, &realtimev1.PublishEventRequest{
		Event: event, PayloadJson: payloadJSON, UserIds: userIDs,
	})
	return err
}

var _ ports.RealtimeClient = (*RealtimeClient)(nil)
