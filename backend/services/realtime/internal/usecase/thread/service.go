// Package thread implements DM thread + message use cases.
package thread

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"qonaqzhai-backend/pkg/errs"

	"qonaqzhai-backend/services/realtime/internal/domain"
	"qonaqzhai-backend/services/realtime/internal/ports"
)

// Deps bundles thread collaborators.
type Deps struct {
	Threads   ports.ThreadRepo
	Auth      ports.AuthClient  // optional — enrich summaries with names
	Publisher ports.Publisher   // ws hub
}

// Service exposes thread operations.
type Service struct{ d Deps }

// New constructs a thread Service.
func New(d Deps) *Service { return &Service{d: d} }

// Ensure idempotently creates a thread for a booking. Called by core-svc gRPC
// when a booking moves to accepted.
func (s *Service) Ensure(ctx context.Context, bookingID, customerID, vendorUserID string) (*domain.Thread, error) {
	if bookingID == "" || customerID == "" || vendorUserID == "" {
		return nil, fmt.Errorf("ensure: %w", errs.ErrInvalidInput)
	}
	return s.d.Threads.EnsureForBooking(ctx, bookingID, customerID, vendorUserID)
}

// Get returns thread + messages, enforcing membership.
func (s *Service) Get(ctx context.Context, userID, threadID string) (*domain.Thread, []*domain.Message, error) {
	t, err := s.d.Threads.FindByID(ctx, threadID)
	if err != nil {
		return nil, nil, err
	}
	if t.CustomerID != userID && t.VendorID != userID {
		return nil, nil, errs.ErrForbidden
	}
	msgs, err := s.d.Threads.ListMessages(ctx, threadID)
	if err != nil {
		return nil, nil, err
	}
	return t, msgs, nil
}

// Send appends a message, enforces membership, and publishes to both peers.
func (s *Service) Send(ctx context.Context, userID, threadID, text string) (*domain.Message, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, fmt.Errorf("empty text: %w", errs.ErrInvalidInput)
	}
	t, err := s.d.Threads.FindByID(ctx, threadID)
	if err != nil {
		return nil, err
	}
	if t.CustomerID != userID && t.VendorID != userID {
		return nil, errs.ErrForbidden
	}
	m, err := s.d.Threads.AddMessage(ctx, &domain.Message{
		ThreadID: t.ID, SenderID: userID, Text: text,
	})
	if err != nil {
		return nil, err
	}
	if s.d.Publisher != nil {
		payload, _ := json.Marshal(m)
		s.d.Publisher.Publish("thread.message", payload, t.CustomerID, t.VendorID)
	}
	return m, nil
}

// Summary enriches a thread with counterpart info, batching auth lookups.
type Summary struct {
	Thread      *domain.Thread `json:"thread"`
	Counterpart string         `json:"counterpart"`
}

// ListSummaries returns the user's threads with the peer's display name.
func (s *Service) ListSummaries(ctx context.Context, userID string) ([]*Summary, error) {
	ts, err := s.d.Threads.ListForUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make([]*Summary, len(ts))
	peerIDs := make([]string, 0, len(ts))
	for i, t := range ts {
		peer := t.VendorID
		if userID == t.VendorID {
			peer = t.CustomerID
		}
		out[i] = &Summary{Thread: t}
		peerIDs = append(peerIDs, peer)
	}
	if s.d.Auth == nil || len(peerIDs) == 0 {
		return out, nil
	}
	users, err := s.d.Auth.GetUsersBatch(ctx, peerIDs)
	if err == nil {
		byID := make(map[string]*ports.ExternalUser, len(users))
		for _, u := range users {
			byID[u.ID] = u
		}
		for i, t := range ts {
			peer := t.VendorID
			if userID == t.VendorID {
				peer = t.CustomerID
			}
			if u, ok := byID[peer]; ok {
				out[i].Counterpart = u.Name
				if out[i].Counterpart == "" {
					out[i].Counterpart = u.Email
				}
			}
		}
	}
	return out, nil
}
