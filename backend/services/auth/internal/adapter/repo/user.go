package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"qonaqzhai-backend/pkg/errs"
	"qonaqzhai-backend/services/auth/internal/domain"
	"qonaqzhai-backend/services/auth/internal/ports"
)

// UserRepo persists users in PostgreSQL.
type UserRepo struct {
	db    *sql.DB
	idGen ports.IDGen
}

// NewUserRepo constructs a user repository bound to db.
func NewUserRepo(db *sql.DB, idGen ports.IDGen) *UserRepo {
	return &UserRepo{db: db, idGen: idGen}
}

// Create inserts a new user, generating an ID if absent. Returns
// errs.ErrAlreadyExists on email conflict.
func (r *UserRepo) Create(ctx context.Context, u *domain.User) (*domain.User, error) {
	if u.ID == "" {
		u.ID = r.idGen.New()
	}
	if u.Status == "" {
		u.Status = domain.UserActive
	}
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO users (id, email, name, password_hash, role, status)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		u.ID, u.Email, u.Name, u.PasswordHash, string(u.Role), string(u.Status),
	)
	if err != nil {
		if isUniqueErr(err) {
			return nil, errs.ErrAlreadyExists
		}
		return nil, fmt.Errorf("insert user: %w", err)
	}
	return r.FindByID(ctx, u.ID)
}

// FindByID retrieves a user by primary key. Returns errs.ErrNotFound when
// absent.
func (r *UserRepo) FindByID(ctx context.Context, id string) (*domain.User, error) {
	return r.queryUser(ctx, `WHERE id = $1`, id)
}

// FindByEmail retrieves a user by unique email.
func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	return r.queryUser(ctx, `WHERE email = $1`, email)
}

// FindByIDs returns the subset of users matching the requested ids, in
// arbitrary order. Missing ids are silently omitted.
func (r *UserRepo) FindByIDs(ctx context.Context, ids []string) ([]*domain.User, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, email, name, password_hash, role, status, created_at
		 FROM users WHERE id = ANY($1)`,
		ids,
	)
	if err != nil {
		return nil, fmt.Errorf("find users by ids: %w", err)
	}
	defer rows.Close()
	return scanUsers(rows)
}

// List returns paginated users, newest first.
func (r *UserRepo) List(ctx context.Context, opts ports.ListUsersOpts) ([]*domain.User, error) {
	q := `SELECT id, email, name, password_hash, role, status, created_at FROM users`
	args := []any{}
	if opts.Role != "" {
		q += ` WHERE role = $1`
		args = append(args, opts.Role)
	}
	q += ` ORDER BY created_at DESC`
	if opts.Limit > 0 {
		q += fmt.Sprintf(` LIMIT $%d`, len(args)+1)
		args = append(args, opts.Limit)
	}
	if opts.Offset > 0 {
		q += fmt.Sprintf(` OFFSET $%d`, len(args)+1)
		args = append(args, opts.Offset)
	}
	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()
	return scanUsers(rows)
}

// UpdateStatus changes a user's lifecycle status.
func (r *UserRepo) UpdateStatus(ctx context.Context, id string, status domain.UserStatus) error {
	res, err := r.db.ExecContext(ctx, `UPDATE users SET status = $1 WHERE id = $2`, string(status), id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errs.ErrNotFound
	}
	return nil
}

// UpdatePasswordHash overwrites the stored password hash.
func (r *UserRepo) UpdatePasswordHash(ctx context.Context, id, hash string) error {
	res, err := r.db.ExecContext(ctx, `UPDATE users SET password_hash = $1 WHERE id = $2`, hash, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errs.ErrNotFound
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
		return nil, errs.ErrNotFound
	}
	return u, err
}

type scanner interface{ Scan(dest ...any) error }

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

func scanUsers(rows *sql.Rows) ([]*domain.User, error) {
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

var _ ports.UserRepo = (*UserRepo)(nil)
