// Package chat implements AI assistant chat sessions with persistence.
// The reply is a deterministic, keyword-driven stub: it recognises a service
// category in the user's message and answers with real approved vendors from
// the catalog. Swap buildReply for a real LLM call without touching
// persistence or the HTTP layer.
package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"

	"qonaqzhai-backend/pkg/errs"
	"qonaqzhai-backend/services/core/internal/domain"
	"qonaqzhai-backend/services/core/internal/ports"
)

// VendorSearch is the slice of vendor querying the chat planner needs.
// vendor.Service satisfies it.
type VendorSearch interface {
	Search(ctx context.Context, q ports.VendorQuery) ([]*domain.Vendor, int, error)
}

// Intent is the structured understanding of a user message — produced by an
// LLM Planner, or by the deterministic keyword fallback.
type Intent struct {
	Category   string `json:"category"`
	Guests     int    `json:"guests"`
	EventDate  string `json:"eventDate"`
	Budget     int64  `json:"budget"`
	City       string `json:"city"`
	WantsOther bool   `json:"wantsOther"`
	Reply      string `json:"reply"`
}

// Planner turns a free-text message into a structured Intent. The Gemini
// adapter implements it; a nil Planner falls back to keyword parsing.
type Planner interface {
	Plan(ctx context.Context, message, lastCategory string) (*Intent, error)
}

// Deps bundles chat collaborators.
type Deps struct {
	Chats   ports.ChatRepo
	Vendors VendorSearch
	Planner Planner
	Logger  *slog.Logger
}

// Service exposes chat operations.
type Service struct{ d Deps }

// New constructs a chat Service. A nil logger falls back to slog.Default().
func New(d Deps) *Service {
	if d.Logger == nil {
		d.Logger = slog.Default()
	}
	return &Service{d: d}
}

// SendResult is returned to the HTTP layer after a turn is persisted.
type SendResult struct {
	ChatID string
	Reply  string
	Blocks json.RawMessage
}

// Send persists the user message, generates + persists an AI reply, and
// returns the reply. An empty chatID starts a new chat owned by userID.
func (s *Service) Send(ctx context.Context, userID, chatID, message string) (*SendResult, error) {
	message = strings.TrimSpace(message)
	if message == "" {
		return nil, fmt.Errorf("message: %w", errs.ErrInvalidInput)
	}

	chat, err := s.resolveChat(ctx, userID, chatID, message)
	if err != nil {
		return nil, err
	}

	// Conversational memory: read what we've already shown in this chat so a
	// follow-up ("show me others") returns fresh vendors and the last category
	// can be reused when the new message doesn't name one.
	shown, lastCategory := s.chatContext(ctx, chat.ID)

	if _, err := s.d.Chats.AddMessage(ctx, &domain.ChatMessage{
		ChatID: chat.ID, Role: domain.ChatRoleUser, Text: message,
	}); err != nil {
		return nil, err
	}

	replyText, blocks := s.buildReply(ctx, message, lastCategory, shown)
	if _, err := s.d.Chats.AddMessage(ctx, &domain.ChatMessage{
		ChatID: chat.ID, Role: domain.ChatRoleAI, Text: replyText, Blocks: blocks,
	}); err != nil {
		return nil, err
	}
	if err := s.d.Chats.TouchChat(ctx, chat.ID); err != nil {
		s.d.Logger.Warn("touch chat failed", "chat", chat.ID, "err", err)
	}

	return &SendResult{ChatID: chat.ID, Reply: replyText, Blocks: blocks}, nil
}

// resolveChat returns the owned chat for chatID, or creates a new one when
// chatID is empty.
func (s *Service) resolveChat(ctx context.Context, userID, chatID, firstMessage string) (*domain.Chat, error) {
	if chatID != "" {
		return s.d.Chats.GetChat(ctx, chatID, userID)
	}
	return s.d.Chats.CreateChat(ctx, &domain.Chat{
		UserID: userID,
		Title:  domain.ChatTitleFromMessage(firstMessage),
	})
}

// List returns the user's chats, most recent first.
func (s *Service) List(ctx context.Context, userID string) ([]*domain.Chat, error) {
	return s.d.Chats.ListChats(ctx, userID)
}

// Detail returns a chat plus its full message history, scoped to userID.
func (s *Service) Detail(ctx context.Context, userID, chatID string) (*domain.Chat, []*domain.ChatMessage, error) {
	chat, err := s.d.Chats.GetChat(ctx, chatID, userID)
	if err != nil {
		return nil, nil, err
	}
	msgs, err := s.d.Chats.ListMessages(ctx, chat.ID)
	if err != nil {
		return nil, nil, err
	}
	return chat, msgs, nil
}

