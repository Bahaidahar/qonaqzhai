package booking_test

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"testing"
	"time"

	"qonaqzhai-backend/pkg/errs"
	"qonaqzhai-backend/services/core/internal/domain"
	"qonaqzhai-backend/services/core/internal/ports"
	"qonaqzhai-backend/services/core/internal/usecase/booking"
)

// --- fakes -------------------------------------------------------------------

type memBookings struct {
	mu   sync.Mutex
	rows map[string]*domain.Booking
}

func newMemBookings() *memBookings { return &memBookings{rows: map[string]*domain.Booking{}} }

func (m *memBookings) Create(_ context.Context, b *domain.Booking) (*domain.Booking, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if b.ID == "" {
		b.ID = "b-" + strconv.Itoa(len(m.rows)+1)
	}
	cp := *b
	m.rows[cp.ID] = &cp
	return &cp, nil
}
func (m *memBookings) Find(_ context.Context, id string) (*domain.Booking, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	b, ok := m.rows[id]
	if !ok {
		return nil, errs.ErrNotFound
	}
	cp := *b
	return &cp, nil
}
func (m *memBookings) ListForCustomer(_ context.Context, customerID string, p ports.Page) ([]*domain.Booking, error) {
	return m.filter(func(b *domain.Booking) bool { return b.CustomerID == customerID }, p), nil
}
func (m *memBookings) ListForVendor(_ context.Context, vendorID string, p ports.Page) ([]*domain.Booking, error) {
	return m.filter(func(b *domain.Booking) bool { return b.VendorID == vendorID }, p), nil
}
func (m *memBookings) ListAll(_ context.Context, p ports.Page) ([]*domain.Booking, error) {
	return m.filter(func(*domain.Booking) bool { return true }, p), nil
}
func (m *memBookings) filter(pred func(*domain.Booking) bool, p ports.Page) []*domain.Booking {
	m.mu.Lock()
	defer m.mu.Unlock()
	p = p.Clamp()
	out := []*domain.Booking{}
	for _, b := range m.rows {
		if pred(b) {
			cp := *b
			out = append(out, &cp)
		}
	}
	if p.Offset >= len(out) {
		return out[:0]
	}
	end := p.Offset + p.Limit
	if end > len(out) {
		end = len(out)
	}
	return out[p.Offset:end]
}
func (m *memBookings) UpdateStatus(_ context.Context, id string, status domain.BookingStatus) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	b, ok := m.rows[id]
	if !ok {
		return errs.ErrNotFound
	}
	b.Status = status
	return nil
}
func (m *memBookings) SetPayment(_ context.Context, id, paymentID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	b, ok := m.rows[id]
	if !ok {
		return errs.ErrNotFound
	}
	b.PaymentID = paymentID
	return nil
}
func (m *memBookings) MarkPaid(_ context.Context, id, paymentID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	b, ok := m.rows[id]
	if !ok {
		return errs.ErrNotFound
	}
	b.PaymentID = paymentID
	b.Status = domain.BookingPaid
	return nil
}
func (m *memBookings) Stats(_ context.Context) (ports.BookingStats, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var s ports.BookingStats
	for _, b := range m.rows {
		s.Total++
		switch b.Status {
		case domain.BookingPending:
			s.Pending++
		case domain.BookingAccepted:
			s.Accepted++
		case domain.BookingPaid:
			s.Paid++
			s.GMV += b.Amount
		}
	}
	return s, nil
}

type memVendors struct {
	mu sync.Mutex
	rows map[string]*domain.Vendor
}

func newMemVendors() *memVendors { return &memVendors{rows: map[string]*domain.Vendor{}} }

func (m *memVendors) seed(id, userID string, status domain.VendorStatus) *domain.Vendor {
	m.mu.Lock()
	defer m.mu.Unlock()
	v := &domain.Vendor{ID: id, UserID: userID, Status: status, Name: "v", Category: "c", City: "Almaty"}
	m.rows[id] = v
	return v
}

