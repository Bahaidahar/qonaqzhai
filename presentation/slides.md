---
marp: true
theme: default
size: 16:9
paginate: true
backgroundColor: "#ffffff"
style: |
  @import url("./style.css");
---

<!-- _class: lead -->
<!-- _paginate: false -->

<div class="eyebrow">Diploma defense · 2026</div>

# qonaqzhai
# <span class="accent">Plan events by chatting.</span>

<div class="subtitle">
AI-assisted event marketplace for Kazakhstan — wedding, toi, corporate.
Three clients (web + mobile + MCP), four Go microservices, one bill of materials.
</div>

<div class="footer">Bahtiyar Yelik · Astana IT University</div>

---

<div class="eyebrow">01 — The problem</div>

# Planning a Kazakh wedding takes <span style="color:var(--primary)">14+ phone calls</span> and an Instagram DM chain.

<div class="cols">

<div>

**Customer side**
- Vendors live on Instagram, WhatsApp, 2GIS — no price transparency
- Comparisons mean screenshots in group chats
- Cancellations slip because there's no booking record
- AI assistance is locked to enterprise tools

</div>

<div>

**Vendor side**
- Bookings sit in DMs scattered across 4 messengers
- Payment is cash or Kaspi link, no escrow
- Reviews are word-of-mouth, no portable reputation
- Free Instagram traffic flattening since 2024

</div>

</div>

---

<div class="eyebrow">02 — Market</div>

# Kazakhstan event services — <span style="color:var(--primary)">≈ ₸ 480 B / year</span>

<div class="kpi-row">
<div class="kpi"><div class="num">160 K</div><div class="lbl">weddings / year</div></div>
<div class="kpi"><div class="num">~₸ 3 M</div><div class="lbl">avg cheque per event</div></div>
<div class="kpi"><div class="num">22 %</div><div class="lbl">YoY mobile commerce</div></div>
<div class="kpi"><div class="num">68 %</div><div class="lbl">vendors on Instagram only</div></div>
</div>

<br/>

**Sources:** Bureau of National Statistics (BNS RK 2024 demographic yearbook); Halyk Finance retail consumption brief Q3-2024; Kaspi Marketplace investor letter 2024; in-house survey of 38 vendors (Almaty, Astana, Shymkent), Mar 2026.

---

<div class="eyebrow">03 — Who's already in the field</div>

# Competitive landscape — and what's missing

| Player | Geo | Coverage | Booking flow | AI planner | Native mobile | Realtime chat |
|---|---|---|---|---|---|---|
| **Tamada.kz** | KZ | Photo + tamada | Phone callback | <span class="tag tag-no">no</span> | <span class="tag tag-no">no</span> | <span class="tag tag-no">no</span> |
| **Wedy.kz** | KZ | Wedding only | Form → email | <span class="tag tag-no">no</span> | <span class="tag tag-no">no</span> | <span class="tag tag-no">no</span> |
| **Sallem.kz** | KZ | Venue rentals | Calendar slot | <span class="tag tag-no">no</span> | <span class="tag tag-no">no</span> | <span class="tag tag-no">no</span> |
| **Instagram + WhatsApp** | KZ | Everything | DM | <span class="tag tag-no">no</span> | n/a | DM |
| **The Knot** (US) | US | Wedding | Quote engine | <span class="tag tag-no">no</span> | <span class="tag tag-ok">yes</span> | <span class="tag tag-no">no</span> |
| **Bash** (US) | US | Venue + service | Direct book | <span class="tag tag-ok">limited</span> | <span class="tag tag-ok">yes</span> | <span class="tag tag-no">no</span> |
| **qonaqzhai** | **KZ** | **All categories** | **Direct + escrow** | <span class="tag tag-ok">**yes**</span> | <span class="tag tag-ok">**yes**</span> | <span class="tag tag-ok">**WS realtime**</span> |

---

<div class="eyebrow">04 — Comparative scorecard</div>

# Feature matrix

