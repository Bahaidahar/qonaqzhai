package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"qonaqzhai-backend/internal/domain"
	"qonaqzhai-backend/internal/usecase"
)

// PhotoRepo persists vendor photo blobs.
type PhotoRepo struct {
	db    *sql.DB
	idGen usecase.IDGen
}

// NewPhotoRepo constructs a photo repository.
func NewPhotoRepo(db *sql.DB, idGen usecase.IDGen) *PhotoRepo {
	return &PhotoRepo{db: db, idGen: idGen}
}

// Create stores a photo and returns its metadata.
func (r *PhotoRepo) Create(ctx context.Context, vendorID, mime string, data []byte) (*domain.Photo, error) {
	id := r.idGen.New()
	if _, err := r.db.ExecContext(ctx,
		`INSERT INTO photos (id, vendor_id, mime, size, data) VALUES (?, ?, ?, ?, ?)`,
		id, vendorID, mime, len(data), data,
	); err != nil {
		return nil, fmt.Errorf("insert photo: %w", err)
	}
	return r.Find(ctx, id)
}

// Find returns the photo (including raw bytes) by id.
func (r *PhotoRepo) Find(ctx context.Context, id string) (*domain.Photo, error) {
	var p domain.Photo
	err := r.db.QueryRowContext(ctx,
		`SELECT id, vendor_id, mime, size, data, created_at FROM photos WHERE id = ?`, id,
	).Scan(&p.ID, &p.VendorID, &p.MIME, &p.Size, &p.Data, &p.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// Delete removes a photo by id (idempotent — missing IDs return nil).
func (r *PhotoRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM photos WHERE id = ?`, id)
	return err
}

// ListIDs returns photo IDs belonging to a vendor, oldest first.
func (r *PhotoRepo) ListIDs(ctx context.Context, vendorID string) ([]string, error) {
	return listPhotoIDs(ctx, r.db, vendorID)
}
