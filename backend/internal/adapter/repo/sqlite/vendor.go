package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"qonaqzhai-backend/internal/domain"
	"qonaqzhai-backend/internal/usecase"
)

// VendorRepo persists vendor profiles.
type VendorRepo struct {
	db    *sql.DB
	idGen usecase.IDGen
}

// NewVendorRepo constructs a vendor repository bound to db.
func NewVendorRepo(db *sql.DB, idGen usecase.IDGen) *VendorRepo {
	return &VendorRepo{db: db, idGen: idGen}
}

// Upsert creates or updates the user's vendor profile.
func (r *VendorRepo) Upsert(ctx context.Context, userID string, in domain.VendorInput) (*domain.Vendor, error) {
	existing, err := r.FindByUserID(ctx, userID)
	if err == nil {
		if _, err := r.db.ExecContext(ctx,
			`UPDATE vendors SET name=?, category=?, city=?, description=?, price_from=?, updated_at=CURRENT_TIMESTAMP WHERE id=?`,
			in.Name, in.Category, in.City, in.Description, in.PriceFrom, existing.ID,
		); err != nil {
			return nil, fmt.Errorf("update vendor: %w", err)
		}
		return r.FindByID(ctx, existing.ID)
	}
	if !errors.Is(err, domain.ErrNotFound) {
		return nil, err
	}
	id := r.idGen.New()
	if _, err := r.db.ExecContext(ctx,
		`INSERT INTO vendors (id, user_id, name, category, city, description, price_from) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		id, userID, in.Name, in.Category, in.City, in.Description, in.PriceFrom,
	); err != nil {
		return nil, fmt.Errorf("insert vendor: %w", err)
	}
	return r.FindByID(ctx, id)
}

// FindByID returns the vendor by primary key (with photo IDs populated).
func (r *VendorRepo) FindByID(ctx context.Context, id string) (*domain.Vendor, error) {
	return r.queryVendor(ctx, `WHERE id = ?`, id)
}

// FindByUserID returns the vendor owned by userID.
func (r *VendorRepo) FindByUserID(ctx context.Context, userID string) (*domain.Vendor, error) {
	return r.queryVendor(ctx, `WHERE user_id = ?`, userID)
}

// Search runs the catalog query with all filters, pagination, sorting.
// Returns the page slice + total count of matching rows.
func (r *VendorRepo) Search(ctx context.Context, q usecase.VendorQuery) ([]*domain.Vendor, int, error) {
	var (
		wheres []string
		args   []any
	)

	useFTS := strings.TrimSpace(q.Q) != ""
	from := `FROM vendors v`
	if useFTS {
		from = `FROM vendors v JOIN vendors_fts f ON v.rowid = f.rowid`
		wheres = append(wheres, `vendors_fts MATCH ?`)
		args = append(args, ftsQuery(q.Q))
	}
	if q.Status != "" {
		wheres = append(wheres, `v.status = ?`)
		args = append(args, string(q.Status))
	}
	if q.Category != "" {
		wheres = append(wheres, `v.category = ?`)
		args = append(args, q.Category)
	}
	if q.City != "" {
		wheres = append(wheres, `v.city = ?`)
		args = append(args, q.City)
	}
	if q.MinPrice > 0 {
		wheres = append(wheres, `v.price_from >= ?`)
		args = append(args, q.MinPrice)
	}
	if q.MaxPrice > 0 {
		wheres = append(wheres, `v.price_from <= ?`)
		args = append(args, q.MaxPrice)
	}
	if q.MinRating > 0 {
		wheres = append(wheres, `v.rating_avg >= ?`)
		args = append(args, q.MinRating)
	}

	whereSQL := ""
	if len(wheres) > 0 {
		whereSQL = " WHERE " + strings.Join(wheres, " AND ")
	}

	// Count
	var total int
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) "+from+whereSQL, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count vendors: %w", err)
	}

	orderBy := " ORDER BY v.created_at DESC"
	switch q.Sort {
	case "price_asc":
		orderBy = " ORDER BY v.price_from ASC"
	case "price_desc":
		orderBy = " ORDER BY v.price_from DESC"
	case "rating_desc":
		orderBy = " ORDER BY v.rating_avg DESC, v.rating_count DESC"
	}

	page := q.Page
	if page < 1 {
		page = 1
	}
	limit := q.Limit
	if limit <= 0 {
		limit = 20
	}
	offset := (page - 1) * limit

	rows, err := r.db.QueryContext(ctx,
		`SELECT v.id, v.user_id, v.name, v.category, v.city, v.description, v.price_from, v.status, v.rating_avg, v.rating_count, v.created_at, v.updated_at `+
			from+whereSQL+orderBy+` LIMIT ? OFFSET ?`,
		append(args, limit, offset)...,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("list vendors: %w", err)
	}
	defer rows.Close()

	out := []*domain.Vendor{}
	for rows.Next() {
		v, err := scanVendor(rows)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	for _, v := range out {
		ids, err := listPhotoIDs(ctx, r.db, v.ID)
		if err != nil {
			return nil, 0, err
		}
		v.PhotoIDs = ids
	}
	return out, total, nil
}

// UpdateStatus moves a vendor to next status.
func (r *VendorRepo) UpdateStatus(ctx context.Context, id string, status domain.VendorStatus) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE vendors SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		string(status), id,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// UpdateRating persists a recomputed rating aggregate.
func (r *VendorRepo) UpdateRating(ctx context.Context, id string, avg float64, count int) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE vendors SET rating_avg = ?, rating_count = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		avg, count, id,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *VendorRepo) queryVendor(ctx context.Context, where string, args ...any) (*domain.Vendor, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, name, category, city, description, price_from, status, rating_avg, rating_count, created_at, updated_at FROM vendors `+where,
		args...,
	)
	v, err := scanVendor(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	ids, err := listPhotoIDs(ctx, r.db, v.ID)
	if err != nil {
		return nil, err
	}
	v.PhotoIDs = ids
	return v, nil
}

func scanVendor(s scanner) (*domain.Vendor, error) {
	var v domain.Vendor
	var status string
	if err := s.Scan(
		&v.ID, &v.UserID, &v.Name, &v.Category, &v.City,
		&v.Description, &v.PriceFrom, &status,
		&v.RatingAvg, &v.RatingCount,
		&v.CreatedAt, &v.UpdatedAt,
	); err != nil {
		return nil, err
	}
	v.Status = domain.VendorStatus(status)
	return &v, nil
}

func listPhotoIDs(ctx context.Context, db *sql.DB, vendorID string) ([]string, error) {
	rows, err := db.QueryContext(ctx, `SELECT id FROM photos WHERE vendor_id = ? ORDER BY created_at ASC`, vendorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ids := []string{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

// ftsQuery wraps the user input into an FTS5 prefix expression on every token.
// Returns a safe form that never crashes on operator characters in user input.
func ftsQuery(raw string) string {
	tokens := strings.Fields(raw)
	for i, t := range tokens {
		// strip FTS5 reserved characters by quoting each token, append * for prefix matching
		t = strings.ReplaceAll(t, `"`, `""`)
		tokens[i] = `"` + t + `"` + "*"
	}
	if len(tokens) == 0 {
		return ""
	}
	return strings.Join(tokens, " ")
}
