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
# <span class="accent">Plan any event<br/>by chatting.</span>

<div class="subtitle">
AI-assisted event services marketplace for Kazakhstan.<br/>
Web · iOS · Android · MCP — one backend, four Go services.
</div>

<div class="footer">Bahtiyar Yelik · Astana IT University</div>

---

<div class="eyebrow">01 — Problem</div>

# 14 phone calls per event.

<div class="cols" style="margin-top: 28px">

<div>
<div class="eyebrow" style="font-size:9px">Customer</div>
<ul>
<li>Vendors live on Instagram, WhatsApp, 2GIS</li>
<li>No prices, no portable reviews</li>
<li>Comparisons happen via group-chat screenshots</li>
<li>Cancellations slip — no booking record</li>
</ul>
</div>

<div>
<div class="eyebrow" style="font-size:9px">Vendor</div>
<ul>
<li>Bookings scattered across four messengers</li>
<li>Cash or Kaspi link — no escrow</li>
<li>Reviews stuck in word-of-mouth</li>
<li>Instagram organic reach flattening since 2024</li>
</ul>
</div>

</div>

---

<div class="eyebrow">02 — Market</div>

# ₸ 480 B / year. <span class="muted">No one platform owns it.</span>

<div class="kpi-row">
<div class="kpi"><div class="num">160 K</div><div class="lbl">weddings + toi / yr</div></div>
<div class="kpi"><div class="num">2.4 M</div><div class="lbl">birthdays / yr</div></div>
<div class="kpi"><div class="num">22 %</div><div class="lbl">YoY mobile commerce</div></div>
<div class="kpi"><div class="num">68 %</div><div class="lbl">vendors on IG only</div></div>
</div>

<div class="subtitle" style="margin-top:20px">
The same photographer shoots Saturday's wedding, Monday's corporate, Sunday's birthday. One platform has to serve <strong>every event type</strong>.
</div>

<div class="footer-note">
BNS RK 2024 · Halyk Finance Q3-2024 · Kaspi investor letter 2024 · in-house survey of 38 vendors (Almaty + Astana + Shymkent, Mar 2026)
</div>

---

<div class="eyebrow">03 — Landscape</div>

# Nobody covers <span style="color:var(--primary)">all events</span> for KZ.

| | Geo | Coverage | Booking | AI | Native app | Realtime chat |
|---|---|---|---|---|---|---|
| Instagram + WhatsApp | KZ | Any (ad-hoc) | DM | ❌ | n/a | DM |
| 2GIS / Yandex Maps | KZ | Catalog + phone | Phone | ❌ | ✅ | ❌ |
| Ticketon.kz | KZ | Tickets only | Buy ticket | ❌ | ✅ | ❌ |
| GigSalad (US) | US | All event types | $20–80 / lead | ❌ | ✅ | ❌ |
| Thumbtack (US) | US | Any pro hire | pay-per-quote | ❌ | ✅ | ❌ |
| Eventbrite | global | Tickets only | Buy ticket | ❌ | ✅ | ❌ |
| **qonaqzhai** | **KZ** | **All events** | **Direct + escrow** | **✅** | **✅** | **✅ WS** |

---

<div class="eyebrow">04 — Feature matrix</div>

# Where the gap actually sits.

| | IG+WA | 2GIS | GigSalad | Thumbtack | Eventbrite | **qonaqzhai** |
|---|---|---|---|---|---|---|
| KZ first | ✅ | ✅ | ❌ | ❌ | ❌ | ✅ |
| All event types | ✅ | ✅ | ✅ | ✅ | ❌ | ✅ |
| Native iOS / Android | ❌ | ✅ | ✅ | ✅ | ✅ | ✅ |
| kk / ru / en | n/a | ⚪ | en | en | en | ✅ |
| AI conversational planner | ❌ | ❌ | ❌ | ❌ | ❌ | ✅ |
| Realtime chat (WS) | ⚪ DM | ❌ | ❌ | ❌ | ❌ | ✅ |
| Escrow payment hold | ❌ | ❌ | ❌ | ❌ | ⚪ | ✅ |
| Flat-fee pricing (no pay-per-lead) | ✅ | ❌ | ❌ | ❌ | n/a | ✅ |
| Programmatic API (MCP) | ❌ | ❌ | ❌ | ❌ | ❌ | ✅ |

---

<div class="eyebrow">05 — Solution</div>

