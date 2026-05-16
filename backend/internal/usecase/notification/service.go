// Package notification implements an async in-app + email + push notifier
// driven by an in-process channel queue and a worker goroutine.
package notification

import (
	"context"
	"log/slog"
	"sync"

	"qonaqzhai-backend/internal/domain"
	"qonaqzhai-backend/internal/usecase"
)

// FCMTokenRepo retrieves device tokens for a user.
type FCMTokenRepo interface {
	TokensForUser(ctx context.Context, userID string) ([]string, error)
	Register(ctx context.Context, userID, token, platform string) error
	Unregister(ctx context.Context, token string) error
}

// Deps bundles notification collaborators.
type Deps struct {
	Notifications usecase.NotificationRepo
	Users         usecase.UserRepo
	FCMTokens     FCMTokenRepo
	Mailer        usecase.Mailer // optional
	Pusher        usecase.Pusher // optional
	Logger        *slog.Logger
	QueueSize     int // default 128
}

// Service queues + dispatches notifications.
type Service struct {
	d      Deps
	queue  chan job
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

type job struct {
	userID string
	n      *domain.Notification
}

// New constructs and starts a worker goroutine.
func New(d Deps) *Service {
	if d.Logger == nil {
		d.Logger = slog.Default()
	}
	if d.QueueSize <= 0 {
		d.QueueSize = 128
	}
	ctx, cancel := context.WithCancel(context.Background())
	s := &Service{
		d:      d,
		queue:  make(chan job, d.QueueSize),
		ctx:    ctx,
		cancel: cancel,
	}
	s.wg.Add(1)
	go s.loop()
	return s
}

// Stop drains the queue and stops the worker. Safe to call multiple times.
func (s *Service) Stop() {
	if s == nil {
		return
	}
	s.cancel()
	close(s.queue)
	s.wg.Wait()
}

// Enqueue persists the notification (status=queued) and schedules delivery.
// The caller fills UserID, Type, Channel, Title, Body.
func (s *Service) Enqueue(ctx context.Context, n *domain.Notification) error {
	created, err := s.d.Notifications.Create(ctx, n)
	if err != nil {
		return err
	}
	select {
	case s.queue <- job{userID: created.UserID, n: created}:
	default:
		// queue full → mark failed so callers can retry via a separate job
		_ = s.d.Notifications.MarkFailed(ctx, created.ID)
		s.d.Logger.Warn("notification queue full", slog.String("id", created.ID))
	}
	return nil
}

// ListForUser returns the in-app inbox for user.
func (s *Service) ListForUser(ctx context.Context, userID string, limit int) ([]*domain.Notification, error) {
	return s.d.Notifications.ListForUser(ctx, userID, limit)
}

func (s *Service) loop() {
	defer s.wg.Done()
	for j := range s.queue {
		s.deliver(j)
	}
}

func (s *Service) deliver(j job) {
	ctx := s.ctx
	n := j.n
	var deliveryErr error

	if n.Channel == domain.ChannelEmail || n.Channel == domain.ChannelBoth {
		if s.d.Mailer != nil {
			u, err := s.d.Users.FindByID(ctx, n.UserID)
			if err == nil {
				if err := s.d.Mailer.Send(ctx, u.Email, n.Title, n.Body); err != nil {
					s.d.Logger.Warn("email send failed", slog.String("err", err.Error()))
					deliveryErr = err
				}
			}
		}
	}
	if n.Channel == domain.ChannelPush || n.Channel == domain.ChannelBoth {
		if s.d.Pusher != nil && s.d.FCMTokens != nil {
			tokens, _ := s.d.FCMTokens.TokensForUser(ctx, n.UserID)
			if len(tokens) > 0 {
				data := map[string]string{"type": string(n.Type)}
				if err := s.d.Pusher.Send(ctx, tokens, n.Title, n.Body, data); err != nil {
					s.d.Logger.Warn("push send failed", slog.String("err", err.Error()))
					deliveryErr = err
				}
			}
		}
	}

	if deliveryErr != nil {
		_ = s.d.Notifications.MarkFailed(ctx, n.ID)
		return
	}
	_ = s.d.Notifications.MarkSent(ctx, n.ID)
}
