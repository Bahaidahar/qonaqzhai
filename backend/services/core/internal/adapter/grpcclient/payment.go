package grpcclient

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	paymentv1 "qonaqzhai-backend/gen/proto/payment/v1"
	"qonaqzhai-backend/services/core/internal/ports"
)

// PaymentClient adapts payment-svc gRPC to ports.PaymentClient.
type PaymentClient struct {
	c    paymentv1.PaymentServiceClient
	conn *grpc.ClientConn
}

// NewPaymentClient dials payment-svc at addr.
func NewPaymentClient(addr string) (*PaymentClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("dial payment %s: %w", addr, err)
	}
	return &PaymentClient{c: paymentv1.NewPaymentServiceClient(conn), conn: conn}, nil
}

// Close releases the underlying connection.
func (p *PaymentClient) Close() error { return p.conn.Close() }

// Charge runs a synchronous payment.
func (p *PaymentClient) Charge(ctx context.Context, in ports.ChargeRequest) (*ports.PaymentResult, error) {
	out, err := p.c.Charge(ctx, &paymentv1.ChargeRequest{
		BookingId: in.BookingID, UserId: in.UserID, CardId: in.CardID,
		Amount: in.Amount, Currency: in.Currency,
	})
	if err != nil {
		return nil, err
	}
	return &ports.PaymentResult{ID: out.GetId(), Status: out.GetStatus()}, nil
}

var _ ports.PaymentClient = (*PaymentClient)(nil)
