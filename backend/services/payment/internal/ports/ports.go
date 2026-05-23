// Package ports defines the interfaces the payment service usecases depend on.
package ports

import (
	"context"
	"time"

	"qonaqzhai-backend/services/payment/internal/domain"
)

// CardRepo persists saved cards.
type CardRepo interface {
	Create(ctx context.Context, c *domain.Card) (*domain.Card, error)
	Find(ctx context.Context, id string) (*domain.Card, error)
	ListByUser(ctx context.Context, userID string) ([]*domain.Card, error)
	Delete(ctx context.Context, id string) error
	SetDefault(ctx context.Context, userID, id string) error
}

// PaymentRepo persists payment attempts.
type PaymentRepo interface {
	Create(ctx context.Context, p *domain.Payment) (*domain.Payment, error)
	Find(ctx context.Context, id string) (*domain.Payment, error)
	FindByBooking(ctx context.Context, bookingID string) (*domain.Payment, error)
	ListByUser(ctx context.Context, userID string) ([]*domain.Payment, error)
	UpdateStatus(ctx context.Context, id string, status domain.PaymentStatus) error
}

// Gateway runs the actual PSP-side capture.
type Gateway interface {
	Charge(ctx context.Context, in ChargeInput) (string, error)
}

// ChargeInput is the payload sent to a PSP gateway adapter.
type ChargeInput struct {
	OrderID  string
	Amount   int64
	Currency string
	Holder   string
	Last4    string
}

// CoreClient lets payment notify core when a booking has been paid.
type CoreClient interface {
	MarkBookingPaid(ctx context.Context, bookingID, paymentID string) error
}

// Clock abstracts the current time.
type Clock interface{ Now() time.Time }

// IDGen generates new opaque entity IDs.
type IDGen interface{ New() string }