// Delete removes a chat owned by userID.
func (s *Service) Delete(ctx context.Context, userID, chatID string) error {
	return s.d.Chats.DeleteChat(ctx, chatID, userID)
}

// Rename updates a chat title owned by userID.
func (s *Service) Rename(ctx context.Context, userID, chatID, title string) error {
	title = strings.TrimSpace(title)
	if title == "" {
		return fmt.Errorf("title: %w", errs.ErrInvalidInput)
	}
	return s.d.Chats.RenameChat(ctx, chatID, userID, title)
}

// buildReply is the deterministic, conversational stub planner.
//   - first ask in a category → a draft plan PLUS matching vendors
//   - "show me others" follow-up → a fresh batch of vendors, excluding the
//     ones already shown in this chat
//   - no category → a draft plan and a nudge to name one
//
// Replace with a real LLM — the persisted shape (text + JSON blocks) is stable.
func (s *Service) buildReply(ctx context.Context, message, lastCategory string, shown []string) (string, json.RawMessage) {
	in := s.resolveIntent(ctx, message, lastCategory)
	category := in.Category
	if category == "" && in.WantsOther {
		category = lastCategory // "show others" without naming a category again
	}

	if category == "" {
		text := in.Reply
		if text == "" {
			text = "Опишите событие и категорию (зал, кейтеринг, декор, фото, музыка, торты) — подберу подрядчиков. Пока вот черновик плана:"
		}
		return text, marshalBlocks(planBlockData(in))
	}

	vs := s.searchVendors(ctx, category, shown, 3)
	if len(vs) == 0 {
		if in.WantsOther {
			return "Это все подходящие варианты в этой категории. Назовите другую категорию или город — поищу ещё.",
				marshalBlocks()
		}
		return "Пока не нашёл подрядчиков в этой категории. Вот черновик плана, чтобы начать:",
			marshalBlocks(planBlockData(in))
	}

	vendors := vendorsBlockData(vs, message, in.Guests, in.EventDate)
	if in.WantsOther {
		text := in.Reply
		if text == "" {
			text = "Вот другие варианты из базы:"
		}
		return text, marshalBlocks(vendors)
	}
	text := in.Reply
	if text == "" {
		text = "Вот черновик плана. Вам могут подойти эти подрядчики — выберите и забронируйте, или попросите другие варианты:"
	}
	return text, marshalBlocks(planBlockData(in), vendors)
}

// resolveIntent asks the LLM Planner to understand the message, falling back to
// deterministic keyword parsing when no planner is configured or it errors.
func (s *Service) resolveIntent(ctx context.Context, message, lastCategory string) Intent {
	if s.d.Planner != nil {
		if it, err := s.d.Planner.Plan(ctx, message, lastCategory); err == nil && it != nil {
			it.Category = normalizeCategory(it.Category)
			if it.City == "" {
				it.City = "Almaty"
			}
			return *it
		} else if err != nil {
			s.d.Logger.Warn("planner failed, using keyword fallback", "err", err)
		}
	}
	guests, eventDate := parseEventContext(message)
	return Intent{
		Category:   detectCategory(message),
		Guests:     guests,
		EventDate:  eventDate,
		Budget:     parseBudget(message),
		City:       parseCity(message),
		WantsOther: detectMoreIntent(message),
	}
}

// normalizeCategory maps an LLM-returned category to a canonical catalog value,
// tolerating case / whitespace drift. Unknown values become "".
func normalizeCategory(c string) string {
	c = strings.TrimSpace(c)
	for _, valid := range []string{"Venue", "Catering", "Music & DJ", "Photo & Video", "Decor & Florists", "Cakes"} {
		if strings.EqualFold(c, valid) {
			return valid
		}
	}
	return ""
}

