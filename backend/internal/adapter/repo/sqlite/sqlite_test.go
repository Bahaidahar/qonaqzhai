package sqlite_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"qonaqzhai-backend/internal/adapter/repo/sqlite"
	"qonaqzhai-backend/internal/domain"
	"qonaqzhai-backend/internal/infra/db"
	"qonaqzhai-backend/internal/infra/idgen"
	"qonaqzhai-backend/internal/usecase"
)

func openDB(t *testing.T) (*sqlite.UserRepo, *sqlite.VendorRepo, *sqlite.BookingRepo, *sqlite.ReviewRepo, *sqlite.RefreshTokenRepo, *sqlite.PasswordResetRepo, *sqlite.PhotoRepo, *sqlite.NotificationRepo) {
	t.Helper()
	conn, err := db.Open(filepath.Join(t.TempDir(), "t.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = conn.Close() })
	ids := idgen.New()
	return sqlite.NewUserRepo(conn, ids),
		sqlite.NewVendorRepo(conn, ids),
		sqlite.NewBookingRepo(conn, ids),
		sqlite.NewReviewRepo(conn, ids),
		sqlite.NewRefreshTokenRepo(conn, ids),
		sqlite.NewPasswordResetRepo(conn, ids),
		sqlite.NewPhotoRepo(conn, ids),
		sqlite.NewNotificationRepo(conn, ids)
}

func TestUserCRUD(t *testing.T) {
	t.Parallel()
	users, _, _, _, _, _, _, _ := openDB(t)
	ctx := context.Background()
	u, err := users.Create(ctx, &domain.User{Email: "a@b.kz", Name: "A", PasswordHash: "hash", Role: domain.RoleCustomer})
	if err != nil {
		t.Fatal(err)
	}
	if u.ID == "" {
		t.Error("no id assigned")
	}
	if _, err := users.Create(ctx, &domain.User{Email: "a@b.kz", Name: "Dup", PasswordHash: "x", Role: domain.RoleCustomer}); err == nil {
		t.Error("duplicate email accepted")
	}
	got, _ := users.FindByID(ctx, u.ID)
	if got.Email != "a@b.kz" {
		t.Errorf("findByID=%+v", got)
	}
	got, _ = users.FindByEmail(ctx, "a@b.kz")
	if got.ID != u.ID {
		t.Errorf("findByEmail mismatch")
	}
	if err := users.UpdateStatus(ctx, u.ID, domain.UserSuspended); err != nil {
		t.Fatal(err)
	}
	got, _ = users.FindByID(ctx, u.ID)
	if got.Status != domain.UserSuspended {
		t.Errorf("status=%s", got.Status)
	}
	if err := users.UpdatePasswordHash(ctx, u.ID, "newhash"); err != nil {
		t.Fatal(err)
	}
	got, _ = users.FindByID(ctx, u.ID)
	if got.PasswordHash != "newhash" {
		t.Errorf("hash=%s", got.PasswordHash)
	}
}

func TestVendorCRUD(t *testing.T) {
	t.Parallel()
	users, vendors, _, _, _, _, _, _ := openDB(t)
	ctx := context.Background()
	u, _ := users.Create(ctx, &domain.User{Email: "v@b.kz", Name: "V", PasswordHash: "h", Role: domain.RoleVendor})
	v, err := vendors.Upsert(ctx, u.ID, domain.VendorInput{Name: "Rixos", Category: "Venue", City: "Almaty", PriceFrom: 1_500_000})
	if err != nil {
		t.Fatal(err)
	}
	if v.Status != domain.VendorPending {
		t.Errorf("status=%s", v.Status)
	}
	// upsert again — same ID, updated name
	v2, _ := vendors.Upsert(ctx, u.ID, domain.VendorInput{Name: "Rixos Plus", Category: "Venue", City: "Almaty", PriceFrom: 2_000_000})
	if v2.ID != v.ID {
		t.Error("upsert created new row")
	}
	if v2.Name != "Rixos Plus" || v2.PriceFrom != 2_000_000 {
		t.Errorf("not updated: %+v", v2)
	}
	if err := vendors.UpdateStatus(ctx, v.ID, domain.VendorApproved); err != nil {
		t.Fatal(err)
	}
	if err := vendors.UpdateRating(ctx, v.ID, 4.6, 12); err != nil {
		t.Fatal(err)
	}
	got, _ := vendors.FindByID(ctx, v.ID)
	if got.RatingAvg != 4.6 || got.RatingCount != 12 || got.Status != domain.VendorApproved {
		t.Errorf("got=%+v", got)
	}
}

func TestVendorSearchFilters(t *testing.T) {
	t.Parallel()
	users, vendors, _, _, _, _, _, _ := openDB(t)
	ctx := context.Background()
	seed := func(email, name, cat, city string, price int64, status domain.VendorStatus) {
		u, _ := users.Create(ctx, &domain.User{Email: email, Name: name, PasswordHash: "h", Role: domain.RoleVendor})
		v, _ := vendors.Upsert(ctx, u.ID, domain.VendorInput{Name: name, Category: cat, City: city, PriceFrom: price, Description: "desc " + name})
		_ = vendors.UpdateStatus(ctx, v.ID, status)
	}
	seed("a@b", "Rixos Almaty", "Venue", "Almaty", 1_500_000, domain.VendorApproved)
	seed("c@d", "Aitu Photo", "Photo", "Almaty", 450_000, domain.VendorApproved)
	seed("e@f", "Aizada Cater", "Catering", "Astana", 800_000, domain.VendorPending)

	list, total, err := vendors.Search(ctx, usecase.VendorQuery{Status: domain.VendorApproved})
	if err != nil {
		t.Fatal(err)
	}
	if total != 2 || len(list) != 2 {
		t.Errorf("approved count=%d total=%d", len(list), total)
	}

	list, _, _ = vendors.Search(ctx, usecase.VendorQuery{Status: domain.VendorApproved, City: "Almaty", Category: "Venue"})
	if len(list) != 1 || list[0].Name != "Rixos Almaty" {
		t.Errorf("filter=%+v", list)
	}

	list, _, _ = vendors.Search(ctx, usecase.VendorQuery{Status: domain.VendorApproved, MaxPrice: 500_000})
	if len(list) != 1 || list[0].Name != "Aitu Photo" {
		t.Errorf("price filter: %+v", list)
	}

	// FTS5
	list, _, _ = vendors.Search(ctx, usecase.VendorQuery{Status: domain.VendorApproved, Q: "Rixos"})
	if len(list) != 1 {
		t.Errorf("fts: got %d", len(list))
	}

	// Sort
	list, _, _ = vendors.Search(ctx, usecase.VendorQuery{Status: domain.VendorApproved, Sort: "price_asc"})
	if list[0].PriceFrom > list[1].PriceFrom {
		t.Error("price_asc broken")
	}
}

func TestBookingCRUD(t *testing.T) {
	t.Parallel()
	users, vendors, bookings, _, _, _, _, _ := openDB(t)
	ctx := context.Background()
	cu, _ := users.Create(ctx, &domain.User{Email: "c@b", Role: domain.RoleCustomer, PasswordHash: "h", Name: "C"})
	vu, _ := users.Create(ctx, &domain.User{Email: "v@b", Role: domain.RoleVendor, PasswordHash: "h", Name: "V"})
	v, _ := vendors.Upsert(ctx, vu.ID, domain.VendorInput{Name: "X", Category: "Y", City: "Z"})

	b, err := bookings.Create(ctx, &domain.Booking{CustomerID: cu.ID, VendorID: v.ID, EventDate: "2026-06-12", GuestCount: 100, Amount: 500_000})
	if err != nil {
		t.Fatal(err)
	}
	if b.Status != domain.BookingPending {
		t.Errorf("status=%s", b.Status)
	}
	if err := bookings.UpdateStatus(ctx, b.ID, domain.BookingAccepted); err != nil {
		t.Fatal(err)
	}
	if err := bookings.SetPayment(ctx, b.ID, "pi_test_123"); err != nil {
		t.Fatal(err)
	}
	got, _ := bookings.Find(ctx, b.ID)
	if got.Status != domain.BookingAccepted || got.PaymentID != "pi_test_123" {
		t.Errorf("update lost: %+v", got)
	}

	list, _ := bookings.ListForCustomer(ctx, cu.ID)
	if len(list) != 1 {
		t.Errorf("customer list=%d", len(list))
	}
}

func TestReviewAggregate(t *testing.T) {
	t.Parallel()
	users, vendors, bookings, reviews, _, _, _, _ := openDB(t)
	ctx := context.Background()
	cu, _ := users.Create(ctx, &domain.User{Email: "c@b", Role: domain.RoleCustomer, PasswordHash: "h", Name: "C"})
	vu, _ := users.Create(ctx, &domain.User{Email: "v@b", Role: domain.RoleVendor, PasswordHash: "h", Name: "V"})
	v, _ := vendors.Upsert(ctx, vu.ID, domain.VendorInput{Name: "X", Category: "Y", City: "Z"})
	b1, _ := bookings.Create(ctx, &domain.Booking{CustomerID: cu.ID, VendorID: v.ID, EventDate: "2026-06-12", Status: domain.BookingCompleted})
	b2, _ := bookings.Create(ctx, &domain.Booking{CustomerID: cu.ID, VendorID: v.ID, EventDate: "2026-06-13", Status: domain.BookingCompleted})

	if _, err := reviews.Create(ctx, &domain.Review{BookingID: b1.ID, CustomerID: cu.ID, VendorID: v.ID, Rating: 5}); err != nil {
		t.Fatal(err)
	}
	if _, err := reviews.Create(ctx, &domain.Review{BookingID: b2.ID, CustomerID: cu.ID, VendorID: v.ID, Rating: 3}); err != nil {
		t.Fatal(err)
	}
	// uniq constraint on booking_id
	if _, err := reviews.Create(ctx, &domain.Review{BookingID: b1.ID, CustomerID: cu.ID, VendorID: v.ID, Rating: 4}); err == nil {
		t.Error("dup review accepted")
	}
	avg, count, err := reviews.AggregateForVendor(ctx, v.ID)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 || avg != 4 {
		t.Errorf("avg=%v count=%d", avg, count)
	}
}

func TestRefreshTokenLifecycle(t *testing.T) {
	t.Parallel()
	users, _, _, _, rts, _, _, _ := openDB(t)
	ctx := context.Background()
	u, _ := users.Create(ctx, &domain.User{Email: "x@y", Role: domain.RoleCustomer, PasswordHash: "h", Name: "X"})
	now := time.Now().UTC()
	tok := &domain.RefreshToken{UserID: u.ID, TokenHash: "hash-abc", ExpiresAt: now.Add(time.Hour)}
	if err := rts.Create(ctx, tok); err != nil {
		t.Fatal(err)
	}
	got, err := rts.FindActiveByHash(ctx, "hash-abc", now)
	if err != nil {
		t.Fatal(err)
	}
	if got.UserID != u.ID {
		t.Errorf("user mismatch")
	}
	if err := rts.Revoke(ctx, got.ID, now); err != nil {
		t.Fatal(err)
	}
	if _, err := rts.FindActiveByHash(ctx, "hash-abc", now); err != domain.ErrNotFound {
		t.Errorf("revoked still returned: %v", err)
	}
}
