// Package chat implements the AI chat use case with fallback behavior.
package chat

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"qonaqzhai-backend/internal/domain"
	"qonaqzhai-backend/internal/usecase"
)

// Deps bundles chat service collaborators.
type Deps struct {
	Vendors usecase.VendorRepo
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

// Generate produces a Reply for the user message.
func (s *Service) Generate(ctx context.Context, message string) (*Reply, error) {
	message = strings.TrimSpace(message)
	if message == "" {
		return nil, domain.ErrInvalidInput
	}
	if s.d.AI != nil {
		refs, err := s.approvedRefs(ctx)
		if err == nil {
			res, err := s.d.AI.Generate(ctx, message, refs)
			if err == nil && res != nil {
				return &Reply{Reply: res.Reply, Blocks: res.Blocks}, nil
			}
			s.d.Logger.Warn("ai generate failed, falling back", slog.String("err", errString(err)))
		}
	}
	return s.fallback(ctx, message), nil
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

func (s *Service) fallback(ctx context.Context, message string) *Reply {
	low := strings.ToLower(message)
	switch {
	case containsAny(low, "budget", "бюджет"):
		return &Reply{
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
		return &Reply{
			Reply: "Top picks from our verified vendors.",
			Blocks: []usecase.ChatBlock{{
				Type: "vendors",
				Data: map[string]any{"query": "Top matches", "items": items},
			}},
		}
	default:
		return &Reply{
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
