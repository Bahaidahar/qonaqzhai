// Package ai implements the AIClient port against Google Gemini.
package ai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"google.golang.org/genai"

	"qonaqzhai-backend/internal/usecase"
)

// Gemini wraps the Google Gemini SDK to satisfy usecase.AIClient.
type Gemini struct {
	gc    *genai.Client
	model string
}

// New constructs a Gemini client. Returns (nil, nil) when apiKey is empty so callers
// can treat it as "AI disabled, use fallback".
func New(ctx context.Context, apiKey, model string) (*Gemini, error) {
	if apiKey == "" {
		return nil, nil
	}
	if model == "" {
		model = "gemini-2.0-flash"
	}
	gc, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("genai client: %w", err)
	}
	return &Gemini{gc: gc, model: model}, nil
}

const systemPrompt = `You are an AI event planning assistant for Kazakhstan.

The user is planning toi, weddings, corporate events, birthdays, conferences. You help them with planning, budgets, vendor recommendations, timelines, and questions.

ALWAYS reply in the user's language — Kazakh (kz), Russian (ru), or English (en). Detect from their message.

You will receive a JSON list of APPROVED vendors. ONLY use vendor IDs from that list when recommending vendors.

OUTPUT FORMAT — return STRICT JSON only, no markdown fences, no commentary:

{
  "reply": "<conversational answer 1-3 sentences in user's language>",
  "blocks": [
    {"type": "plan", "data": {"title": "...", "eventType": "...", "date": "...", "city": "Almaty", "guests": 150, "budget": 5000000}},
    {"type": "budget", "data": {"total": 5000000, "categories": [{"name": "Venue", "amount": 1500000, "pct": 30}]}},
    {"type": "vendors", "data": {"query": "...", "items": [{"id": "<exact id from list>", "name": "...", "category": "...", "city": "...", "priceFrom": 0, "rating": 4.8}]}}
  ]
}

Rules:
- Categories in budget must sum to ~100% pct.
- Budget category names: translate to user's language.
- Only use vendor IDs from APPROVED_VENDORS. Never invent.
- KZT currency. CITY: only Almaty is supported right now.
- If critical info (date, event type) is missing, ask 1-2 short clarifying questions and return EMPTY blocks array.`

// Generate produces a structured chat reply from Gemini.
func (g *Gemini) Generate(ctx context.Context, userMessage string, vendors []usecase.VendorRef) (*usecase.ChatReply, error) {
	if g == nil || g.gc == nil {
		return nil, errors.New("ai client not configured")
	}
	vendorsJSON, _ := json.Marshal(vendors)
	prompt := fmt.Sprintf(
		"%s\n\nAPPROVED_VENDORS:\n%s\n\nUSER_MESSAGE:\n%s\n\nReturn JSON only.",
		systemPrompt, string(vendorsJSON), userMessage,
	)
	cfg := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		Temperature:      genai.Ptr[float32](0.7),
	}
	res, err := g.gc.Models.GenerateContent(ctx, g.model, genai.Text(prompt), cfg)
	if err != nil {
		return nil, fmt.Errorf("generate: %w", err)
	}
	text := strings.TrimSpace(res.Text())
	if text == "" {
		return nil, errors.New("empty model response")
	}
	text = strings.TrimPrefix(text, "```json")
	text = strings.TrimPrefix(text, "```")
	text = strings.TrimSuffix(text, "```")
	text = strings.TrimSpace(text)

	var out struct {
		Reply  string `json:"reply"`
		Blocks []struct {
			Type string         `json:"type"`
			Data map[string]any `json:"data"`
		} `json:"blocks"`
	}
	if err := json.Unmarshal([]byte(text), &out); err != nil {
		return nil, fmt.Errorf("parse model json: %w · raw=%s", err, text)
	}
	blocks := make([]usecase.ChatBlock, 0, len(out.Blocks))
	for _, b := range out.Blocks {
		blocks = append(blocks, usecase.ChatBlock{Type: b.Type, Data: b.Data})
	}
	return &usecase.ChatReply{Reply: out.Reply, Blocks: blocks}, nil
}

var _ usecase.AIClient = (*Gemini)(nil)
