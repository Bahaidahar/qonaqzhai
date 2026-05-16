# Qonaqzhai — Diploma Implementation Plan

**Project:** Event/Wedding services marketplace (Kazakhstan)
**Stack:** Go backend + Next.js 16 frontend + Flutter mobile (future)
**Target:** Diploma defense (KZ university)

---

## 0. Current State

- Backend: Go + SQLite, layered (`api/app/auth/model/store`), JWT auth, Gemini AI chat, vendor/booking/admin flows, e2e covered.
- Frontend: Next.js 16 + Turbopack, i18n, Playwright e2e, custom structure.
- Roles: customer / vendor / admin.
- Deployed locally only.

---

## 1. Architecture

### 1.1 Backend — Clean Architecture (Uncle Bob)

**Why:** Academic defensibility (citation: Martin, R. "Clean Architecture", 2017). Clear layer boundaries. Easy to draw diagrams for пояснительная записка.

**Target layout:**

```
backend/
├── cmd/
│   └── qonaqzhai/
│       └── main.go              # composition root only
├── internal/
│   ├── domain/                  # entities + business rules (no deps)
│   │   ├── user.go
│   │   ├── vendor.go
│   │   ├── booking.go
│   │   ├── review.go
│   │   └── errors.go
│   ├── usecase/                 # application logic, depends on domain + ports
│   │   ├── auth/
│   │   ├── vendor/
│   │   ├── booking/
│   │   ├── review/
│   │   ├── search/
│   │   └── ports.go             # repository + service interfaces
│   ├── adapter/
│   │   ├── http/                # delivery layer (handlers, routes, DTOs)
│   │   │   ├── handler/
│   │   │   ├── middleware/
│   │   │   ├── dto/
│   │   │   └── router.go
│   │   ├── repository/          # data access impl (sqlite/postgres)
│   │   │   ├── user_repo.go
│   │   │   ├── vendor_repo.go
│   │   │   └── ...
│   │   ├── ai/                  # gemini client impl
│   │   ├── mail/                # SMTP impl
│   │   └── push/                # Firebase Cloud Messaging impl
│   ├── infra/
│   │   ├── db/                  # connection, migrations runner
│   │   ├── config/              # env loading
│   │   ├── logger/
│   │   └── token/               # JWT impl
│   └── app/
│       └── app.go               # DI wiring
├── migrations/                  # sql files (goose or golang-migrate)
├── docs/
│   └── swagger/                 # generated OpenAPI
└── tests/
    ├── integration/
    └── e2e/
```

**Layer rules:**
- `domain` imports nothing from project
- `usecase` imports only `domain`
- `adapter` imports `usecase` + `domain`
- `cmd` + `app` wire everything

**Migration steps from current code:**
1. Extract entities from `internal/model/types.go` → `internal/domain/`
2. Create port interfaces in `internal/usecase/ports.go` (UserRepo, VendorRepo, etc.)
3. Move SQLite code from `internal/store/sqlite.go` → `internal/adapter/repository/*` implementing ports
4. Split `internal/api/*.go` handlers → thin HTTP layer calling usecases
5. Move business logic out of handlers into `usecase/*`
6. Move `internal/ai/gemini.go` → `internal/adapter/ai/`
7. Wire in `internal/app/app.go`

### 1.2 Frontend — Feature-Sliced Design (FSD)

**Why:** Industry-recognized methodology (feature-sliced.design). Maps cleanly onto Next.js App Router. Defensible architecture choice.

**Target layout:**

```
frontend/src/
├── app/                 # Next.js app router (FSD: app layer)
│   ├── (routes)/
│   ├── layout.tsx
│   └── providers.tsx
├── pages/               # FSD pages layer (page compositions)
├── widgets/             # complex UI blocks (Header, VendorCard, ChatPanel)
├── features/            # user-facing features (auth, booking-flow, search-filter, review-submit)
│   └── <feature>/
│       ├── ui/
│       ├── model/       # zustand store / hooks
│       ├── api/         # API calls
│       └── index.ts     # public api
├── entities/            # domain entities (User, Vendor, Booking, Review)
│   └── <entity>/
│       ├── ui/
│       ├── model/
│       └── api/
└── shared/              # reusable, no business logic
    ├── ui/              # button, input, modal
    ├── api/             # axios/fetch client, interceptors
    ├── lib/             # utils, hooks
    ├── config/          # env, constants
    └── i18n/
```