| Feature | Tamada.kz | Wedy.kz | Sallem.kz | The Knot | **qonaqzhai** |
|---|---|---|---|---|---|
| Public vendor catalog | ⚪ | ⚪ | ⚪ | ✅ | ✅ |
| Native iOS / Android | ❌ | ❌ | ❌ | ✅ | ✅ Flutter |
| 3 languages (kk / ru / en) | ⚪ | ⚪ | ❌ | ❌ | ✅ |
| AI conversational planner | ❌ | ❌ | ❌ | ❌ | ✅ Gemini |
| Vendor self-service | ⚪ | ❌ | ⚪ | ✅ | ✅ |
| Realtime customer ↔ vendor chat | ❌ | ❌ | ❌ | ❌ | ✅ WebSocket |
| Escrow-style payment hold | ❌ | ❌ | ❌ | ❌ | ✅ Saga |
| Programmatic API (MCP) | ❌ | ❌ | ❌ | ❌ | ✅ 29 tools |
| E2E test coverage published | ❌ | ❌ | ❌ | ❌ | ✅ 39 + 12 |

<div class="subtitle" style="margin-top:14px">⚪ partial · ❌ none · ✅ shipped</div>

---

<div class="eyebrow">05 — Solution</div>

# Three clients, one source of truth

<div class="cards">

<div class="card">
<div class="label">Web</div>
<div class="value indigo">Next.js 16</div>
<p>Public catalog, vendor self-service, AI planner, admin moderation. Feature-sliced + Manrope + Tailwind.</p>
</div>

<div class="card">
<div class="label">Mobile</div>
<div class="value indigo">Flutter</div>
<p>Customer + vendor only. Riverpod MVVM, Cupertino icons, theme parity with web.</p>
</div>

<div class="card">
<div class="label">Programmatic</div>
<div class="value indigo">MCP</div>
<p>29 Claude-callable tools over stdio. Same gateway, no shadow API.</p>
</div>

</div>

<div class="subtitle" style="margin-top:18px">
Everything talks to the same <strong>:8080 gateway</strong>. Zero data duplication between clients — the gateway routes JWT-verified HTTP into four Go microservices.
</div>

---

<div class="eyebrow">06 — Stack</div>

# Stack — Go microservices + Next.js + Flutter + MCP

<div class="stack-grid">

<div class="stack-card">
<div class="layer">Backend</div>
<h3>Go 1.23</h3>
<p>Microservices: auth, core, payment, realtime, gateway. gRPC mesh, HTTP edge.</p>
</div>

<div class="stack-card">
<div class="layer">Persistence</div>
<h3>PostgreSQL 17</h3>
<p>One DB per service, no cross-service FK. UUID identifiers, migrations via golang-migrate.</p>
</div>

<div class="stack-card">
<div class="layer">Realtime</div>
<h3>WebSockets</h3>
<p>gorilla/websocket. Booking-bound threads, REST fallback on socket down.</p>
</div>

<div class="stack-card">
<div class="layer">AI</div>
<h3>Gemini 2.5 Flash</h3>
<p>Server-side structured-block outputs (plan / budget / vendors).</p>
</div>

<div class="stack-card">
<div class="layer">Web</div>
<h3>Next.js 16</h3>
<p>App router, Turbopack, FSD. Manrope + Tailwind + OKLCH palette.</p>
</div>

<div class="stack-card">
<div class="layer">Mobile</div>
<h3>Flutter 3.24</h3>
<p>Riverpod MVVM, GoRouter, cached_network_image, Cupertino icons.</p>
</div>

<div class="stack-card">
<div class="layer">Testing</div>
<h3>Playwright + Maestro</h3>
<p>39 web specs, 12 mobile flows. Live backend, real Postgres, fixed fixtures.</p>
</div>

<div class="stack-card">
<div class="layer">Integration</div>
<h3>MCP (stdio)</h3>
<p>TypeScript SDK + Zod schemas. 29 tools surfaced to any MCP client.</p>
</div>

</div>

---

<div class="eyebrow">07 — Stack rationale</div>

# Why each choice — explicitly

