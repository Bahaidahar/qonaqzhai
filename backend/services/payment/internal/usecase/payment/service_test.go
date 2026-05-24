package payment_test

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"testing"
	"time"

	"qonaqzhai-backend/pkg/errs"
	"qonaqzhai-backend/services/payment/internal/domain"
	"qonaqzhai-backend/services/payment/internal/ports"
	"qonaqzhai-backend/services/payment/internal/usecase/payment"
)

type memPayments struct {
	mu   sync.Mutex
	rows []*domain.Payment
	seq  int
}

func (m *memPayments) Create(_ context.Context, p *domain.Payment) (*domain.Payment, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, e := range m.rows {
		if e.BookingID == p.BookingID {
			return nil, errs.ErrAlreadyExists
		}
	}
	m.seq++
	p.ID = "p" + strconv.Itoa(m.seq)
	cp := *p
	m.rows = append(m.rows, &cp)
	return &cp, nil
}
func (m *memPayments) Find(_ context.Context, id string) (*domain.Payment, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, p := range m.rows {
		if p.ID == id {
			cp := *p
			return &cp, nil
		}
	}
	return nil, errs.ErrNotFound
}
func (m *memPayments) FindByBooking(_ context.Context, id string) (*domain.Payment, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, p := range m.rows {
		if p.BookingID == id {
			cp := *p
			return &cp, nil
		}
	}
	return nil, errs.ErrNotFound
}
func (m *memPayments) ListByUser(_ context.Context, uid string) ([]*domain.Payment, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := []*domain.Payment{}
	for _, p := range m.rows {
		if p.UserID == uid {
			cp := *p
			out = append(out, &cp)
		}
	}
	return out, nil
}
func (m *memPayments) UpdateStatus(_ context.Context, id string, s domain.PaymentStatus) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, p := range m.rows {
		if p.ID == id {
			p.Status = s
			return nil
		}
	}
	return errs.ErrNotFound
}

type memCards struct{ owner string }

func (m *memCards) Create(context.Context, *domain.Card) (*domain.Card, error) { return nil, nil }
func (m *memCards) Find(_ context.Context, id string) (*domain.Card, error) {
	return &domain.Card{ID: id, UserID: m.owner}, nil
}
func (m *memCards) ListByUser(context.Context, string) ([]*domain.Card, error) { return nil, nil }
func (m *memCards) Delete(context.Context, string) error                       { return nil }
func (m *memCards) SetDefault(context.Context, string, string) error           { return nil }

type alwaysOK struct{}

func (alwaysOK) Charge(context.Context, ports.ChargeInput) (string, error) { return "ref-ok", nil }

type alwaysFail struct{}

func (alwaysFail) Charge(context.Context, ports.ChargeInput) (string, error) {
	return "", errors.New("gateway declined")
}

type stubCore struct{ marked int; mu sync.Mutex }

func (s *stubCore) MarkBookingPaid(context.Context, string, string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.marked++
	return nil
}

type fixedClock struct{}

func (fixedClock) Now() time.Time { return time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC) }

func newSvc(gw ports.Gateway, core ports.CoreClient) (*payment.Service, *memPayments) {
	pays := &memPayments{}
	svc := payment.New(payment.Deps{
		Payments: pays, Cards: &memCards{owner: "u"}, Gateway: gw, Core: core, Clock: fixedClock{},
	})
	return svc, pays
}

func TestCharge_HappyPath_NotifiesCore(t *testing.T) {
	core := &stubCore{}
	svc, _ := newSvc(alwaysOK{}, core)
	p, err := svc.Charge(context.Background(), payment.ChargeInput{
		BookingID: "b1", UserID: "u", CardID: "c1", Amount: 100,
	})
	if err != nil {
		t.Fatal(err)
	}
	if p.Status != domain.PaymentCaptured || p.ProviderRef == "" {
		t.Fatalf("unexpected payment: %+v", p)
	}
	if core.marked != 1 {
		t.Fatalf("expected MarkBookingPaid called once, got %d", core.marked)
	}
}

func TestCharge_GatewayFailure_PersistsAsFailed(t *testing.T) {
	svc, _ := newSvc(alwaysFail{}, nil)
	p, err := svc.Charge(context.Background(), payment.ChargeInput{
		BookingID: "b1", UserID: "u", CardID: "c1", Amount: 100,
	})
	if err == nil {
		t.Fatal("expected wrapped error")
	}
	if p == nil || p.Status != domain.PaymentFailed {
		t.Fatalf("expected failed row persisted, got %+v err=%v", p, err)
	}
}

func TestCharge_Idempotent_ReturnsExisting(t *testing.T) {
	svc, _ := newSvc(alwaysOK{}, nil)
	first, err := svc.Charge(context.Background(), payment.ChargeInput{
		BookingID: "b1", UserID: "u", CardID: "c1", Amount: 100,
	})
	if err != nil {
		t.Fatal(err)
	}
	second, err := svc.Charge(context.Background(), payment.ChargeInput{
		BookingID: "b1", UserID: "u", CardID: "c1", Amount: 100,
	})
	if !errors.Is(err, errs.ErrAlreadyExists) {
		t.Fatalf("expected already-exists on duplicate, got %v", err)
	}
	if second.ID != first.ID {
		t.Fatalf("expected same payment, got different id")
	}
}

func TestCharge_WrongCardOwner_Forbidden(t *testing.T) {
	pays := &memPayments{}
	cards := &memCards{owner: "other"}
	svc := payment.New(payment.Deps{Payments: pays, Cards: cards, Gateway: alwaysOK{}, Clock: fixedClock{}})
	_, err := svc.Charge(context.Background(), payment.ChargeInput{
		BookingID: "b1", UserID: "u", CardID: "c1", Amount: 100,
	})
	if !errors.Is(err, errs.ErrForbidden) {
		t.Fatalf("want forbidden, got %v", err)
	}
}

func TestCharge_InvalidInputs(t *testing.T) {
	svc, _ := newSvc(alwaysOK{}, nil)
	for _, in := range []payment.ChargeInput{
		{BookingID: "", UserID: "u", CardID: "c", Amount: 1},
		{BookingID: "b", UserID: "", CardID: "c", Amount: 1},
		{BookingID: "b", UserID: "u", CardID: "", Amount: 1},
		{BookingID: "b", UserID: "u", CardID: "c", Amount: 0},
	} {
		_, err := svc.Charge(context.Background(), in)
		if !errors.Is(err, errs.ErrInvalidInput) {
			t.Fatalf("input %+v should be invalid, got %v", in, err)
		}
	}
}

func TestRefund_RequiresCaptured(t *testing.T) {
	svc, _ := newSvc(alwaysFail{}, nil)
	p, _ := svc.Charge(context.Background(), payment.ChargeInput{
		BookingID: "b1", UserID: "u", CardID: "c1", Amount: 100,
	})
	if _, err := svc.Refund(context.Background(), p.ID); !errors.Is(err, errs.ErrConflict) {
		t.Fatalf("want conflict refunding failed payment, got %v", err)
	}
}

func TestRefund_CapturedFlipsToRefunded(t *testing.T) {
	svc, _ := newSvc(alwaysOK{}, nil)
	p, _ := svc.Charge(context.Background(), payment.ChargeInput{
		BookingID: "b1", UserID: "u", CardID: "c1", Amount: 100,
	})
	got, err := svc.Refund(context.Background(), p.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != domain.PaymentRefunded {
		t.Fatalf("expected refunded, got %s", got.Status)
	}
}
