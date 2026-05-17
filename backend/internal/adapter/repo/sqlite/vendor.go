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
			`UPDATE vendors SET name=$1, category=$2, city=$3, description=$4, price_from=$5, updated_at=now() WHERE id=$6`,
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
		`INSERT INTO vendors (id, user_id, name, category, city, description, price_from)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		id, userID, in.Name, in.Category, in.City, in.Description, in.PriceFrom,
	); err != nil {
		return nil, fmt.Errorf("insert vendor: %w", err)
	}
	return r.FindByID(ctx, id)
}

// FindByID returns the vendor by primary key (with photo IDs populated).
func (r *VendorRepo) FindByID(ctx context.Context, id string) (*domain.Vendor, error) {
	return r.queryVendor(ctx, `WHERE id = $1`, id)
}

// FindByUserID returns the vendor owned by userID.
func (r *VendorRepo) FindByUserID(ctx context.Context, userID string) (*domain.Vendor, error) {
	return r.queryVendor(ctx, `WHERE user_id = $1`, userID)
}

// Search runs the catalog query with filters, pagination, and sorting.
// Full-text search uses the `search_tsv` generated tsvector column (GIN-indexed).
func (r *VendorRepo) Search(ctx context.Context, q usecase.VendorQuery) ([]*domain.Vendor, int, error) {
	var (
		wheres []string
		args   []any
	)
	add := func(template string, val any) {
		args = append(args, val)
		wheres = append(wheres, strings.Replace(template, "$?", fmt.Sprintf("$%d", len(args)), 1))
	}

	if s := strings.TrimSpace(q.Q); s != "" {
		add(`v.search_tsv @@ plainto_tsquery('simple', $?)`, s)
	}
	if q.Status != "" {
		add(`v.status = $?`, string(q.Status))
	}
	if q.Category != "" {
		add(`v.category = $?`, q.Category)
	}
	if q.City != "" {
		add(`v.city = $?`, q.City)
	}
	if q.MinPrice > 0 {
		add(`v.price_from >= $?`, q.MinPrice)
	}
	if q.MaxPrice > 0 {
		add(`v.price_from <= $?`, q.MaxPrice)
	}
	if q.MinRating > 0 {
		add(`v.rating_avg >= $?`, q.MinRating)
	}

	whereSQL := ""
	if len(wheres) > 0 {
		whereSQL = " WHERE " + strings.Join(wheres, " AND ")
	}

	var total int
	if err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM vendors v`+whereSQL, args...,
	).Scan(&total); err != nil {
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

	args = append(args, limit, offset)
	limitOffset := fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)-1, len(args))

	rows, err := r.db.QueryContext(ctx,
		`SELECT v.id, v.user_id, v.name, v.category, v.city, v.description, v.price_from, v.status, v.rating_avg, v.rating_count, v.created_at, v.updated_at
		 FROM vendors v`+whereSQL+orderBy+limitOffset,
		args...,
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
		`UPDATE vendors SET status = $1, updated_at = now() WHERE id = $2`,
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
		`UPDATE vendors SET rating_avg = $1, rating_count = $2, updated_at = now() WHERE id = $3`,
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
	rows, err := db.QueryContext(ctx,
		`SELECT id FROM photos WHERE vendor_id = $1 ORDER BY created_at ASC`, vendorID)
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