**FSD import rules (enforced via eslint-plugin-boundaries):**
- `shared` → nothing
- `entities` → `shared`
- `features` → `entities`, `shared`
- `widgets` → `features`, `entities`, `shared`
- `pages` → `widgets`, `features`, `entities`, `shared`
- `app` → all

**Migration:** keep Next.js App Router under `src/app/`, move components into `widgets/features/entities/shared` per FSD rules.

---

## 2. Database

**Decision:** Keep SQLite for diploma defense (zero-ops, runs anywhere). Postgres = stretch goal.

**Tasks:**
- Introduce migrations via `golang-migrate` or `goose`
- Add ERD diagram (dbdiagram.io) for пояснительная записка
- Tables to add: `reviews`, `categories`, `vendor_photos` (if not normalized), `notifications`, `password_reset_tokens`, `refresh_tokens`

**Migration files:**
```
migrations/
├── 0001_init.up.sql / 0001_init.down.sql
├── 0002_reviews.up.sql
├── 0003_refresh_tokens.up.sql
├── 0004_password_reset.up.sql
└── 0005_notifications.up.sql
```

---

## 3. Features to Implement

### 3.1 Search & Filters
- Backend: `GET /api/vendors?category=&city=&price_min=&price_max=&rating_min=&q=&sort=&page=&limit=`
- Full-text search on name/description (SQLite FTS5 module)
- Filter: category, city, price range, min rating
- Sort: price asc/desc, rating desc, newest
- Pagination with `total`/`page`/`limit` metadata

**Frontend:**
- `features/vendor-search` + `features/vendor-filter`
- URL-synced filters (Next.js searchParams)
- Debounced query
- Loading skeletons

### 3.2 Reviews & Ratings
- Entity: `Review { id, vendor_id, customer_id, booking_id, rating(1-5), text, created_at }`
- Rule: only customer with completed booking can review (once per booking)
- Backend: `POST /api/vendors/{id}/reviews`, `GET /api/vendors/{id}/reviews`
- Vendor average rating cached on vendor row, recomputed on review insert/delete
- Frontend: `features/review-submit`, `features/review-list`, star widget in `shared/ui`

### 3.3 Notifications
- **Email (Gmail SMTP)** — booking confirmations, password reset, vendor approval
- **Push (Firebase Cloud Messaging)** — booking status changes for Flutter mobile app
- Tasks:
  - `internal/adapter/mail/smtp.go` — wrap `net/smtp`
  - `internal/adapter/push/fcm.go` — Firebase Cloud Messaging client (HTTP v1 API)
  - Notification queue (in-memory channel + worker) — async, no request blocking
  - Templates in `templates/email/*.html`
- DB table `notifications` for in-app inbox

### 3.4 Admin Dashboard & Analytics
- Metrics: total users / vendors / bookings / revenue / DAU / approval funnel
- Endpoints: `GET /api/admin/stats`, `GET /api/admin/stats/timeseries?metric=&from=&to=`
- Frontend: Recharts or Chart.js inside `widgets/admin-dashboard`
- Cards: KPI tiles + line chart (bookings/day) + bar chart (top categories) + funnel

### 3.5 Payments (TBD — stretch)

Stripe **не работает в KZ**. Локальные варианты:

| Провайдер | Sandbox без юр.лица | KZ-native | Заметки |
|---|---|---|---|
| **Freedom Pay (PayBox.money)** | да, test-merchant | да | Агрегатор: Kaspi QR + Halyk + карты. Рекомендую. |
| Kaspi Pay | нет (нужен ИП/ТОО + договор) | да | Долго, не подойдёт для диплома |
| Halyk Bank Epay | нет (договор с банком) | да | Долго |
| CloudPayments KZ | да | частично | Резерв |
| Mock checkout (свой) | да | n/a | Fallback, оформить как stub до подписания PSP |

**Recommendation:** Freedom Pay (PayBox) test mode.
- Endpoint `POST /api/bookings/{id}/pay` → creates PayBox payment intent
- Redirect customer на PayBox checkout
- Callback `POST /api/webhooks/paybox` → verify signature → mark booking paid
- Поддержка Kaspi QR через агрегатор PayBox → демо для защиты

**Fallback path:** свой mock checkout (`/api/payments/mock`) с фейк-успехом, имитирующий PSP-флоу. Защищать как "stub-имплементация для MVP до интеграции с production PSP". Webhook-логика идентична — потом меняется только адаптер.

---

## 4. Security

