package inmem

import (
	"context"
	"sort"
	"sync"
	"time"

	"qonaqzhai-backend/internal/domain"
)

// CardRepo is an in-memory card store for tests.
type CardRepo struct {
	mu    sync.Mutex
	byID  map[string]*domain.PaymentCard
	idGen func() string
}

func NewCardRepo(idGen func() string) *CardRepo {
	return &CardRepo{byID: map[string]*domain.PaymentCard{}, idGen: idGen}
}

func (r *CardRepo) Create(_ context.Context, c *domain.PaymentCard) (*domain.PaymentCard, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *c
	if cp.ID == "" {
		cp.ID = r.idGen()
	}
	if cp.CreatedAt.IsZero() {
		cp.CreatedAt = time.Now().UTC()
	}
	count := 0
	for _, ex := range r.byID {
		if ex.UserID == cp.UserID {
			count++
		}
	}
	if count == 0 {
		cp.IsDefault = true
	}
	if cp.IsDefault {
		for _, ex := range r.byID {
			if ex.UserID == cp.UserID {
				ex.IsDefault = false
			}
		}
	}
	r.byID[cp.ID] = &cp
	out := cp
	return &out, nil
}

func (r *CardRepo) FindByID(_ context.Context, id string) (*domain.PaymentCard, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	c, ok := r.byID[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	cp := *c
	return &cp, nil
}

func (r *CardRepo) ListForUser(_ context.Context, userID string) ([]*domain.PaymentCard, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := []*domain.PaymentCard{}
	for _, c := range r.byID {
		if c.UserID == userID {
			cp := *c
			out = append(out, &cp)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].IsDefault != out[j].IsDefault {
			return out[i].IsDefault
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

func (r *CardRepo) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	c, ok := r.byID[id]
	if !ok {
		return domain.ErrNotFound
	}
	wasDefault := c.IsDefault
	userID := c.UserID
	delete(r.byID, id)
	if wasDefault {
		var newest *domain.PaymentCard
		for _, ex := range r.byID {
			if ex.UserID == userID && (newest == nil || ex.CreatedAt.After(newest.CreatedAt)) {
				newest = ex
			}
		}
		if newest != nil {
			newest.IsDefault = true
		}
	}
	return nil
}

func (r *CardRepo) SetDefault(_ context.Context, userID, cardID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	target, ok := r.byID[cardID]
	if !ok {
		return domain.ErrNotFound
	}
	if target.UserID != userID {
		return domain.ErrForbidden
	}
	for _, ex := range r.byID {
		if ex.UserID == userID {
			ex.IsDefault = false
		}
	}
	target.IsDefault = true
	return nil
}
