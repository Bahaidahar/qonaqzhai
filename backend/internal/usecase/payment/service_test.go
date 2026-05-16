package payment_test

import (
	"context"
	"errors"
	"testing"

	"qonaqzhai-backend/internal/domain"
	"qonaqzhai-backend/internal/usecase"
	"qonaqzhai-backend/internal/usecase/inmem"
	"qonaqzhai-backend/internal/usecase/payment"
)

type stubGateway struct {
	intentURL string
	intentID  string
	verifyOK  bool
	wantAmt   int64
	intentErr error
}

func (s *stubGateway) CreatePayment(_ context.Context, in usecase.PaymentIntent) (usecase.PaymentRedirect, error) {
	if s.intentErr != nil {
		return usecase.PaymentRedirect{}, s.intentErr
	}
	s.wantAmt = in.Amount
	return usecase.PaymentRedirect{URL: s.intentURL, TransactionID: s.intentID}, nil
}

func (s *stubGateway) VerifyCallback(form map[string]string) (usecase.CallbackResult, error) {
	if !s.verifyOK {
		return usecase.CallbackResult{}, errors.New("invalid sig")
	}
	return usecase.CallbackResult{
		OrderID:       form["pg_order_id"],
		TransactionID: form["pg_payment_id"],
		Amount:        1500,
		Success:       true,
	}, nil
}

func newSvc(t *testing.T, gw *stubGateway) (*payment.Service, *inmem.UserRepo, *inmem.BookingRepo) {
	t.Helper()
	users := inmem.NewUserRepo(inmem.NewSeqIDGen("u-").New)
	bookings := inmem.NewBookingRepo(inmem.NewSeqIDGen("b-").New)
	svc := payment.New(payment.Deps{
		Bookings: bookings,
		Users:    users,
		Gateway:  gw,
		BaseURL:  "https://app",
	})
	return svc, users, bookings
}

func TestStartIntentRequiresOwnership(t *testing.T) {
	t.Parallel()
	gw := &stubGateway{intentURL: "https://gw/pay", intentID: "pi_1"}
	svc, users, bookings := newSvc(t, gw)
	ctx := context.Background()
	u, _ := users.Create(ctx, &domain.User{Email: "x@y.kz", Role: domain.RoleCustomer})
	b, _ := bookings.Create(ctx, &domain.Booking{CustomerID: u.ID, VendorID: "v1", Amount: 1500, Status: domain.BookingPending})

	_, err := svc.StartIntent(ctx, "other", b.ID)
	if !errors.Is(err, domain.ErrForbidden) {
		t.Errorf("err=%v want ErrForbidden", err)
	}
}

func TestStartIntentRejectsZeroAmount(t *testing.T) {
	t.Parallel()
	gw := &stubGateway{}
	svc, users, bookings := newSvc(t, gw)
	ctx := context.Background()
	u, _ := users.Create(ctx, &domain.User{Email: "x@y.kz", Role: domain.RoleCustomer})
	b, _ := bookings.Create(ctx, &domain.Booking{CustomerID: u.ID, VendorID: "v1", Amount: 0, Status: domain.BookingPending})

	if _, err := svc.StartIntent(ctx, u.ID, b.ID); !errors.Is(err, domain.ErrInvalidInput) {
		t.Errorf("err=%v", err)
	}
}

func TestStartIntentReturnsRedirect(t *testing.T) {
	t.Parallel()
	gw := &stubGateway{intentURL: "https://gw/pay", intentID: "pi_42"}
	svc, users, bookings := newSvc(t, gw)
	ctx := context.Background()
	u, _ := users.Create(ctx, &domain.User{Email: "x@y.kz", Role: domain.RoleCustomer})
	b, _ := bookings.Create(ctx, &domain.Booking{CustomerID: u.ID, VendorID: "v1", Amount: 1500, Status: domain.BookingPending})

	url, err := svc.StartIntent(ctx, u.ID, b.ID)
	if err != nil {
		t.Fatal(err)
	}
	if url != "https://gw/pay" {
		t.Errorf("url=%s", url)
	}
	updated, _ := bookings.Find(ctx, b.ID)
	if updated.PaymentID != "pi_42" {
		t.Errorf("paymentId=%s", updated.PaymentID)
	}
}

func TestHandleCallbackMarksPaid(t *testing.T) {
	t.Parallel()
	gw := &stubGateway{verifyOK: true}
	svc, users, bookings := newSvc(t, gw)
	ctx := context.Background()
	u, _ := users.Create(ctx, &domain.User{Email: "x@y.kz", Role: domain.RoleCustomer})
	b, _ := bookings.Create(ctx, &domain.Booking{CustomerID: u.ID, VendorID: "v1", Amount: 1500, Status: domain.BookingPending})

	form := map[string]string{
		"pg_order_id":   b.ID,
		"pg_payment_id": "pi_callback",
	}
	if err := svc.HandleCallback(ctx, form); err != nil {
		t.Fatal(err)
	}
	updated, _ := bookings.Find(ctx, b.ID)
	if updated.Status != domain.BookingPaid {
		t.Errorf("status=%s", updated.Status)
	}
	if updated.PaymentID != "pi_callback" {
		t.Errorf("paymentId=%s", updated.PaymentID)
	}
}

func TestHandleCallbackRejectsBadSignature(t *testing.T) {
	t.Parallel()
	gw := &stubGateway{verifyOK: false}
	svc, _, _ := newSvc(t, gw)
	if err := svc.HandleCallback(context.Background(), map[string]string{}); err == nil {
		t.Error("invalid sig accepted")
	}
}