### 4.1 Rate Limiting
- Per-IP limiter on auth endpoints (login/signup): 10 req / min
- Per-user limiter on chat: 30 req / min
- Lib: `golang.org/x/time/rate` + middleware

### 4.2 Refresh Tokens
- Access token: short TTL (15 min), JWT
- Refresh token: long TTL (30 days), opaque, stored hashed in DB
- Endpoints: `POST /api/auth/refresh`, `POST /api/auth/logout` (revoke refresh)
- Rotation on use (single-use refresh tokens)

### 4.3 Password Reset
- `POST /api/auth/forgot-password { email }` → emit email with one-time token (TTL 1h, hashed in DB)
- `POST /api/auth/reset-password { token, new_password }` → verify + update + invalidate token
- Rate limit on forgot-password endpoint

### 4.4 General Hardening
- bcrypt cost 12 for passwords
- CORS: explicit allowed origins
- HTTPS only in prod (via nginx)
- Secure cookies for refresh token (HttpOnly, Secure, SameSite=Lax)
- Input validation: `go-playground/validator` on all DTOs
- SQL: parameterized queries only (already done)

---

## 5. Deployment (Docker)

```
deploy/
├── docker-compose.yml
├── backend.Dockerfile         # multi-stage: golang:alpine → distroless
├── frontend.Dockerfile        # multi-stage: node:alpine → standalone next
├── nginx/
│   ├── nginx.conf             # reverse proxy + HTTPS
│   └── certbot/               # Let's Encrypt
└── .env.example
```

**Services:**
- `backend` — Go binary on `:8080`
- `frontend` — Next.js standalone on `:3000`
- `nginx` — reverse proxy on `:80/:443`, routes `/api/*` → backend, rest → frontend
- (Postgres later if added)

**VPS deploy target:** any KZ provider (PS.kz, Hoster.kz) or DigitalOcean droplet.
Domain + Let's Encrypt cert via certbot sidecar.

---

## 6. CI/CD (GitHub Actions)

```
.github/workflows/
├── backend.yml          # on push to backend/**
│   ├── go fmt / go vet
│   ├── golangci-lint
│   ├── go test -race -cover
│   └── docker build & push
├── frontend.yml         # on push to frontend/**
│   ├── pnpm lint
│   ├── pnpm typecheck
│   ├── pnpm test
│   ├── playwright e2e
│   └── docker build & push
└── deploy.yml           # on tag v*
    └── ssh to VPS + docker compose pull && up -d
```

**Quality gates:** PR cannot merge if lint/test fails.

---

## 7. API Documentation (Swagger / OpenAPI)

- Tool: `swaggo/swag` (Go annotations → OpenAPI 3) OR hand-written `openapi.yaml`
- Annotations on each handler
- Generate `docs/swagger.json` + `swagger.yaml`
- Serve Swagger UI at `/api/docs` (dev only)
- Export `openapi.yaml` for пояснительная записка appendix

---

## 8. Documentation (for Defense)

Required artifacts for пояснительная записка:

- **Введение** — актуальность (KZ event market gap), цель, задачи
- **Аналог анализ** — comparison table: GoSwana, Wezoom, Marry.kz vs Qonaqzhai
- **ER-диаграмма** — dbdiagram.io export
- **UML:**
  - Use Case (customer / vendor / admin)
  - Class diagram (domain entities)
  - Sequence (booking flow, auth flow, chat flow)
  - Component diagram (C4 level 2)
- **Архитектура:** Clean Architecture diagram + FSD diagram
- **Скриншоты** UI всех ключевых экранов
- **API reference** (Swagger export)
- **Тестирование** — coverage report, e2e scenarios list
- **Экономическая часть** (если требуется кафедрой)
- **Охрана труда** (если требуется)

Tools:
- diagrams: draw.io / Mermaid / PlantUML
- ER: dbdiagram.io

---

## 9. Research Topics (научная новизна)

Pick 1-2 for defense — diplomas love "research component":

1. **Bench AI prompts:** measure Gemini response quality on event-planning prompts (Kazakh/Russian). Latency + relevance scoring. Result: prompt template choice justified by data.
2. **Vendor matching algorithm:** weighted scoring (category match + price fit + rating + locality) — describe formula, compare with naïve filter.
3. **Search relevance:** SQLite FTS5 vs naïve LIKE — benchmark precision/recall on 50 query corpus.
4. **i18n approach:** runtime vs build-time localization in Next.js 16 — measure bundle size impact.

---

## 10. Implementation Phases & Order

