package thread_test

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"testing"
	"time"

	"qonaqzhai-backend/pkg/errs"
	"qonaqzhai-backend/services/realtime/internal/domain"
	"qonaqzhai-backend/services/realtime/internal/ports"
	"qonaqzhai-backend/services/realtime/internal/usecase/thread"
)

type memThreads struct {
	mu       sync.Mutex
	rows     map[string]*domain.Thread
	messages map[string][]*domain.Message
	seq      int
}

func newMem() *memThreads {
	return &memThreads{rows: map[string]*domain.Thread{}, messages: map[string][]*domain.Message{}}
}

func (m *memThreads) EnsureForBooking(_ context.Context, bid, cid, vid string) (*domain.Thread, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, t := range m.rows {
		if t.BookingID == bid {
			cp := *t
			return &cp, nil
		}
	}
	m.seq++
	t := &domain.Thread{ID: "t" + strconv.Itoa(m.seq), BookingID: bid, CustomerID: cid, VendorID: vid,
		CreatedAt: time.Now(), UpdatedAt: time.Now()}
	m.rows[t.ID] = t
	cp := *t
	return &cp, nil
}
func (m *memThreads) FindByID(_ context.Context, id string) (*domain.Thread, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	t, ok := m.rows[id]
	if !ok {
		return nil, errs.ErrNotFound
	}
	cp := *t
	return &cp, nil
}
func (m *memThreads) ListForUser(_ context.Context, uid string) ([]*domain.Thread, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := []*domain.Thread{}
	for _, t := range m.rows {
		if t.CustomerID == uid || t.VendorID == uid {
			cp := *t
			out = append(out, &cp)
		}
	}
	return out, nil
}
func (m *memThreads) AddMessage(_ context.Context, msg *domain.Message) (*domain.Message, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.seq++
	msg.ID = "m" + strconv.Itoa(m.seq)
	msg.CreatedAt = time.Now()
	cp := *msg
	m.messages[msg.ThreadID] = append(m.messages[msg.ThreadID], &cp)
	return &cp, nil
}
func (m *memThreads) ListMessages(_ context.Context, tid string) ([]*domain.Message, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := []*domain.Message{}
	for _, msg := range m.messages[tid] {
		cp := *msg
		out = append(out, &cp)
	}
	return out, nil
}

type stubAuth struct{ users map[string]*ports.ExternalUser }

func (s *stubAuth) GetUsersBatch(_ context.Context, ids []string) ([]*ports.ExternalUser, error) {
	out := []*ports.ExternalUser{}
	for _, id := range ids {
		if u, ok := s.users[id]; ok {
			out = append(out, u)
		}
	}
	return out, nil
}

type stubPublisher struct {
	mu     sync.Mutex
	events int
	last   []string
}

func (p *stubPublisher) Publish(_ string, _ []byte, userIDs ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.events++
	p.last = append([]string(nil), userIDs...)
}

func newSvc() (*thread.Service, *memThreads, *stubPublisher, *stubAuth) {
	threads := newMem()
	pub := &stubPublisher{}
	auth := &stubAuth{users: map[string]*ports.ExternalUser{
		"vendor": {ID: "vendor", Name: "Vendor One"},
		"cust":   {ID: "cust", Email: "c@x.kz"},
	}}
	svc := thread.New(thread.Deps{Threads: threads, Publisher: pub, Auth: auth})
	return svc, threads, pub, auth
}

// --- tests -------------------------------------------------------------------

func TestEnsure_Idempotent(t *testing.T) {
	svc, _, _, _ := newSvc()
	a, err := svc.Ensure(context.Background(), "b1", "cust", "vendor")
	if err != nil {
		t.Fatal(err)
	}
	b, err := svc.Ensure(context.Background(), "b1", "cust", "vendor")
	if err != nil {
		t.Fatal(err)
	}
	if a.ID != b.ID {
		t.Fatalf("Ensure should be idempotent, got %s vs %s", a.ID, b.ID)
	}
}

func TestEnsure_RequiresAllFields(t *testing.T) {
	svc, _, _, _ := newSvc()
	if _, err := svc.Ensure(context.Background(), "", "c", "v"); !errors.Is(err, errs.ErrInvalidInput) {
		t.Fatalf("want invalid input, got %v", err)
	}
}

func TestSend_OwnerAuthorised_BroadcastsToBothPeers(t *testing.T) {
	svc, _, pub, _ := newSvc()
	tr, _ := svc.Ensure(context.Background(), "b1", "cust", "vendor")
	msg, err := svc.Send(context.Background(), "cust", tr.ID, "hello")
	if err != nil {
		t.Fatal(err)
	}
	if msg.Text != "hello" || msg.SenderID != "cust" {
		t.Fatalf("message wrong: %+v", msg)
	}
	if pub.events != 1 || len(pub.last) != 2 {
		t.Fatalf("publish should target both peers, got events=%d targets=%v", pub.events, pub.last)
	}
}

func TestSend_NonMember_Forbidden(t *testing.T) {
	svc, _, _, _ := newSvc()
	tr, _ := svc.Ensure(context.Background(), "b1", "cust", "vendor")
	_, err := svc.Send(context.Background(), "stranger", tr.ID, "spy")
	if !errors.Is(err, errs.ErrForbidden) {
		t.Fatalf("want forbidden, got %v", err)
	}
}

func TestSend_EmptyText_Invalid(t *testing.T) {
	svc, _, _, _ := newSvc()
	tr, _ := svc.Ensure(context.Background(), "b1", "cust", "vendor")
	_, err := svc.Send(context.Background(), "cust", tr.ID, "  ")
	if !errors.Is(err, errs.ErrInvalidInput) {
		t.Fatalf("want invalid input, got %v", err)
	}
}

func TestGet_AuthZ_MemberOnly(t *testing.T) {
	svc, _, _, _ := newSvc()
	tr, _ := svc.Ensure(context.Background(), "b1", "cust", "vendor")
	if _, _, err := svc.Get(context.Background(), "cust", tr.ID); err != nil {
		t.Fatalf("member should read: %v", err)
	}
	if _, _, err := svc.Get(context.Background(), "vendor", tr.ID); err != nil {
		t.Fatalf("member should read: %v", err)
	}
	if _, _, err := svc.Get(context.Background(), "stranger", tr.ID); !errors.Is(err, errs.ErrForbidden) {
		t.Fatalf("stranger must be forbidden, got %v", err)
	}
}

func TestListSummaries_EnrichesCounterpart(t *testing.T) {
	svc, _, _, _ := newSvc()
	_, _ = svc.Ensure(context.Background(), "b1", "cust", "vendor")
	summaries, err := svc.ListSummaries(context.Background(), "cust")
	if err != nil {
		t.Fatal(err)
	}
	if len(summaries) != 1 || summaries[0].Counterpart != "Vendor One" {
		t.Fatalf("counterpart enrichment failed: %+v", summaries)
	}
}
