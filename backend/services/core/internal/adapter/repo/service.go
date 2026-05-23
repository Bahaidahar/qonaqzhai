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

// ServiceRepo persists vendor service menus.
type ServiceRepo struct {
	db    *sql.DB
	idGen ports.IDGen
}

// NewServiceRepo constructs a service repository.
func NewServiceRepo(db *sql.DB, idGen ports.IDGen) *ServiceRepo {
	return &ServiceRepo{db: db, idGen: idGen}
}

const serviceCols = `id, vendor_id, name, description, price, unit, is_active, created_at, updated_at`

// Create inserts a service.
func (r *ServiceRepo) Create(ctx context.Context, s *domain.Service) (*domain.Service, error) {
	if s.ID == "" {
		s.ID = r.idGen.New()
	}
	if s.Unit == "" {
		s.Unit = domain.UnitFixed
	}
	if _, err := r.db.ExecContext(ctx,
		`INSERT INTO services (id, vendor_id, name, description, price, unit, is_active)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		s.ID, s.VendorID, s.Name, s.Description, s.Price, string(s.Unit), s.IsActive,
	); err != nil {
		return nil, fmt.Errorf("insert service: %w", err)
	}
	return r.Find(ctx, s.ID)
}

// Update mutates an existing service. Nil IsActive leaves the column alone.
func (r *ServiceRepo) Update(ctx context.Context, id string, in domain.ServiceInput) (*domain.Service, error) {
	if in.IsActive == nil {
		if _, err := r.db.ExecContext(ctx,
			`UPDATE services SET name=$1, description=$2, price=$3, unit=$4, updated_at=now() WHERE id=$5`,
			in.Name, in.Description, in.Price, string(in.Unit), id,
		); err != nil {
			return nil, fmt.Errorf("update service: %w", err)
		}
	} else {
		if _, err := r.db.ExecContext(ctx,
			`UPDATE services SET name=$1, description=$2, price=$3, unit=$4, is_active=$5, updated_at=now() WHERE id=$6`,
			in.Name, in.Description, in.Price, string(in.Unit), *in.IsActive, id,
		); err != nil {
			return nil, fmt.Errorf("update service: %w", err)
		}
	}
	return r.Find(ctx, id)
}

// Delete removes a service.
func (r *ServiceRepo) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM services WHERE id = $1`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errs.ErrNotFound
	}
	return nil
}

// Find returns a service by primary key.
func (r *ServiceRepo) Find(ctx context.Context, id string) (*domain.Service, error) {
	row := r.db.QueryRowContext(ctx, `SELECT `+serviceCols+` FROM services WHERE id = $1`, id)
	s, err := scanService(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errs.ErrNotFound
	}
	return s, err
}

// ListByVendor returns all services owned by vendorID.
func (r *ServiceRepo) ListByVendor(ctx context.Context, vendorID string) ([]*domain.Service, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT `+serviceCols+` FROM services WHERE vendor_id = $1 ORDER BY created_at ASC`,
		vendorID,
	)
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

func scanService(s scanner) (*domain.Service, error) {
	var v domain.Service
	var unit string
	if err := s.Scan(&v.ID, &v.VendorID, &v.Name, &v.Description, &v.Price, &unit, &v.IsActive, &v.CreatedAt, &v.UpdatedAt); err != nil {
		return nil, err
	}
	v.Unit = domain.ServiceUnit(unit)
	return &v, nil
}

var _ ports.ServiceRepo = (*ServiceRepo)(nil)
