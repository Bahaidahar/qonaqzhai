// Package payment implements the charge / refund lifecycle.
package payment

import (
	"context"
	"fmt"

	"qonaqzhai-backend/pkg/errs"
	"qonaqzhai-backend/services/payment/internal/domain"
	"qonaqzhai-backend/services/payment/internal/ports"
)

// Deps bundles payment usecase collaborators.
type Deps struct {
	Payments ports.PaymentRepo
	Cards    ports.CardRepo
	Gateway  ports.Gateway
	Core     ports.CoreClient // optional — when nil, MarkBookingPaid is skipped
	Clock    ports.Clock
}

// Service exposes payment operations.
type Service struct{ d Deps }

// New constructs a payment Service.
func New(d Deps) *Service { return &Service{d: d} }

// ChargeInput is the payload for a synchronous capture.
type ChargeInput struct {
	BookingID string
	UserID    string
	CardID    string
	Amount    int64
	Currency  string
}

// Charge runs a synchronous capture against the PSP and persists a payment row.
// On success it also notifies core-svc so the booking moves to paid.
func (s *Service) Charge(ctx context.Context, in ChargeInput) (*domain.Payment, error) {
	if in.BookingID == "" || in.UserID == "" || in.CardID == "" || in.Amount <= 0 {
		return nil, fmt.Errorf("charge fields: %w", errs.ErrInvalidInput)
	}
	if existing, err := s.d.Payments.FindByBooking(ctx, in.BookingID); err == nil {
		return existing, errs.ErrAlreadyExists
	}
	c, err := s.d.Cards.Find(ctx, in.CardID)
	if err != nil {
		return nil, err
	}
	if c.UserID != in.UserID {
		return nil, errs.ErrForbidden
	}
	if in.Currency == "" {
		in.Currency = "KZT"
	}
	ref, err := s.d.Gateway.Charge(ctx, ports.ChargeInput{
		OrderID: in.BookingID, Amount: in.Amount, Currency: in.Currency,
		Holder: c.Holder, Last4: c.Last4,
	})
	status := domain.PaymentCaptured
	if err != nil {
		status = domain.PaymentFailed
	}
	p, perr := s.d.Payments.Create(ctx, &domain.Payment{
		BookingID: in.BookingID, UserID: in.UserID, CardID: in.CardID,
		Amount: in.Amount, Currency: in.Currency,
		Status: status, ProviderRef: ref,
	})
	if perr != nil {
		return nil, perr
	}
	if err != nil {
		return p, fmt.Errorf("gateway: %w", err)
	}
	if s.d.Core != nil {
		_ = s.d.Core.MarkBookingPaid(ctx, in.BookingID, p.ID)
	}
	return p, nil
}

// Refund flips a captured payment to refunded. Real PSP refund is out of scope
// for the diploma; the row is updated only.
func (s *Service) Refund(ctx context.Context, paymentID string) (*domain.Payment, error) {
	p, err := s.d.Payments.Find(ctx, paymentID)
	if err != nil {
		return nil, err
	}
	if p.Status != domain.PaymentCaptured {
		return nil, fmt.Errorf("status %s: %w", p.Status, errs.ErrConflict)
	}
	if err := s.d.Payments.UpdateStatus(ctx, paymentID, domain.PaymentRefunded); err != nil {
		return nil, err
	}
	return s.d.Payments.Find(ctx, paymentID)
}

// Find returns a payment by id.
func (s *Service) Find(ctx context.Context, id string) (*domain.Payment, error) {
	return s.d.Payments.Find(ctx, id)
}

// ListForUser returns the caller's payments.
func (s *Service) ListForUser(ctx context.Context, userID string) ([]*domain.Payment, error) {
	return s.d.Payments.ListByUser(ctx, userID)
}
