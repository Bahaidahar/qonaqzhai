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

// ChatRepo persists AI chats and messages.
type ChatRepo struct {
	db    *sql.DB
	idGen ports.IDGen
}

// NewChatRepo constructs a chat repository.
func NewChatRepo(db *sql.DB, idGen ports.IDGen) *ChatRepo {
	return &ChatRepo{db: db, idGen: idGen}
}

const chatCols = `id, user_id, title, created_at, updated_at`

// CreateChat inserts a new chat owned by c.UserID.
func (r *ChatRepo) CreateChat(ctx context.Context, c *domain.Chat) (*domain.Chat, error) {
	if c.ID == "" {
		c.ID = r.idGen.New()
	}
	if _, err := r.db.ExecContext(ctx,
		`INSERT INTO chats (id, user_id, title) VALUES ($1, $2, $3)`,
		c.ID, c.UserID, c.Title,
	); err != nil {
		return nil, fmt.Errorf("insert chat: %w", err)
	}
	return r.GetChat(ctx, c.ID, c.UserID)
}

// GetChat returns a chat by id scoped to userID.
func (r *ChatRepo) GetChat(ctx context.Context, id, userID string) (*domain.Chat, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT `+chatCols+` FROM chats WHERE id = $1 AND user_id = $2`, id, userID)
	var c domain.Chat
	if err := row.Scan(&c.ID, &c.UserID, &c.Title, &c.CreatedAt, &c.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, fmt.Errorf("get chat: %w", err)
	}
	return &c, nil
}

// ListChats returns a user's chats, most recently updated first.
func (r *ChatRepo) ListChats(ctx context.Context, userID string) ([]*domain.Chat, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT `+chatCols+` FROM chats WHERE user_id = $1 ORDER BY updated_at DESC`, userID)
	if err != nil {
		return nil, fmt.Errorf("list chats: %w", err)
	}
	defer rows.Close()
	out := []*domain.Chat{}
	for rows.Next() {
		var c domain.Chat
		if err := rows.Scan(&c.ID, &c.UserID, &c.Title, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, &c)
	}
	return out, rows.Err()
}

// DeleteChat removes a chat (and its messages via cascade) owned by userID.
func (r *ChatRepo) DeleteChat(ctx context.Context, id, userID string) error {
	res, err := r.db.ExecContext(ctx,
		`DELETE FROM chats WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return fmt.Errorf("delete chat: %w", err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errs.ErrNotFound
	}
	return nil
}

// RenameChat updates a chat title owned by userID.
func (r *ChatRepo) RenameChat(ctx context.Context, id, userID, title string) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE chats SET title = $1, updated_at = now() WHERE id = $2 AND user_id = $3`,
		title, id, userID)
	if err != nil {
		return fmt.Errorf("rename chat: %w", err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errs.ErrNotFound
	}
	return nil
}

// TouchChat bumps a chat's updated_at so it sorts to the top of the list.
func (r *ChatRepo) TouchChat(ctx context.Context, id string) error {
	if _, err := r.db.ExecContext(ctx,
		`UPDATE chats SET updated_at = now() WHERE id = $1`, id); err != nil {
		return fmt.Errorf("touch chat: %w", err)
	}
	return nil
}

// AddMessage appends a message to a chat. Blocks default to an empty array.
func (r *ChatRepo) AddMessage(ctx context.Context, m *domain.ChatMessage) (*domain.ChatMessage, error) {
	if m.ID == "" {
		m.ID = r.idGen.New()
	}
	blocks := m.Blocks
	if len(blocks) == 0 {
		blocks = []byte("[]")
	}
	row := r.db.QueryRowContext(ctx,
		`INSERT INTO chat_messages (id, chat_id, role, text, blocks)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING `+msgCols,
		m.ID, m.ChatID, string(m.Role), m.Text, []byte(blocks),
	)
	return scanMessage(row)
}

// ListMessages returns a chat's messages oldest-first.
func (r *ChatRepo) ListMessages(ctx context.Context, chatID string) ([]*domain.ChatMessage, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT `+msgCols+` FROM chat_messages WHERE chat_id = $1 ORDER BY created_at ASC`, chatID)
	if err != nil {
		return nil, fmt.Errorf("list messages: %w", err)
	}
	defer rows.Close()
	out := []*domain.ChatMessage{}
	for rows.Next() {
		m, err := scanMessage(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

const msgCols = `id, chat_id, role, text, blocks, created_at`

func scanMessage(s scanner) (*domain.ChatMessage, error) {
	var m domain.ChatMessage
	var role string
	var blocks []byte
	if err := s.Scan(&m.ID, &m.ChatID, &role, &m.Text, &blocks, &m.CreatedAt); err != nil {
		return nil, err
	}
	m.Role = domain.ChatRole(role)
	m.Blocks = blocks
	return &m, nil
}

var _ ports.ChatRepo = (*ChatRepo)(nil)
