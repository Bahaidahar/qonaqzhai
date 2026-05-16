// Package inmem provides in-memory port implementations for use in unit tests.
// Not safe for production — single process, no durability.
package inmem

import (
	"context"
	"errors"
	"sort"
	"strings"
	"sync"
	"time"

	"qonaqzhai-backend/internal/domain"
	"qonaqzhai-backend/internal/usecase"
)

// --- User repo ---

type UserRepo struct {
	mu    sync.Mutex
	byID  map[string]*domain.User
	idGen func() string
}

func NewUserRepo(idGen func() string) *UserRepo {
	return &UserRepo{byID: map[string]*domain.User{}, idGen: idGen}
}

func (r *UserRepo) Create(_ context.Context, u *domain.User) (*domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, existing := range r.byID {
		if existing.Email == u.Email {
			return nil, domain.ErrAlreadyExists
		}
	}
	cp := *u
	if cp.ID == "" {
		cp.ID = r.idGen()
	}
	if cp.CreatedAt.IsZero() {
		cp.CreatedAt = time.Now().UTC()
	}
	if cp.Status == "" {
		cp.Status = domain.UserActive
	}
	r.byID[cp.ID] = &cp
	out := cp
	return &out, nil
}

func (r *UserRepo) FindByID(_ context.Context, id string) (*domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	u, ok := r.byID[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	cp := *u
	return &cp, nil
}

func (r *UserRepo) FindByEmail(_ context.Context, email string) (*domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, u := range r.byID {
		if u.Email == email {
			cp := *u
			return &cp, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (r *UserRepo) List(_ context.Context) ([]*domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]*domain.User, 0, len(r.byID))
	for _, u := range r.byID {
		cp := *u
		out = append(out, &cp)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.After(out[j].CreatedAt) })
	return out, nil
}

func (r *UserRepo) UpdateStatus(_ context.Context, id string, status domain.UserStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	u, ok := r.byID[id]
	if !ok {
		return domain.ErrNotFound
	}
	u.Status = status
	return nil
}

func (r *UserRepo) UpdatePasswordHash(_ context.Context, id, hash string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	u, ok := r.byID[id]
	if !ok {
		return domain.ErrNotFound
	}
	u.PasswordHash = hash
	return nil
}

// --- Vendor repo ---

type VendorRepo struct {
	mu      sync.Mutex
	byID    map[string]*domain.Vendor
	byUser  map[string]string
	idGen   func() string
	photos  map[string][]string
}

func NewVendorRepo(idGen func() string) *VendorRepo {
	return &VendorRepo{
		byID:   map[string]*domain.Vendor{},
		byUser: map[string]string{},
		idGen:  idGen,
		photos: map[string][]string{},
	}
}

func (r *VendorRepo) Upsert(_ context.Context, userID string, in domain.VendorInput) (*domain.Vendor, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now().UTC()
	if id, ok := r.byUser[userID]; ok {
		v := r.byID[id]
		v.Name = in.Name
		v.Category = in.Category
		v.City = in.City
		v.Description = in.Description
		v.PriceFrom = in.PriceFrom
		v.UpdatedAt = now
		cp := *v
		return &cp, nil
	}
	v := &domain.Vendor{
		ID:          r.idGen(),
		UserID:      userID,
		Name:        in.Name,
		Category:    in.Category,
		City:        in.City,
		Description: in.Description,
		PriceFrom:   in.PriceFrom,
		Status:      domain.VendorPending,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	r.byID[v.ID] = v
	r.byUser[userID] = v.ID
	cp := *v
	return &cp, nil
}

func (r *VendorRepo) FindByID(_ context.Context, id string) (*domain.Vendor, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	v, ok := r.byID[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	cp := *v
	cp.PhotoIDs = append([]string{}, r.photos[id]...)
	return &cp, nil
}

func (r *VendorRepo) FindByUserID(_ context.Context, userID string) (*domain.Vendor, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	id, ok := r.byUser[userID]
	if !ok {
		return nil, domain.ErrNotFound
	}
	cp := *r.byID[id]
	cp.PhotoIDs = append([]string{}, r.photos[id]...)
	return &cp, nil
}

func (r *VendorRepo) Search(_ context.Context, q usecase.VendorQuery) ([]*domain.Vendor, int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	all := make([]*domain.Vendor, 0, len(r.byID))
	for _, v := range r.byID {
		if q.Status != "" && v.Status != q.Status {
			continue
		}
		if q.Category != "" && v.Category != q.Category {
			continue
		}
		if q.City != "" && v.City != q.City {
			continue
		}
		if q.MinPrice > 0 && v.PriceFrom < q.MinPrice {
			continue
		}
		if q.MaxPrice > 0 && v.PriceFrom > q.MaxPrice {
			continue
		}
		if q.MinRating > 0 && v.RatingAvg < q.MinRating {
			continue
		}
		if q.Q != "" {
			needle := strings.ToLower(q.Q)
			if !strings.Contains(strings.ToLower(v.Name), needle) &&
				!strings.Contains(strings.ToLower(v.Description), needle) {
				continue
			}
		}
		cp := *v
		cp.PhotoIDs = append([]string{}, r.photos[v.ID]...)
		all = append(all, &cp)
	}
	total := len(all)
	switch q.Sort {
	case "price_asc":
		sort.SliceStable(all, func(i, j int) bool { return all[i].PriceFrom < all[j].PriceFrom })
	case "price_desc":
		sort.SliceStable(all, func(i, j int) bool { return all[i].PriceFrom > all[j].PriceFrom })
	case "rating_desc":
		sort.SliceStable(all, func(i, j int) bool { return all[i].RatingAvg > all[j].RatingAvg })
	default:
		sort.SliceStable(all, func(i, j int) bool { return all[i].CreatedAt.After(all[j].CreatedAt) })
	}
	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit <= 0 {
		q.Limit = 20
	}
	start := (q.Page - 1) * q.Limit
	if start >= len(all) {
		return []*domain.Vendor{}, total, nil
	}
	end := start + q.Limit
	if end > len(all) {
		end = len(all)
	}
	return all[start:end], total, nil
}

func (r *VendorRepo) UpdateStatus(_ context.Context, id string, status domain.VendorStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	v, ok := r.byID[id]
	if !ok {
		return domain.ErrNotFound
	}
	v.Status = status
	v.UpdatedAt = time.Now().UTC()
	return nil
}

func (r *VendorRepo) UpdateRating(_ context.Context, id string, avg float64, count int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	v, ok := r.byID[id]
	if !ok {
		return domain.ErrNotFound
	}
	v.RatingAvg = avg
	v.RatingCount = count
	return nil
}

// AddPhotoID is a test helper to associate photo IDs to a vendor.
func (r *VendorRepo) AddPhotoID(vendorID, photoID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.photos[vendorID] = append(r.photos[vendorID], photoID)
}

// --- Photo repo ---

type PhotoRepo struct {
	mu     sync.Mutex
	byID   map[string]*domain.Photo
	idGen  func() string
	vendor *VendorRepo
}

func NewPhotoRepo(idGen func() string, vendor *VendorRepo) *PhotoRepo {
	return &PhotoRepo{byID: map[string]*domain.Photo{}, idGen: idGen, vendor: vendor}
}

func (r *PhotoRepo) Create(_ context.Context, vendorID, mime string, data []byte) (*domain.Photo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	p := &domain.Photo{
		ID:        r.idGen(),
		VendorID:  vendorID,
		MIME:      mime,
		Size:      int64(len(data)),
		Data:      append([]byte{}, data...),
		CreatedAt: time.Now().UTC(),
	}
	r.byID[p.ID] = p
	if r.vendor != nil {
		r.vendor.AddPhotoID(vendorID, p.ID)
	}
	cp := *p
	return &cp, nil
}

func (r *PhotoRepo) Find(_ context.Context, id string) (*domain.Photo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	p, ok := r.byID[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	cp := *p
	return &cp, nil
}

func (r *PhotoRepo) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.byID, id)
	return nil
}

func (r *PhotoRepo) ListIDs(_ context.Context, vendorID string) ([]string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	ids := []string{}
	for _, p := range r.byID {
		if p.VendorID == vendorID {
			ids = append(ids, p.ID)
		}
	}
	return ids, nil
}

// --- Booking repo ---

type BookingRepo struct {
	mu    sync.Mutex
	byID  map[string]*domain.Booking
	idGen func() string
}

func NewBookingRepo(idGen func() string) *BookingRepo {
	return &BookingRepo{byID: map[string]*domain.Booking{}, idGen: idGen}
}

func (r *BookingRepo) Create(_ context.Context, b *domain.Booking) (*domain.Booking, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *b
	if cp.ID == "" {
		cp.ID = r.idGen()
	}
	if cp.Status == "" {
		cp.Status = domain.BookingPending
	}
	if cp.CreatedAt.IsZero() {
		cp.CreatedAt = time.Now().UTC()
	}
	r.byID[cp.ID] = &cp
	out := cp
	return &out, nil
}

func (r *BookingRepo) Find(_ context.Context, id string) (*domain.Booking, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	b, ok := r.byID[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	cp := *b
	return &cp, nil
}

func (r *BookingRepo) ListForCustomer(_ context.Context, customerID string) ([]*domain.Booking, error) {
	return r.filter(func(b *domain.Booking) bool { return b.CustomerID == customerID }), nil
}

func (r *BookingRepo) ListForVendor(_ context.Context, vendorID string) ([]*domain.Booking, error) {
	return r.filter(func(b *domain.Booking) bool { return b.VendorID == vendorID }), nil
}

func (r *BookingRepo) ListAll(_ context.Context) ([]*domain.Booking, error) {
	return r.filter(func(_ *domain.Booking) bool { return true }), nil
}

func (r *BookingRepo) filter(keep func(*domain.Booking) bool) []*domain.Booking {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := []*domain.Booking{}
	for _, b := range r.byID {
		if keep(b) {
			cp := *b
			out = append(out, &cp)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.After(out[j].CreatedAt) })
	return out
}

func (r *BookingRepo) UpdateStatus(_ context.Context, id string, status domain.BookingStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	b, ok := r.byID[id]
	if !ok {
		return domain.ErrNotFound
	}
	b.Status = status
	return nil
}

func (r *BookingRepo) SetPayment(_ context.Context, id, paymentID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	b, ok := r.byID[id]
	if !ok {
		return domain.ErrNotFound
	}
	b.PaymentID = paymentID
	return nil
}

// --- Review repo ---

type ReviewRepo struct {
	mu      sync.Mutex
	byID    map[string]*domain.Review
	byBook  map[string]string
	idGen   func() string
}

func NewReviewRepo(idGen func() string) *ReviewRepo {
	return &ReviewRepo{
		byID:   map[string]*domain.Review{},
		byBook: map[string]string{},
		idGen:  idGen,
	}
}

func (r *ReviewRepo) Create(_ context.Context, in *domain.Review) (*domain.Review, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.byBook[in.BookingID]; exists {
		return nil, domain.ErrAlreadyExists
	}
	cp := *in
	if cp.ID == "" {
		cp.ID = r.idGen()
	}
	if cp.CreatedAt.IsZero() {
		cp.CreatedAt = time.Now().UTC()
	}
	r.byID[cp.ID] = &cp
	r.byBook[cp.BookingID] = cp.ID
	out := cp
	return &out, nil
}

func (r *ReviewRepo) FindByID(_ context.Context, id string) (*domain.Review, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	rv, ok := r.byID[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	cp := *rv
	return &cp, nil
}

func (r *ReviewRepo) FindByBooking(_ context.Context, bookingID string) (*domain.Review, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	id, ok := r.byBook[bookingID]
	if !ok {
		return nil, domain.ErrNotFound
	}
	cp := *r.byID[id]
	return &cp, nil
}

func (r *ReviewRepo) ListByVendor(_ context.Context, vendorID string) ([]*domain.Review, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := []*domain.Review{}
	for _, rv := range r.byID {
		if rv.VendorID == vendorID {
			cp := *rv
			out = append(out, &cp)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.After(out[j].CreatedAt) })
	return out, nil
}

func (r *ReviewRepo) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	rv, ok := r.byID[id]
	if !ok {
		return domain.ErrNotFound
	}
	delete(r.byBook, rv.BookingID)
	delete(r.byID, id)
	return nil
}

func (r *ReviewRepo) AggregateForVendor(_ context.Context, vendorID string) (float64, int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var sum, count int
	for _, rv := range r.byID {
		if rv.VendorID == vendorID {
			sum += rv.Rating
			count++
		}
	}
	if count == 0 {
		return 0, 0, nil
	}
	return float64(sum) / float64(count), count, nil
}

// --- Refresh token repo ---

type RefreshTokenRepo struct {
	mu       sync.Mutex
	byHash   map[string]*domain.RefreshToken
	idGen    func() string
}

func NewRefreshTokenRepo(idGen func() string) *RefreshTokenRepo {
	return &RefreshTokenRepo{byHash: map[string]*domain.RefreshToken{}, idGen: idGen}
}

func (r *RefreshTokenRepo) Create(_ context.Context, t *domain.RefreshToken) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *t
	if cp.ID == "" {
		cp.ID = r.idGen()
	}
	if cp.CreatedAt.IsZero() {
		cp.CreatedAt = time.Now().UTC()
	}
	r.byHash[cp.TokenHash] = &cp
	return nil
}

func (r *RefreshTokenRepo) FindActiveByHash(_ context.Context, hash string, now time.Time) (*domain.RefreshToken, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	t, ok := r.byHash[hash]
	if !ok {
		return nil, domain.ErrNotFound
	}
	if !t.Active(now) {
		return nil, domain.ErrNotFound
	}
	cp := *t
	return &cp, nil
}

func (r *RefreshTokenRepo) Revoke(_ context.Context, id string, at time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, t := range r.byHash {
		if t.ID == id {
			t.RevokedAt = &at
			return nil
		}
	}
	return domain.ErrNotFound
}

func (r *RefreshTokenRepo) RevokeAllForUser(_ context.Context, userID string, at time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, t := range r.byHash {
		if t.UserID == userID && t.RevokedAt == nil {
			t.RevokedAt = &at
		}
	}
	return nil
}

// --- Password reset repo ---

type PasswordResetRepo struct {
	mu     sync.Mutex
	byHash map[string]*domain.PasswordResetToken
	idGen  func() string
}

func NewPasswordResetRepo(idGen func() string) *PasswordResetRepo {
	return &PasswordResetRepo{byHash: map[string]*domain.PasswordResetToken{}, idGen: idGen}
}

func (r *PasswordResetRepo) Create(_ context.Context, t *domain.PasswordResetToken) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *t
	if cp.ID == "" {
		cp.ID = r.idGen()
	}
	if cp.CreatedAt.IsZero() {
		cp.CreatedAt = time.Now().UTC()
	}
	r.byHash[cp.TokenHash] = &cp
	return nil
}

func (r *PasswordResetRepo) FindByHash(_ context.Context, hash string) (*domain.PasswordResetToken, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	t, ok := r.byHash[hash]
	if !ok {
		return nil, domain.ErrNotFound
	}
	cp := *t
	return &cp, nil
}

func (r *PasswordResetRepo) MarkUsed(_ context.Context, id string, at time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, t := range r.byHash {
		if t.ID == id {
			t.UsedAt = &at
			return nil
		}
	}
	return domain.ErrNotFound
}

// --- Notification repo ---

type NotificationRepo struct {
	mu    sync.Mutex
	byID  map[string]*domain.Notification
	idGen func() string
}

func NewNotificationRepo(idGen func() string) *NotificationRepo {
	return &NotificationRepo{byID: map[string]*domain.Notification{}, idGen: idGen}
}

func (r *NotificationRepo) Create(_ context.Context, n *domain.Notification) (*domain.Notification, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *n
	if cp.ID == "" {
		cp.ID = r.idGen()
	}
	if cp.CreatedAt.IsZero() {
		cp.CreatedAt = time.Now().UTC()
	}
	if cp.Status == "" {
		cp.Status = "queued"
	}
	r.byID[cp.ID] = &cp
	out := cp
	return &out, nil
}

func (r *NotificationRepo) ListForUser(_ context.Context, userID string, limit int) ([]*domain.Notification, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := []*domain.Notification{}
	for _, n := range r.byID {
		if n.UserID == userID {
			cp := *n
			out = append(out, &cp)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.After(out[j].CreatedAt) })
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out, nil
}

func (r *NotificationRepo) MarkSent(_ context.Context, id string) error {
	return r.mark(id, "sent")
}

func (r *NotificationRepo) MarkFailed(_ context.Context, id string) error {
	return r.mark(id, "failed")
}

func (r *NotificationRepo) mark(id, status string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	n, ok := r.byID[id]
	if !ok {
		return domain.ErrNotFound
	}
	n.Status = status
	return nil
}

// --- Hasher (deterministic, for tests) ---

type PlainHasher struct{}

func (PlainHasher) Hash(plain string) (string, error)        { return "h:" + plain, nil }
func (PlainHasher) Verify(hash, plain string) error {
	if hash == "h:"+plain {
		return nil
	}
	return errors.New("mismatch")
}

// --- Token issuer (deterministic, for tests) ---

type FakeIssuer struct {
	Issued map[string]*domain.User
}

func NewFakeIssuer() *FakeIssuer { return &FakeIssuer{Issued: map[string]*domain.User{}} }

func (f *FakeIssuer) Issue(u *domain.User, _ time.Duration) (string, error) {
	tok := "tok-" + u.ID
	f.Issued[tok] = u
	return tok, nil
}

func (f *FakeIssuer) Parse(raw string) (usecase.Claims, error) {
	u, ok := f.Issued[raw]
	if !ok {
		return usecase.Claims{}, domain.ErrUnauthorized
	}
	return usecase.Claims{UserID: u.ID, Email: u.Email, Role: u.Role}, nil
}

// --- Clock / IDGen ---

type FixedClock struct{ T time.Time }

func (c *FixedClock) Now() time.Time { return c.T }
func (c *FixedClock) Advance(d time.Duration) { c.T = c.T.Add(d) }

type SeqIDGen struct {
	mu   sync.Mutex
	N    int
	Prefix string
}

func NewSeqIDGen(prefix string) *SeqIDGen { return &SeqIDGen{Prefix: prefix} }

func (g *SeqIDGen) New() string {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.N++
	return g.Prefix + idToString(g.N)
}

// Func adapts a function to the IDGen interface.
type Func func() string

func (f Func) New() string { return f() }

func idToString(n int) string {
	// avoid strconv import to keep this file dependency-light
	if n == 0 {
		return "0"
	}
	digits := []byte{}
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}
