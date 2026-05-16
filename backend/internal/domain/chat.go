package domain

import "time"

// ChatRole is the speaker role on a chat message.
type ChatRole string

const (
	ChatRoleUser ChatRole = "user"
	ChatRoleAI   ChatRole = "ai"
)

// Chat is a single conversation belonging to a user.
type Chat struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ChatMessage is one message inside a chat. AI replies may include
// structured blocks serialised as raw JSON for forward compatibility.
type ChatMessage struct {
	ID         string    `json:"id"`
	ChatID     string    `json:"chatId"`
	Role       ChatRole  `json:"role"`
	Text       string    `json:"text"`
	BlocksJSON string    `json:"-"`            // raw JSON, surfaced via Blocks
	Blocks     []any     `json:"blocks,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
}
