// Package chat implements the AI chat use case with fallback behavior + history persistence.
package chat

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"strings"

	"qonaqzhai-backend/internal/domain"
	"qonaqzhai-backend/internal/usecase"
)

// Deps bundles chat service collaborators.
type Deps struct {
	Vendors usecase.VendorRepo
	Chats   usecase.ChatRepo // optional — when nil, history is not persisted
	AI      usecase.AIClient // optional — nil disables Gemini, fallback only
	Logger  *slog.Logger
}

// Service generates assistant replies using AI when available, falling back to
// deterministic keyword responses when not.
type Service struct{ d Deps }

// New constructs a chat Service.
func New(d Deps) *Service {
	if d.Logger == nil {
		d.Logger = slog.Default()
	}
	return &Service{d: d}
}

// Reply combines the conversational text with structured UI blocks.
type Reply struct {
	Reply  string
	Blocks []usecase.ChatBlock
}

// SendOutcome is the result of sending a message:
// the AI reply plus the chat id (newly created or reused).
type SendOutcome struct {
	ChatID string
	Reply  string
	Blocks []usecase.ChatBlock
}

// Generate produces a Reply for the user message without touching history.
// Used by stateless callers (e.g. tests).
func (s *Service) Generate(ctx context.Context, message string) (*Reply, error) {
	message = strings.TrimSpace(message)
	if message == "" {
		return nil, domain.ErrInvalidInput
	}
	r := s.generate(ctx, message)
	return &Reply{Reply: r.Reply, Blocks: r.Blocks}, nil
}

// Send persists the user message + AI reply against a chat owned by userID.
// When chatID is empty, a fresh chat is created and its id is returned.
func (s *Service) Send(ctx context.Context, userID, chatID, message string) (*SendOutcome, error) {
	message = strings.TrimSpace(message)
	if message == "" {
		return nil, domain.ErrInvalidInput
	}
	if s.d.Chats == nil {
		// History disabled — fall back to stateless behavior.
		r := s.generate(ctx, message)
		return &SendOutcome{ChatID: "", Reply: r.Reply, Blocks: r.Blocks}, nil
	}

	// Resolve / create chat.
	var chat *domain.Chat
	if chatID != "" {
		c, err := s.d.Chats.FindByID(ctx, chatID)
		if err != nil {
			return nil, err
		}
		if c.UserID != userID {
			return nil, domain.ErrForbidden
		}
		chat = c
	} else {
		c, err := s.d.Chats.Create(ctx, userID, deriveTitle(message))
		if err != nil {
			return nil, err
		}
		chat = c
	}

	// Persist user message.
	if _, err := s.d.Chats.AddMessage(ctx, &domain.ChatMessage{
		ChatID: chat.ID, Role: domain.ChatRoleUser, Text: message,
	}); err != nil {
		return nil, err
	}

	// Generate + persist AI reply.
	r := s.generate(ctx, message)
	var blocksJSON string
	if len(r.Blocks) > 0 {
		raw, err := json.Marshal(r.Blocks)
		if err == nil {
			blocksJSON = string(raw)
		}
	}
	if _, err := s.d.Chats.AddMessage(ctx, &domain.ChatMessage{
		ChatID:     chat.ID,
		Role:       domain.ChatRoleAI,
		Text:       r.Reply,
		BlocksJSON: blocksJSON,
	}); err != nil {
		return nil, err
	}

	// Bump updatedAt + title if the chat title was empty.
	if chat.Title == "" {
		_ = s.d.Chats.UpdateTitle(ctx, chat.ID, deriveTitle(message))
	} else {
		_ = s.d.Chats.Touch(ctx, chat.ID)
	}

	return &SendOutcome{ChatID: chat.ID, Reply: r.Reply, Blocks: r.Blocks}, nil
}

// ListChats returns the user's chats, newest first.
func (s *Service) ListChats(ctx context.Context, userID string) ([]*domain.Chat, error) {
	if s.d.Chats == nil {
		return []*domain.Chat{}, nil
	}
	return s.d.Chats.ListForUser(ctx, userID, 100)
}

// GetChat returns a chat + its messages, enforcing ownership.
func (s *Service) GetChat(ctx context.Context, userID, chatID string) (*domain.Chat, []*domain.ChatMessage, error) {
	if s.d.Chats == nil {
		return nil, nil, domain.ErrNotFound
	}
	c, err := s.d.Chats.FindByID(ctx, chatID)
	if err != nil {
		return nil, nil, err
	}
	if c.UserID != userID {
		return nil, nil, domain.ErrForbidden
	}
	msgs, err := s.d.Chats.ListMessages(ctx, chatID)
	if err != nil {
		return nil, nil, err
	}
	return c, msgs, nil
}