func (m *memVendors) Upsert(context.Context, string, domain.VendorInput) (*domain.Vendor, error) {
	return nil, errors.New("not used in test")
}
func (m *memVendors) FindByID(_ context.Context, id string) (*domain.Vendor, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	v, ok := m.rows[id]
	if !ok {
		return nil, errs.ErrNotFound
	}
	cp := *v
	return &cp, nil
}
func (m *memVendors) FindByUserID(_ context.Context, userID string) (*domain.Vendor, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, v := range m.rows {
		if v.UserID == userID {
			cp := *v
			return &cp, nil
		}
	}
	return nil, errs.ErrNotFound
}
func (m *memVendors) FindByIDs(context.Context, []string) ([]*domain.Vendor, error) {
	return nil, nil
}
func (m *memVendors) Search(context.Context, ports.VendorQuery) ([]*domain.Vendor, int, error) {
	return nil, 0, nil
}
func (m *memVendors) UpdateStatus(_ context.Context, id string, status domain.VendorStatus) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	v, ok := m.rows[id]
	if !ok {
		return errs.ErrNotFound
	}
	v.Status = status
	return nil
}
func (m *memVendors) UpdateRating(_ context.Context, id string, avg float64, count int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	v, ok := m.rows[id]
	if !ok {
		return errs.ErrNotFound
	}
	v.RatingAvg = avg
	v.RatingCount = count
	return nil
}
func (m *memVendors) CountByStatus(context.Context) (map[domain.VendorStatus]int, error) {
	return nil, nil
}

type memNotifs struct{ count int; mu sync.Mutex }

func (m *memNotifs) Enqueue(_ context.Context, n *domain.Notification) (*domain.Notification, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.count++
	n.ID = "n-" + strconv.Itoa(m.count)
	return n, nil
}
func (m *memNotifs) ListForUser(context.Context, string, ports.Page) ([]*domain.Notification, error) {
	return nil, nil
}
func (m *memNotifs) MarkSent(context.Context, string) error { return nil }

type stubPayments struct {
	chargeErr   error
	chargeStatus string
}

func (s *stubPayments) Charge(_ context.Context, in ports.ChargeRequest) (*ports.PaymentResult, error) {
	if s.chargeErr != nil {
		return nil, s.chargeErr
	}
	status := s.chargeStatus
	if status == "" {
		status = "captured"
	}
	return &ports.PaymentResult{ID: "pay-" + in.BookingID, Status: status}, nil
}

type stubRealtime struct{ ensured, published int; mu sync.Mutex }

func (s *stubRealtime) EnsureThread(_ context.Context, _, _, _ string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ensured++
	return nil
}
func (s *stubRealtime) Publish(_ context.Context, _ string, _ []byte, _ ...string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.published++
	return nil
}

// --- helpers -----------------------------------------------------------------

func newSvc(t *testing.T, payments ports.PaymentClient, realtime ports.RealtimeClient) (*booking.Service, *memBookings, *memVendors) {
	t.Helper()
	bookings := newMemBookings()
	vendors := newMemVendors()
	return booking.New(booking.Deps{
		Bookings: bookings, Vendors: vendors, Notifications: &memNotifs{},
		Payments: payments, Realtime: realtime,
	}), bookings, vendors
}

