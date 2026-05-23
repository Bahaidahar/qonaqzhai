package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"qonaqzhai-backend/pkg/errs"
	"qonaqzhai-backend/services/payment/internal/domain"
	"qonaqzhai-backend/services/payment/internal/ports"
)

// CardRepo persists saved cards.
type CardRepo struct {
	db    *sql.DB
	idGen ports.IDGen
}

// NewCardRepo constructs a card repository.
func NewCardRepo(db *sql.DB, idGen ports.IDGen) *CardRepo {
	return &CardRepo{db: db, idGen: idGen}
}

const cardCols = `id, user_id, brand, last4, exp_month, exp_year, holder, is_default, created_at`

// Create inserts a card. If first card for the user, it is marked default.
func (r *CardRepo) Create(ctx context.Context, c *domain.Card) (*domain.Card, error) {
	if c.ID == "" {
		c.ID = r.idGen.New()
	}
	var count int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM cards WHERE user_id = $1`, c.UserID).Scan(&count); err != nil {
		return nil, err
	}
	if count == 0 {
		c.IsDefault = true
	}
	if _, err := r.db.ExecContext(ctx,
		`INSERT INTO cards (id, user_id, brand, last4, exp_month, exp_year, holder, is_default)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		c.ID, c.UserID, c.Brand, c.Last4, c.ExpMonth, c.ExpYear, c.Holder, c.IsDefault,
	); err != nil {
		return nil, fmt.Errorf("insert card: %w", err)
	}
	return r.Find(ctx, c.ID)
}

// Find returns a card by id.
func (r *CardRepo) Find(ctx context.Context, id string) (*domain.Card, error) {
	row := r.db.QueryRowContext(ctx, `SELECT `+cardCols+` FROM cards WHERE id = $1`, id)
	c, err := scanCard(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errs.ErrNotFound
	}
	return c, err
}

// ListByUser returns the user's cards (default first).
func (r *CardRepo) ListByUser(ctx context.Context, userID string) ([]*domain.Card, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT `+cardCols+` FROM cards WHERE user_id = $1 ORDER BY is_default DESC, created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("list cards: %w", err)
	}
	defer rows.Close()
	out := []*domain.Card{}
	for rows.Next() {
		c, err := scanCard(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// Delete removes a card. If the deleted card was default, promotes the next.
func (r *CardRepo) Delete(ctx context.Context, id string) error {
	c, err := r.Find(ctx, id)
	if err != nil {
		return err
	}
	if _, err := r.db.ExecContext(ctx, `DELETE FROM cards WHERE id = $1`, id); err != nil {
		return err
	}
	if c.IsDefault {
		_, _ = r.db.ExecContext(ctx,
			`UPDATE cards SET is_default = TRUE
			 WHERE id = (SELECT id FROM cards WHERE user_id = $1 ORDER BY created_at ASC LIMIT 1)`,
			c.UserID,
		)
	}
	return nil
}

// SetDefault flips the default flag for the given card and demotes siblings.
func (r *CardRepo) SetDefault(ctx context.Context, userID, id string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, `UPDATE cards SET is_default = FALSE WHERE user_id = $1`, userID); err != nil {
		return err
	}
	res, err := tx.ExecContext(ctx, `UPDATE cards SET is_default = TRUE WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errs.ErrNotFound
	}
	return tx.Commit()
}

func scanCard(s scanner) (*domain.Card, error) {
	var c domain.Card
	if err := s.Scan(&c.ID, &c.UserID, &c.Brand, &c.Last4, &c.ExpMonth, &c.ExpYear, &c.Holder, &c.IsDefault, &c.CreatedAt); err != nil {
		return nil, err
	}
	return &c, nil
}

var _ ports.CardRepo = (*CardRepo)(nil)
