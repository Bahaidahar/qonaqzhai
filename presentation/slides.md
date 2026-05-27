---
marp: true
size: 16:9
paginate: true
theme: default
style: |
  /* ─── Base — palette mirrors frontend/src/app/globals.css ── */
  :root {
    --brand: #5B47F4;          /* electric indigo — site primary */
    --brand-dark: #4538C7;
    --brand-soft: #EEEBFF;     /* primary-tint surface */
    --ink: #0F1024;            /* foreground oklch(0.14 0.018 275) */
    --ink-soft: #1B1D3A;
    --muted: #6B6C82;
    --bg: #FFFFFF;             /* editorial white */
    --bg-muted: #F6F6FA;
    --card: #FFFFFF;
    --line: #E5E5EE;
  }
  @import url('https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700;800;900&family=JetBrains+Mono:wght@400;500&display=swap');

  section {
    font-family: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
    font-size: 22px;
    color: var(--ink);
    background: var(--bg);
    padding: 60px 70px 50px 70px;
    position: relative;
    letter-spacing: -0.01em;
  }

  /* accent stripe top-right — single indigo accent, no gradient */
  section::before {
    content: '';
    position: absolute;
    top: 0; right: 0;
    width: 240px; height: 4px;
    background: var(--brand);
  }

  /* brand mark bottom-left */
  section::after {
    content: 'QONAQZHAI';
    position: absolute;
    bottom: 22px; left: 70px;
    font-size: 12px;
    font-weight: 800;
    letter-spacing: 0.25em;
    color: var(--brand);
  }

  /* page number top-right */
  section[data-marpit-pagination]::after,
  section:not(.lead):not(.section) {
    counter-increment: page;
  }
  /* page number — use a dedicated element */
  section:not(.lead):not(.section)[data-marpit-pagination] {
    /* default pagination already handled by Marp */
  }

  /* ─── Headings ─────────────────────────────────────── */
  h1 {
    font-size: 52px;
    font-weight: 800;
    color: var(--ink);
    margin: 0 0 22px 0;
    letter-spacing: -0.035em;
    line-height: 1.05;
    padding-bottom: 0;
    border-bottom: none;
    display: block;
  }
  h2 {
    font-size: 26px;
    font-weight: 600;
    color: var(--muted);
    margin: 8px 0 16px 0;
    letter-spacing: -0.018em;
    line-height: 1.3;
  }
  h3 {
    font-size: 13px;
    font-weight: 700;
    color: var(--brand);
    margin: 14px 0 8px 0;
    text-transform: uppercase;
    letter-spacing: 0.18em;
  }
  h4 { font-size: 19px; font-weight: 700; color: var(--ink); margin: 10px 0 4px; }

  p { line-height: 1.55; margin: 6px 0; }
  strong { color: var(--ink); font-weight: 700; }
  em { color: var(--brand); font-style: normal; font-weight: 600; }

  ul, ol { margin: 4px 0 4px 20px; padding: 0; }
  li { margin: 4px 0; line-height: 1.5; }
  li::marker { color: var(--brand); }

  /* ─── Code ─────────────────────────────────────────── */
  code {
    font-family: 'JetBrains Mono', monospace;
    background: var(--brand-soft);
    color: var(--brand-dark);
    padding: 2px 8px;
    border-radius: 4px;
    font-size: 18px;
    font-weight: 500;
  }
  pre {
    background: var(--ink);
    color: #E2E5F0;
    padding: 22px 26px;
    border-radius: 12px;
    font-size: 15px;
    line-height: 1.45;
    box-shadow: 0 10px 30px -10px rgba(15, 16, 36, 0.35);
    margin: 12px 0;
    overflow-x: auto;
  }
  pre code {
    background: transparent;
    color: inherit;
    padding: 0;
    font-size: 15px;
    font-weight: 400;
  }

  /* ─── Tables ───────────────────────────────────────── */
  table {
    border-collapse: separate;
    border-spacing: 0;
    width: 100%;
    font-size: 17px;
    margin: 10px 0;
    border-radius: 10px;
    overflow: hidden;
    border: 1px solid var(--line);
  }
  thead th {
    background: var(--bg-muted);
    color: var(--ink);
    padding: 12px 14px;
    text-align: left;
    font-weight: 700;
    font-size: 13px;
    text-transform: uppercase;
    letter-spacing: 0.1em;
    border-bottom: 1px solid var(--line);
  }
  td {
    background: var(--card);
    padding: 11px 14px;
    border-top: 1px solid var(--line);
    vertical-align: top;
  }
  tbody tr:first-child td { border-top: none; }
  tbody tr:hover td { background: var(--brand-soft); }

  /* ─── Links ────────────────────────────────────────── */
  a { color: var(--brand-dark); text-decoration: none; border-bottom: 1.5px solid var(--brand); }

  /* ─── Pagination ──────────────────────────────────── */
  section::after {
    /* override above for pagination */
  }
  /* pagination styled by marp natively */

  /* ─── LEAD (title / closing) ──────────────────────── */
  section.lead {
    background:
      radial-gradient(circle at 18% 22%, rgba(91,71,244,0.30) 0%, transparent 55%),
      radial-gradient(circle at 82% 78%, rgba(91,71,244,0.18) 0%, transparent 60%),
      linear-gradient(160deg, #0F1024 0%, #15163B 60%, #1B1D3A 100%);
    color: #fff;
    padding: 80px 90px;
  }
  section.lead::before {
    width: 100%;
    height: 4px;
    background: var(--brand);
  }
  section.lead::after {
    color: rgba(255,255,255,0.5);
    bottom: 40px;
  }
  section.lead h1 {
    font-size: 96px;
    color: #fff;
    border: none;
    padding: 0;
    margin: 0 0 20px 0;
    font-weight: 900;
    letter-spacing: -0.045em;
    line-height: 0.95;
  }
  section.lead h2 {
    color: rgba(255,255,255,0.92);
    font-size: 32px;
    font-weight: 400;
    margin: 0 0 40px 0;
    max-width: 80%;
    letter-spacing: -0.015em;
  }
  section.lead h3 {
    color: var(--brand);
    font-size: 15px;
    letter-spacing: 0.2em;
    margin-top: 60px;
    font-weight: 700;
  }
  section.lead p { color: rgba(255,255,255,0.85); font-size: 22px; }
  section.lead strong { color: var(--brand); }
  section.lead a { color: var(--brand); border-color: var(--brand); }

  /* ─── SECTION DIVIDER ─────────────────────────────── */
  section.section {
    background: var(--ink);
    color: #fff;
    padding: 100px 90px;
    display: flex;
    flex-direction: column;
    justify-content: center;
  }
  section.section::before {
    width: 240px;
    height: 4px;
    background: var(--brand);
    top: 0; left: 0; right: auto;
  }
  section.section::after {
    color: rgba(255,255,255,0.45);
  }
  section.section .num {
    color: var(--brand);
    font-size: 140px;
    font-weight: 900;
    line-height: 0.85;
    letter-spacing: -0.05em;
    margin-bottom: 24px;
    opacity: 0.95;
  }
  section.section h1 {
    color: #fff;
    font-size: 64px;
    border: none;
    padding: 0;
    margin: 0;
    font-weight: 800;
    letter-spacing: -0.035em;
  }
  section.section h2 {
    color: rgba(255,255,255,0.72);
    font-size: 24px;
    font-weight: 400;
    margin-top: 16px;
    max-width: 70%;
    letter-spacing: -0.015em;
  }

  /* ─── DARK content slide ──────────────────────────── */
  section.dark {
    background: var(--ink);
    color: #E2E5F0;
  }
  section.dark::before { background: var(--brand); }
  section.dark h1 { color: #fff; border-bottom-color: var(--brand); }
  section.dark h2 { color: rgba(255,255,255,0.85); }
  section.dark h3 { color: var(--brand); }
  section.dark strong { color: #fff; }
  section.dark em { color: var(--brand); }
  section.dark td { background: rgba(255,255,255,0.04); color: #E2E5F0; border-top-color: rgba(255,255,255,0.08); }
  section.dark tbody tr:hover td { background: rgba(91,71,244,0.12); }
  section.dark::after { color: var(--brand); }
  section.dark code { background: rgba(91,71,244,0.22); color: #C7C0FF; }

  /* ─── Helpers ─────────────────────────────────────── */
  .columns-2 { display: grid; grid-template-columns: 1fr 1fr; gap: 32px; margin-top: 16px; }
  .columns-3 { display: grid; grid-template-columns: 1fr 1fr 1fr; gap: 22px; margin-top: 16px; }

  .card {
    background: var(--card);
    border: 1px solid var(--line);
    border-left: 3px solid var(--brand);
    border-radius: 10px;
    padding: 18px 22px;
    box-shadow: 0 1px 3px rgba(15, 16, 36, 0.04);
  }
  .card.tint {
    background: var(--brand-soft);
    border-color: rgba(91,71,244,0.18);
    border-left-color: var(--brand);
  }
  .card.dark {
    background: var(--ink);
    color: #fff;
    border-color: var(--ink-soft);
    border-left-color: var(--brand);
  }
  .card.dark strong { color: var(--brand); }
  .card h4 {
    margin-top: 0;
    color: var(--brand);
    text-transform: uppercase;
    font-size: 13px;
    letter-spacing: 0.12em;
    font-weight: 700;
  }
  section.dark .card { background: rgba(255,255,255,0.04); border-color: rgba(255,255,255,0.08); color: #E2E5F0; }
  section.dark .card h4 { color: var(--brand); }
  section.dark .card strong { color: #fff; }

  .pill {
    display: inline-block;
    background: var(--ink);
    color: #fff;
    padding: 5px 14px;
    border-radius: 6px;
    font-size: 13px;
    font-weight: 600;
    margin: 3px 4px 3px 0;
    letter-spacing: 0.02em;
  }
  .pill.brand { background: var(--brand); }
  .pill.tint { background: var(--brand-soft); color: var(--brand-dark); }
  .pill.outline { background: transparent; color: var(--ink); border: 1.5px solid var(--line); }

  .stat-grid {
    display: grid;
    grid-template-columns: repeat(5, 1fr);
    gap: 18px;
    margin: 20px 0;
  }
  .stat {
    background: var(--card);
    border: 1px solid var(--line);
    border-radius: 12px;
    padding: 22px 18px;
    text-align: center;
    box-shadow: 0 1px 3px rgba(15, 16, 36, 0.04);
  }
  .stat .n {
    display: block;
    color: var(--brand);
    font-size: 46px;
    font-weight: 900;
    line-height: 1;
    letter-spacing: -0.035em;
  }
  .stat .l {
    display: block;
    color: var(--muted);
    font-size: 13px;
    margin-top: 8px;
    text-transform: uppercase;
    letter-spacing: 0.08em;
    font-weight: 600;
  }

  .kicker {
    display: inline-block;
    color: var(--brand);
    font-size: 13px;
    font-weight: 700;
    letter-spacing: 0.2em;
    text-transform: uppercase;
    margin-bottom: 10px;
    padding: 4px 0;
    border-bottom: 2px solid var(--brand);
  }

  .lead-stat {
    font-size: 92px;
    font-weight: 900;
    color: var(--brand);
    line-height: 1;
    letter-spacing: -0.045em;
  }
  .lead-stat .unit { color: var(--ink); font-size: 36px; font-weight: 700; margin-left: 6px; }

  .timeline {
    display: grid;
    grid-template-columns: 30px 1fr;
    gap: 14px;
    margin: 10px 0;
  }
  .timeline .num {
    background: var(--brand);
    color: #fff;
    width: 30px; height: 30px;
    border-radius: 50%;
    display: flex; align-items: center; justify-content: center;
    font-weight: 800;
    font-size: 14px;
  }
---

<!-- _class: lead -->
<!-- _paginate: false -->

<span class="kicker" >Diploma Defense · May 2026</span>

# Qonaqzhai

## AI-powered event services marketplace for Kazakhstan.
## Connecting customers, vendors, and traditions.

### Bahaidahar · Astana IT University

---

<!-- _class: section -->

<div class="num">01</div>

# The Problem

## Event planning in Kazakhstan is fragmented, untrustworthy, and ignores local culture.

---

# Status quo is broken

<div class="columns-2">

<div class="card">
<h4>How it works today</h4>

- Discovery via Instagram DMs
- WhatsApp chains for quotes
- No verified reviews
- Deposits to private cards = trust risk
- Customer juggles 5–10 vendors
- No platform supports Kazakh first-class
- Global tools ignore traditions
</div>

<div class="card tint">
<h4>Market size</h4>

<div class="lead-stat">$1.5<span class="unit">B</span></div>

KZ events sector per year — fragmented across Instagram, fl.kz, weddingsalon.kz. **No platform owns the full journey.**
</div>

</div>

---

<!-- _class: section -->

<div class="num">02</div>

# The Solution

## Three clients, one platform. AI-first. Bilingual KZ / RU + EN.

---

# Qonaqzhai at a glance

```
                Customer                          Vendor
                   │                                │
                   ▼                                ▼
        ┌──────────────────────────────────────────────┐
        │  AI Planner    Search       Reviews          │
        │  (Claude)                                    │
        │  Booking       Calendar     Analytics        │
        │  Payment       Chat         Dashboard        │
        │  (PayBox)      (WS)                          │
        └──────────────────────────────────────────────┘
                              │
                              ▼
                         Qonaqzhai
                  Web  ·  iOS  ·  Android
```

<div class="columns-3">
<div><span class="pill brand">Trilingual KZ / RU / EN</span></div>
<div><span class="pill brand">AI event planner</span></div>
<div><span class="pill brand">PayBox · KZT-native</span></div>
</div>

---

# By the numbers

<div class="stat-grid">
  <div class="stat"><span class="n">5</span><span class="l">Go services</span></div>
  <div class="stat"><span class="n">4</span><span class="l">Postgres DBs</span></div>
  <div class="stat"><span class="n">16</span><span class="l">gRPC methods</span></div>
  <div class="stat"><span class="n">68</span><span class="l">HTTP endpoints</span></div>
  <div class="stat"><span class="n">16</span><span class="l">DB tables</span></div>
</div>

<div class="stat-grid">
  <div class="stat"><span class="n">13</span><span class="l">Mobile features</span></div>
  <div class="stat"><span class="n">16</span><span class="l">Web routes</span></div>
  <div class="stat"><span class="n">2,033</span><span class="l">Go test LOC</span></div>
  <div class="stat"><span class="n">12</span><span class="l">E2E specs</span></div>
  <div class="stat"><span class="n">3</span><span class="l">Languages</span></div>
</div>

---

<!-- _class: section -->

<div class="num">03</div>

# The Stack

## Picked layer by layer against measured 2025-2026 industry data.

---

# Technology decisions

<div class="columns-2">

<div class="card">
<h4>Backend</h4>

<span class="pill">Go 1.24</span> <span class="pill brand">gRPC</span> <span class="pill">Postgres 14</span>

5 microservices · gateway / auth / core / payment / realtime
Protobuf contracts inter-service · HTTP/JSON at edge
JWT · PayBox PSP
</div>

<div class="card">
<h4>Web</h4>

<span class="pill">Next.js 16</span> <span class="pill brand">React 19</span> <span class="pill">FSD</span>

App Router · Cache Components · Partial pre-rendering
Tailwind v4 · Playwright E2E
</div>

<div class="card">
<h4>Mobile</h4>

<span class="pill">Flutter 3.24</span> <span class="pill brand">Riverpod</span> <span class="pill">MVVM</span>

Impeller engine · go_router · Dio + auth interceptor
Firebase Messaging · Secure storage
</div>

<div class="card tint">
<h4>AI Layer</h4>

<span class="pill brand">Claude Sonnet 4.6</span> <span class="pill">MCP</span>

Tool-use event planner
Agentic dev workflow with MCP servers
</div>

</div>

---

<!-- _class: section -->

<div class="num">04</div>

# Architecture

## Domain-driven · service-isolated · contract-first.

---

# Domain model — four bounded contexts

```
┌──────────┐   ┌─────────────────────────┐   ┌──────────┐   ┌──────────┐
│   Auth   │   │          Core           │   │ Payment  │   │ Realtime │
├──────────┤   ├─────────────────────────┤   ├──────────┤   ├──────────┤
│ User     │   │ Vendor                  │   │ Card     │   │ Thread   │
│ Refresh  │   │ Service                 │   │ Payment  │   │ Message  │
│ Reset    │   │ Booking (state machine) │   │          │   │          │
│          │   │ Review · Photo          │   │          │   │          │
│          │   │ Notification · FCM      │   │          │   │          │
└──────────┘   └─────────────────────────┘   └──────────┘   └──────────┘
```

**Booking lifecycle** — `pending → accepted → completed → paid`
Branches — `declined` · `cancelled`. Explicit state machine in `core/usecase/booking/`.

---

<!-- _class: dark -->

# Backend topology

```
┌──────────┐
│  client  │  mobile + web
└────┬─────┘
     │ HTTP
┌────▼──────────────────────────────────────────────────┐
│                  gateway :8080                        │
│   verifies JWT once · routes by prefix                │
└──┬──────────┬───────────────┬──────────────┬─────────┘
   │          │               │              │
   ▼          ▼               ▼              ▼
┌────────┐ ┌──────────┐ ┌────────────┐ ┌──────────────┐
│ auth   │ │  core    │ │  payment   │ │   realtime   │
│ :8081  │ │  :8082   │ │  :8083     │ │   :8084 WS   │
└───┬────┘ └──┬───────┘ └──┬─────────┘ └──┬───────────┘
    │         │           │              │
┌───▼───┐ ┌──▼─────┐  ┌───▼──────┐ ┌─────▼───────┐
│auth-db│ │core-db │  │payment-db│ │realtime-db  │
└───────┘ └────────┘  └──────────┘ └─────────────┘
```

**One DB per service** · **gRPC inter-service** · **HTTP at edge**

---

# Service contracts — 16 gRPC methods

```
core ─────► auth      GetUser · GetUsersBatch · VerifyToken
core ─────► payment   Charge · Refund
core ─────► realtime  EnsureThread · PublishEvent
payment ──► core      MarkBookingPaid          (saga callback)
realtime ─► auth      GetUsersBatch            (peer-name enrich)
gateway ──► auth      VerifyToken              (once per request)
```

<div class="card tint">
<h4>Payment saga — synchronous</h4>

`core.booking.Pay` → `payment.Charge` → `core.MarkBookingPaid`

Diploma scale — synchronous wins on simplicity.
Production migration path — event bus (NATS / Kafka).
</div>

---

# Clean Architecture per service

```
backend/services/core/internal/
├── domain/         Pure types — Vendor, Booking, Review (no I/O)
├── ports/          Interfaces — VendorRepo, PaymentClient, Notifier
├── usecase/        Business logic — vendor/, booking/, review/, photo/
└── adapter/
    ├── http/       chi router · JSON marshalling · middleware
    ├── grpc/       gRPC server impl (core.proto)
    ├── grpcclient/ Outbound to auth, payment, realtime
    ├── repo/       PostgreSQL (sqlx)
    └── push/       FCM HTTP v1 sender
```

Layers depend inward only — `adapter → usecase → domain`. **Dependency Inversion** keeps usecase pure and testable. Table-driven Go tests stub the ports.

---

# Web — Feature-Sliced Design

<div class="columns-2">

<div>

```
frontend/src/
├── app/        Routes — vendors, bookings,
│               threads, admin
├── features/   ai-chat · auth · booking-flow
│               messaging · reviews · services
│               theme · vendor-search
├── entities/   Shared domain types
├── widgets/    Composed UI blocks
└── shared/     UI kit · API · i18n · hooks
```

</div>

<div>

<div class="card">
<h4>Why FSD</h4>

Strict horizontal isolation. Dependency rule:
`shared → entities → features → widgets → app`
</div>

<div class="card tint" style="margin-top:14px">
<h4>Next.js 16 features used</h4>

- Cache Components with `"use cache"`
- Partial pre-rendering on vendor lists
- React Server Components on vendor detail
</div>

</div>

</div>

---

# Mobile — Flutter · Riverpod · Clean Arch

```
mobile/lib/
├── core/         network · router · theme · i18n · DI
├── features/    13 modules — auth · vendor_catalog · booking
│                 ai_chat · payment · cards · reviews · messaging
│                 notifications · vendor_self · admin · onboarding
│                 settings
│   └── <feature>/
│       ├── data/         DTOs · datasources · repository impl
│       ├── domain/       Entities · abstract repo · use cases
│       └── presentation/ ViewModels (Riverpod) · screens
```

<div class="columns-3">
  <div><span class="pill brand">Riverpod — compile-safe DI</span></div>
  <div><span class="pill brand">go_router — deep links</span></div>
  <div><span class="pill brand">Dio — auto JWT refresh</span></div>
</div>

---

# Data model — 16 tables, 4 databases

| Database | Tables | Highlights |
|---|---|---|
| `auth-db` | `users`, `refresh_tokens`, `password_reset_tokens` | role enum · hashed tokens |
| `core-db` | `vendors`, `services`, `photos`, `bookings`, `reviews`, `notifications`, `fcm_tokens` | full-text GIN · rating aggregation · BYTEA blobs |
| `payment-db` | `cards`, `payments` | last4 + brand · PAN never persisted |
| `realtime-db` | `threads`, `thread_messages` | one thread per booking · append-only |

**No cross-service foreign keys.** `core.bookings.customer_id` is a plain UUID — `auth-svc` owns the actual `users` row. Cross-refs resolved via batched gRPC `auth.GetUsersBatch`.

---

<!-- _class: dark -->

# Real-time chat — WebSocket hub

```
mobile ──WS──► gateway :8080 ──upgrade──► realtime :8084
                                                │
                                                ▼
                                        ┌────────────┐
                                        │   Hub      │
                                        │ per-user   │
                                        │  channels  │
                                        └─────┬──────┘
                                              │
                                              ▼
                                    PostgreSQL append-only
                                    threads + thread_messages
```

One thread per booking · Hub fans out and persists · JWT verified at upgrade · FCM deep-link fallback when peer offline

---

<!-- _class: section -->

<div class="num">05</div>

# Why This Stack

## Every choice backed by 2025–2026 industry data.

---

# Stack justification — comparative

| Layer | Choice | Alternative | Why we picked it |
|---|---|---|---|
| Backend | **Go 1.24** | Node.js · Java | ~2.6× CPU-bound · mature gRPC tooling |
| Inter-svc | **gRPC + Protobuf** | REST / JSON | 5–10× throughput · 70–90% smaller payloads |
| Web | **Next.js 16 + React 19** | Remix · SvelteKit | Cache Components + PPR stable in v16 |
| Web arch | **Feature-Sliced Design** | Atomic · flat | Hard horizontal isolation |
| Mobile | **Flutter 3.24 Impeller** | React Native | 46% vs 35% share · 60–120 fps |
| Mobile state | **Riverpod** | Bloc · Provider | Compile-safe DI · no `BuildContext` |
| Database | **Postgres per service** | Shared DB | Service autonomy · independent scaling |
| AI | **Claude API + MCP** | OpenAI raw | MCP is cross-vendor standard (Linux Foundation) |
| E2E | **Playwright** | Cypress | Parallel cross-browser · TS-native |

---

<!-- _class: section -->

<div class="num">06</div>

# Competitive Landscape

## Where Qonaqzhai sits among global · CIS · Kazakh players.

---

# Direct competitor matrix

| Platform | Region | Model | AI | Mobile | KZ Traditions |
|---|---|---|---|---|---|
| The Bash | US | Commission | No | Yes | No |
| GigSalad | US | Freemium | Limited | Yes | No |
| Peerspace | US | Commission + escrow | Recs | Yes | No |
| The Knot | US | Listing + ads | Checklist | Yes | No |
| Wedding.ru | RU | Paid listings | No | Web-first | No |
| Eventie.kz | KZ | Paid listings | No | Web | Partial |
| weddingsalon.kz | KZ | Directory | No | Web | Partial |
| **Qonaqzhai** | **KZ** | **Commission + escrow** | **Tool-use planner** | **Native Flutter** | **Yes** |

---

# What competitors miss — what we built

<div class="columns-2">

<div class="card">
<h4>Gaps we fill</h4>

- Kazakh-first UX (KZ / RU / EN equal)
- Local PayBox · KZT-native
- AI planner with real **tool use**
- Mobile-native, not web-first
- Tradition bundles — тұсаукесер · беташар · шашу
</div>

<div class="card tint">
<h4>Inspiration ported in</h4>

- Peerspace visual-first listings
- GigSalad freemium tier for vendors
- The Knot budget tracker
- Eventbrite promoted slots
- Peerspace escrow with milestone release
</div>

</div>

---

<!-- _class: section -->

<div class="num">07</div>

# Signature Feature

## AI Event Planner with tool use.

---

<!-- _class: dark -->

# Claude calls real backend tools

```
User: "Plan corporate event for 50 in Astana, 20 Aug, budget 1.5M ₸"
                              │
                              ▼
              Claude Sonnet 4.6 (system prompt + tools)
                              │
        ┌────────────┬────────┴────────┬─────────────┐
        ▼            ▼                 ▼             ▼
  search_vendors  check_availability  draft_booking  estimate_total
   (core gRPC)     (core gRPC)         (core gRPC)    (local calc)
                              │
                              ▼
                    Structured plan in chat
        • Venue × 1            180 k ₸
        • Catering × 50        900 k ₸
        • Photographer × 1     200 k ₸
        • Music DJ × 1         150 k ₸
        ─────────────────────────────────
        Total: 1.43 M ₸     [Book all]  [Customize]
```

Anthropic Messages API · `tools` parameter · JSON schemas · SSE streaming from `core` service.

---

<!-- _class: section -->

<div class="num">08</div>

# Security & Quality

## Defense in depth · table-driven tests · E2E gate.

---

# Security model — defense in depth

| Layer | Mechanism |
|---|---|
| Transport | HTTPS at gateway · Let's Encrypt |
| AuthN | JWT HS256 · 15 min access + 7 day refresh · hashed |
| AuthZ | RBAC middleware · per-route guards · ownership checks |
| Edge | Gateway calls `auth.VerifyToken` on every request |
| Service | Each service re-verifies JWT — no implicit trust |
| Rate limit | Per-IP 20–40 req / 10 s on every service |
| Password reset | One-time short-TTL tokens · SMTP delivery |
| PCI scope | PAN never persisted — PayBox PSP holds the scope |
| Photo uploads | MIME whitelist · 5 MB cap · content-type sniffed |
| Cross-DB joins | **Forbidden** — no shared schema, no leaked FKs |

---

# Testing & quality

| Component | Files | LOC | Approach |
|---|---|---|---|
| Auth service | 1 | 386 | Signup · login · refresh · password reset |
| Core service | 4 | 1,007 | Vendor CRUD · booking state machine · review aggregation |
| Payment service | 2 | 436 | Card validation · charge · refund · PSP mock |
| Realtime service | 1 | 204 | Thread ensure · message send · WS connection |
| **Backend total** | **8** | **2,033** | Table-driven Go tests · `-race` enabled |
| Frontend E2E | 12 | ~800 | Playwright — auth · booking · chat · admin · vendor |
| Mobile | smoke | growing | Riverpod `ProviderContainer` mocks |

**CI gate** — backend tests + Playwright E2E must pass before merge.

---

<!-- _class: section -->

<div class="num">09</div>

# Agentic Engineering

## MCP, Claude, and how this diploma was actually built.

---

# How this diploma was built

<div class="columns-2">

<div class="card">
<h4>MCP servers used in development</h4>

- **GitHub MCP** — branch / PR / issue automation
- **Filesystem MCP** — repo-wide refactors
- **Postgres MCP** — schema introspection
- **Context7 MCP** — Next.js 16 + Riverpod docs
- **Exa / WebSearch** — competitor research
</div>

<div class="card tint">
<h4>Why this matters in 2026</h4>

- MCP donated to **Linux Foundation Agentic AI Foundation** (2025)
- **10K+** MCP servers live · **97M+** SDK installs
- Stack Overflow 2025 — **84%** of developers use AI tools
- IDC — **$36.2B** AI dev-tool market by end-2026
</div>

</div>

This diploma mirrors how engineering is actually done today.

---

<!-- _class: section -->

<div class="num">10</div>

# Demo & Roadmap

## What you'll see today, and what ships next.

---

# Live demo — 8 minutes

<div class="timeline"><div class="num">1</div><div><strong>Localization</strong> · 30s · Open mobile · switch KZ → RU → EN</div></div>
<div class="timeline"><div class="num">2</div><div><strong>AI Planner</strong> · 2 min · <em>"Plan corporate event for 50 in Astana, 20 Aug, 1.5M ₸"</em> · Claude returns 4-vendor plan with reasoning</div></div>
<div class="timeline"><div class="num">3</div><div><strong>Booking + Calendar</strong> · 1 min · Tap venue · view availability · confirm Aug 20</div></div>
<div class="timeline"><div class="num">4</div><div><strong>Payment</strong> · 1 min · PayBox test card · push notification · deep-link into booking</div></div>
<div class="timeline"><div class="num">5</div><div><strong>Vendor side</strong> · 1 min · Log in as vendor · see new booking in analytics · revenue chart updates</div></div>
<div class="timeline"><div class="num">6</div><div><strong>Admin moderation</strong> · 30s · Web admin moderates flagged review</div></div>
<div class="timeline"><div class="num">7</div><div><strong>Architecture trace</strong> · 1 min · <code>docker compose logs</code> — saga <code>core → payment.Charge → core.MarkBookingPaid</code></div></div>
<div class="timeline"><div class="num">8</div><div><strong>Q & A</strong> · 1 min</div></div>

---

# Roadmap

<div class="columns-3">

<div class="card">
<h4>Tier 1 — ships before defense</h4>

- AI event planner with tool use
- Vendor analytics dashboard
- KZ / RU / EN i18n completion
- Push notifications with deep-links
- Chat polish — typing · read receipts
</div>

<div class="card tint">
<h4>Tier 2 — next iteration</h4>

- Smart vendor recommendations
- Dynamic pricing (peak / off-peak)
- Google Calendar sync
- Referral · promo codes
- Verified-booking review badge
</div>

<div class="card dark">
<h4>Tier 3 — research</h4>

- Escrow with milestone release
- Voice input for AI planner
- ML matching via embeddings
- KYC via egov.kz public API
- AI-generated event runbook
</div>

</div>

---

# Lessons learned — honest retrospective

| What worked | What I would change |
|---|---|
| One DB per service — clear ownership | Synchronous payment saga is fragile · publish events instead (NATS / Kafka) |
| Riverpod for Flutter — fewer rebuilds | Hand-rolled DTOs · should generate from OpenAPI from day 1 |
| FSD on web — features stay isolated | Tailwind v4 + Next.js 16 still rough · hydration warnings remain |
| Claude + MCP — 2–3× productivity | Need stricter prompt discipline · early sessions drifted scope |
| Per-service Go modules — independent CI | Workspace `replace` directives are subtle · broke CI twice |
| gRPC for internal — no API drift | Browser debugging painful · we keep REST at the gateway for a reason |

---

<!-- _class: lead -->
<!-- _paginate: false -->

<span class="kicker" >Thank you</span>

# Questions?

## Repo · github.com/Bahaidahar/diploma
## Architecture deep-dive · presentation/research/code-map.md
## Stack justification · presentation/research/stack-justification.md
## Comparative analysis · presentation/research/competitors.md

### Bahaidahar · Astana IT University · May 2026