// chatContext scans a chat's prior AI vendor blocks to recover which vendors
// were already shown and the last category discussed.
func (s *Service) chatContext(ctx context.Context, chatID string) ([]string, string) {
	msgs, err := s.d.Chats.ListMessages(ctx, chatID)
	if err != nil {
		return nil, ""
	}
	var shown []string
	lastCategory := ""
	for _, m := range msgs {
		if m.Role != domain.ChatRoleAI || len(m.Blocks) == 0 {
			continue
		}
		var blocks []struct {
			Type string `json:"type"`
			Data struct {
				Items []struct {
					ID       string `json:"id"`
					Category string `json:"category"`
				} `json:"items"`
			} `json:"data"`
		}
		if err := json.Unmarshal(m.Blocks, &blocks); err != nil {
			continue
		}
		for _, b := range blocks {
			if b.Type != "vendors" {
				continue
			}
			for _, it := range b.Data.Items {
				shown = append(shown, it.ID)
				if it.Category != "" {
					lastCategory = it.Category
				}
			}
		}
	}
	return shown, lastCategory
}

// searchVendors returns up to n approved vendors of a category, skipping any id
// in exclude. It over-fetches so exclusions still leave a full batch.
func (s *Service) searchVendors(ctx context.Context, category string, exclude []string, n int) []*domain.Vendor {
	if s.d.Vendors == nil {
		return nil
	}
	vs, _, err := s.d.Vendors.Search(ctx, ports.VendorQuery{
		Category: category, Sort: "rating_desc", Limit: n + len(exclude) + 4,
	})
	if err != nil {
		s.d.Logger.Warn("chat vendor search failed", "category", category, "err", err)
		return nil
	}
	skip := make(map[string]struct{}, len(exclude))
	for _, id := range exclude {
		skip[id] = struct{}{}
	}
	out := make([]*domain.Vendor, 0, n)
	for _, v := range vs {
		if _, seen := skip[v.ID]; seen {
			continue
		}
		out = append(out, v)
		if len(out) == n {
			break
		}
	}
	return out
}

// vendorsBlockData renders a vendors block. Event context (guests/date) rides
// along so the chat cards can book directly.
func vendorsBlockData(vs []*domain.Vendor, query string, guests int, eventDate string) map[string]any {
	items := make([]map[string]any, 0, len(vs))
	for _, v := range vs {
		items = append(items, map[string]any{
			"id":        v.ID,
			"name":      v.Name,
			"category":  v.Category,
			"rating":    v.RatingAvg,
			"priceFrom": v.PriceFrom,
			"city":      v.City,
		})
	}
	return map[string]any{"type": "vendors", "data": map[string]any{
		"query":     query,
		"items":     items,
		"guests":    guests,
		"eventDate": eventDate,
	}}
}

// marshalBlocks marshals zero or more block maps into the JSON array the client
// expects.
func marshalBlocks(blocks ...map[string]any) json.RawMessage {
	if blocks == nil {
		blocks = []map[string]any{}
	}
	raw, _ := json.Marshal(blocks)
	return raw
}

// detectMoreIntent recognises a "give me other options" follow-up.
func detectMoreIntent(message string) bool {
	m := strings.ToLower(message)
	for _, kw := range []string{
		"не нрав", "не подход", "друг", "ещё", "еще", "альтернатив", "замен",
		"other", "another", "more", "else", "instead",
	} {
		if strings.Contains(m, kw) {
			return true
		}
	}
	return false
}

var (
	reGuests = regexp.MustCompile(`(\d{1,5})\s*(?:гост|чел|адам|қонақ|guest|person|pax)`)
	reAnyNum = regexp.MustCompile(`\d{1,5}`)
	reISO    = regexp.MustCompile(`\b(\d{4})-(\d{2})-(\d{2})\b`)
	reDMY    = regexp.MustCompile(`\b(\d{1,2})\.(\d{1,2})(?:\.(\d{2,4}))?\b`)
	reDMonth = regexp.MustCompile(`\b(\d{1,2})\s+(январ|феврал|март|апрел|ма[йя]|июн|июл|август|сентябр|октябр|ноябр|декабр|jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec)`)
)

var monthNum = map[string]int{
	"январ": 1, "jan": 1, "феврал": 2, "feb": 2, "март": 3, "mar": 3,
	"апрел": 4, "apr": 4, "ма": 5, "may": 5, "июн": 6, "jun": 6,
	"июл": 7, "jul": 7, "август": 8, "aug": 8, "сентябр": 9, "sep": 9,
	"октябр": 10, "oct": 10, "ноябр": 11, "nov": 11, "декабр": 12, "dec": 12,
}

