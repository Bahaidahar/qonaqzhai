// Package thread implements customer↔vendor DM threads scoped to one booking.
// A thread becomes available only after the vendor accepts the booking.
package thread

import (
	"context"
	"fmt"
	"strings"

	"qonaqzhai-backend/internal/domain"
	"qonaqzhai-backend/internal/usecase"
)

// Notifier emits per-message notifications. Optional.
type Notifier interface {
	Enqueue(ctx context.Context, n *domain.Notification) error
}

// Publisher broadcasts realtime events to connected users (e.g. WebSocket hub).
type Publisher interface {
	Publish(event string, payload any, userIDs ...string)
}

// Deps bundles thread service collaborators.
type Deps struct {
	Threads   usecase.ThreadRepo
	Bookings  usecase.BookingRepo
	Vendors   usecase.VendorRepo
	Users     usecase.UserRepo
	Notifier  Notifier  // optional
	Publisher Publisher // optional — realtime fan-out
}

// Summary enriches a thread with the linked booking + counterpart info, so the
// UI does not need to fan-out per-thread lookups to render the inbox.
type Summary struct {
	Thread      *domain.BookingThread `json:"thread"`
	BookingID   string                `json:"bookingId"`
	EventDate   string                `json:"eventDate"`
	GuestCount  int                   `json:"guestCount"`
	Amount      int64                 `json:"amount"`
	Status      string                `json:"status"`
	VendorName  string                `json:"vendorName"`
	Counterpart string                `json:"counterpart"`
}

// Service exposes thread operations.
type Service struct{ d Deps }

// New constructs a thread Service.
func New(d Deps) *Service { return &Service{d: d} }

// EnsureForBooking opens (or returns existing) a thread for an accepted booking.
// Called by the booking lifecycle when status transitions to accepted.
func (s *Service) EnsureForBooking(ctx context.Context, b *domain.Booking) (*domain.BookingThread, error) {
	v, err := s.d.Vendors.FindByID(ctx, b.VendorID)
	if err != nil {
		return nil, err
	}
	return s.d.Threads.CreateForBooking(ctx, b.ID, b.CustomerID, v.UserID)
}

// ListForUser returns threads visible to the user (own customer + vendor sides).
func (s *Service) ListForUser(ctx context.Context, userID string) ([]*domain.BookingThread, error) {
	return s.d.Threads.ListForUser(ctx, userID)
}

// ListSummariesForUser returns threads with denormalized booking + counterpart
// info. Missing related rows fall back to empty strings instead of failing the
// whole listing — best-effort enrichment.
func (s *Service) ListSummariesForUser(ctx context.Context, userID string) ([]*Summary, error) {
	threads, err := s.d.Threads.ListForUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make([]*Summary, 0, len(threads))
	for _, t := range threads {
		sum := &Summary{Thread: t}
		if b, err := s.d.Bookings.Find(ctx, t.BookingID); err == nil && b != nil {
			sum.BookingID = b.ID
			sum.EventDate = b.EventDate
			sum.GuestCount = b.GuestCount
			sum.Amount = b.Amount
			sum.Status = string(b.Status)
			if v, err := s.d.Vendors.FindByID(ctx, b.VendorID); err == nil && v != nil {
				sum.VendorName = v.Name
			}
		}
		counterpartUserID := t.VendorID
		if userID == t.VendorID {
			counterpartUserID = t.CustomerID
		}
		if s.d.Users != nil {
			if u, err := s.d.Users.FindByID(ctx, counterpartUserID); err == nil && u != nil {
				sum.Counterpart = u.Name
				if sum.Counterpart == "" {
					sum.Counterpart = u.Email
				}
			}
		}
		out = append(out, sum)
	}
	return out, nil
}

// Get returns thread + messages, enforcing membership.
func (s *Service) Get(ctx context.Context, userID, threadID string) (*domain.BookingThread, []*domain.ThreadMessage, error) {
	t, err := s.d.Threads.FindByID(ctx, threadID)
	if err != nil {
		return nil, nil, err
	}
	if t.CustomerID != userID && t.VendorID != userID {
		return nil, nil, domain.ErrForbidden
	}
	msgs, err := s.d.Threads.ListMessages(ctx, threadID)
	if err != nil {
		return nil, nil, err
	}
	return t, msgs, nil
}

// Send appends a message to the thread, enforcing membership.
func (s *Service) Send(ctx context.Context, userID, threadID, text string) (*domain.ThreadMessage, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, fmt.Errorf("empty text: %w", domain.ErrInvalidInput)
	}
	t, err := s.d.Threads.FindByID(ctx, threadID)
	if err != nil {
		return nil, err
	}
	if t.CustomerID != userID && t.VendorID != userID {
		return nil, domain.ErrForbidden
	}
	m, err := s.d.Threads.AddMessage(ctx, &domain.ThreadMessage{
		ThreadID: t.ID, SenderID: userID, Text: text,
	})
	if err != nil {
		return nil, err
	}
	// Realtime fan-out to both participants (sender + peer). The sender's other
	// devices benefit from the echo too.
	if s.d.Publisher != nil {
		s.d.Publisher.Publish("thread.message", m, t.CustomerID, t.VendorID)
	}
	s.notifyOther(ctx, t, userID, text)
	return m, nil
}

func (s *Service) notifyOther(ctx context.Context, t *domain.BookingThread, sender, text string) {
	if s.d.Notifier == nil {
		return
	}
	other := t.VendorID
	if sender == t.VendorID {
		other = t.CustomerID
	}
	preview := text
	if len(preview) > 80 {
		preview = preview[:80] + "…"
	}
	_ = s.d.Notifier.Enqueue(ctx, &domain.Notification{
		UserID:  other,
		Type:    "thread.message",
		Channel: domain.ChannelPush,
		Title:   "New message",
		Body:    preview,
	})
}
