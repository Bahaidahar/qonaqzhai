# Qonaqzhai — Diploma Documentation Plan (Revised)

**Old plan** = microservices, Kafka, ELK, restaurants + rentals split, multi-DB. Overscoped, doesn't match implementation.
**Revised plan** = matches actual stack (Go monolith Clean Arch + Next.js FSD + Flutter mobile + SQLite/Postgres + Gemini AI + Firebase FCM + Gmail SMTP + Freedom Pay (PayBox)).

---

## 1. Use Case List (revised)

**Customer (User):**
- Register / Login (with refresh token)
- Reset password
- Create event request
- Use AI planner (Gemini chat)
- Browse vendors (search + filter by category, city, price, rating)
- Book vendor service
- Pay online (Freedom Pay / PayBox test mode — Kaspi QR, Halyk, cards)
- Leave review + rating (after completed booking)
- View bookings inbox / status
- Receive notifications (email + push)

**Vendor (Business partner):**
- Register / Login as vendor
- Manage vendor profile (name, category, city, price, description, photos)
- Receive booking requests
- Accept / Decline / Mark completed bookings
- View own reviews
- Receive notifications

**Admin:**
- Moderate vendors (approve / reject / suspend)
- Manage users (suspend / unsuspend)
- View platform stats + analytics (KPI tiles, time series, top categories)
- Moderate reviews (delete abuse)

**Removed from old plan:**
- "Rental items" / "Restaurant" split → unified `vendor` entity with `category` field (Venue / Catering / Photo / Decor / Music / etc.). Cleaner data model.
- "Chat with support" → out of scope (AI chat covers planning, not support).
- "Process reports" admin action → covered by review moderation.

---

## 2. Component Diagram (revised, matches reality)

```
┌─────────────────────────────────────────────────────────────┐
│  Client Layer                                               │
│  ┌──────────────┐    ┌──────────────┐                       │
│  │ Next.js Web  │    │ Flutter App  │                       │
│  │  (FSD)       │    │  (mobile)    │                       │
│  └──────┬───────┘    └──────┬───────┘                       │
└─────────┼────────────────────┼─────────────────────────────-┘
          │ HTTPS              │ HTTPS
          ▼                    ▼
┌─────────────────────────────────────────────────────────────┐
│  Reverse Proxy (nginx + Let's Encrypt)                      │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│  Backend Monolith (Go, Clean Architecture)                  │
│                                                             │
│  ┌─────────────────────────────────────────────────────┐    │
│  │ HTTP Layer  (handlers, middleware, DTOs)            │    │
│  │   ├─ AuthHandler                                    │    │
│  │   ├─ VendorHandler                                  │    │
│  │   ├─ BookingHandler                                 │    │
│  │   ├─ ReviewHandler                                  │    │
│  │   ├─ ChatHandler                                    │    │
│  │   ├─ AdminHandler                                   │    │
│  │   └─ Middleware: JWT, RateLimit, CORS, Logger       │    │
│  ├─────────────────────────────────────────────────────┤    │
│  │ Use Cases   (application logic)                     │    │
│  │   AuthUC · VendorUC · BookingUC · ReviewUC          │    │
│  │   SearchUC · ChatUC · AdminUC · NotificationUC      │    │
│  ├─────────────────────────────────────────────────────┤    │
│  │ Domain      (entities + business rules)             │    │
│  │   User · Vendor · Booking · Review · Notification   │    │
│  ├─────────────────────────────────────────────────────┤    │
│  │ Adapters    (interface implementations)             │    │
│  │   SQLiteRepo · GeminiClient · SMTPMailer ·          │    │
│  │   Firebase FCMClient · PayBoxClient · JWTIssuer            │    │
│  └─────────────────────────────────────────────────────┘    │
└──────┬─────────────┬──────────┬──────────┬──────────┬───────┘
       │             │          │          │          │
       ▼             ▼          ▼          ▼          ▼
   ┌────────┐   ┌─────────┐  ┌──────┐  ┌──────┐  ┌────────┐
   │SQLite/ │   │ Gemini  │  │Gmail │  │Firebase FCM │  │PayBox  │
   │Postgres│   │ AI API  │  │SMTP  │  │ API  │  │Test API│
   └────────┘   └─────────┘  └──────┘  └──────┘  └────────┘
```