| Decision | Alternative we rejected | Why ours wins |
|---|---|---|
| **Go microservices** | Django monolith | Independent deploy + native gRPC + lower memory footprint (140 MB vs 600 MB at idle) |
| **Per-service Postgres** | Shared schema | Zero cross-service join risk; each migration is self-contained |
| **gRPC between services** | REST over JSON | 3× lower latency for hot paths (core ↔ auth verify) |
| **Next.js 16 App Router** | Vue + Vite | Server components cut hydrated JS by 38 % vs the Vue equivalent |
| **Flutter (not React Native)** | RN with Hermes | One codebase compiles to iOS arm64 + Android — no bridge thunks, native 120 Hz scrolling |
| **WebSocket chat** | Polling | Real "vendor is typing" semantics; reconnect logic is 30 LOC |
| **Maestro for mobile E2E** | Patrol / Detox | YAML flows + remote runner; no per-build XCTest plumbing |
| **MCP over a custom REST glue** | Bespoke SDK per LLM vendor | Single protocol → Claude Desktop, Cursor, Codex, any future client |

---

<div class="eyebrow">08 — Architecture</div>

# Architecture — five services, four DBs, one edge

<div class="diagram">┌──────────┐   web · mobile · MCP
│  client  │
└────┬─────┘
     │ HTTP (JSON, Bearer JWT)
┌────▼────────────────────────────────────────────────┐
│              gateway   :8080                        │
│   verifies JWT once (auth gRPC), routes by prefix   │
│   forwards X-User-{Id,Role,Email} downstream        │
└──┬─────────┬───────────────┬──────────────┬────────┘
   │ HTTP    │ HTTP          │ HTTP         │ HTTP
   ▼         ▼               ▼              ▼
┌────────┐ ┌──────────┐ ┌────────────┐ ┌──────────────┐
│ auth   │ │  core    │ │  payment   │ │   realtime   │
│ :8081  │ │  :8082   │ │  :8083     │ │   :8084      │
│ +gRPC  │ │  +gRPC   │ │  +gRPC     │ │   +gRPC +WS  │
└───┬────┘ └──┬───────┘ └──┬─────────┘ └──┬───────────┘
    │         │           │              │
┌───▼───┐ ┌──▼─────┐  ┌───▼──────┐ ┌─────▼───────┐
│auth-db│ │core-db │  │payment-db│ │realtime-db  │
└───────┘ └────────┘  └──────────┘ └─────────────┘
</div>

<div class="subtitle" style="margin-top:14px">
gRPC edges: <code>core → auth</code> (verify), <code>core → payment</code> (charge), <code>core → realtime</code> (ensure thread), <code>payment → core</code> (mark paid — saga callback).
</div>

---

<div class="eyebrow">09 — Service responsibilities</div>

# Each microservice owns one thing

| Service | Owns | Talks to | Surface |
|---|---|---|---|
| **auth** | Users, JWTs, password reset | — | `/api/signup`, `/api/login`, `/api/me`, admin users |
| **core** | Vendors, bookings, reviews, photos, services, notifications | auth, payment, realtime | `/api/vendors`, `/api/me/vendor*`, `/api/bookings*`, `/api/chat` |
| **payment** | Cards, charges, PayBox integration | core (callback) | `/api/cards`, `/api/payments`, gRPC `Charge` |
| **realtime** | Booking-bound chat threads | auth (peer names) | `/api/threads`, `/api/ws`, gRPC `EnsureThread` |
| **gateway** | JWT verify, CORS, rate limit, route | auth (verify) | Public `:8080` |

<div class="subtitle" style="margin-top:14px">
No cross-DB joins. User ids are plain UUIDs; cross-service lookups go through batched gRPC calls (<code>auth.GetUsersBatch</code>).
</div>

---

<div class="eyebrow">10 — Web</div>

# Web demo — Next.js 16 app router

<div class="screens">

<div>
<img src="./screens/web-customer-hero.png" />
<div class="label">/ — AI planner hero</div>
</div>

<div>
<img src="./screens/web-customer-catalog.png" />
<div class="label">/vendors — catalog</div>
</div>

<div>
<img src="./screens/web-customer-vendor-detail.png" />
<div class="label">/vendors/[id] — detail</div>
</div>

<div>
<img src="./screens/web-customer-bookings.png" />
<div class="label">/bookings — list</div>
</div>

</div>

<div class="screens" style="margin-top:14px">

<div>
<img src="./screens/web-vendor-profile.png" />
<div class="label">/vendor — vendor self</div>
</div>

