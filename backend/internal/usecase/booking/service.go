// Package booking implements booking creation and lifecycle transitions.
package booking

import (
	"context"
	"fmt"
	"strings"

	"qonaqzhai-backend/internal/domain"
	"qonaqzhai-backend/internal/usecase"
)

// Notifier emits booking-related notifications. Optional collaborator;
// usecase tolerates nil for unit tests + when notifications are disabled.
type Notifier interface {
	Enqueue(ctx context.Context, n *domain.Notification) error
}

// Deps bundles booking service collaborators.
type Deps struct {
	Bookings usecase.BookingRepo
	Vendors  usecase.VendorRepo
	Clock    usecase.Clock
	Notifier Notifier // optional
}

// Service exposes booking operations.
type Service struct{ d Deps }

// New constructs a booking Service.
func New(d Deps) *Service { return &Service{d: d} }

// CreateInput captures user-supplied fields for a new booking.
type CreateInput struct {
	VendorID   string
	EventDate  string
	GuestCount int
	Note       string
	Amount     int64
}

// Create validates input, verifies vendor existence and inserts a pending booking.
func (s *Service) Create(ctx context.Context, customerID string, in CreateInput) (*domain.Booking, error) {
	in.VendorID = strings.TrimSpace(in.VendorID)
	in.EventDate = strings.TrimSpace(in.EventDate)
	if in.VendorID == "" || in.EventDate == "" {
		return nil, fmt.Errorf("vendorId and eventDate required: %w", domain.ErrInvalidInput)
	}
	if in.GuestCount < 0 {
		return nil, fmt.Errorf("guestCount: %w", domain.ErrInvalidInput)
	}
	if _, err := s.d.Vendors.FindByID(ctx, in.VendorID); err != nil {
		return nil, err
	}
	b := &domain.Booking{
		CustomerID: customerID,
		VendorID:   in.VendorID,
		EventDate:  in.EventDate,
		GuestCount: in.GuestCount,
		Note:       strings.TrimSpace(in.Note),
		Status:     domain.BookingPending,
		Amount:     in.Amount,
		CreatedAt:  s.d.Clock.Now(),
	}
	created, err := s.d.Bookings.Create(ctx, b)
	if err != nil {
		return nil, err
	}
	if v, err := s.d.Vendors.FindByID(ctx, created.VendorID); err == nil && s.d.Notifier != nil {
		_ = s.d.Notifier.Enqueue(ctx, &domain.Notification{
			UserID:  v.UserID,
			Type:    domain.NotifBookingCreated,
			Channel: domain.ChannelBoth,
			Title:   "New booking request",
			Body:    "<p>You have a new booking request — please review it in your inbox.</p>",
		})
	}
	return created, nil
}

// Find returns a booking by id.
func (s *Service) Find(ctx context.Context, id string) (*domain.Booking, error) {
	return s.d.Bookings.Find(ctx, id)
}

// VendorTransition moves the booking to next, enforcing vendor ownership.
func (s *Service) VendorTransition(ctx context.Context, vendorUserID, bookingID string, next domain.BookingStatus) (*domain.Booking, error) {
	b, err := s.d.Bookings.Find(ctx, bookingID)
	if err != nil {
		return nil, err
	}
	v, err := s.d.Vendors.FindByUserID(ctx, vendorUserID)
	if err != nil {
		return nil, domain.ErrForbidden
	}
	if v.ID != b.VendorID {
		return nil, domain.ErrForbidden
	}
	if !b.VendorMayTransition(next) {
		return nil, fmt.Errorf("invalid vendor transition %s→%s: %w", b.Status, next, domain.ErrInvalidInput)
	}
	if err := s.d.Bookings.UpdateStatus(ctx, b.ID, next); err != nil {
		return nil, err
	}
	s.notifyVendorTransition(ctx, b, next)
	return s.d.Bookings.Find(ctx, b.ID)
}

func (s *Service) notifyVendorTransition(ctx context.Context, b *domain.Booking, next domain.BookingStatus) {
	if s.d.Notifier == nil {
		return
	}
	var notifType domain.NotificationType
	var title, body string
	switch next {
	case domain.BookingAccepted:
		notifType = domain.NotifBookingAccepted
		title = "Booking accepted"
		body = "<p>Your booking request was accepted by the vendor.</p>"
	case domain.BookingDeclined:
		notifType = domain.NotifBookingDeclined
		title = "Booking declined"
		body = "<p>Your booking request was declined.</p>"
	default:
		return
	}
	_ = s.d.Notifier.Enqueue(ctx, &domain.Notification{
		UserID:  b.CustomerID,
		Type:    notifType,
		Channel: domain.ChannelBoth,
		Title:   title,
		Body:    body,
	})
}

// CustomerTransition moves the booking to next, enforcing ownership.
func (s *Service) CustomerTransition(ctx context.Context, customerID, bookingID string, next domain.BookingStatus) (*domain.Booking, error) {
	b, err := s.d.Bookings.Find(ctx, bookingID)
	if err != nil {
		return nil, err
	}
	if b.CustomerID != customerID {
		return nil, domain.ErrForbidden
	}
	if !b.CustomerMayTransition(next) {
		return nil, fmt.Errorf("invalid customer transition %s→%s: %w", b.Status, next, domain.ErrInvalidInput)
	}
	if err := s.d.Bookings.UpdateStatus(ctx, b.ID, next); err != nil {
		return nil, err
	}
	return s.d.Bookings.Find(ctx, b.ID)
}

// AdminTransition forces any valid status — moderation override.
func (s *Service) AdminTransition(ctx context.Context, bookingID string, next domain.BookingStatus) (*domain.Booking, error) {
	if !next.Valid() {
		return nil, domain.ErrInvalidInput
	}
	if err := s.d.Bookings.UpdateStatus(ctx, bookingID, next); err != nil {
		return nil, err
	}
	return s.d.Bookings.Find(ctx, bookingID)
}

// ListForCustomer returns the bookings made by customerID.
func (s *Service) ListForCustomer(ctx context.Context, customerID string) ([]*domain.Booking, error) {
	return s.d.Bookings.ListForCustomer(ctx, customerID)
}

// ListForVendor returns bookings against vendor's profile.
// Returns empty slice (not error) if the vendor has no profile yet.
func (s *Service) ListForVendor(ctx context.Context, vendorUserID string) ([]*domain.Booking, error) {
	v, err := s.d.Vendors.FindByUserID(ctx, vendorUserID)
	if err != nil {
		return []*domain.Booking{}, nil
	}
	return s.d.Bookings.ListForVendor(ctx, v.ID)
}

// ListAll returns every booking (admin only).
func (s *Service) ListAll(ctx context.Context) ([]*domain.Booking, error) {
	return s.d.Bookings.ListAll(ctx)
}