**Diff from old:**
- No API Gateway (single backend, nginx is enough)
- No microservices (monolith with Clean Arch layers)
- No Kafka/RabbitMQ (in-process notification queue via Go channels)
- No ELK / Prometheus / Grafana (slog structured logging + optional Uptime Kuma for monitoring demo)
- No Auth Server (JWT issued by backend itself)
- No Redis (caching not needed at MVP scale; vendor rating cached in DB column)

---

## 3. Package Diagram (revised)

**Backend packages (Clean Arch):**
```
cmd.qonaqzhai           → main, composition root
internal.domain         → entities, value objects, domain errors
internal.usecase        → application services + port interfaces
internal.adapter.http   → handlers, middleware, DTOs, router
internal.adapter.repo   → SQLite/Postgres repositories
internal.adapter.ai     → Gemini client
internal.adapter.mail   → SMTP client
internal.adapter.push   → Firebase FCM client
internal.adapter.pay    → PayBox client
internal.infra.db       → connection, migrations
internal.infra.token    → JWT issuer
internal.infra.config   → env loading
internal.infra.logger   → slog setup
internal.app            → DI wiring
```

**Frontend packages (FSD):**
```
src.app                 → Next.js App Router
src.pages               → page compositions
src.widgets             → complex UI blocks
src.features            → user-facing features
src.entities            → domain entities (UI + model)
src.shared.ui           → reusable components
src.shared.api          → HTTP client
src.shared.lib          → utils, hooks
src.shared.config       → env, constants
src.shared.i18n         → translations
```

**Diff from old:** "Monitoring", "Security" not separate packages → cross-cutting (middleware + slog).

---

## 4. Object Diagram (sample instances at runtime)

Example state during a booking flow:

```
user: User
  id      = "u_42"
  email   = "aigerim@example.kz"
  role    = "customer"
  status  = "active"

vendor: Vendor
  id          = "v_07"
  ownerId     = "u_91"
  name        = "Rixos Almaty Ballroom"
  category    = "Venue"
  city        = "Almaty"
  priceFrom   = 1500000
  status      = "approved"
  ratingAvg   = 4.7
  ratingCount = 23

booking: Booking
  id        = "b_1024"
  customerId= user.id
  vendorId  = vendor.id
  date      = 2026-06-12
  status    = "accepted"
  amount    = 1800000
  paymentId = "pi_3PqXyz..."

review: Review
  id         = "r_551"
  bookingId  = booking.id
  customerId = user.id
  vendorId   = vendor.id
  rating     = 5
  text       = "Шикарный зал, всё на уровне"

aiPlan: ChatMessage[]
  [
    { role: "user",      content: "свадьба на 200 человек в Алматы, бюджет 8 млн" },
    { role: "assistant", content: "Рекомендую: Rixos Almaty Ballroom (Venue, ~1.5M)..." }
  ]

notification: Notification
  id      = "n_88"
  userId  = vendor.ownerId
  type    = "booking.created"
  channel = "email+push"
  status  = "sent"
```

**Diff from old:** removed standalone `rentalItems` / `restaurant` objects → both unified as `vendor` with category.

---

## 5. Sequence Diagram — Event Creation & Booking Workflow