func seed(t *testing.T, vendors *memVendors, bookings *memBookings, status domain.BookingStatus) *domain.Booking {
	t.Helper()
	vendors.seed("v1", "vendor-user", domain.VendorApproved)
	b, err := bookings.Create(context.Background(), &domain.Booking{
		CustomerID: "cust", VendorID: "v1", EventDate: "2026-06-01",
		GuestCount: 50, Amount: 200000, Status: status,
	})
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// --- tests -------------------------------------------------------------------

func TestCreate_HappyPath(t *testing.T) {
	svc, _, vendors := newSvc(t, nil, nil)
	vendors.seed("v1", "vendor-user", domain.VendorApproved)
	b, err := svc.Create(context.Background(), "cust", booking.CreateInput{
		VendorID: "v1", EventDate: "2026-06-01", GuestCount: 50, Amount: 100000,
	})
	if err != nil {
		t.Fatal(err)
	}
	if b.Status != domain.BookingPending {
		t.Fatalf("expected pending, got %s", b.Status)
	}
}

func TestCreate_VendorNotPublicForbidden(t *testing.T) {
	svc, _, vendors := newSvc(t, nil, nil)
	vendors.seed("v1", "vendor-user", domain.VendorPending)
	_, err := svc.Create(context.Background(), "cust", booking.CreateInput{
		VendorID: "v1", EventDate: "2026-06-01", Amount: 1,
	})
	if !errors.Is(err, errs.ErrForbidden) {
		t.Fatalf("want forbidden, got %v", err)
	}
}

func TestVendorTransition_Accept_EnsuresThread(t *testing.T) {
	rt := &stubRealtime{}
	svc, bookings, vendors := newSvc(t, nil, rt)
	b := seed(t, vendors, bookings, domain.BookingPending)
	got, err := svc.VendorTransition(context.Background(), "vendor-user", b.ID, domain.BookingAccepted)
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != domain.BookingAccepted {
		t.Fatalf("expected accepted, got %s", got.Status)
	}
	if rt.ensured != 1 {
		t.Fatalf("expected EnsureThread called once, got %d", rt.ensured)
	}
}

func TestVendorTransition_NotOwner_Forbidden(t *testing.T) {
	svc, bookings, vendors := newSvc(t, nil, nil)
	b := seed(t, vendors, bookings, domain.BookingPending)
	_, err := svc.VendorTransition(context.Background(), "someone-else", b.ID, domain.BookingAccepted)
	if !errors.Is(err, errs.ErrForbidden) {
		t.Fatalf("want forbidden, got %v", err)
	}
}

func TestVendorTransition_InvalidState_Conflict(t *testing.T) {
	svc, bookings, vendors := newSvc(t, nil, nil)
	b := seed(t, vendors, bookings, domain.BookingDeclined)
	_, err := svc.VendorTransition(context.Background(), "vendor-user", b.ID, domain.BookingAccepted)
	if !errors.Is(err, errs.ErrConflict) {
		t.Fatalf("want conflict, got %v", err)
	}
}

func TestCustomerCancel_OnlyOwner(t *testing.T) {
	svc, bookings, vendors := newSvc(t, nil, nil)
	b := seed(t, vendors, bookings, domain.BookingPending)
	_, err := svc.CustomerCancel(context.Background(), "wrong", b.ID)
	if !errors.Is(err, errs.ErrForbidden) {
		t.Fatalf("want forbidden, got %v", err)
	}
	got, err := svc.CustomerCancel(context.Background(), "cust", b.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != domain.BookingCancelled {
		t.Fatalf("expected cancelled, got %s", got.Status)
	}
}

func TestPay_RequiresAcceptedStatus(t *testing.T) {
	svc, bookings, vendors := newSvc(t, &stubPayments{}, nil)
	b := seed(t, vendors, bookings, domain.BookingPending)
	_, err := svc.Pay(context.Background(), "cust", b.ID, "card1", "KZT")
	if !errors.Is(err, errs.ErrConflict) {
		t.Fatalf("want conflict, got %v", err)
	}
}

func TestPay_AtomicMarkPaid(t *testing.T) {
	pay := &stubPayments{}
	svc, bookings, vendors := newSvc(t, pay, nil)
	b := seed(t, vendors, bookings, domain.BookingAccepted)
	got, err := svc.Pay(context.Background(), "cust", b.ID, "card1", "KZT")
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != domain.BookingPaid {
		t.Fatalf("expected paid, got %s", got.Status)
	}
	if got.PaymentID == "" {
		t.Fatal("expected payment id set")
	}
}

func TestPay_GatewayDeclined_KeepsAccepted(t *testing.T) {
	pay := &stubPayments{chargeStatus: "failed"}
	svc, bookings, vendors := newSvc(t, pay, nil)
	b := seed(t, vendors, bookings, domain.BookingAccepted)
	_, err := svc.Pay(context.Background(), "cust", b.ID, "card1", "KZT")
	if !errors.Is(err, errs.ErrConflict) {
		t.Fatalf("want conflict, got %v", err)
	}
	got, _ := bookings.Find(context.Background(), b.ID)
	if got.Status != domain.BookingAccepted {
		t.Fatalf("failed payment must not flip status, got %s", got.Status)
	}
}

func TestPay_NoPaymentsClient_Upstream(t *testing.T) {
	svc, bookings, vendors := newSvc(t, nil, nil)
	b := seed(t, vendors, bookings, domain.BookingAccepted)
	_, err := svc.Pay(context.Background(), "cust", b.ID, "card1", "KZT")
	if !errors.Is(err, errs.ErrUpstream) {
		t.Fatalf("want upstream, got %v", err)
	}
}

func TestMarkPaid_GRPCEntrypoint(t *testing.T) {
	svc, bookings, vendors := newSvc(t, nil, nil)
	b := seed(t, vendors, bookings, domain.BookingAccepted)
	got, err := svc.MarkPaid(context.Background(), b.ID, "external-ref")
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != domain.BookingPaid || got.PaymentID != "external-ref" {
		t.Fatalf("mark paid did not commit: %+v", got)
	}
}

func TestFind_Authz(t *testing.T) {
	svc, bookings, vendors := newSvc(t, nil, nil)
	b := seed(t, vendors, bookings, domain.BookingPending)
	if _, err := svc.Find(context.Background(), "cust", b.ID); err != nil {
		t.Fatalf("owner should see: %v", err)
	}
	if _, err := svc.Find(context.Background(), "vendor-user", b.ID); err != nil {
		t.Fatalf("vendor should see: %v", err)
	}
	if _, err := svc.Find(context.Background(), "stranger", b.ID); !errors.Is(err, errs.ErrForbidden) {
		t.Fatalf("stranger must be forbidden, got %v", err)
	}
}

func TestStats(t *testing.T) {
	svc, bookings, vendors := newSvc(t, nil, nil)
	vendors.seed("v1", "u", domain.VendorApproved)
	for i, st := range []domain.BookingStatus{
		domain.BookingPending, domain.BookingAccepted, domain.BookingPaid, domain.BookingPaid,
	} {
		_, _ = bookings.Create(context.Background(), &domain.Booking{
			ID: "b" + strconv.Itoa(i), VendorID: "v1", CustomerID: "c", Status: st, Amount: 100,
		})
	}
	s, err := svc.Stats(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if s.Total != 4 || s.Paid != 2 || s.GMV != 200 {
		t.Fatalf("stats wrong: %+v", s)
	}
}

func TestIsAccepted(t *testing.T) {
	svc, bookings, vendors := newSvc(t, nil, nil)
	b := seed(t, vendors, bookings, domain.BookingAccepted)
	_, vendorUser, ok, err := svc.IsAccepted(context.Background(), b.ID)
	if err != nil || !ok {
		t.Fatalf("expected accepted, got %v %v", ok, err)
	}
	if vendorUser != "vendor-user" {
		t.Fatalf("expected vendor user id, got %s", vendorUser)
	}
}

func TestListPagination(t *testing.T) {
	svc, bookings, vendors := newSvc(t, nil, nil)
	vendors.seed("v1", "u", domain.VendorApproved)
	for i := 0; i < 5; i++ {
		_, _ = bookings.Create(context.Background(), &domain.Booking{
			VendorID: "v1", CustomerID: "c", Status: domain.BookingPending, Amount: 1,
		})
	}
	got, err := svc.ListForCustomer(context.Background(), "c", ports.Page{Limit: 2})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

// quiet timestamps so test runner reports stay short.
var _ = time.Now
