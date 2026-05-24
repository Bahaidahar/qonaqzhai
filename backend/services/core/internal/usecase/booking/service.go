// Package booking implements the booking lifecycle.
package booking

import (
	"context"
	"encoding/json"
	"fmt"

	"qonaqzhai-backend/pkg/errs"
	"qonaqzhai-backend/services/core/internal/domain"
	"qonaqzhai-backend/services/core/internal/ports"
)

// Deps bundles booking usecase collaborators.
type Deps struct {
	Bookings      ports.BookingRepo
	Vendors       ports.VendorRepo
	Notifications ports.NotificationRepo
	Payments      ports.PaymentClient
	Realtime      ports.RealtimeClient
}

// Service exposes booking operations.
type Service struct{ d Deps }

// New constructs a booking Service.
func New(d Deps) *Service { return &Service{d: d} }

// CreateInput captures fields supplied by the customer.
type CreateInput struct {
	VendorID   string
	ServiceID  string
	EventDate  string
	GuestCount int
	Note       string
	Amount     int64
}

// Create makes a pending booking against vendorID.
func (s *Service) Create(ctx context.Context, customerID string, in CreateInput) (*domain.Booking, error) {
	if in.VendorID == "" || in.EventDate == "" {
		return nil, fmt.Errorf("vendor_id + event_date required: %w", errs.ErrInvalidInput)
	}
	v, err := s.d.Vendors.FindByID(ctx, in.VendorID)
	if err != nil {
		return nil, err
	}
	if !v.IsPublic() {
		return nil, fmt.Errorf("vendor not bookable: %w", errs.ErrForbidden)
	}
	b := &domain.Booking{
		CustomerID: customerID,
		VendorID:   in.VendorID,
		ServiceID:  in.ServiceID,
		EventDate:  in.EventDate,
		GuestCount: in.GuestCount,
		Note:       in.Note,
		Amount:     in.Amount,
		Status:     domain.BookingPending,
	}
	created, err := s.d.Bookings.Create(ctx, b)
	if err != nil {
		return nil, err
	}
	s.notify(ctx, v.UserID, "booking.created", "New booking request", in.Note)
	return created, nil
}

// VendorTransition is called by the vendor to accept / decline / complete.
func (s *Service) VendorTransition(ctx context.Context, vendorUserID, bookingID string, next domain.BookingStatus) (*domain.Booking, error) {
	b, err := s.authorizeVendor(ctx, vendorUserID, bookingID)
	if err != nil {
		return nil, err
	}
	if !b.VendorMayTransition(next) {
		return nil, fmt.Errorf("%s -> %s: %w", b.Status, next, errs.ErrConflict)
	}
	if err := s.d.Bookings.UpdateStatus(ctx, b.ID, next); err != nil {
		return nil, err
	}
	if next == domain.BookingAccepted && s.d.Realtime != nil {
		v, _ := s.d.Vendors.FindByID(ctx, b.VendorID)
		if v != nil {
			_ = s.d.Realtime.EnsureThread(ctx, b.ID, b.CustomerID, v.UserID)
		}
	}
	s.notify(ctx, b.CustomerID, "booking."+string(next), "Booking "+string(next), "")
	return s.d.Bookings.Find(ctx, b.ID)
}

// CustomerCancel lets the customer cancel pending/accepted bookings.
func (s *Service) CustomerCancel(ctx context.Context, customerID, bookingID string) (*domain.Booking, error) {
	b, err := s.d.Bookings.Find(ctx, bookingID)
	if err != nil {
		return nil, err
	}
	if b.CustomerID != customerID {
		return nil, errs.ErrForbidden
	}
	if !b.CustomerMayTransition(domain.BookingCancelled) {
		return nil, errs.ErrConflict
	}
	if err := s.d.Bookings.UpdateStatus(ctx, b.ID, domain.BookingCancelled); err != nil {
		return nil, err
	}
	v, _ := s.d.Vendors.FindByID(ctx, b.VendorID)
	if v != nil {
		s.notify(ctx, v.UserID, "booking.cancelled", "Booking cancelled", "")
	}
	return s.d.Bookings.Find(ctx, b.ID)
}