```
Customer    Next.js     Backend     Gemini    PayBox   Vendor    SMTP/Firebase FCM
   │           │           │           │         │       │           │
   │ open chat │           │           │         │       │           │
   ├──────────▶│           │           │         │       │           │
   │           │ POST /chat│           │         │       │           │
   │           ├──────────▶│           │         │       │           │
   │           │           │ prompt    │         │       │           │
   │           │           ├──────────▶│         │       │           │
   │           │           │◀──────────┤         │       │           │
   │           │           │ inject vendors      │       │           │
   │           │◀──────────┤           │         │       │           │
   │           │ AI plan   │           │         │       │           │
   │◀──────────┤           │           │         │       │           │
   │           │           │           │         │       │           │
   │ search    │           │           │         │       │           │
   ├──────────▶│ GET /vendors?filters  │         │       │           │
   │           ├──────────▶│           │         │       │           │
   │           │◀──────────┤           │         │       │           │
   │ pick venue│           │           │         │       │           │
   │ POST /bookings        │           │         │       │           │
   ├──────────▶├──────────▶│           │         │       │           │
   │           │           │ INSERT booking(pending)     │           │
   │           │           │           │         │       │           │
   │           │           │ notify vendor       │       │           │
   │           │           ├─────────────────────────────────────────▶│
   │           │◀──────────┤           │         │       │           │
   │ POST /bookings/{id}/pay           │         │       │           │
   ├──────────▶├──────────▶│ create Checkout Session    │           │
   │           │           ├────────────────────▶│       │           │
   │           │           │◀────────────────────┤       │           │
   │           │◀──────────┤ session URL         │       │           │
   │◀──────────┤           │           │         │       │           │
   │ redirect to Freedom Pay checkout  │         │       │           │
   ├──────────────────────────────────────────────▶      │           │
   │ pay              │           │         │       │           │
   │◀────────────────────────────────────────────┤       │           │
   │           │           │ webhook payment.succeeded   │           │
   │           │           │◀────────────────────┤       │           │
   │           │           │ UPDATE booking(paid)│       │           │
   │           │           │ notify customer + vendor    │           │
   │           │           ├─────────────────────────────────────────▶│
   │           │           │           │         │       │           │
   │           │           │ Vendor accepts via /bookings/{id}/accept│
   │           │           │◀──────────────────────────────┤         │
   │           │           │ notify customer             │           │
   │           │           ├─────────────────────────────────────────▶│
```

---

## 6. Activity Diagram — User Plans Event

```
       (start)
          │
          ▼
   ┌──────────────┐
   │ Login/Signup │
   └──────┬───────┘
          ▼
   ┌──────────────────┐
   │ Open AI planner  │
   └──────┬───────────┘
          ▼
   ┌──────────────────────────────────┐
   │ Describe event (date, guests,    │
   │ budget, style, city)             │
   └──────┬───────────────────────────┘
          ▼
   ┌─────────────────────────────────┐
   │ AI returns vendor suggestions   │
   │ (Gemini + real vendor inject)   │
   └──────┬──────────────────────────┘
          ▼
   ┌────────────────────────┐
   │ Refine via filters     │
   │ (category/price/city)  │
   └──────┬─────────────────┘
          ▼
   ┌────────────────────┐
   │ Open vendor page   │
   │ + read reviews     │
   └──────┬─────────────┘
          ▼
   ┌─────────────────┐
   │ Create booking  │
   └──────┬──────────┘
          ▼
       ◇ payment? ◇
       │         │
      yes        no
       ▼         ▼
  ┌──────┐  ┌─────────────────┐
  │PayBox│  │ pending payment │
  └──┬───┘  └─────────────────┘
     ▼
 ┌──────────────────┐
 │ Vendor accepts   │
 │ /declines        │
 └──────┬───────────┘
        ▼
 ┌──────────────────────┐
 │ Notify customer      │
 │ (email + push)       │
 └──────┬───────────────┘
        ▼
 ┌────────────────────┐
 │ Event takes place  │
 │ → booking completed│
 └──────┬─────────────┘
        ▼
 ┌──────────────┐
 │ Leave review │
 └──────┬───────┘
        ▼
      (end)
```

---

## 7. High-Level Architecture (revised, no overpromised infra)

```
┌──────────────────────────────────────────────────────────┐
│ Client Layer                                             │
│   Next.js 16 (web, FSD)   ·   Flutter (mobile)           │
└─────────────────────┬────────────────────────────────────┘
                      │ HTTPS
┌─────────────────────▼────────────────────────────────────┐
│ Edge Layer                                               │
│   nginx (reverse proxy + TLS via Let's Encrypt)          │
└─────────────────────┬────────────────────────────────────┘
                      │
┌─────────────────────▼────────────────────────────────────┐
│ Application Layer (Go monolith, Clean Architecture)      │
│   HTTP handlers → Use cases → Domain → Adapters          │
│   In-process notification worker (goroutine)             │
└──┬──────────┬───────────┬─────────┬──────────┬───────────┘
   │          │           │         │          │
┌──▼─────┐ ┌──▼──────┐ ┌──▼────┐ ┌──▼─────┐ ┌──▼─────────┐
│ DB     │ │ Gemini  │ │ Gmail │ │ Firebase FCM  │ │ PayBox     │
│ SQLite/│ │ AI      │ │ SMTP  │ │ Push   │ │ (test mode)│
│ Postgres│ │         │ │       │ │        │ │            │
└────────┘ └─────────┘ └───────┘ └────────┘ └────────────┘

┌──────────────────────────────────────────────────────────┐
│ Cross-cutting                                            │
│   slog structured logs · uptime-kuma (optional) ·        │
│   GitHub Actions CI/CD · Docker Compose deploy           │
└──────────────────────────────────────────────────────────┘
```

