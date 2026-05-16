package notification_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"qonaqzhai-backend/internal/domain"
	"qonaqzhai-backend/internal/usecase/inmem"
	"qonaqzhai-backend/internal/usecase/notification"
)

type stubMailer struct {
	mu    sync.Mutex
	calls int
	err   error
}

func (s *stubMailer) Send(_ context.Context, _, _, _ string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.calls++
	return s.err
}

type stubPusher struct {
	mu     sync.Mutex
	tokens []string
}

func (s *stubPusher) Send(_ context.Context, tokens []string, _, _ string, _ map[string]string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tokens = append(s.tokens, tokens...)
	return nil
}

type fakeFCMTokens struct{ map_ map[string][]string }

func (f *fakeFCMTokens) TokensForUser(_ context.Context, userID string) ([]string, error) {
	return f.map_[userID], nil
}
func (f *fakeFCMTokens) Register(_ context.Context, userID, token, _ string) error {
	f.map_[userID] = append(f.map_[userID], token)
	return nil
}
func (f *fakeFCMTokens) Unregister(_ context.Context, _ string) error { return nil }

func waitFor(t *testing.T, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatal("condition not met within timeout")
}

func TestEnqueueDeliversEmail(t *testing.T) {
	t.Parallel()
	users := inmem.NewUserRepo(inmem.NewSeqIDGen("u-").New)
	notifs := inmem.NewNotificationRepo(inmem.NewSeqIDGen("n-").New)
	u, _ := users.Create(context.Background(), &domain.User{Email: "x@y.com", Role: domain.RoleCustomer})
	mailer := &stubMailer{}
	svc := notification.New(notification.Deps{
		Notifications: notifs,
		Users:         users,
		Mailer:        mailer,
	})
	defer svc.Stop()

	if err := svc.Enqueue(context.Background(), &domain.Notification{
		UserID:  u.ID,
		Type:    domain.NotifSignupWelcome,
		Channel: domain.ChannelEmail,
		Title:   "Welcome",
		Body:    "<p>hi</p>",
	}); err != nil {
		t.Fatal(err)
	}
	waitFor(t, func() bool {
		mailer.mu.Lock()
		defer mailer.mu.Unlock()
		return mailer.calls == 1
	})
	list, _ := notifs.ListForUser(context.Background(), u.ID, 10)
	if len(list) != 1 {
		t.Fatalf("got %d notifications", len(list))
	}
	waitFor(t, func() bool {
		list, _ := notifs.ListForUser(context.Background(), u.ID, 10)
		return len(list) > 0 && list[0].Status == "sent"
	})
}

func TestEnqueueDeliversPush(t *testing.T) {
	t.Parallel()
	users := inmem.NewUserRepo(inmem.NewSeqIDGen("u-").New)
	notifs := inmem.NewNotificationRepo(inmem.NewSeqIDGen("n-").New)
	tokens := &fakeFCMTokens{map_: map[string][]string{}}
	pusher := &stubPusher{}

	u, _ := users.Create(context.Background(), &domain.User{Email: "x@y.com", Role: domain.RoleCustomer})
	_ = tokens.Register(context.Background(), u.ID, "tok-1", "android")

	svc := notification.New(notification.Deps{
		Notifications: notifs,
		Users:         users,
		FCMTokens:     tokens,
		Pusher:        pusher,
	})
	defer svc.Stop()

	_ = svc.Enqueue(context.Background(), &domain.Notification{
		UserID:  u.ID,
		Type:    domain.NotifBookingCreated,
		Channel: domain.ChannelPush,
		Title:   "New booking",
		Body:    "Aigerim sent you a request",
	})
	waitFor(t, func() bool {
		pusher.mu.Lock()
		defer pusher.mu.Unlock()
		return len(pusher.tokens) == 1
	})
}

func TestEnqueueMarksFailedOnMailerError(t *testing.T) {
	t.Parallel()
	users := inmem.NewUserRepo(inmem.NewSeqIDGen("u-").New)
	notifs := inmem.NewNotificationRepo(inmem.NewSeqIDGen("n-").New)
	u, _ := users.Create(context.Background(), &domain.User{Email: "x@y.com", Role: domain.RoleCustomer})
	mailer := &stubMailer{err: errors.New("smtp down")}
	svc := notification.New(notification.Deps{Notifications: notifs, Users: users, Mailer: mailer})
	defer svc.Stop()

	_ = svc.Enqueue(context.Background(), &domain.Notification{
		UserID: u.ID, Type: domain.NotifPasswordReset, Channel: domain.ChannelEmail,
		Title: "Reset", Body: "tok",
	})
	waitFor(t, func() bool {
		list, _ := notifs.ListForUser(context.Background(), u.ID, 10)
		return len(list) == 1 && list[0].Status == "failed"
	})
}
