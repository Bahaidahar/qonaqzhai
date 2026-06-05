// Package gemini adapts Google's Generative Language API into the chat
// planner interface. It performs natural-language understanding only — vendor
// data always comes from the database, so the assistant never hallucinates
// bookable vendors.
package gemini

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"qonaqzhai-backend/services/core/internal/usecase/chat"
)

const endpoint = "https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s"

// Client calls Gemini to parse an event request into a structured intent.
type Client struct {
	apiKey string
	model  string
	http   *http.Client
}

// New constructs a Gemini client. Returns nil when apiKey is empty so callers
// can transparently fall back to the keyword planner.
func New(apiKey, model string) *Client {
	if apiKey == "" {
		return nil
	}
	if model == "" {
		model = "gemini-2.5-flash"
	}
	return &Client{
		apiKey: apiKey,
		model:  model,
		http:   &http.Client{Timeout: 12 * time.Second},
	}
}

const systemPrompt = `Ты — ассистент планирования событий в Казахстане (города Алматы, Астана, Шымкент).
Категории услуг ровно такие: "Venue", "Catering", "Music & DJ", "Photo & Video", "Decor & Florists", "Cakes".
По запросу пользователя извлеки параметры события и верни JSON по схеме.
Правила:
- category: одна строка из списка выше, или "" если категория не ясна.
- guests: число гостей (0 если не указано).
- eventDate: формат YYYY-MM-DD, год по умолчанию 2026, "" если даты нет.
- budget: бюджет в тенге числом (0 если не указан); "5 млн" = 5000000, "500 тыс" = 500000.
- city: "Almaty" | "Astana" | "Shymkent" (по умолчанию Almaty).
- wantsOther: true, если пользователь просит другие/другие варианты/не нравится/ещё/замену.
- reply: короткий дружелюбный ответ на языке пользователя (1-2 предложения), скажи что подобрал подрядчиков из базы и предложи выбрать или попросить другие варианты.
Последняя обсуждаемая категория в этом чате: %q.
Запрос пользователя: %q`

type geminiRequest struct {
	Contents         []content        `json:"contents"`
	GenerationConfig generationConfig `json:"generationConfig"`
}

type content struct {
	Parts []part `json:"parts"`
}

type part struct {
	Text string `json:"text"`
}

type generationConfig struct {
	ResponseMimeType string          `json:"responseMimeType"`
	ResponseSchema   json.RawMessage `json:"responseSchema"`
	Temperature      float64         `json:"temperature"`
}

var responseSchema = json.RawMessage(`{
  "type": "object",
  "properties": {
    "category":   {"type": "string"},
    "guests":     {"type": "integer"},
    "eventDate":  {"type": "string"},
    "budget":     {"type": "integer"},
    "city":       {"type": "string"},
    "wantsOther": {"type": "boolean"},
    "reply":      {"type": "string"}
  },
  "required": ["category","guests","eventDate","budget","city","wantsOther","reply"]
}`)

type geminiResponse struct {
	Candidates []struct {
		Content content `json:"content"`
	} `json:"candidates"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

// Plan calls Gemini and returns the parsed intent. Any transport, status, or
// decode error is returned so the caller can fall back to the keyword planner.
func (c *Client) Plan(ctx context.Context, message, lastCategory string) (*chat.Intent, error) {
	reqBody := geminiRequest{
		Contents: []content{{Parts: []part{{Text: fmt.Sprintf(systemPrompt, lastCategory, message)}}}},
		GenerationConfig: generationConfig{
			ResponseMimeType: "application/json",
			ResponseSchema:   responseSchema,
			Temperature:      0.4,
		},
	}
	raw, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf(endpoint, c.model, c.apiKey)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var gr geminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&gr); err != nil {
		return nil, err
	}
	if gr.Error != nil {
		return nil, fmt.Errorf("gemini: %s", gr.Error.Message)
	}
	if len(gr.Candidates) == 0 || len(gr.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("gemini: empty response")
	}

	var intent chat.Intent
	if err := json.Unmarshal([]byte(gr.Candidates[0].Content.Parts[0].Text), &intent); err != nil {
		return nil, fmt.Errorf("gemini: decode intent: %w", err)
	}
	return &intent, nil
}

var _ chat.Planner = (*Client)(nil)
