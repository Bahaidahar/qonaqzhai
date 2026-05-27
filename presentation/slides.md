---
marp: true
size: 16:9
paginate: true
theme: default
style: |
  :root {
    --brand: #5B47F4;
    --brand-dark: #4538C7;
    --brand-soft: #EEEBFF;
    --ink: #0F1024;
    --ink-soft: #1B1D3A;
    --muted: #6B6C82;
    --bg: #FFFFFF;
    --bg-muted: #F6F6FA;
    --line: #E5E5EE;
  }
  @import url('https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700;800;900&family=JetBrains+Mono:wght@400;500&display=swap');

  section {
    font-family: 'Inter', -apple-system, sans-serif;
    font-size: 24px;
    color: var(--ink);
    background: var(--bg);
    padding: 70px 90px 60px 90px;
    position: relative;
    letter-spacing: -0.01em;
  }

  /* small mark in corner — restrained */
  section::after {
    content: attr(data-marpit-pagination) ' / 18 · Qonaqzhai';
    position: absolute;
    bottom: 28px; right: 90px;
    font-size: 12px;
    font-weight: 500;
    color: var(--muted);
    letter-spacing: 0.05em;
  }

  h1 {
    font-size: 54px;
    font-weight: 800;
    color: var(--ink);
    margin: 0 0 24px 0;
    letter-spacing: -0.04em;
    line-height: 1.05;
  }
  h2 {
    font-size: 26px;
    font-weight: 400;
    color: var(--muted);
    margin: 0 0 22px 0;
    letter-spacing: -0.018em;
    line-height: 1.4;
    max-width: 78%;
  }
  h3 {
    font-size: 13px;
    font-weight: 700;
    color: var(--brand);
    margin: 0 0 18px 0;
    text-transform: uppercase;
    letter-spacing: 0.22em;
  }
  h4 { font-size: 19px; font-weight: 700; color: var(--ink); margin: 12px 0 4px; }

  p { line-height: 1.55; margin: 8px 0; max-width: 90%; }
  strong { color: var(--ink); font-weight: 700; }
  em { color: var(--brand); font-style: normal; font-weight: 600; }

  ul { margin: 8px 0 8px 22px; padding: 0; }
  li { margin: 7px 0; line-height: 1.5; }
  li::marker { color: var(--brand); }

  code {
    font-family: 'JetBrains Mono', monospace;
    background: var(--brand-soft);
    color: var(--brand-dark);
    padding: 2px 7px;
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
    line-height: 1.5;
    margin: 14px 0;
  }
  pre code {
    background: transparent;
    color: inherit;
    padding: 0;
    font-size: 15px;
  }

  table {
    border-collapse: separate;
    border-spacing: 0;
    width: 100%;
    font-size: 18px;
    margin: 12px 0;
    border-radius: 8px;
    overflow: hidden;
    border: 1px solid var(--line);
  }
  thead th {
    background: var(--bg-muted);
    color: var(--ink);
    padding: 12px 16px;
    text-align: left;
    font-weight: 700;
    font-size: 13px;
    text-transform: uppercase;
    letter-spacing: 0.1em;
  }
  td {
    padding: 12px 16px;
    border-top: 1px solid var(--line);
    vertical-align: top;
  }
  tbody tr:first-child td { border-top: none; }

  a { color: var(--brand-dark); border-bottom: 1.5px solid var(--brand); text-decoration: none; }

  blockquote {
    border: none;
    padding: 0;
    margin: 28px 0;
    font-size: 38px;
    font-weight: 600;
    line-height: 1.25;
    color: var(--ink);
    letter-spacing: -0.025em;
    max-width: 88%;
  }
  blockquote::before {
    content: '“';
    color: var(--brand);
    font-size: 80px;
    line-height: 0;
    vertical-align: -28px;
    margin-right: 6px;
    font-weight: 800;
  }
  blockquote p { margin: 0; max-width: 100%; }

  /* ─── COVER ────────────────────────────────────────── */
  section.cover {
    background: var(--bg);
    padding: 100px 110px;
  }
  section.cover::after { display: none; }
  section.cover .stripe {
    width: 80px; height: 4px; background: var(--brand);
    margin-bottom: 36px;
  }
  section.cover h1 {
    font-size: 110px;
    letter-spacing: -0.055em;
    line-height: 0.92;
    margin: 0 0 28px 0;
  }
  section.cover h2 {
    font-size: 30px;
    color: var(--ink);
    font-weight: 500;
    max-width: 70%;
    line-height: 1.3;
    margin: 0 0 64px 0;
  }
  section.cover .meta {
    font-size: 16px;
    color: var(--muted);
    letter-spacing: 0.06em;
  }
  section.cover .meta strong {
    color: var(--ink);
    font-weight: 700;
  }

  /* ─── DARK ─────────────────────────────────────────── */
  section.dark {
    background: var(--ink);
    color: #E2E5F0;
  }
  section.dark::after { color: rgba(255,255,255,0.4); }
  section.dark h1 { color: #fff; }
  section.dark h2 { color: rgba(255,255,255,0.7); }
  section.dark h3 { color: var(--brand); }
  section.dark strong { color: #fff; }
  section.dark code { background: rgba(91,71,244,0.22); color: #C7C0FF; }
  section.dark td { color: #E2E5F0; border-top-color: rgba(255,255,255,0.08); }
  section.dark thead th { background: rgba(255,255,255,0.05); color: #fff; }
  section.dark blockquote { color: #fff; }

  /* ─── QUOTE / BREATHER ─────────────────────────────── */
  section.quote {
    background: var(--bg-muted);
    padding: 110px 130px;
    display: flex;
    flex-direction: column;
    justify-content: center;
  }
  section.quote::after { color: var(--muted); }
  section.quote .stripe {
    width: 60px; height: 3px; background: var(--brand);
    margin-bottom: 30px;
  }
  section.quote blockquote {
    font-size: 46px;
    line-height: 1.2;
    margin: 0 0 24px 0;
  }
  section.quote .sub {
    font-size: 18px;
    color: var(--muted);
    letter-spacing: 0.04em;
  }

  /* ─── BIG NUMBER ───────────────────────────────────── */
  section.big {
    background: var(--bg);
    padding: 80px 110px;
  }
  section.big .huge {
    font-size: 220px;
    font-weight: 900;
    color: var(--brand);
    line-height: 0.9;
    letter-spacing: -0.06em;
    margin: 30px 0 10px 0;
  }
  section.big .huge .unit {
    color: var(--ink);
    font-size: 84px;
    font-weight: 700;
    margin-left: 12px;
  }
  section.big .note {
    font-size: 22px;
    color: var(--muted);
    max-width: 70%;
    line-height: 1.5;
    margin-top: 18px;
  }

  /* Helpers */
  .columns-2 { display: grid; grid-template-columns: 1fr 1fr; gap: 50px; margin-top: 8px; }
  .columns-3 { display: grid; grid-template-columns: 1fr 1fr 1fr; gap: 36px; margin-top: 8px; }
  .lede { font-size: 28px; line-height: 1.4; max-width: 80%; color: var(--ink); font-weight: 400; }
  .lede em { color: var(--brand); font-weight: 600; }
  .small { font-size: 16px; color: var(--muted); margin-top: 14px; }
  .accent { color: var(--brand); font-weight: 700; }
---

<!-- _class: cover -->
<!-- _paginate: false -->

<div class="stripe"></div>

# Qonaqzhai

## A marketplace for booking Kazakh weddings, corporate events, and traditional ceremonies — built with Go, Flutter, and an AI planner that actually plans.

<div class="meta">

<strong>Bahaidahar</strong> · Astana IT University · May 2026

</div>

---

### Where this started

# I tried to plan a wedding.

It took three months, six WhatsApp chats, and one deposit I almost didn't get back.

That's why this exists.

---

<!-- _class: quote -->

<div class="stripe"></div>

> In Kazakhstan, planning an event still means scrolling Instagram, asking your aunt, and trusting a stranger with a card number.

<div class="sub">— Roughly every couple I asked.</div>

---

# What I actually built

<div class="columns-2">

<div>

A platform that lets a customer **describe an event in plain Russian or Kazakh** and walk away with a booked, paid plan.

Three clients sharing one Go backend. A Claude-powered planner that calls real APIs instead of hallucinating prices.

</div>

<div>

- **Web** — Next.js 16 for desktop browsing
- **Mobile** — Flutter app, KZ / RU / EN
- **Backend** — five Go microservices behind one gateway
- **AI** — Claude Sonnet 4.6 with tool use, not just chat

</div>

</div>

---

<!-- _class: big -->

<div style="font-size:14px;color:var(--muted);letter-spacing:0.2em;text-transform:uppercase;font-weight:700">Backend, end to end</div>

<div class="huge">2,033<span class="unit">LOC of tests</span></div>

<div class="note">Eight test files across the four services I actually own. Table-driven, <code>-race</code> enabled, run on every push. The number isn't impressive on its own — it's that I wrote them at the same time as the code, not after.</div>

---

# The architecture, in one breath

```
client (mobile + web)
    │ HTTP
    ▼
 gateway :8080  ─── verifies JWT once, forwards everything else
    │
    ├──► auth      users, JWT, password reset
    ├──► core      vendors, bookings, reviews, photos, FCM
    ├──► payment   cards, PayBox saga
    └──► realtime  WS hub, threads, messages

each owns its own Postgres DB
no shared schema, no cross-DB joins
```

When I split this out from a monolith in phase 9, I didn't believe it would feel cleaner. It does. Each service has one job and one schema.

---

# Why I split it

<div class="columns-2">

<div>

I could have shipped this as a Go monolith. Honestly, for diploma traffic, I should have.

I split it anyway because I wanted to learn the *real* operational pain of microservices — gRPC contracts that drift, distributed transactions, cold-start chains.

</div>

<div>

Things I learned the hard way:

- A workspace `replace` directive can break CI silently
- "One DB per service" sounds clean until you need a user's name in three places
- gRPC is great until you try to debug it from a browser

</div>

</div>

---

# The payment saga

I almost wrote distributed transactions. I'm glad I didn't.

```
core.booking.Pay
    │
    ▼
payment.Charge    (gRPC, sync)
    │
    ▼  on success
core.MarkBookingPaid    (gRPC, sync)
```

Synchronous. Two-step. No event bus, no outbox table.

If the second call fails, the payment row is still captured and a follow-up job reconciles. It's not the textbook answer. It is the answer that ships.

---

# Mobile, where most users live

```
mobile/lib/features/
  auth/           vendor_catalog/    booking/
  ai_chat/        payment/           cards/
  reviews/        messaging/         notifications/
  vendor_self/    admin/             onboarding/
  settings/
```

Thirteen feature folders. Each one is `data → domain → presentation`, Riverpod for state, Dio with a 401-refresh interceptor for HTTP, `go_router` for deep links.

I picked Flutter over React Native because I wanted **one Dart codebase shipping native binaries** — not a JS bridge. Six months in, no regrets, but I still miss hot-reload between Dart and the Go server.

---

# The AI part — the slide I'm asked about most

It's not a chatbot wrapping Claude. It's Claude with **four tools** it can actually call:

```
search_vendors(category, city, budget, date)
check_availability(vendor_id, date)
draft_booking(vendor_id, service_id, guests, date)
estimate_total(items[])
```

When a user says *"corporate event, 50 people, Astana, August 20, 1.5M ₸"* —
Claude calls `search_vendors`, picks the top three by rating, calls `check_availability` for each, drafts a plan, and streams it back. No vendor list invented from thin air.

---

# Stack — opinions, not features

| Layer | What I shipped | Why, in one sentence |
|---|---|---|
| Backend | Go 1.24 + gRPC | Fast, boring, statically linked — runs anywhere |
| Web | Next.js 16 + React 19 | Cache Components finally make App Router feel finished |
| Mobile | Flutter 3.24 | One codebase, native binaries, real animations |
| State | Riverpod | Compile-safe DI without `BuildContext` headaches |
| DB | Postgres per service | Forces real service boundaries, no shared-schema laziness |
| AI | Claude API + MCP | Tool use is the part that turns chat into a product |

Every row is a choice I'd have to defend, so I'm defending them.

---

# Where it sits in the market

| Platform | Where | What it does | What it misses |
|---|---|---|---|
| The Knot · WeddingWire | US | Big directories, no booking | Kazakh, KZT, traditions |
| GigSalad · The Bash | US | Vendor commission | No AI, no escrow |
| Peerspace | US | Venue-only, escrow | Not events, not KZ |
| Wedding.ru · FlyBride | RU | Paid listings | No transaction flow |
| Eventie.kz · weddingsalon.kz | KZ | Directories | Web-only, no AI, no payment |

So I'm not building the tenth event marketplace. I'm building the first one that speaks Kazakh, takes KZT, and actually books the event.

---

# The hard stuff I'd do differently

A short honest list, written after the bugs.

- **Sync payment saga.** Works at this scale; will break at 10× scale. Should be event-driven.
- **Hand-rolled DTOs.** I should have generated Dart and TypeScript clients from OpenAPI from day one. Three classes of bugs would have disappeared.
- **Photos in Postgres as `BYTEA`.** Easy demo, wrong for production. Belongs in object storage.
- **No feature flags.** Every change ships to everyone. Fine for diploma; not fine for real users.
- **Test coverage is uneven.** Backend is solid, mobile is sparse. Mobile needs proper ViewModel tests.

---

<!-- _class: quote -->

<div class="stripe"></div>

> The diploma isn't the code. The diploma is the decisions I can explain.

<div class="sub">— What I keep telling myself when I'm tempted to add more features.</div>

---

# What I'd build next

Things I'd add if this kept going past defense.

- **Vendor calendar sync** with Google Calendar so double-booking becomes impossible.
- **Verified-booking badge** — only customers who actually completed an event can leave reviews. Cuts review-spam to zero.
- **Embedding-based vendor recommendations** — replace today's weighted-score SQL with a real similarity model.
- **Voice input for the AI planner**, because typing on mobile while wedding-planning is brutal.
- **A proper escrow flow** with milestone release. Customers stop being scared.

---

# Demo

I'll show you the parts that matter, in this order.

1. Switch the mobile app between Kazakh, Russian, English.
2. Ask the AI to plan an event. Watch it actually call the API.
3. Book one of the vendors it suggested. Pay with a test card.
4. Catch the push notification, deep-link into the booking detail.
5. Log in as the vendor. The booking is there, the dashboard updated.
6. Show the `docker compose logs` of the saga in action.

Eight minutes if I don't ramble. Ten if I do.

---

<!-- _class: dark -->

# What I want you to walk away with

<div class="lede">

This wasn't a tutorial project. I made decisions I'd defend in a real codebase — about service boundaries, about which tool the AI actually needs, about what to ship synchronously and what to leave to a follow-up.

<br>

If a single demo is enough to convince you of one thing, let it be that **the AI planner isn't a wrapper**. It calls real code.

</div>

---

<!-- _class: cover -->
<!-- _paginate: false -->

<div class="stripe"></div>

# Questions?

## I'd rather you ask the awkward ones now than discover them in production.

<div class="meta">

<strong>github.com/Bahaidahar/diploma</strong> · katiev2802@icloud.com

</div>
