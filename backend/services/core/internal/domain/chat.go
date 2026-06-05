package domain

import (
	"encoding/json"
	"time"
)

// ChatRole identifies who authored a chat message.
type ChatRole string

const (
	ChatRoleUser ChatRole = "user"
	ChatRoleAI   ChatRole = "ai"
)

// Chat is an AI assistant conversation owned by a user.
type Chat struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ChatMessage is one turn in a chat. Blocks holds structured UI payloads
// (plan / budget / vendor cards) as raw JSON so the schema can evolve without
// a migration.
type ChatMessage struct {
	ID        string          `json:"id"`
	ChatID    string          `json:"chatId"`
	Role      ChatRole        `json:"role"`
	Text      string          `json:"text"`
	Blocks    json.RawMessage `json:"blocks,omitempty"`
	CreatedAt time.Time       `json:"createdAt"`
}

// ChatTitleFromMessage derives a short chat title from its first message.
func ChatTitleFromMessage(msg string) string {
	const max = 48
	runes := []rune(msg)
	if len(runes) <= max {
		return msg
	}
	return string(runes[:max]) + "…"
}