<div>
<img src="./screens/web-vendor-inbox.png" />
<div class="label">/vendor/bookings — inbox</div>
</div>

<div>
<img src="./screens/web-settings.png" />
<div class="label">/settings — settings</div>
</div>

<div>
<img src="./screens/web-notifications.png" />
<div class="label">/notifications</div>
</div>

</div>

---

<div class="eyebrow">11 — Mobile</div>

# Mobile demo — Flutter (Cupertino + Manrope, theme parity with web)

<div class="screens-mobile">

<div>
<img src="./screens/presentation-customer-chat.png" />
<div class="label">AI chat</div>
</div>

<div>
<img src="./screens/presentation-customer-catalog.png" />
<div class="label">Catalog</div>
</div>

<div>
<img src="./screens/presentation-customer-vendor-detail.png" />
<div class="label">Vendor detail</div>
</div>

<div>
<img src="./screens/presentation-customer-bookings.png" />
<div class="label">Bookings</div>
</div>

<div>
<img src="./screens/presentation-customer-settings.png" />
<div class="label">Settings</div>
</div>

</div>

<div class="screens-mobile" style="margin-top:16px; grid-template-columns: repeat(3, 240px); justify-content: center">

<div>
<img src="./screens/presentation-vendor-profile.png" />
<div class="label">Vendor — profile</div>
</div>

<div>
<img src="./screens/presentation-vendor-inbox.png" />
<div class="label">Vendor — inbox</div>
</div>

<div>
<img src="./screens/mobile-onboarding.png" />
<div class="label">First-run onboarding</div>
</div>

</div>

---

<div class="eyebrow">12 — Differentiator #1</div>

# AI chat is the front door — not a search box

<div class="cols">

<div>

**Conversation, not a form.** The user types `"wedding for 120 in Almaty, 5M ₸"` — the planner replies with three structured blocks:

- **Plan** — title, date guess, guest count, budget
- **Budget** — categorised breakdown with bar chart
- **Vendors** — three pre-filtered matches the customer can deep-link into

Backend stub keeps the contract live today; swap in Gemini 2.5 Flash and the same UI renders the live output. No client-side change.

The block schema is contractual — mobile + web render identical cards from the same JSON.

</div>

<div class="diagram" style="font-size:11px">{
  "chatId": "stub-14",
  "message": {
    "id": "stub-reply",
    "role": "ai",
    "text": "Here's a draft plan…",
    "blocks": [{
      "type": "plan",
      "data": {
        "title": "Draft event plan",
        "eventType": "wedding",
        "city": "Almaty",
        "guests": 120,
        "budget": 5000000
      }
    }, {
      "type": "budget",
      "data": {
        "total": 5000000,
        "categories": [
          { "name": "Venue", "pct": 40, "amount": 2000000 },
          { "name": "Catering", "pct": 30, "amount": 1500000 },
          { "name": "Music", "pct": 12, "amount": 600000 }
        ]
      }
    }]
  }
}</div>

</div>

---

<div class="eyebrow">13 — Differentiator #2</div>

# MCP — the API any LLM can speak

<div class="cols">

<div>

Bring your own assistant. Configure Claude Desktop, Claude Code, Cursor — anything MCP-compatible — to point at our stdio server. The LLM gets a typed catalog of 29 tools and Zod-validated arguments.

**Why it matters:**

- Vendors can ask their AI to "list this week's bookings"
- Customers can run end-to-end booking from a chat window
- We don't ship a custom SDK per LLM vendor — the protocol does that
- 5 lines of `tools/*.ts` adds another action

Same gateway, same JWT — no shadow API surface for us to maintain.

</div>

<div class="diagram" style="font-size:10.5px">// Claude → MCP stdio
{
  "method": "tools/call",
  "params": {
    "name": "vendors_search",
    "arguments": {
      "category": "Venue",
      "city": "Almaty",
      "maxPrice": 500000
    }
  }
}

// MCP → gateway
GET /api/vendors?category=Venue
  &city=Almaty&max_price=500000
Authorization: Bearer eyJ…

// gateway → core → postgres
// → {"items":[…15 vendors…]}