# Three clients, one source of truth.

<div class="cards">

<div class="card">
<div class="label">Web</div>
<div class="value indigo">Next.js 16</div>
<p>Public catalog, vendor self-service, AI planner, admin moderation.</p>
</div>

<div class="card">
<div class="label">Mobile</div>
<div class="value indigo">Flutter</div>
<p>Customer + vendor. Riverpod MVVM. Theme parity with web.</p>
</div>

<div class="card">
<div class="label">Programmatic</div>
<div class="value indigo">MCP</div>
<p>29 Claude-callable tools over stdio. Same gateway, no shadow API.</p>
</div>

</div>

<div class="subtitle" style="margin-top:18px">
Everything hits the same <code>:8080</code> gateway. JWT verified once at the edge, forwarded as <code>X-User-*</code> headers to four Go microservices.
</div>

---

<div class="eyebrow">06 — Stack</div>

# Stack.

<div class="stack-grid">

<div class="stack-card">
<div class="layer">Backend</div>
<h3>Go 1.23</h3>
<p>5 microservices · gRPC mesh · HTTP edge</p>
</div>

<div class="stack-card">
<div class="layer">Persistence</div>
<h3>PostgreSQL 17</h3>
<p>One DB per service · UUID ids · no cross-FK</p>
</div>

<div class="stack-card">
<div class="layer">Realtime</div>
<h3>WebSockets</h3>
<p>gorilla/websocket · REST fallback</p>
</div>

<div class="stack-card">
<div class="layer">AI</div>
<h3>Gemini 2.5</h3>
<p>Structured blocks (plan / budget / vendors)</p>
</div>

<div class="stack-card">
<div class="layer">Web</div>
<h3>Next.js 16</h3>
<p>App router · Turbopack · FSD · OKLCH palette</p>
</div>

<div class="stack-card">
<div class="layer">Mobile</div>
<h3>Flutter 3.24</h3>
<p>Riverpod · GoRouter · Cupertino icons</p>
</div>

<div class="stack-card">
<div class="layer">Testing</div>
<h3>Playwright + Maestro</h3>
<p>39 web specs · 12 mobile flows · live backend</p>
</div>

<div class="stack-card">
<div class="layer">Integration</div>
<h3>MCP (stdio)</h3>
<p>TypeScript SDK · Zod schemas · 29 tools</p>
</div>

</div>

---

<div class="eyebrow">07 — Stack rationale</div>

# Why these, not the obvious alternatives.

| Decision | Alternative rejected | Why |
|---|---|---|
| Go microservices | Django monolith | Independent deploy · native gRPC · ~4× lower idle memory |
| Postgres per service | Shared schema | No cross-service join risk · isolated migrations |
| gRPC between services | REST/JSON | Lower latency on hot paths (core ↔ auth verify) |
| Next.js App Router | Vue + Vite | Server components cut hydrated JS by ~38 % |
| Flutter | React Native | One codebase compiles native · no bridge thunks |
| WebSocket chat | Long-polling | Real "vendor is typing" · reconnect in 30 LOC |
| Maestro | Patrol / Detox | YAML flows · no per-build XCTest plumbing |
| MCP | Bespoke SDK per LLM | One protocol → Claude, Cursor, Codex, anything |

---

<div class="eyebrow">08 — Architecture</div>

# 5 services, 4 databases, 1 edge.

<div class="diagram">┌──────────┐  web · mobile · MCP
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

<div class="subtitle" style="margin-top:10px">
gRPC mesh: <code>core → auth</code> (verify) · <code>core → payment</code> (charge saga) · <code>core → realtime</code> (ensure thread) · <code>payment → core</code> (mark paid callback).
</div>

---

<div class="eyebrow">09 — Service map</div>

# Each service owns one thing.

| Service | Owns | gRPC out | HTTP surface |
|---|---|---|---|
| auth | Users, JWTs, password reset | — | `/api/signup` · `/api/login` · `/api/me` · admin users |
| core | Vendors, bookings, reviews, photos, services, notifications | auth · payment · realtime | `/api/vendors*` · `/api/me/vendor*` · `/api/bookings*` · `/api/chat` |
| payment | Cards, charges, PayBox | core (callback) | `/api/cards` · `/api/payments` · `Charge` |
| realtime | Booking-bound chat | auth (peer names) | `/api/threads` · `/api/ws` · `EnsureThread` |
| gateway | JWT verify, CORS, rate limit | auth (verify) | public `:8080` |

