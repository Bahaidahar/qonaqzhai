package card_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"qonaqzhai-backend/internal/domain"
	"qonaqzhai-backend/internal/usecase/card"
	"qonaqzhai-backend/internal/usecase/inmem"
)

func newSvc(t *testing.T) *card.Service {
	t.Helper()
	i := 0
	gen := func() string {
		i++
		return "c" + string(rune('0'+i))
	}
	repo := inmem.NewCardRepo(gen)
	return card.New(card.Deps{
		Cards: repo,
		Now:   func() time.Time { return time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC) },
	})
}

func TestAddCard_FirstIsDefault(t *testing.T) {
	s := newSvc(t)
	c, err := s.Add(context.Background(), "u1", domain.CardInput{
		Number: "4242 4242 4242 4242", ExpMonth: 12, ExpYear: 2027, Holder: "TEST USER",
	}, false)
	if err != nil {
		t.Fatalf("add: %v", err)
	}
	if c.Last4 != "4242" || c.Brand != "visa" || !c.IsDefault {
		t.Fatalf("bad card: %+v", c)
	}
}

func TestAddCard_RejectsExpired(t *testing.T) {
	s := newSvc(t)
	_, err := s.Add(context.Background(), "u1", domain.CardInput{
		Number: "4242424242424242", ExpMonth: 12, ExpYear: 2024,
	}, false)
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("want ErrInvalidInput, got %v", err)
	}
}

func TestAddCard_RejectsBadPAN(t *testing.T) {
	s := newSvc(t)
	_, err := s.Add(context.Background(), "u1", domain.CardInput{
		Number: "abcd", ExpMonth: 12, ExpYear: 2027,
	}, false)
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("want ErrInvalidInput, got %v", err)
	}
}

func TestSetDefault_OnlyOwner(t *testing.T) {
	s := newSvc(t)
	c, _ := s.Add(context.Background(), "u1", domain.CardInput{
		Number: "5555555555554444", ExpMonth: 6, ExpYear: 2030,
	}, false)
	if err := s.SetDefault(context.Background(), "u2", c.ID); !errors.Is(err, domain.ErrForbidden) {
		t.Fatalf("want forbidden, got %v", err)
	}
}

func TestDelete_PromotesNewDefault(t *testing.T) {
	s := newSvc(t)
	c1, _ := s.Add(context.Background(), "u1", domain.CardInput{Number: "4242424242424242", ExpMonth: 1, ExpYear: 2030}, false)
	c2, _ := s.Add(context.Background(), "u1", domain.CardInput{Number: "5555555555554444", ExpMonth: 1, ExpYear: 2030}, true)
	if !c2.IsDefault {
		t.Fatalf("c2 should be default")
	}
	if err := s.Delete(context.Background(), "u1", c2.ID); err != nil {
		t.Fatal(err)
	}
	def, err := s.DefaultFor(context.Background(), "u1")
	if err != nil {
		t.Fatal(err)
	}
	if def.ID != c1.ID {
		t.Fatalf("want %s, got %s", c1.ID, def.ID)
	}
}

func TestDefaultFor_NoCards(t *testing.T) {
	s := newSvc(t)
	_, err := s.DefaultFor(context.Background(), "u1")
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("want ErrNotFound, got %v", err)
	}
}
