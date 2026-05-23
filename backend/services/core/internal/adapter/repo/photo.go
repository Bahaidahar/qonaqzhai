package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"qonaqzhai-backend/pkg/errs"
	"qonaqzhai-backend/services/core/internal/domain"
	"qonaqzhai-backend/services/core/internal/ports"
)

// PhotoRepo persists vendor photos as bytea.
type PhotoRepo struct {
	db    *sql.DB
	idGen ports.IDGen
}

// NewPhotoRepo constructs a photo repository.
func NewPhotoRepo(db *sql.DB, idGen ports.IDGen) *PhotoRepo {
	return &PhotoRepo{db: db, idGen: idGen}
}

// Insert persists a new photo.
func (r *PhotoRepo) Insert(ctx context.Context, p *domain.Photo) (*domain.Photo, error) {
	if p.ID == "" {
		p.ID = r.idGen.New()
	}
	if _, err := r.db.ExecContext(ctx,
		`INSERT INTO photos (id, vendor_id, mime, size, data) VALUES ($1, $2, $3, $4, $5)`,
		p.ID, p.VendorID, p.MIME, p.Size, p.Data,
	); err != nil {
		return nil, fmt.Errorf("insert photo: %w", err)
	}
	return r.Find(ctx, p.ID)
}

// Find returns a photo by primary key (data included).
func (r *PhotoRepo) Find(ctx context.Context, id string) (*domain.Photo, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, vendor_id, mime, size, data, created_at FROM photos WHERE id = $1`, id,
	)
	var p domain.Photo
	if err := row.Scan(&p.ID, &p.VendorID, &p.MIME, &p.Size, &p.Data, &p.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}
	return &p, nil
}

// Delete removes a photo.
func (r *PhotoRepo) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM photos WHERE id = $1`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errs.ErrNotFound
	}
	return nil
}

// ListByVendor returns metadata (no bytes) for vendor's photos.
func (r *PhotoRepo) ListByVendor(ctx context.Context, vendorID string) ([]*domain.Photo, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, vendor_id, mime, size, created_at FROM photos WHERE vendor_id = $1 ORDER BY created_at ASC`,
		vendorID,
	)
	if err != nil {
		return nil, fmt.Errorf("list photos: %w", err)
	}
	defer rows.Close()
	out := []*domain.Photo{}
	for rows.Next() {
		var p domain.Photo
		if err := rows.Scan(&p.ID, &p.VendorID, &p.MIME, &p.Size, &p.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, &p)
	}
	return out, rows.Err()
}

var _ ports.PhotoRepo = (*PhotoRepo)(nil)
