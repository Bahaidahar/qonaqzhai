package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"qonaqzhai-backend/internal/domain"
	"qonaqzhai-backend/internal/usecase"
)

// ServiceRepo persists per-vendor services (menu items).
type ServiceRepo struct {
	db    *sql.DB
	idGen usecase.IDGen
}

// NewServiceRepo constructs a service repository.
func NewServiceRepo(db *sql.DB, idGen usecase.IDGen) *ServiceRepo {
	return &ServiceRepo{db: db, idGen: idGen}
}

// Create inserts a new service.
func (r *ServiceRepo) Create(ctx context.Context, vendorID string, in domain.ServiceInput) (*domain.Service, error) {
	active := 1
	if in.IsActive != nil && !*in.IsActive {
		active = 0
	}
	id := r.idGen.New()
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO services (id, vendor_id, name, description, price, unit, is_active)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		id, vendorID, in.Name, in.Description, in.Price, string(in.Unit), active,
	)
	if err != nil {
		return nil, fmt.Errorf("insert service: %w", err)
	}
	return r.FindByID(ctx, id)
}

// Update modifies an existing service.
func (r *ServiceRepo) Update(ctx context.Context, id string, in domain.ServiceInput) (*domain.Service, error) {
	args := []any{in.Name, in.Description, in.Price, string(in.Unit)}
	stmt := `UPDATE services SET name=?, description=?, price=?, unit=?, updated_at=CURRENT_TIMESTAMP`
	if in.IsActive != nil {
		stmt += `, is_active=?`
		v := 0
		if *in.IsActive {
			v = 1
		}
		args = append(args, v)
	}
	stmt += ` WHERE id=?`
	args = append(args, id)

	res, err := r.db.ExecContext(ctx, stmt, args...)
	if err != nil {
		return nil, fmt.Errorf("update service: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return nil, domain.ErrNotFound
	}
	return r.FindByID(ctx, id)
}

// FindByID returns a service by primary key.
func (r *ServiceRepo) FindByID(ctx context.Context, id string) (*domain.Service, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, vendor_id, name, description, price, unit, is_active, created_at, updated_at
		 FROM services WHERE id = ?`, id,
	)
	s, err := scanService(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return s, err
}

// ListByVendor returns services for a vendor, ordered by creation.
func (r *ServiceRepo) ListByVendor(ctx context.Context, vendorID string, activeOnly bool) ([]*domain.Service, error) {
	q := `SELECT id, vendor_id, name, description, price, unit, is_active, created_at, updated_at
	      FROM services WHERE vendor_id = ?`
	if activeOnly {
		q += ` AND is_active = 1`
	}
	q += ` ORDER BY created_at ASC`
	rows, err := r.db.QueryContext(ctx, q, vendorID)
	if err != nil {
		return nil, fmt.Errorf("list services: %w", err)
	}
	defer rows.Close()
	out := []*domain.Service{}
	for rows.Next() {
		s, err := scanService(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

// Delete removes a service.
func (r *ServiceRepo) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM services WHERE id = ?`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// MinActivePrice returns the lowest active service price for a vendor.
// Returns 0 when the vendor has no active services.
func (r *ServiceRepo) MinActivePrice(ctx context.Context, vendorID string) (int64, error) {
	var price sql.NullInt64
	if err := r.db.QueryRowContext(ctx,
		`SELECT MIN(price) FROM services WHERE vendor_id = ? AND is_active = 1`, vendorID,
	).Scan(&price); err != nil {
		return 0, err
	}
	if !price.Valid {
		return 0, nil
	}
	return price.Int64, nil
}

func scanService(s scanner) (*domain.Service, error) {
	var srv domain.Service
	var unit string
	var active int
	if err := s.Scan(
		&srv.ID, &srv.VendorID, &srv.Name, &srv.Description,
		&srv.Price, &unit, &active, &srv.CreatedAt, &srv.UpdatedAt,
	); err != nil {
		return nil, err
	}
	srv.Unit = domain.ServiceUnit(unit)
	srv.IsActive = active != 0
	return &srv, nil
}

var _ usecase.ServiceRepo = (*ServiceRepo)(nil)