// Pay charges a saved card on payment-svc and, on success, flips the booking
// to paid in a single atomic UPDATE (see BookingRepo.MarkPaid). The saga is
// still not strictly two-phase: if the capture succeeds but MarkPaid fails the
// caller sees an error and the customer keeps the captured payment. Customer
// retries Pay → payment.Charge returns ErrAlreadyExists (unique on booking_id)
// → we recover the existing payment id via Payments.Charge response and finish
// the local commit.
func (s *Service) Pay(ctx context.Context, customerID, bookingID, cardID, currency string) (*domain.Booking, error) {
	b, err := s.d.Bookings.Find(ctx, bookingID)
	if err != nil {
		return nil, err
	}
	if b.CustomerID != customerID {
		return nil, errs.ErrForbidden
	}
	if b.Status != domain.BookingAccepted {
		return nil, fmt.Errorf("status %s: %w", b.Status, errs.ErrConflict)
	}
	if s.d.Payments == nil {
		return nil, fmt.Errorf("payments unavailable: %w", errs.ErrUpstream)
	}
	if currency == "" {
		currency = "KZT"
	}
	res, err := s.d.Payments.Charge(ctx, ports.ChargeRequest{
		BookingID: b.ID, UserID: b.CustomerID, CardID: cardID,
		Amount: b.Amount, Currency: currency,
	})
	if err != nil {
		return nil, err
	}
	if res.Status != "captured" {
		return nil, fmt.Errorf("payment %s: %w", res.Status, errs.ErrConflict)
	}
	if err := s.d.Bookings.MarkPaid(ctx, b.ID, res.ID); err != nil {
		return nil, fmt.Errorf("mark paid (payment %s captured): %w", res.ID, err)
	}
	v, _ := s.d.Vendors.FindByID(ctx, b.VendorID)
	if v != nil {
		s.notify(ctx, v.UserID, "booking.paid", "Booking paid", "")
	}
	return s.d.Bookings.Find(ctx, b.ID)
}

// MarkPaid is the gRPC entrypoint payment-svc uses to push back a webhook-style
// confirmation when the customer paid via redirect rather than direct Charge.
// One atomic UPDATE.
func (s *Service) MarkPaid(ctx context.Context, bookingID, paymentID string) (*domain.Booking, error) {
	if err := s.d.Bookings.MarkPaid(ctx, bookingID, paymentID); err != nil {
		return nil, err
	}
	return s.d.Bookings.Find(ctx, bookingID)
}

// Find returns a booking visible to the caller (customer or owning vendor).
func (s *Service) Find(ctx context.Context, callerID, bookingID string) (*domain.Booking, error) {
	b, err := s.d.Bookings.Find(ctx, bookingID)
	if err != nil {
		return nil, err
	}
	if b.CustomerID == callerID {
		return b, nil
	}
	v, err := s.d.Vendors.FindByID(ctx, b.VendorID)
	if err != nil {
		return nil, err
	}
	if v.UserID != callerID {
		return nil, errs.ErrForbidden
	}
	return b, nil
}

// GetRaw is the gRPC entrypoint — no caller check.
func (s *Service) GetRaw(ctx context.Context, id string) (*domain.Booking, error) {
	return s.d.Bookings.Find(ctx, id)
}

// IsAccepted answers the realtime-svc gRPC query for thread-creation.
func (s *Service) IsAccepted(ctx context.Context, id string) (*domain.Booking, string, bool, error) {
	b, err := s.d.Bookings.Find(ctx, id)
	if err != nil {
		return nil, "", false, err
	}
	if b.Status != domain.BookingAccepted && b.Status != domain.BookingPaid && b.Status != domain.BookingCompleted {
		return b, "", false, nil
	}
	v, err := s.d.Vendors.FindByID(ctx, b.VendorID)
	if err != nil {
		return b, "", true, nil
	}
	return b, v.UserID, true, nil
}

// ListForCustomer returns paginated bookings made by customer.
func (s *Service) ListForCustomer(ctx context.Context, customerID string, p ports.Page) ([]*domain.Booking, error) {
	return s.d.Bookings.ListForCustomer(ctx, customerID, p)
}

// ListForVendor returns paginated bookings against vendorID.
func (s *Service) ListForVendor(ctx context.Context, vendorUserID string, p ports.Page) ([]*domain.Booking, error) {
	v, err := s.d.Vendors.FindByUserID(ctx, vendorUserID)
	if err != nil {
		return nil, err
	}
	return s.d.Bookings.ListForVendor(ctx, v.ID, p)
}

// Stats returns aggregate booking counts + GMV.
func (s *Service) Stats(ctx context.Context) (ports.BookingStats, error) {
	return s.d.Bookings.Stats(ctx)
}

func (s *Service) authorizeVendor(ctx context.Context, vendorUserID, bookingID string) (*domain.Booking, error) {
	b, err := s.d.Bookings.Find(ctx, bookingID)
	if err != nil {
		return nil, err
	}
	v, err := s.d.Vendors.FindByID(ctx, b.VendorID)
	if err != nil {
		return nil, err
	}
	if v.UserID != vendorUserID {
		return nil, errs.ErrForbidden
	}
	return b, nil
}

func (s *Service) notify(ctx context.Context, userID, kind, title, body string) {
	if s.d.Notifications == nil {
		return
	}
	n := &domain.Notification{
		UserID: userID, Type: kind, Channel: domain.ChannelPush,
		Title: title, Body: body,
	}
	if _, err := s.d.Notifications.Enqueue(ctx, n); err != nil {
		return
	}
	if s.d.Realtime != nil {
		payload, _ := json.Marshal(map[string]string{"type": kind, "title": title, "body": body})
		_ = s.d.Realtime.Publish(ctx, kind, payload, userID)
	}
}
