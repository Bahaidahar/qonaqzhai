// Package card implements saved-card (mock) use cases — list/add/delete/setDefault.
// Cards are mock-only. PANs are never persisted; only last4 + brand are stored.
package card

import (
	"context"
	"time"

	"qonaqzhai-backend/internal/domain"
	"qonaqzhai-backend/internal/usecase"
)

// Deps bundles collaborators for the card service.
type Deps struct {
	Cards usecase.CardRepo
	Now   func() time.Time
}

// Service exposes card management use cases.
type Service struct{ d Deps }

// New constructs a Service.
func New(d Deps) *Service {
	if d.Now == nil {
		d.Now = func() time.Time { return time.Now().UTC() }
	}
	return &Service{d: d}
}

// List returns the user's saved cards.
func (s *Service) List(ctx context.Context, userID string) ([]*domain.PaymentCard, error) {
	return s.d.Cards.ListForUser(ctx, userID)
}

// Add validates the input and stores a new card (last4 + brand only).
func (s *Service) Add(ctx context.Context, userID string, in domain.CardInput, makeDefault bool) (*domain.PaymentCard, error) {
	in.Normalize()
	if err := in.Validate(s.d.Now()); err != nil {
		return nil, err
	}
	year := in.ExpYear
	if year < 100 {
		year += 2000
	}
	c := &domain.PaymentCard{
		UserID:    userID,
		Brand:     domain.DetectBrand(in.Number),
		Last4:     in.Last4(),
		ExpMonth:  in.ExpMonth,
		ExpYear:   year,
		Holder:    in.Holder,
		IsDefault: makeDefault,
	}
	return s.d.Cards.Create(ctx, c)
}

// Delete removes a card after verifying ownership.
func (s *Service) Delete(ctx context.Context, userID, cardID string) error {
	c, err := s.d.Cards.FindByID(ctx, cardID)
	if err != nil {
		return err
	}
	if c.UserID != userID {
		return domain.ErrForbidden
	}
	return s.d.Cards.Delete(ctx, cardID)
}

// SetDefault marks a card as the user's default payment method.
func (s *Service) SetDefault(ctx context.Context, userID, cardID string) error {
	return s.d.Cards.SetDefault(ctx, userID, cardID)
}

// DefaultFor returns the user's default card (or first if none flagged).
// Returns ErrNotFound when the user has no cards.
func (s *Service) DefaultFor(ctx context.Context, userID string) (*domain.PaymentCard, error) {
	cards, err := s.d.Cards.ListForUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(cards) == 0 {
		return nil, domain.ErrNotFound
	}
	for _, c := range cards {
		if c.IsDefault {
			return c, nil
		}
	}
	return cards[0], nil
}
