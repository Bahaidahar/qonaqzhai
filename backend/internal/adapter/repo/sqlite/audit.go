package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"qonaqzhai-backend/internal/domain"
	"qonaqzhai-backend/internal/usecase"
)

// AuditRepo persists audit log entries.
type AuditRepo struct {
	db    *sql.DB
	idGen usecase.IDGen
}

// NewAuditRepo constructs an audit repository.
func NewAuditRepo(db *sql.DB, idGen usecase.IDGen) *AuditRepo {
	return &AuditRepo{db: db, idGen: idGen}
}

// Create inserts an audit entry.
func (r *AuditRepo) Create(ctx context.Context, e *domain.AuditEntry) error {
	if e.ID == "" {
		e.ID = r.idGen.New()
	}
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO audit_log (id, actor_id, actor_email, action, target_type, target_id, meta)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		e.ID, e.ActorID, e.ActorEmail, e.Action, e.TargetType, e.TargetID, e.Meta,
	)
	if err != nil {
		return fmt.Errorf("insert audit: %w", err)
	}
	return nil
}

// List returns the latest entries (newest first).
func (r *AuditRepo) List(ctx context.Context, limit int) ([]*domain.AuditEntry, error) {
	if limit <= 0 {
		limit = 100
	}
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, actor_id, actor_email, action, target_type, target_id, meta, created_at
		 FROM audit_log ORDER BY created_at DESC LIMIT ?`, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []*domain.AuditEntry{}
	for rows.Next() {
		var e domain.AuditEntry
		if err := rows.Scan(
			&e.ID, &e.ActorID, &e.ActorEmail, &e.Action,
			&e.TargetType, &e.TargetID, &e.Meta, &e.CreatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, &e)
	}
	return out, rows.Err()
}

var _ usecase.AuditRepo = (*AuditRepo)(nil)
