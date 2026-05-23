// Package grpc exposes payment service operations over gRPC.
package grpc

import (
	"context"

	paymentv1 "qonaqzhai-backend/gen/proto/payment/v1"
	"qonaqzhai-backend/pkg/grpcutil"

	"qonaqzhai-backend/services/payment/internal/domain"
	"qonaqzhai-backend/services/payment/internal/usecase/card"
	"qonaqzhai-backend/services/payment/internal/usecase/payment"
)

// Server implements paymentv1.PaymentServiceServer.
type Server struct {
	paymentv1.UnimplementedPaymentServiceServer
	payments *payment.Service
	cards    *card.Service
}

// New constructs the gRPC server.
func New(pmt *payment.Service, crd *card.Service) *Server {
	return &Server{payments: pmt, cards: crd}
}

// Charge runs a synchronous capture against the PSP.
func (s *Server) Charge(ctx context.Context, req *paymentv1.ChargeRequest) (*paymentv1.Payment, error) {
	p, err := s.payments.Charge(ctx, payment.ChargeInput{
		BookingID: req.GetBookingId(), UserID: req.GetUserId(),
		CardID: req.GetCardId(), Amount: req.GetAmount(), Currency: req.GetCurrency(),
	})
	if err != nil {
		return nil, grpcutil.ToStatus(err)
	}
	return paymentProto(p), nil
}

// Refund marks the payment refunded.
func (s *Server) Refund(ctx context.Context, req *paymentv1.RefundRequest) (*paymentv1.Payment, error) {
	p, err := s.payments.Refund(ctx, req.GetPaymentId())
	if err != nil {
		return nil, grpcutil.ToStatus(err)
	}
	return paymentProto(p), nil
}

// GetPayment returns a payment by id.
func (s *Server) GetPayment(ctx context.Context, req *paymentv1.GetPaymentRequest) (*paymentv1.Payment, error) {
	p, err := s.payments.Find(ctx, req.GetPaymentId())
	if err != nil {
		return nil, grpcutil.ToStatus(err)
	}
	return paymentProto(p), nil
}

// ListCardsByUser returns the user's saved cards.
func (s *Server) ListCardsByUser(ctx context.Context, req *paymentv1.ListCardsByUserRequest) (*paymentv1.ListCardsByUserResponse, error) {
	cs, err := s.cards.List(ctx, req.GetUserId())
	if err != nil {
		return nil, grpcutil.ToStatus(err)
	}
	out := make([]*paymentv1.Card, len(cs))
	for i, c := range cs {
		out[i] = cardProto(c)
	}
	return &paymentv1.ListCardsByUserResponse{Cards: out}, nil
}

func paymentProto(p *domain.Payment) *paymentv1.Payment {
	if p == nil {
		return nil
	}
	return &paymentv1.Payment{
		Id:          p.ID,
		BookingId:   p.BookingID,
		UserId:      p.UserID,
		CardId:      p.CardID,
		Amount:      p.Amount,
		Currency:    p.Currency,
		Status:      string(p.Status),
		ProviderRef: p.ProviderRef,
		CreatedAt:   p.CreatedAt.Unix(),
	}
}

func cardProto(c *domain.Card) *paymentv1.Card {
	if c == nil {
		return nil
	}
	return &paymentv1.Card{
		Id:        c.ID,
		UserId:    c.UserID,
		Brand:     c.Brand,
		Last4:     c.Last4,
		ExpMonth:  int32(c.ExpMonth),
		ExpYear:   int32(c.ExpYear),
		Holder:    c.Holder,
		IsDefault: c.IsDefault,
		CreatedAt: c.CreatedAt.Unix(),
	}
}