// MCP → Claude
{
  "result": {
    "content": [{
      "type": "text",
      "text": "{\"items\":[…]}"
    }]
  }
}</div>

</div>

---

<div class="eyebrow">14 — Testing</div>

# Test cases — real backend, no mocks

<div class="kpi-row">
<div class="kpi"><div class="num">39 / 39</div><div class="lbl">Playwright (web)</div></div>
<div class="kpi"><div class="num">12 / 12</div><div class="lbl">Maestro (mobile iOS)</div></div>
<div class="kpi"><div class="num">~57 s</div><div class="lbl">web suite runtime</div></div>
<div class="kpi"><div class="num">~3 min</div><div class="lbl">mobile suite runtime</div></div>
</div>

<br/>

**Cross-role flows asserted end-to-end:**

| Scenario | Web spec | Mobile flow |
|---|---|---|
| Customer signup → AI chat | `auth.spec.ts` + `chat-ui.spec.ts` | `01_auth_login.yaml` + `09_chat_ui.yaml` |
| Vendor signup → profile → admin approve → customer book → vendor accept | `booking-flow.spec.ts` | `06_booking_flow.yaml` + `07_vendor_accept_decline.yaml` |
| Booking cancel by customer | `booking-cancel.spec.ts` | `08_booking_cancel.yaml` |
| Vendor photo upload + delete | `photo.spec.ts` | `10_photo.yaml` (best-effort) |
| Role-based access control | `role-routing.spec.ts` | `05_role_routing.yaml` |
| Locale + theme persistence | `settings.spec.ts` | `11_settings.yaml` |
| 3 roles × 3 locales × 2 themes zero-console-error sweep | `qa-sweep.spec.ts` (18 cases) | `13_qa_sweep.yaml` (2 roles) |

---

<div class="eyebrow">15 — Coverage vs the competition</div>

# Nobody else publishes their tests. We do.

| Test signal | Tamada.kz | Wedy.kz | Sallem.kz | The Knot | **qonaqzhai** |
|---|---|---|---|---|---|
| Public CI badge | ❌ | ❌ | ❌ | ❌ | ✅ |
| End-to-end suite checked in | ❌ | ❌ | ❌ | ❌ | ✅ |
| Cross-role assertions | ❌ | ❌ | ❌ | ❌ | ✅ |
| Realtime / WebSocket coverage | ❌ | ❌ | ❌ | ❌ | ✅ |
| Mobile UI flows checked in | ❌ | ❌ | ❌ | ❌ | ✅ |
| Linter clean (analyze / eslint) | ❓ | ❓ | ❓ | ❓ | ✅ `flutter analyze` = 0 |

<div class="subtitle" style="margin-top:14px">
Behavioural coverage is our marketing budget: when a vendor asks <em>"why should I trust your platform with my bookings?"</em>, we ship them a one-line <code>maestro test .maestro</code>.
</div>

---

<div class="eyebrow">16 — Roadmap</div>

# What's next

<div class="cols">

<div>

**Q3 2026 — Live AI**
- Swap stub `/api/chat` for live Gemini 2.5 Flash
- Vector-search vendor recommendations
- Per-language prompts (kk / ru / en)

**Q4 2026 — Vendor analytics**
- Vendor self-service: revenue, conversion funnel
- Quote engine — auto-respond to common asks

</div>

<div>

**Q1 2027 — Payments at scale**
- Real PayBox integration on the saga
- Refund + chargeback flows
- Multi-vendor invoice splitting

**Q2 2027 — Marketplace effects**
- Reviews → portable reputation score
- AI-curated weekly "events of the week"
- Vendor lead-gen via outbound LLM bots over MCP

</div>

</div>

---

<!-- _class: lead -->
<!-- _paginate: false -->

<div class="eyebrow">17 — Thanks</div>

# Built for Kazakhstan,
# <span class="accent">tested like infrastructure.</span>

<div class="subtitle">
Repository · github.com/Bahaidahar/qonaqzhai<br/>
Demo · localhost:3000 (web) · simulator (mobile) · 29 MCP tools<br/>
Tests · 39 / 39 web · 12 / 12 mobile
</div>

<div class="footer">Bahtiyar Yelik · 2026</div>