<div class="subtitle" style="margin-top:14px">
Zero cross-DB joins. User ids are plain UUIDs. Cross-service lookups batch through gRPC (<code>auth.GetUsersBatch</code>).
</div>

---

<div class="eyebrow">10 — Web</div>

# Web — Next.js 16 + Manrope + indigo.

<div class="screens">
<div><img src="./screens/web-customer-hero.png" /><div class="label">AI planner hero</div></div>
<div><img src="./screens/web-customer-catalog.png" /><div class="label">Catalog</div></div>
<div><img src="./screens/web-customer-vendor-detail.png" /><div class="label">Vendor detail</div></div>
<div><img src="./screens/web-customer-bookings.png" /><div class="label">Bookings</div></div>
</div>

<div class="screens" style="margin-top:12px">
<div><img src="./screens/web-vendor-profile.png" /><div class="label">Vendor profile</div></div>
<div><img src="./screens/web-vendor-inbox.png" /><div class="label">Vendor inbox</div></div>
<div><img src="./screens/web-settings.png" /><div class="label">Settings</div></div>
<div><img src="./screens/web-notifications.png" /><div class="label">Notifications</div></div>
</div>

---

<div class="eyebrow">11 — Mobile</div>

# Mobile — Flutter, Cupertino icons, theme parity.

<div class="screens-mobile">
<div><img src="./screens/presentation-customer-chat.png" /><div class="label">AI chat</div></div>
<div><img src="./screens/presentation-customer-catalog.png" /><div class="label">Catalog</div></div>
<div><img src="./screens/presentation-customer-vendor-detail.png" /><div class="label">Vendor detail</div></div>
<div><img src="./screens/presentation-customer-bookings.png" /><div class="label">Bookings</div></div>
<div><img src="./screens/presentation-customer-settings.png" /><div class="label">Settings</div></div>
</div>

<div class="screens-mobile" style="margin-top:10px; grid-template-columns: repeat(2, 130px); justify-content: center">
<div><img src="./screens/presentation-vendor-profile.png" /><div class="label">Vendor — profile</div></div>
<div><img src="./screens/presentation-vendor-inbox.png" /><div class="label">Vendor — inbox</div></div>
</div>

---

<div class="eyebrow">12 — Differentiator</div>

# AI is the front door — not a search bar.

<div class="cols" style="margin-top: 16px">

<div>

The user types <code>"toi for 120 in Almaty, 5M ₸"</code>. The planner replies with three structured blocks:

- **Plan** — title, date guess, guests, budget
- **Budget** — bar-charted categorical breakdown
- **Vendors** — three deep-linkable matches

Block schema is contractual. Web and mobile render identical cards from the same JSON. Backend stub today, Gemini swap-in tomorrow — no client change.

</div>

<div class="diagram" style="font-size:9.5px">{
  "chatId": "stub-14",
  "message": {
    "role": "ai",
    "text": "Here's a draft plan…",
    "blocks": [{
      "type": "plan",
      "data": {
        "title": "Draft event plan",
        "eventType": "toi",
        "city": "Almaty",
        "guests": 120,
        "budget": 5000000
      }
    }, {
      "type": "budget",
      "data": {
        "total": 5000000,
        "categories": [
          { "name": "Venue",    "pct": 40, "amount": 2000000 },
          { "name": "Catering", "pct": 30, "amount": 1500000 },
          { "name": "Music",    "pct": 12, "amount":  600000 }
        ]
      }
    }]
  }
}</div>

</div>

---

<div class="eyebrow">13 — Differentiator</div>

# MCP — the API any LLM can speak.

<div class="cols" style="margin-top: 16px">

<div>

Bring-your-own assistant. Configure Claude Desktop, Cursor, Codex — anything MCP-compatible — to point at our stdio server. The LLM gets 29 typed tools with Zod-validated arguments.

