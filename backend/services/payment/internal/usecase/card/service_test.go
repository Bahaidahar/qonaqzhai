package card_test

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"testing"
	"time"

	"qonaqzhai-backend/pkg/errs"
	"qonaqzhai-backend/services/payment/internal/domain"
	"qonaqzhai-backend/services/payment/internal/usecase/card"
)

type memCards struct {
	mu   sync.Mutex
	rows []*domain.Card
	seq  int
}

func newCards() *memCards { return &memCards{rows: []*domain.Card{}} }

func (m *memCards) Create(_ context.Context, c *domain.Card) (*domain.Card, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.seq++
	c.ID = "c" + strconv.Itoa(m.seq)
	// first card auto-default
	hasDefault := false
	for _, e := range m.rows {
		if e.UserID == c.UserID {
			hasDefault = hasDefault || e.IsDefault
		}
	}
	if !hasDefault {
		c.IsDefault = true
	}
	cp := *c
	m.rows = append(m.rows, &cp)
	return &cp, nil
}
func (m *memCards) Find(_ context.Context, id string) (*domain.Card, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, c := range m.rows {
		if c.ID == id {
			cp := *c
			return &cp, nil
		}
	}
	return nil, errs.ErrNotFound
}
func (m *memCards) ListByUser(_ context.Context, uid string) ([]*domain.Card, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := []*domain.Card{}
	for _, c := range m.rows {
		if c.UserID == uid {
			cp := *c
			out = append(out, &cp)
		}
	}
	return out, nil
}
func (m *memCards) Delete(_ context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, c := range m.rows {
		if c.ID == id {
			wasDefault, uid := c.IsDefault, c.UserID
			m.rows = append(m.rows[:i], m.rows[i+1:]...)
			if wasDefault {
				for _, e := range m.rows {
					if e.UserID == uid {
						e.IsDefault = true
						break
					}
				}
			}
			return nil
		}
	}
	return errs.ErrNotFound
}
func (m *memCards) SetDefault(_ context.Context, uid, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	found := false
	for _, c := range m.rows {
		if c.UserID != uid {
			continue
		}
		if c.ID == id {
			c.IsDefault = true
			found = true
		} else {
			c.IsDefault = false
		}
	}
	if !found {
		return errs.ErrNotFound
	}
	return nil
}

type fixedClock struct{ t time.Time }

func (c *fixedClock) Now() time.Time { return c.t }

func newSvc() (*card.Service, *memCards) {
	cards := newCards()
	svc := card.New(card.Deps{
		Cards: cards,
		Clock: &fixedClock{t: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)},
	})
	return svc, cards
}

func validInput() domain.CardInput {
	return domain.CardInput{Number: "4111 1111 1111 1111", ExpMonth: 6, ExpYear: 30, Holder: "TEST"}
}

func TestAdd_HappyPath(t *testing.T) {
	svc, _ := newSvc()
	c, err := svc.Add(context.Background(), "u", validInput())
	if err != nil {
		t.Fatal(err)
	}
	if c.Last4 != "1111" || c.Brand != "visa" {
		t.Fatalf("brand/last4 wrong: %+v", c)
	}
	if !c.IsDefault {
		t.Fatal("first card must be default")
	}
}

func TestAdd_RejectsExpired(t *testing.T) {
	svc, _ := newSvc()
	in := validInput()
	in.ExpYear = 20
	_, err := svc.Add(context.Background(), "u", in)
	if !errors.Is(err, errs.ErrInvalidInput) {
		t.Fatalf("want invalid input, got %v", err)
	}
}

func TestAdd_RejectsNonDigits(t *testing.T) {
	svc, _ := newSvc()
	in := validInput()
	in.Number = "4111-1111-1111-XYZA"
	_, err := svc.Add(context.Background(), "u", in)
	if !errors.Is(err, errs.ErrInvalidInput) {
		t.Fatalf("want invalid input, got %v", err)
	}
}

func TestDelete_AuthorisesOwner(t *testing.T) {
	svc, _ := newSvc()
	c, _ := svc.Add(context.Background(), "owner", validInput())
	if err := svc.Delete(context.Background(), "stranger", c.ID); !errors.Is(err, errs.ErrForbidden) {
		t.Fatalf("stranger must be forbidden, got %v", err)
	}
	if err := svc.Delete(context.Background(), "owner", c.ID); err != nil {
		t.Fatalf("owner delete failed: %v", err)
	}
}

func TestDelete_PromotesNextDefault(t *testing.T) {
	svc, cards := newSvc()
	a, _ := svc.Add(context.Background(), "owner", validInput())
	in := validInput()
	in.Number = "5500000000000004"
	b, _ := svc.Add(context.Background(), "owner", in)

	if err := svc.Delete(context.Background(), "owner", a.ID); err != nil {
		t.Fatal(err)
	}
	got, err := cards.Find(context.Background(), b.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !got.IsDefault {
		t.Fatalf("after default deletion the survivor must be promoted")
	}
}

func TestSetDefault(t *testing.T) {
	svc, _ := newSvc()
	a, _ := svc.Add(context.Background(), "owner", validInput())
	in := validInput()
	in.Number = "5500000000000004"
	b, _ := svc.Add(context.Background(), "owner", in)
	if err := svc.SetDefault(context.Background(), "owner", b.ID); err != nil {
		t.Fatal(err)
	}
	list, _ := svc.List(context.Background(), "owner")
	for _, c := range list {
		if c.ID == b.ID && !c.IsDefault {
			t.Fatalf("b should be default")
		}
		if c.ID == a.ID && c.IsDefault {
			t.Fatalf("a should no longer be default")
		}
	}
}

func TestSetDefault_StrangerForbidden(t *testing.T) {
	svc, _ := newSvc()
	c, _ := svc.Add(context.Background(), "owner", validInput())
	if err := svc.SetDefault(context.Background(), "stranger", c.ID); !errors.Is(err, errs.ErrForbidden) {
		t.Fatalf("want forbidden, got %v", err)
	}
}