**Removed from old plan:**
- Kafka / RabbitMQ → in-process channel queue (sufficient for MVP)
- ELK stack → slog + log files (justify by scale)
- Prometheus/Grafana → optional, replace with simple `/metrics` JSON endpoint + Uptime Kuma if needed
- Separate Auth Server → JWT signed by backend (monolith)
- Redis → not needed; if added, only for rate-limit token bucket persistence

**Justification line for записка:**
> Microservices, Kafka, ELK stack were considered but rejected at MVP scale. Distributed infrastructure introduces operational overhead disproportionate to a 2-region pilot with < 10k MAU. Clean Architecture monolith provides equivalent modularity at the code level with O(1) deployment complexity. Future migration path preserved: each use case package is a microservice candidate.

---

## 8. Diploma Structure (revised, aligned with actual deliverables)

### Chapter 1 — Theoretical Foundations

1.1 Event services market in Kazakhstan (gap analysis)
1.2 Marketplace platforms: business models, monetization
1.3 AI in event planning: LLM-based recommendation
1.4 Existing solutions analysis (GoSwana, Wezoom, Marry.kz, international: The Bash, WeddingWire) — comparison table with feature gaps
1.5 Methodologies: Clean Architecture (Martin), Feature-Sliced Design

### Chapter 2 — System Analysis and Design

2.1 Functional requirements (customer / vendor / admin)
2.2 Non-functional requirements (security, performance, scalability, observability, i18n)
2.3 Use Case diagram
2.4 ER diagram (User, Vendor, Booking, Review, Notification, RefreshToken, PasswordResetToken)
2.5 Class diagram (domain entities)
2.6 Sequence diagrams (signup+login, booking flow, AI chat flow, payment flow)
2.7 Activity diagram (user plans event)
2.8 Component diagram (Clean Architecture layers)
2.9 Package diagram (backend + frontend)
2.10 High-level deployment architecture

### Chapter 3 — Implementation, Integration, Testing

3.1 Backend implementation (Go, Clean Architecture, layer-by-layer walkthrough)
3.2 Frontend implementation (Next.js 16 + FSD)
3.3 Mobile implementation (Flutter, shared OpenAPI client)
3.4 Database schema + migrations (golang-migrate)
3.5 API design (REST, Swagger/OpenAPI 3 spec)
3.6 External integrations:
  - Gemini AI (prompt engineering, vendor injection)
  - Gmail SMTP (transactional email)
  - Firebase FCM (push notifications)
  - Freedom Pay (PayBox) mode (payments)
3.7 Security: JWT + refresh tokens, password reset flow, rate limiting, bcrypt
3.8 Testing:
  - Unit tests (Go testing, table-driven)
  - Integration tests (repo + HTTP)
  - E2E tests (Playwright)
  - Coverage report (target 80%+)

### Chapter 4 — Deployment, Operations, Impact

4.1 Containerization (multi-stage Dockerfiles)
4.2 Docker Compose orchestration (backend + frontend + nginx)
4.3 Reverse proxy + HTTPS (nginx + certbot)
4.4 CI/CD pipeline (GitHub Actions: lint → test → build → deploy)
4.5 Logging strategy (slog structured, JSON in prod)
4.6 Monitoring (health endpoint + optional Uptime Kuma)
4.7 Economic impact
4.8 Social impact
4.9 Limitations + future improvements
4.10 Research component results (see §10)

---

## 9. Functional & Non-Functional Requirements (cleaned)

**Functional:**
- Auth: register/login/logout, refresh, password reset
- Customer: create event request, AI chat, search/filter vendors, book, pay, review
- Vendor: manage profile, manage bookings, view reviews
- Admin: moderate vendors/reviews, manage users, view analytics
- Notifications: email + push for key events
- i18n: Kazakh / Russian / English