<ul style="margin-top:8px">
<li>Vendor asks AI "list this week's bookings"</li>
<li>Customer runs end-to-end booking from a chat</li>
<li>5 lines of <code>tools/*.ts</code> adds another action</li>
<li>Same gateway, same JWT — no shadow API</li>
</ul>

</div>

<div class="diagram" style="font-size:9.5px">// Claude → MCP (stdio)
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
GET /api/vendors?
   category=Venue&city=Almaty
   &max_price=500000
Authorization: Bearer eyJ…

// gateway → core → postgres
// → {"items":[…15 vendors…]}

// MCP → Claude (tool result)
"Rixos Almaty Ballroom from
 500k ₸, Maestro Studio
 from 450k ₸…"
</div>

</div>

---

<div class="eyebrow">14 — Testing</div>

# Real backend. No mocks.

<div class="kpi-row">
<div class="kpi"><div class="num">39/39</div><div class="lbl">Playwright (web)</div></div>
<div class="kpi"><div class="num">12/12</div><div class="lbl">Maestro (mobile iOS)</div></div>
<div class="kpi"><div class="num">~57 s</div><div class="lbl">web runtime</div></div>
<div class="kpi"><div class="num">~3 min</div><div class="lbl">mobile runtime</div></div>
</div>

| Scenario | Web | Mobile |
|---|---|---|
| Customer signup → AI chat | `auth.spec.ts` + `chat-ui.spec.ts` | `01_auth_login` + `09_chat_ui` |
| Vendor profile → admin approve → customer book → vendor accept | `booking-flow.spec.ts` | `06_booking_flow` + `07_vendor_accept` |
| Booking cancel (customer) | `booking-cancel.spec.ts` | `08_booking_cancel` |
| Photo upload / delete | `photo.spec.ts` | `10_photo` |
| Role-based access | `role-routing.spec.ts` | `05_role_routing` |
| Locale + theme persistence | `settings.spec.ts` | `11_settings` |
| 3 roles × 3 locales × 2 themes sweep | `qa-sweep.spec.ts` (18) | `13_qa_sweep` |

---

<div class="eyebrow">15 — Coverage vs the rest</div>

# Nobody else publishes tests. We do.

| | IG+WA | 2GIS / Yandex | GigSalad | Thumbtack | **qonaqzhai** |
|---|---|---|---|---|---|
| Public CI badge | ❌ | ❌ | ❌ | ❌ | ✅ |
| E2E suite in the repo | n/a | ❌ | ❌ | ❌ | ✅ |
| Cross-role assertions | ❌ | ❌ | ❌ | ❌ | ✅ |
| WebSocket coverage | ❌ | ❌ | ❌ | ❌ | ✅ |
| Mobile UI flows in the repo | ❌ | ❌ | ❌ | ❌ | ✅ |
| Linter clean | ❓ | ❓ | ❓ | ❓ | ✅ `flutter analyze = 0` |

<div class="subtitle" style="margin-top:14px">
Behavioural coverage is our marketing budget. Anyone can run <code>maestro test .maestro</code> and watch the full booking flow on their own simulator.
</div>

---

<div class="eyebrow">16 — Roadmap</div>

# What's next.

<div class="cols">

<div>
<div class="eyebrow" style="font-size:9px">Q3 2026 — Live AI</div>
<ul>
<li>Swap stub <code>/api/chat</code> for live Gemini 2.5</li>
<li>Vector-search vendor recommendations</li>
<li>Per-language prompts (kk / ru / en)</li>
</ul>

<div class="eyebrow" style="font-size:9px; margin-top:14px">Q4 2026 — Vendor analytics</div>
<ul>
<li>Revenue + conversion funnel dashboards</li>
<li>Auto-reply quote engine for common asks</li>
</ul>
</div>

<div>
<div class="eyebrow" style="font-size:9px">Q1 2027 — Payments at scale</div>
<ul>
<li>Real PayBox on the saga</li>
<li>Refund + chargeback</li>
<li>Multi-vendor invoice split</li>
</ul>

<div class="eyebrow" style="font-size:9px; margin-top:14px">Q2 2027 — Network effects</div>
<ul>
<li>Reviews → portable reputation</li>
<li>AI-curated weekly events</li>
<li>Vendor lead-gen via outbound LLM bots over MCP</li>
</ul>
</div>

</div>

---

<!-- _class: lead -->
<!-- _paginate: false -->

<div class="eyebrow">17 — Thanks</div>

# Built for KZ,<br/><span class="accent">tested like infrastructure.</span>

<div class="subtitle">
Repository — <code>github.com/Bahaidahar/qonaqzhai</code><br/>
Demo — <code>localhost:3000</code> (web) · simulator (mobile) · 29 MCP tools<br/>
Tests — 39 / 39 web · 12 / 12 mobile
</div>

<div class="footer">Bahtiyar Yelik · 2026</div>
