package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"qonaqzhai-backend/internal/domain"
	"qonaqzhai-backend/internal/usecase"
)

// CardRepo persists saved payment cards (mock — last4 only).
type CardRepo struct {
	db    *sql.DB
	idGen usecase.IDGen
}

// NewCardRepo constructs a card repository.
func NewCardRepo(db *sql.DB, idGen usecase.IDGen) *CardRepo {
	return &CardRepo{db: db, idGen: idGen}
}

// Create inserts a new card. If user has no cards yet, the new one is marked default.
func (r *CardRepo) Create(ctx context.Context, c *domain.PaymentCard) (*domain.PaymentCard, error) {
	if c.ID == "" {
		c.ID = r.idGen.New()
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin: %w", err)
	}
	defer tx.Rollback()
	var count int
	if err := tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM payment_cards WHERE user_id = $1`, c.UserID).Scan(&count); err != nil {
		return nil, fmt.Errorf("count cards: %w", err)
	}
	isDefault := count == 0 || c.IsDefault
	if isDefault {
		if _, err := tx.ExecContext(ctx, `UPDATE payment_cards SET is_default = FALSE WHERE user_id = $1`, c.UserID); err != nil {
			return nil, fmt.Errorf("clear default: %w", err)
		}
	}
	if _, err := tx.ExecContext(ctx,
		`INSERT INTO payment_cards (id, user_id, brand, last4, exp_month, exp_year, holder, is_default)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		c.ID, c.UserID, c.Brand, c.Last4, c.ExpMonth, c.ExpYear, c.Holder, isDefault,
	); err != nil {
		return nil, fmt.Errorf("insert card: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}
	return r.FindByID(ctx, c.ID)
}

// FindByID looks up a card by primary key.
func (r *CardRepo) FindByID(ctx context.Context, id string) (*domain.PaymentCard, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, brand, last4, exp_month, exp_year, holder, is_default, created_at
		 FROM payment_cards WHERE id = $1`, id,
	)
	return scanCard(row)
}

// ListForUser returns all cards for a user, newest first, default first.
func (r *CardRepo) ListForUser(ctx context.Context, userID string) ([]*domain.PaymentCard, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, user_id, brand, last4, exp_month, exp_year, holder, is_default, created_at
		 FROM payment_cards WHERE user_id = $1
		 ORDER BY is_default DESC, created_at DESC`, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("list cards: %w", err)
	}
	defer rows.Close()
	out := []*domain.PaymentCard{}
	for rows.Next() {
		c, err := scanCard(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// Delete removes a card. If it was default, a remaining card is promoted.
func (r *CardRepo) Delete(ctx context.Context, id string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin: %w", err)
	}
	defer tx.Rollback()
	var userID string
	var isDefault bool
	if err := tx.QueryRowContext(ctx, `SELECT user_id, is_default FROM payment_cards WHERE id = $1`, id).Scan(&userID, &isDefault); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrNotFound
		}
		return fmt.Errorf("lookup card: %w", err)
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM payment_cards WHERE id = $1`, id); err != nil {
		return fmt.Errorf("delete card: %w", err)
	}
	if isDefault {
		if _, err := tx.ExecContext(ctx,
			`UPDATE payment_cards SET is_default = TRUE WHERE id = (
			   SELECT id FROM payment_cards WHERE user_id = $1 ORDER BY created_at DESC LIMIT 1
			 )`, userID,
		); err != nil {
			return fmt.Errorf("promote default: %w", err)
		}
	}
	return tx.Commit()
}

// SetDefault marks the card as default for the user, clearing others.
func (r *CardRepo) SetDefault(ctx context.Context, userID, cardID string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin: %w", err)
	}
	defer tx.Rollback()
	var owner string
	if err := tx.QueryRowContext(ctx, `SELECT user_id FROM payment_cards WHERE id = $1`, cardID).Scan(&owner); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrNotFound
		}
		return fmt.Errorf("lookup card: %w", err)
	}
	if owner != userID {
		return domain.ErrForbidden
	}
	if _, err := tx.ExecContext(ctx, `UPDATE payment_cards SET is_default = FALSE WHERE user_id = $1`, userID); err != nil {
		return fmt.Errorf("clear default: %w", err)
	}
	if _, err := tx.ExecContext(ctx, `UPDATE payment_cards SET is_default = TRUE WHERE id = $1`, cardID); err != nil {
		return fmt.Errorf("set default: %w", err)
	}
	return tx.Commit()
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanCard(row rowScanner) (*domain.PaymentCard, error) {
	var c domain.PaymentCard
	if err := row.Scan(&c.ID, &c.UserID, &c.Brand, &c.Last4, &c.ExpMonth, &c.ExpYear, &c.Holder, &c.IsDefault, &c.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &c, nil
}

var _ usecase.CardRepo = (*CardRepo)(nil)