### Phase 1 — Architecture refactor (week 1)
- Backend Clean Architecture migration
- Frontend FSD restructure
- Migrations system

### Phase 2 — Features (week 2)
- Reviews + ratings
- Search + filters
- Refresh tokens + password reset
- Rate limiting

### Phase 3 — Notifications + Admin (week 3)
- Gmail SMTP
- Firebase Cloud Messaging integration
- Admin dashboard with charts

### Phase 4 — Payments (week 4, if scoped)
- Freedom Pay / PayBox test integration

### Phase 5 — DevOps (week 4-5)
- Dockerfiles + compose
- nginx + HTTPS
- GitHub Actions CI/CD
- VPS deploy

### Phase 6 — Docs (week 5-6)
- Swagger
- ER + UML diagrams
- Пояснительная записка draft
- Research component writeup

### Phase 7 — Flutter mobile (parallel, post-defense-prep)

**Architecture:** simplified Clean Architecture, feature-first layout, Riverpod + MVVM (no BLoC).

```
mobile/lib/
├── core/
│   ├── network/        # Dio client, interceptors, error mapping
│   ├── di/             # Riverpod providers root
│   ├── router/         # go_router config
│   ├── theme/
│   └── utils/
├── features/
│   ├── auth/
│   │   ├── data/
│   │   │   ├── datasources/auth_remote_datasource.dart
│   │   │   ├── models/auth_dto.dart            # JSON DTOs
│   │   │   └── repositories/auth_repository_impl.dart
│   │   ├── domain/
│   │   │   ├── entities/user.dart              # pure Dart, no JSON
│   │   │   ├── repositories/auth_repository.dart  # abstract
│   │   │   └── usecases/login_usecase.dart
│   │   └── presentation/
│   │       ├── viewmodels/login_viewmodel.dart # AsyncNotifier (Riverpod)
│   │       ├── screens/login_screen.dart
│   │       └── widgets/
│   ├── vendor_catalog/
│   ├── booking/
│   ├── ai_chat/
│   ├── payment/
│   ├── reviews/
│   └── notifications/
└── main.dart
```

**Layer rules:**
- `domain` — pure Dart, no Flutter imports, no JSON, no Dio. Entities + abstract repos + use cases
- `data` — implements `domain` repos, talks to Dio, maps DTO ↔ entity
- `presentation` — ViewModels (Riverpod `AsyncNotifier` / `Notifier`), screens, widgets

**Dependencies:**
- `flutter_riverpod` — state management + DI
- `dio` + `retrofit` (or hand-written client) — HTTP
- `freezed` + `json_serializable` — DTOs and entities
- `go_router` — navigation
- `firebase_messaging` — FCM push
- `flutter_localizations` + `intl` — i18n (KZ / RU / EN)
- `cached_network_image` — vendor photos
- `flutter_secure_storage` — refresh token storage

**MVVM convention:**
- ViewModel = Riverpod `AsyncNotifier<State>` per screen
- View = `ConsumerWidget` reading the provider, dispatching ViewModel methods
- State = immutable freezed class (loading / data / error)

**Codegen client:**
- Generate Dart client from backend `openapi.yaml` via `openapi-generator-cli`
- Wrap generated client in repository implementations to keep `domain` clean

---

## 11. Open Decisions

| Topic | Status | Notes |
|-------|--------|-------|
| Postgres vs SQLite | SQLite default, Postgres = stretch | depends on time |
| Payments | TBD | Freedom Pay / PayBox test mode recommended (Stripe не работает в KZ) |
| Push provider | Firebase Cloud Messaging | free, нативная интеграция с Flutter |
| Research topic | TBD | pick one from §9 |
| VPS provider | TBD | PS.kz / Hoster.kz / DO |

---

## 12. Acceptance Criteria for Defense

- [ ] Clean Architecture backend, layer boundaries enforced
- [ ] FSD frontend, eslint boundaries plugin green
- [ ] Migrations system, ERD diagram
- [ ] Reviews + search + filters working e2e
- [ ] Notifications: email + push
- [ ] Admin dashboard with charts
- [ ] Refresh tokens + password reset + rate limit
- [ ] Swagger UI live
- [ ] Docker compose runs full stack
- [ ] CI green on main
- [ ] Deployed on public URL with HTTPS
- [ ] 80%+ test coverage backend
- [ ] Playwright e2e covers golden paths
- [ ] Пояснительная записка drafted with all diagrams
- [ ] Research component documented
