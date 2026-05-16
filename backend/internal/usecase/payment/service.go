// Package payment implements the booking payment use case (PayBox integration).
package payment

import (
	"context"
	"fmt"

	"qonaqzhai-backend/internal/domain"
	"qonaqzhai-backend/internal/usecase"
)

// Notifier emits payment-related notifications. Optional.
type Notifier interface {
	Enqueue(ctx context.Context, n *domain.Notification) error
}

// Deps bundles payment service collaborators.
type Deps struct {
	Bookings usecase.BookingRepo
	Vendors  usecase.VendorRepo // optional — used to notify the vendor on payment
	Users    usecase.UserRepo
	Gateway  usecase.PaymentGateway
	Notifier Notifier // optional
	BaseURL  string   // returned to PayBox as success / failure redirect base
}

// Service exposes payment creation and webhook handling.
type Service struct{ d Deps }

// New constructs a payment Service.
func New(d Deps) *Service { return &Service{d: d} }

// StartIntent initialises payment for a booking and returns the redirect URL.
func (s *Service) StartIntent(ctx context.Context, customerID, bookingID string) (string, error) {
	if s.d.Gateway == nil {
		return "", fmt.Errorf("payment gateway not configured: %w", domain.ErrInvalidInput)
	}
	b, err := s.d.Bookings.Find(ctx, bookingID)
	if err != nil {
		return "", err
	}
	if b.CustomerID != customerID {
		return "", domain.ErrForbidden
	}
	if b.Amount <= 0 {
		return "", fmt.Errorf("booking amount must be positive: %w", domain.ErrInvalidInput)
	}
	u, err := s.d.Users.FindByID(ctx, customerID)
	if err != nil {
		return "", err
	}
	intent := usecase.PaymentIntent{
		OrderID:       b.ID,
		Amount:        b.Amount,
		Currency:      "KZT",
		Description:   fmt.Sprintf("Qonaqzhai booking %s", b.ID),
		CustomerEmail: u.Email,
		SuccessURL:    s.d.BaseURL + "/payments/success?booking=" + b.ID,
		FailureURL:    s.d.BaseURL + "/payments/failure?booking=" + b.ID,
	}
	res, err := s.d.Gateway.CreatePayment(ctx, intent)
	if err != nil {
		return "", err
	}
	if err := s.d.Bookings.SetPayment(ctx, b.ID, res.TransactionID); err != nil {
		return "", err
	}
	return res.URL, nil
}

// HandleCallback verifies a PayBox callback and transitions the booking.
func (s *Service) HandleCallback(ctx context.Context, form map[string]string) error {
	if s.d.Gateway == nil {
		return domain.ErrInvalidInput
	}
	result, err := s.d.Gateway.VerifyCallback(form)
	if err != nil {
		return err
	}
	if !result.Success {
		return nil // silent ack
	}
	b, err := s.d.Bookings.Find(ctx, result.OrderID)
	if err != nil {
		return err
	}
	if result.Amount != b.Amount {
		return fmt.Errorf("amount mismatch: paid=%d expected=%d: %w",
			result.Amount, b.Amount, domain.ErrConflict)
	}
	if err := s.d.Bookings.SetPayment(ctx, b.ID, result.TransactionID); err != nil {
		return err
	}
	if err := s.d.Bookings.UpdateStatus(ctx, b.ID, domain.BookingPaid); err != nil {
		return err
	}
	s.notifyPaid(ctx, b)
	return nil
}

func (s *Service) notifyPaid(ctx context.Context, b *domain.Booking) {
	if s.d.Notifier == nil {
		return
	}
	_ = s.d.Notifier.Enqueue(ctx, &domain.Notification{
		UserID:  b.CustomerID,
		Type:    domain.NotifBookingPaid,
		Channel: domain.ChannelBoth,
		Title:   "Payment confirmed",
		Body:    "<p>Your booking is now paid. The vendor has been notified.</p>",
	})
	if s.d.Vendors != nil {
		if v, err := s.d.Vendors.FindByID(ctx, b.VendorID); err == nil {
			_ = s.d.Notifier.Enqueue(ctx, &domain.Notification{
				UserID:  v.UserID,
				Type:    domain.NotifBookingPaid,
				Channel: domain.ChannelBoth,
				Title:   "Booking paid",
				Body:    "<p>Customer has paid for the upcoming booking.</p>",
			})
		}
	}
}
