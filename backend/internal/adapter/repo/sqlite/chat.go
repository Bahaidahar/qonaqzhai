package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"qonaqzhai-backend/internal/domain"
	"qonaqzhai-backend/internal/usecase"
)

// ChatRepo persists chats and chat messages.
type ChatRepo struct {
	db    *sql.DB
	idGen usecase.IDGen
}

// NewChatRepo constructs a chat repository.
func NewChatRepo(db *sql.DB, idGen usecase.IDGen) *ChatRepo {
	return &ChatRepo{db: db, idGen: idGen}
}

// Create inserts a new chat row.
func (r *ChatRepo) Create(ctx context.Context, userID, title string) (*domain.Chat, error) {
	id := r.idGen.New()
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO chats (id, user_id, title) VALUES ($1, $2, $3)`,
		id, userID, title,
	)
	if err != nil {
		return nil, fmt.Errorf("insert chat: %w", err)
	}
	return r.FindByID(ctx, id)
}

// FindByID returns a chat by primary key.
func (r *ChatRepo) FindByID(ctx context.Context, id string) (*domain.Chat, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, title, created_at, updated_at FROM chats WHERE id = $1`, id,
	)
	var c domain.Chat
	if err := row.Scan(&c.ID, &c.UserID, &c.Title, &c.CreatedAt, &c.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &c, nil
}

// ListForUser returns chats owned by user, newest first.
func (r *ChatRepo) ListForUser(ctx context.Context, userID string, limit int) ([]*domain.Chat, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, user_id, title, created_at, updated_at FROM chats
		 WHERE user_id = $1 ORDER BY updated_at DESC LIMIT $2`, userID, limit,
	)
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

// UpdateTitle changes the chat title.
func (r *ChatRepo) UpdateTitle(ctx context.Context, id, title string) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE chats SET title = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`,
		title, id,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// Touch bumps updated_at so the chat moves to the top of the sidebar.
func (r *ChatRepo) Touch(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE chats SET updated_at = CURRENT_TIMESTAMP WHERE id = $1`, id,
	)
	return err
}

// Delete removes a chat (cascades to messages).
func (r *ChatRepo) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM chats WHERE id = $1`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// AddMessage inserts a chat message.
func (r *ChatRepo) AddMessage(ctx context.Context, m *domain.ChatMessage) (*domain.ChatMessage, error) {
	if m.ID == "" {
		m.ID = r.idGen.New()
	}
	// Serialise blocks → blocks_json. Keep empty string when there are no blocks.
	blocksJSON := m.BlocksJSON
	if blocksJSON == "" && len(m.Blocks) > 0 {
		raw, err := json.Marshal(m.Blocks)
		if err != nil {
			return nil, fmt.Errorf("marshal blocks: %w", err)
		}
		blocksJSON = string(raw)
	}
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO chat_messages (id, chat_id, role, text, blocks_json) VALUES ($1, $2, $3, $4, $5)`,
		m.ID, m.ChatID, string(m.Role), m.Text, blocksJSON,
	)
	if err != nil {
		return nil, fmt.Errorf("insert chat message: %w", err)
	}
	m.BlocksJSON = blocksJSON
	return m, nil
}

// ListMessages returns messages for a chat ordered by creation time.
func (r *ChatRepo) ListMessages(ctx context.Context, chatID string) ([]*domain.ChatMessage, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, chat_id, role, text, blocks_json, created_at FROM chat_messages
		 WHERE chat_id = $1 ORDER BY created_at ASC`, chatID,
	)
	if err != nil {
		return nil, fmt.Errorf("list chat messages: %w", err)
	}
	defer rows.Close()
	out := []*domain.ChatMessage{}
	for rows.Next() {
		var m domain.ChatMessage
		var role string
		if err := rows.Scan(&m.ID, &m.ChatID, &role, &m.Text, &m.BlocksJSON, &m.CreatedAt); err != nil {
			return nil, err
		}
		m.Role = domain.ChatRole(role)
		if m.BlocksJSON != "" {
			var blocks []any
			if err := json.Unmarshal([]byte(m.BlocksJSON), &blocks); err == nil {
				m.Blocks = blocks
			}
		}
		out = append(out, &m)
	}
	return out, rows.Err()
}

var _ usecase.ChatRepo = (*ChatRepo)(nil)