// parseEventContext best-effort extracts guest count and an ISO event date
// from free text. Either may be zero/empty when not present.
func parseEventContext(text string) (int, string) {
	m := strings.ToLower(text)

	guests := 0
	if g := reGuests.FindStringSubmatch(m); g != nil {
		guests = atoiSafe(g[1])
	} else if n := reAnyNum.FindString(m); n != "" {
		guests = atoiSafe(n)
	}

	date := ""
	const defaultYear = 2026
	switch {
	case reISO.MatchString(m):
		p := reISO.FindStringSubmatch(m)
		date = fmt.Sprintf("%s-%s-%s", p[1], p[2], p[3])
	case reDMY.MatchString(m):
		p := reDMY.FindStringSubmatch(m)
		year := defaultYear
		if p[3] != "" {
			year = atoiSafe(p[3])
			if year < 100 {
				year += 2000
			}
		}
		date = fmt.Sprintf("%04d-%02d-%02d", year, atoiSafe(p[2]), atoiSafe(p[1]))
	case reDMonth.MatchString(m):
		p := reDMonth.FindStringSubmatch(m)
		if mn, ok := monthNum[p[2]]; ok {
			date = fmt.Sprintf("%04d-%02d-%02d", defaultYear, mn, atoiSafe(p[1]))
		}
	}
	return guests, date
}

func atoiSafe(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

// planBlockData builds a draft-plan block from a resolved Intent, with sensible
// defaults for anything the user didn't specify.
func planBlockData(in Intent) map[string]any {
	guests := in.Guests
	if guests == 0 {
		guests = 100
	}
	budget := in.Budget
	if budget == 0 {
		budget = 3000000
	}
	city := in.City
	if city == "" {
		city = "Almaty"
	}
	return map[string]any{
		"type": "plan",
		"data": map[string]any{
			"title":     "Черновик плана события",
			"eventType": "event",
			"city":      city,
			"date":      in.EventDate,
			"guests":    guests,
			"budget":    budget,
		},
	}
}

var (
	reMln       = regexp.MustCompile(`(\d[\d.,\s]*)\s*(?:млн|миллион|million)`)
	reTys       = regexp.MustCompile(`(\d[\d.,\s]*)\s*(?:тыс|тысяч|k\b)`)
	reBudgetRaw = regexp.MustCompile(`бюджет\D*(\d[\d.,\s]*\d|\d)`)
)

func cleanInt(s string) int64 {
	digits := strings.NewReplacer(" ", "", ".", "", ",", "").Replace(s)
	n, _ := strconv.ParseInt(digits, 10, 64)
	return n
}

// parseBudget pulls an approximate budget. A million/thousand unit is taken
// anywhere ("5 млн"); a bare number counts only right after "бюджет" so guest
// counts aren't mistaken for money. Returns 0 when nothing matches.
func parseBudget(message string) int64 {
	m := strings.ToLower(message)
	if g := reMln.FindStringSubmatch(m); g != nil {
		return cleanInt(g[1]) * 1000000
	}
	if g := reTys.FindStringSubmatch(m); g != nil {
		return cleanInt(g[1]) * 1000
	}
	if g := reBudgetRaw.FindStringSubmatch(m); g != nil {
		return cleanInt(g[1])
	}
	return 0
}

// parseCity recognises the three MVP cities; defaults to Almaty.
func parseCity(message string) string {
	m := strings.ToLower(message)
	switch {
	case strings.Contains(m, "астан") || strings.Contains(m, "astana") || strings.Contains(m, "нур-султан"):
		return "Astana"
	case strings.Contains(m, "шымкент") || strings.Contains(m, "shymkent"):
		return "Shymkent"
	default:
		return "Almaty"
	}
}

// detectCategory maps free-text keywords (ru / en / kk) to a catalog category.
func detectCategory(message string) string {
	m := strings.ToLower(message)
	type rule struct {
		category string
		keywords []string
	}
	rules := []rule{
		{"Venue", []string{"зал", "venue", "мекен", "банкет", "площадк", "ресторан", "той хана", "тойхана"}},
		{"Catering", []string{"кейтеринг", "catering", "еда", "меню", "дастархан", "dastarkhan", "ас"}},
		{"Music & DJ", []string{"музык", "music", "dj", "диджей", "тамада", "mc", "ведущ", "артист"}},
		{"Photo & Video", []string{"фото", "видео", "photo", "video", "съёмк", "съемк", "оператор"}},
		{"Decor & Florists", []string{"декор", "decor", "цвет", "флорист", "оформлен", "шар", "гүл"}},
		{"Cakes", []string{"торт", "cake", "десерт", "сладк", "кондитер"}},
	}
	for _, r := range rules {
		for _, kw := range r.keywords {
			if strings.Contains(m, kw) {
				return r.category
			}
		}
	}
	return ""
}