// DeleteChat removes a chat owned by userID (cascades to messages).
func (s *Service) DeleteChat(ctx context.Context, userID, chatID string) error {
	if s.d.Chats == nil {
		return domain.ErrNotFound
	}
	c, err := s.d.Chats.FindByID(ctx, chatID)
	if err != nil {
		return err
	}
	if c.UserID != userID {
		return domain.ErrForbidden
	}
	return s.d.Chats.Delete(ctx, chatID)
}

// RenameChat updates a chat's title.
func (s *Service) RenameChat(ctx context.Context, userID, chatID, title string) error {
	if s.d.Chats == nil {
		return domain.ErrNotFound
	}
	c, err := s.d.Chats.FindByID(ctx, chatID)
	if err != nil {
		return err
	}
	if c.UserID != userID {
		return domain.ErrForbidden
	}
	return s.d.Chats.UpdateTitle(ctx, chatID, strings.TrimSpace(title))
}

// generate is the LLM call with fallback. Internal use only.
func (s *Service) generate(ctx context.Context, message string) Reply {
	if s.d.AI != nil {
		refs, err := s.approvedRefs(ctx)
		if err == nil {
			res, err := s.d.AI.Generate(ctx, message, refs)
			if err == nil && res != nil {
				return Reply{Reply: res.Reply, Blocks: res.Blocks}
			}
			s.d.Logger.Warn("ai generate failed, falling back", slog.String("err", errString(err)))
		}
	}
	return s.fallback(ctx, message)
}

func (s *Service) approvedRefs(ctx context.Context) ([]usecase.VendorRef, error) {
	vendors, _, err := s.d.Vendors.Search(ctx, usecase.VendorQuery{Status: domain.VendorApproved, Limit: 100})
	if err != nil {
		return nil, err
	}
	refs := make([]usecase.VendorRef, 0, len(vendors))
	for _, v := range vendors {
		refs = append(refs, usecase.VendorRef{
			ID:        v.ID,
			Name:      v.Name,
			Category:  v.Category,
			City:      v.City,
			PriceFrom: v.PriceFrom,
		})
	}
	return refs, nil
}

func (s *Service) fallback(ctx context.Context, message string) Reply {
	low := strings.ToLower(message)
	switch {
	case containsAny(low, "budget", "бюджет"):
		return Reply{
			Reply: "Here's a typical KZ split for your event.",
			Blocks: []usecase.ChatBlock{{
				Type: "budget",
				Data: map[string]any{
					"total": 5_000_000,
					"categories": []map[string]any{
						{"name": "Venue", "amount": 1_500_000, "pct": 30},
						{"name": "Catering", "amount": 1_800_000, "pct": 36},
						{"name": "Decor", "amount": 600_000, "pct": 12},
						{"name": "Photo & Video", "amount": 450_000, "pct": 9},
						{"name": "Music & DJ", "amount": 250_000, "pct": 5},
						{"name": "Reserve", "amount": 400_000, "pct": 8},
					},
				},
			}},
		}
	case containsAny(low, "vendor", "вендор", "photo", "фотограф"):
		vendors, _, _ := s.d.Vendors.Search(ctx, usecase.VendorQuery{Status: domain.VendorApproved, Limit: 5})
		items := make([]map[string]any, 0, len(vendors))
		for _, v := range vendors {
			items = append(items, map[string]any{
				"id":        v.ID,
				"name":      v.Name,
				"category":  v.Category,
				"city":      v.City,
				"priceFrom": v.PriceFrom,
				"rating":    v.RatingAvg,
			})
		}
		return Reply{
			Reply: "Top picks from our verified vendors.",
			Blocks: []usecase.ChatBlock{{
				Type: "vendors",
				Data: map[string]any{"query": "Top matches", "items": items},
			}},
		}
	default:
		return Reply{
			Reply: "Drafted your plan. Ask for budget, vendors, or timeline next.",
			Blocks: []usecase.ChatBlock{{
				Type: "plan",
				Data: map[string]any{
					"title":     "Aigerim & Daulet — Toi",
					"eventType": "Wedding / Toi",
					"date":      "Aug 12, 2026",
					"city":      "Almaty",
					"guests":    150,
					"budget":    5_000_000,
				},
			}},
		}
	}
}

func containsAny(s string, subs ...string) bool {
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

func errString(err error) string {
	if err == nil {
		return ""
	}
	if errors.Is(err, context.Canceled) {
		return "canceled"
	}
	return err.Error()
}

// deriveTitle creates a short title from the first user message.
func deriveTitle(message string) string {
	t := strings.Join(strings.Fields(message), " ")
	if len(t) > 60 {
		t = t[:60] + "…"
	}
	return t
}
