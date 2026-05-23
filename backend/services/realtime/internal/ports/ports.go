// Package ports defines the interfaces the realtime usecases depend on.
package ports

import (
	"context"

	"qonaqzhai-backend/services/realtime/internal/domain"
)

// ThreadRepo persists DM threads + messages.
type ThreadRepo interface {
	EnsureForBooking(ctx context.Context, bookingID, customerID, vendorID string) (*domain.Thread, error)
	FindByID(ctx context.Context, id string) (*domain.Thread, error)
	ListForUser(ctx context.Context, userID string) ([]*domain.Thread, error)
	AddMessage(ctx context.Context, m *domain.Message) (*domain.Message, error)
	ListMessages(ctx context.Context, threadID string) ([]*domain.Message, error)
}

// AuthClient lets realtime enrich thread summaries with user data.
type AuthClient interface {
	GetUsersBatch(ctx context.Context, userIDs []string) ([]*ExternalUser, error)
}

// ExternalUser is the subset of auth-svc User we need here.
type ExternalUser struct {
	ID    string
	Email string
	Name  string
}

// Publisher pushes events to connected WS clients. Implemented by the hub.
type Publisher interface {
	Publish(event string, payloadJSON []byte, userIDs ...string)
}

// IDGen generates new opaque entity IDs.
type IDGen interface{ New() string }
