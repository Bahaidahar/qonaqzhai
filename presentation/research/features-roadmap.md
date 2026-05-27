# Qonaqzhai — Feature Roadmap for Defense

Prioritization framework: **Impact** (user value + demo wow-factor) × **Effort** (engineering days) × **Risk** (integration complexity).

Legend: 🟢 already in repo (need polish) · 🟡 partial · 🔴 net-new

---

## Tier 1 — Add before defense (high impact, low-medium effort)

### 1. AI Event Planner Agent with Tool Use 🟡
**Current state**: `mobile/lib/features/ai_chat/` + backend chat repository exist as plain Q&A.
**Upgrade**: switch to **Claude Sonnet 4.6 with tool use** — the AI can call backend endpoints (`search_vendors`, `check_availability`, `draft_booking`) and return structured results. End-to-end demo: *"Plan a 100-guest wedding in Almaty for August 15 under 2,000,000 ₸"* → AI returns curated vendor shortlist, draft schedule, total estimate.
**Effort**: 3-4 days. **Impact**: ⭐⭐⭐⭐⭐ (signature feature for diploma).
**Tech**: Anthropic Messages API tool-use, JSON schema per tool, streaming SSE from `core` service.

### 2. Vendor Analytics Dashboard 🔴
**Current state**: vendors only see bookings list (`mobile/lib/features/vendor_self/`).
**Upgrade**: dashboard with revenue curves, booking funnel (views → contact → booked), top services, rating trend.
**Effort**: 2-3 days frontend, 1 day backend SQL aggregations.
**Tech**: `core` service materialized views, Recharts on web, `fl_chart` on Flutter.

### 3. Multi-language UI completion (EN/RU/KZ) 🟡
**Current state**: `mobile/lib/core/i18n/i18n.dart` exists with EN/RU/KZ keys; frontend i18n partial.
**Upgrade**: complete dictionary coverage, language switcher in settings, locale-aware date/currency, Kazakh script polish.
**Effort**: 1-2 days. **Impact**: huge for KZ market positioning.

### 4. Push notifications with deep-links 🟡
**Current state**: Firebase Messaging wired but flow incomplete.
**Upgrade**: send push on `booking_accepted`, `payment_succeeded`, `new_message`. Tap → opens exact screen via `go_router` deep-link.
**Effort**: 1-2 days. **Impact**: live demo factor.

### 5. Real-time chat polish 🟡
**Current state**: WS hub in `realtime` service works; UI exists.
**Upgrade**: typing indicators, read receipts, image attachments, online presence.
**Effort**: 2 days.

---

## Tier 2 — Nice-to-have (high impact, medium effort)

### 6. Smart vendor recommendations 🔴
Personalized vendor ranking based on user's prior bookings, viewed vendors, budget signals. Initial impl: weighted scoring SQL (no ML required). Future: embedding similarity.
**Effort**: 2 days. **Tech**: `core` service `/api/vendors/recommended` endpoint, score = w1·rating + w2·category_match + w3·budget_fit.

### 7. Dynamic pricing 🔴
Peak/off-peak pricing per vendor (Saturday weddings cost more than Tuesday corporate). Vendor-defined rules in admin panel.
**Effort**: 2-3 days.

### 8. Calendar integration 🔴
Vendor availability synced from Google Calendar (read-only). Prevents double-booking.
**Effort**: 2 days. **Tech**: Google Calendar API, OAuth2.

### 9. Referral / promo codes 🔴
Invite friend → both get 5% off first booking. Promo entry in checkout.
**Effort**: 1-2 days. **Tech**: new table `promo_codes` in core-db, validation gRPC call from payment.

### 10. Review with photos + verified-booking badge 🟡
**Current state**: reviews exist.
**Upgrade**: only let users review if `booking.status = completed`; allow photo upload via existing `image_picker`.
**Effort**: 1 day.

---

## Tier 3 — Stretch / future work (mention in "next steps" slide)

### 11. ML-based vendor matching with embeddings
Use sentence-transformers to embed vendor descriptions + customer query, rank by cosine similarity.

### 12. Voice input for AI planner
`speech_to_text` package → AI chat. Useful for hands-free planning on mobile.

### 13. Escrow payments
Hold customer funds until event date. PayBox supports holds; payment-service saga extension.

### 14. Vendor on-boarding KYC
Auto-verify business registration via egov.kz public API. Adds trust badge.

### 15. Group bookings / split payments
Multiple guests pay separately for a single event (group dinners, bachelor parties).

### 16. Aigystyn (Kazakh-traditions add-on)
Curated package builder for traditional ceremonies: тұсаукесер, беташар, шашу, наурыз — predefined service bundles.

### 17. Vendor live-streaming previews
Vendor can broadcast a tour of their venue from the app; customers watch live.

### 18. AI-generated event timeline / runbook
After booking confirmation, AI generates minute-by-minute script for the day. Exportable to PDF.

---

## Suggested Defense Demo Flow (8 minutes)

1. Open mobile app → switch language KZ → EN (i18n) — **30 s**
2. Open AI chat → "Plan corporate event for 50 people in Astana, budget 1.5M ₸" → AI returns 3 vendor cards with reasoning — **2 min**
3. Tap vendor → see availability calendar (Google Calendar sync) → book Aug 20 — **1 min**
4. Pay with PayBox (test mode) → receive push notification → tap → deep-link into booking detail — **1 min**
5. Switch to vendor account → see new booking in dashboard, revenue chart updates — **1 min**
6. Open admin web → moderate flagged review — **30 s**
7. Show architecture diagram, gRPC trace in Grafana / logs — **1 min**
8. Q&A buffer — **1 min**

---

## Effort Budget (1-2 weeks)

| Week | Items |
|------|-------|
| Week 1 (research + presentation) | Comparative analysis, stack justification, architecture diagrams, presentation MD, features doc — **done in parallel** |
| Week 2 day 1-2 | Tier 1 #3 i18n completion + #4 push deep-links |
| Week 2 day 3-5 | Tier 1 #1 AI tool-use planner |
| Week 2 day 6-7 | Tier 1 #2 vendor analytics dashboard, polish |

Goal: ship Tier 1 #1, #2, #3, #4 — gives strong demo without scope risk. Tier 2 features pitched as "next iteration" if cut.
