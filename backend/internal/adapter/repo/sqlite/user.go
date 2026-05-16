// Package sqlite implements the persistence ports against a SQLite database.
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

// UserRepo persists users in SQLite.
type UserRepo struct {
	db    *sql.DB
	idGen usecase.IDGen
}

// NewUserRepo constructs a user repository bound to db.
func NewUserRepo(db *sql.DB, idGen usecase.IDGen) *UserRepo {
	return &UserRepo{db: db, idGen: idGen}
}

// Create inserts a new user, generating an ID if one is not set.
func (r *UserRepo) Create(ctx context.Context, u *domain.User) (*domain.User, error) {
	if u.ID == "" {
		u.ID = r.idGen.New()
	}
	if u.Status == "" {
		u.Status = domain.UserActive
	}
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO users (id, email, name, password_hash, role, status) VALUES (?, ?, ?, ?, ?, ?)`,
		u.ID, u.Email, u.Name, u.PasswordHash, string(u.Role), string(u.Status),
	)
	if err != nil {
		if isUniqueErr(err) {
			return nil, domain.ErrAlreadyExists
		}
		return nil, fmt.Errorf("insert user: %w", err)
	}
	return r.FindByID(ctx, u.ID)
}

// FindByID retrieves a user by primary key.
func (r *UserRepo) FindByID(ctx context.Context, id string) (*domain.User, error) {
	return r.queryUser(ctx, `WHERE id = ?`, id)
}

// FindByEmail retrieves a user by (unique) email.
func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	return r.queryUser(ctx, `WHERE email = ?`, email)
}

// List returns all users ordered by created_at descending.
func (r *UserRepo) List(ctx context.Context) ([]*domain.User, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, email, name, password_hash, role, status, created_at FROM users ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()
	out := []*domain.User{}
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, rows.Err()
}

// UpdateStatus changes a user's lifecycle status.
func (r *UserRepo) UpdateStatus(ctx context.Context, id string, status domain.UserStatus) error {
	res, err := r.db.ExecContext(ctx, `UPDATE users SET status = ? WHERE id = ?`, string(status), id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// UpdatePasswordHash overwrites the stored password hash.
func (r *UserRepo) UpdatePasswordHash(ctx context.Context, id, hash string) error {
	res, err := r.db.ExecContext(ctx, `UPDATE users SET password_hash = ? WHERE id = ?`, hash, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *UserRepo) queryUser(ctx context.Context, where string, args ...any) (*domain.User, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, email, name, password_hash, role, status, created_at FROM users `+where,
		args...,
	)
	u, err := scanUser(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return u, err
}

type scanner interface {
	Scan(dest ...any) error
}

func scanUser(s scanner) (*domain.User, error) {
	var u domain.User
	var role, status string
	if err := s.Scan(&u.ID, &u.Email, &u.Name, &u.PasswordHash, &role, &status, &u.CreatedAt); err != nil {
		return nil, err
	}
	u.Role = domain.Role(role)
	u.Status = domain.UserStatus(status)
	return &u, nil
}

func isUniqueErr(err error) bool {
	return err != nil && strings.Contains(err.Error(), "UNIQUE constraint failed")
}
