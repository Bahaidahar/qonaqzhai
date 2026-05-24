// Package notification exposes notification reads + FCM token registration.
package notification

import (
	"context"

	"qonaqzhai-backend/services/core/internal/domain"
	"qonaqzhai-backend/services/core/internal/ports"
)

// Deps bundles notification collaborators.
type Deps struct {
	Notifications ports.NotificationRepo
	FCMTokens     ports.FCMTokenRepo
}

// Service exposes notification ops.
type Service struct{ d Deps }

// New constructs a notification Service.
func New(d Deps) *Service { return &Service{d: d} }

// ListForUser returns the paginated notifications for userID.
func (s *Service) ListForUser(ctx context.Context, userID string, p ports.Page) ([]*domain.Notification, error) {
	return s.d.Notifications.ListForUser(ctx, userID, p)
}

// RegisterToken upserts an FCM token bound to userID.
func (s *Service) RegisterToken(ctx context.Context, userID, token, platform string) error {
	return s.d.FCMTokens.Upsert(ctx, &domain.FCMToken{
		UserID: userID, Token: token, Platform: platform,
	})
}

// UnregisterToken removes a token.
func (s *Service) UnregisterToken(ctx context.Context, token string) error {
	return s.d.FCMTokens.Delete(ctx, token)
}
