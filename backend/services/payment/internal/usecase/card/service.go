// Package card implements saved card management.
package card

import (
	"context"
	"fmt"

	"qonaqzhai-backend/pkg/errs"
	"qonaqzhai-backend/services/payment/internal/domain"
	"qonaqzhai-backend/services/payment/internal/ports"
)

// Deps bundles card usecase collaborators.
type Deps struct {
	Cards ports.CardRepo
	Clock ports.Clock
}

// Service exposes card operations.
type Service struct{ d Deps }

// New constructs a card Service.
func New(d Deps) *Service { return &Service{d: d} }

// Add validates the input + stores last4/brand (PAN is dropped).
func (s *Service) Add(ctx context.Context, userID string, in domain.CardInput) (*domain.Card, error) {
	in.Normalize()
	if !in.Validate(s.d.Clock.Now()) {
		return nil, fmt.Errorf("card: %w", errs.ErrInvalidInput)
	}
	return s.d.Cards.Create(ctx, &domain.Card{
		UserID:   userID,
		Brand:    domain.DetectBrand(in.Number),
		Last4:    in.Last4(),
		ExpMonth: in.ExpMonth,
		ExpYear:  in.ExpYear,
		Holder:   in.Holder,
	})
}

// List returns the caller's cards.
func (s *Service) List(ctx context.Context, userID string) ([]*domain.Card, error) {
	return s.d.Cards.ListByUser(ctx, userID)
}

// Delete removes a card owned by callerID.
func (s *Service) Delete(ctx context.Context, callerID, id string) error {
	c, err := s.d.Cards.Find(ctx, id)
	if err != nil {
		return err
	}
	if c.UserID != callerID {
		return errs.ErrForbidden
	}
	return s.d.Cards.Delete(ctx, id)
}

// SetDefault flips the default flag for callerID's card.
func (s *Service) SetDefault(ctx context.Context, callerID, id string) error {
	c, err := s.d.Cards.Find(ctx, id)
	if err != nil {
		return err
	}
	if c.UserID != callerID {
		return errs.ErrForbidden
	}
	return s.d.Cards.SetDefault(ctx, callerID, id)
}

// Find returns a card by id without authz (gRPC entrypoint).
func (s *Service) Find(ctx context.Context, id string) (*domain.Card, error) {
	return s.d.Cards.Find(ctx, id)
}