**Non-functional:**
- **Security:** OWASP Top 10 mitigations, bcrypt cost 12, JWT short TTL + refresh rotation, rate limiting, HTTPS, parameterized SQL, input validation
- **Performance:** p95 API latency < 300ms (excluding AI calls), search query < 100ms with indexes
- **Reliability:** stateless backend → restart-safe, DB ACID transactions, idempotent webhooks
- **Scalability:** stateless monolith ready for horizontal scaling behind load balancer; DB extractable to Postgres
- **Observability:** structured logs, /healthz endpoint, request tracing via correlation ID middleware
- **Extensibility:** Clean Architecture boundaries → swap adapters without touching usecases
- **Cross-platform:** web (Next.js) + mobile (Flutter, iOS + Android)
- **Internationalization:** runtime locale switching, KZ/RU/EN

**Cut from old:** "high performance" (vague → quantified above), "cloud deployment" (specified → Docker on VPS).

---

## 10. Research Component (научная новизна) — pick 1

Old plan listed 7 broad areas. Realistic for diploma timeframe:

**Recommended:** *"AI-assisted vendor recommendation: comparing zero-shot prompt vs prompt-injected vendor catalog"*
- Method: 2 prompt strategies × 50 test queries (KZ/RU)
- Metrics: relevance (manual scoring 1-5), latency, token cost
- Result: table + chart proving injected-context prompt > zero-shot
- Why it works: directly tied to actual `ChatUC` code; can ship as graphs in chapter 3

Alternative narrower options if time short:
- Search relevance: FTS5 vs LIKE benchmark on vendor names/descriptions
- i18n bundle size impact in Next.js 16 (runtime vs build-time)

---

## 11. Economic & Social Impact (concrete, defensible)

**Economic:**
- SME digitization: vendors gain online channel without building own site
- Commission model (5-10% per booking) → marketplace revenue projection
- Cost reduction for customers: 1 platform vs calling 10 vendors manually (time saved ~15h per event)
- Local payment processor integration (Freedom Pay / PayBox aggregating Kaspi QR + Halyk + cards) keeps money in-region

**Social:**
- Lowers barrier to organizing cultural events (toi, beshik-toi, kyz uzatu)
- Supports preservation of national traditions via category for traditional services (national music, baursak catering, shanyrak decoration)
- Accessibility: mobile-first → reaches users without desktops
- Trilingual (KZ/RU/EN) → inclusion across language groups
- Transparency: rating system improves vendor quality competition

---

## 12. Future Research Directions (defense Q&A buffer)

- Generative AI for event scenarios (LLM-generated full event plan with timeline)
- Predictive analytics: demand forecasting per category/season
- Vendor fraud detection: anomaly scoring on review patterns
- Explainable AI: surface why each vendor was recommended
- Federated reviews: cross-platform reputation portability
- Behavioral analytics: customer journey funnel optimization

---

## 13. Mapping: Diploma Doc ↔ Code Artifact

| Doc Section | Code Reference |
|---|---|
| Use Case diagram | `internal/adapter/http/router.go` endpoints |
| ER diagram | `migrations/*.sql` |
| Class diagram | `internal/domain/*.go` |
| Sequence: auth | `usecase/auth`, `adapter/http/handler/auth.go` |
| Sequence: booking | `usecase/booking`, `adapter/http/handler/booking.go` |
| Component diagram | `internal/*` tree |
| Package diagram | Go packages + FSD layers |
| API spec | `docs/swagger.json` |
| Tests | `tests/`, `e2e/` |
| Deploy | `deploy/docker-compose.yml` |
| CI/CD | `.github/workflows/*.yml` |

---

## 14. Removed / Out-of-Scope (state plainly in записка)

- Microservices (monolith chosen — justified by scale)
- Kafka/RabbitMQ (in-process queue sufficient)
- ELK / Prometheus / Grafana (slog + Uptime Kuma)
- Redis (not needed at MVP scale)
- Separate Auth Server (JWT in monolith)
- Restaurant / Rental items as separate entities (unified `vendor` with `category`)
- Customer-support chat (AI planner only)

State explicitly that these were considered and rejected → shows architectural maturity at defense.
