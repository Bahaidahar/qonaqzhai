# Microservices Split — Handoff

Status as of this commit: **phases 1–3 complete; auth-svc fully operational.**
Phases 4–10 are scoped and ready to execute in follow-up sessions. This
document is intentionally self-contained so a fresh session can resume cold.

## Done

### Phase 1 — proto contracts (`backend/proto/`)
- `auth.v1`: VerifyToken, GetUser, GetUsersBatch, AdminListUsers, AdminSetStatus
- `core.v1`: GetVendor, GetVendorByUser, ListVendorsByIDs, GetBooking,
  IsBookingAccepted, MarkBookingPaid, AdminStats
- `payment.v1`: Charge, Refund, GetPayment, ListCardsByUser
- `realtime.v1`: EnsureThread, PublishEvent
- Generated Go in `backend/gen/proto/` (own go.mod).
- `make proto` regenerates everything.

### Phase 2 — shared module (`backend/pkg/`)
- `pkg/errs` — sentinel errors (ErrNotFound, ErrUnauthorized, …).
- `pkg/httpx` — JSON helpers, CORS, recover, access log, rate limiter,
  client-IP extraction, error-to-status mapper.
- `pkg/grpcutil` — errs↔gRPC code conversion + logging/recover interceptors.
- `pkg/logger` — slog JSON setup tagged per service.
- `pkg/config` — env helpers (EnvOr, MustEnv, DurationEnv, BoolEnv, RandomHex).
- `pkg/auth` — Claims, gRPC Verifier client (dial auth-svc), HTTP middleware
  (Required/Optional/RequireRole), HS256 JWTSigner.
- `backend/go.work` lists all modules; local imports resolved via `replace`
  directives in each `go.mod` so each service builds standalone too.

### Phase 3 — auth-svc (`backend/services/auth/`)
- Own go.mod, own database (`AUTH_DATABASE_URL`).
- Owns: `users`, `refresh_tokens`, `password_reset_tokens`. No other service
  touches this database.
- HTTP `/api/signup|login|refresh|logout|forgot-password|reset-password|me|health`.
- gRPC: implements every AuthService method from `auth.v1`.
- bcrypt hasher, SMTP mailer (nil-safe), UUIDv4 IDs, system clock.
- Admin seeding through `ADMIN_EMAIL` + `ADMIN_PASSWORD` env vars
  (idempotent — does not silently overwrite). Replaces the previous
  hardcoded `"admin12345"` in `internal/app/app.go`.
- Unit tests for signup/login/refresh, suspended user rejection, full
  password-reset flow, single-use reset tokens, admin seed idempotency.
- Run: `cd services/auth && go run ./cmd/auth`.

## Outstanding (in priority order)

### Phase 4 — core-svc (largest piece)
Owns: vendors, services-menu, photos, bookings, reviews, notifications,
fcm_tokens, audit_log. No FK to `users` — `customer_id` / vendor's `user_id`
are plain UUIDs sourced from auth-svc.

Skeleton already exists at `backend/services/core/` (just `go.mod` +
placeholder package). Build it out following the auth-svc template:

```
services/core/
├── cmd/core/main.go
├── internal/
│   ├── domain/         # vendor.go, booking.go, service.go, photo.go,
│   │                   #   review.go, notification.go, audit.go
│   ├── ports/ports.go  # interfaces for repos + auth-client + payment-client +
│   │                   #   realtime-client + clock/idgen
│   ├── usecase/
│   │   ├── vendor/     # Submit, Update, Approve/Reject, FindPublic, Search
│   │   ├── booking/    # Create, Accept, Decline, Cancel, Pay, List…
│   │   ├── review/     # Submit (only after BookingCompleted/Paid)
│   │   ├── photo/      # Upload, Serve (5MB max, image/* MIME)
│   │   ├── service/    # vendor services menu CRUD
│   │   ├── search/     # FTS over vendors.search_tsv
│   │   ├── notification/  # in-app + push fanout
│   │   └── admin/      # Stats (fan-out gRPC to auth + payment for totals)
│   └── adapter/
│       ├── repo/       # postgres.go (Open + Migrate), vendor.go, booking.go,
│       │               #   service.go, photo.go, review.go, notification.go,
│       │               #   fcm_token.go, audit.go; migrations/0001_init.up.sql
│       ├── http/       # handler.go (or split by domain), router.go
│       ├── grpc/       # server.go implementing core.v1
│       ├── grpcclient/ # auth.go (calls auth.GetUser/GetUsersBatch),
│       │               #   payment.go (Charge on booking.Pay),
│       │               #   realtime.go (EnsureThread on booking.Accept,
│       │               #   PublishEvent for notification fanout)
│       ├── push/fcm.go # port: services/realtime + this share a notifier?
│       │               # Decide: fcm_tokens are core's because notifications
│       │               # for bookings/reviews are emitted here.
│       ├── ai/gemini.go  # only if you keep AI inside core; the
│       │               # diploma's AI chat is being moved to realtime-svc.
│       ├── clock/clock.go
│       └── idgen/uuid.go
└── go.mod (already exists)
```

Migration plan:
- `0001_init.up.sql`: vendors (with `search_tsv` generated column + GIN
  index), services, photos (bytea), bookings, reviews, notifications,
  fcm_tokens, audit_log. NO `users` table, no `REFERENCES users(id)`.
- Foreign keys WITHIN core (booking.vendor_id → vendors, review.booking_id →
  bookings) are fine and should be preserved.

Cross-service flows:
- `core.usecase.booking.Accept` → calls `realtime.EnsureThread` (gRPC) so
  the chat thread exists when the customer opens the booking.
- `core.usecase.booking.Pay` → calls `payment.Charge`; on success persists
  `payment_id` and flips status to `paid`.
- `admin.Stats` → fan-out: `auth.AdminListUsers(limit=1)` for users count
  (extend proto with `Count` if you want only totals), `payment.*` for
  GMV, plus local counts on vendors/bookings.
- `core.adapter.http` middleware uses `pkg/auth.NewMiddleware` with a
  `pkg/auth.NewVerifier(AUTH_GRPC_ADDR)` so every request validates
  against auth-svc.

Defaults: HTTP `:8082`, gRPC `:9082`, env `CORE_DATABASE_URL` /
`CORE_HTTP_ADDR` / `CORE_GRPC_ADDR` / `AUTH_GRPC_ADDR` /
`PAYMENT_GRPC_ADDR` / `REALTIME_GRPC_ADDR`.

### Phase 5 — payment-svc
Owns: `cards`, `payments`. Wraps PayBox provider (logic already in
`internal/adapter/pay/paybox.go` — copy verbatim, swap imports).

Migrations `0001_init.up.sql`:
```sql
CREATE TABLE cards (
  id TEXT PRIMARY KEY, user_id TEXT NOT NULL, brand TEXT NOT NULL,
  last4 TEXT NOT NULL, exp_month INT NOT NULL, exp_year INT NOT NULL,
  holder TEXT NOT NULL DEFAULT '', is_default BOOL NOT NULL DEFAULT FALSE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_cards_user ON cards(user_id);

CREATE TABLE payments (
  id TEXT PRIMARY KEY, booking_id TEXT NOT NULL UNIQUE, user_id TEXT NOT NULL,
  card_id TEXT NOT NULL, amount BIGINT NOT NULL, currency TEXT NOT NULL DEFAULT 'KZT',
  status TEXT NOT NULL CHECK (status IN ('pending','captured','failed','refunded')),
  provider_ref TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_payments_user ON payments(user_id);
CREATE INDEX idx_payments_booking ON payments(booking_id);
```

HTTP `/api/cards`, `/api/payments` for client. gRPC implements `payment.v1`.

Saga: core calls `payment.Charge` synchronously when a booking moves to
`paid`. On gRPC failure, core does NOT flip the booking. No two-phase
commit — accept best-effort: if charge succeeds but the follow-up
`MarkBookingPaid` returns/UpdateStatus fails, the next refresh of the
booking shows it as paid via payments lookup. Add a reconciliation cron
later if needed.

Defaults: HTTP `:8083`, gRPC `:9083`.

### Phase 6 — realtime-svc
Owns: `threads`, `thread_messages`, `chats` (AI), `chat_messages` (AI).
WebSocket hub for client push (existing `internal/adapter/ws/hub.go` ports
cleanly).

Migration creates the four tables (no FK on customer_id/vendor_id —
those live in auth/core).

HTTP: `/api/threads`, `/api/threads/{id}/messages`, `/api/chats`,
`/api/chats/{id}`, `/api/ws`. gRPC implements `realtime.v1.EnsureThread`
(idempotent — uses unique index on booking_id) and `PublishEvent`
(no-DB write — pure hub broadcast).

Cross-service:
- On `EnsureThread`, store the booking_id + customer_id + vendor_id passed
  by core. No need to call back.
- For per-thread enrichment (counterpart name, vendor name, booking
  status), the realtime HTTP `/api/threads/summaries` endpoint batches
  `auth.GetUsersBatch` and `core.ListVendorsByIDs` calls.

If you keep AI chat here, also move `internal/adapter/ai/gemini.go` and
`internal/usecase/chat/service.go`. The AI chat doesn't need any
cross-service call — it's user-scoped only.

Defaults: HTTP `:8084`, gRPC `:9084`.

### Phase 7 — gateway
Pure reverse proxy in front of the four services. Routes:

```
/api/auth/*        → auth-svc HTTP (rewrite /api/auth/login → /api/login)
/api/payments/*    → payment-svc
/api/cards/*       → payment-svc
/api/threads/*     → realtime-svc
/api/chats/*       → realtime-svc
/api/ws            → realtime-svc (HTTP upgrade-aware proxy)
/api/*             → core-svc (catch-all for vendors/bookings/reviews/etc.)
```

JWT verification once at the edge via `pkg/auth.Middleware` so backends
can trust the `X-User-Id`, `X-User-Role` headers it forwards. Rate limit
per-IP with stricter caps on auth routes. CORS lives here too — backends
respond without any CORS headers.

Defaults: HTTP `:8080`, no DB. Env: `AUTH_HTTP_URL`, `CORE_HTTP_URL`,
`PAYMENT_HTTP_URL`, `REALTIME_HTTP_URL`, `AUTH_GRPC_ADDR`.

### Phase 8 — e2e suite (`backend/tests/e2e/`)
testcontainers-compose spinning up 4 Postgres + 5 services. Critical flows
to cover end-to-end:
- signup → login → refresh
- create vendor profile → admin approves
- customer creates booking → vendor accepts → thread auto-created →
  customer + vendor exchange messages (via WS)
- customer adds card → pays booking → status flips to paid → leaves review
- admin lists vendors/bookings, sees stats

Existing tests under `backend/tests/` use a single-process harness — they
mostly need import rewrites (point at the gateway URL) and the helpers
should boot via compose instead of in-process `app.New()`.

### Phase 9 — delete monolith
Remove (only after phases 4–8 are green):
- `backend/cmd/qonaqzhai/`
- `backend/internal/` (everything)
- `backend/services/auth-svc/`, `services/core-svc/`, `services/realtime-svc/`,
  `services/gateway/` (the *old* monolithic copies)
- `backend/go.mod` (root) — replace with a one-line file that just declares
  `module qonaqzhai-backend` if anything outside services still depends on
  it; otherwise delete entirely (workspace doesn't need a root module).
- Old generated artifacts already replaced in Phase 1.
- The `qonaqzhai.db*` SQLite leftovers.

After deletion, run `go work sync && make build && make test`.

### Phase 10 — deploy + README
- `deploy/docker-compose.yml`: 4 Postgres services (auth_db, core_db,
  payment_db, realtime_db), 5 service containers, depends_on + healthchecks.
- Per-service Dockerfile (multi-stage with distroless base, copying only
  its own `services/<name>/` plus `pkg/` + `gen/proto/`).
- Update `backend/README.md` and root `README.md` with the architecture
  diagram + per-service envs.
- Mobile + frontend point at `gateway:8080` only — they should not need
  any other changes.

## Cross-cutting decisions to revisit

- **Notifications/FCM:** decided to live in core (booking events drive
  most notifications). If realtime-svc also wants to send pushes, it
  should call core via gRPC instead of duplicating the FCM client.
- **Audit log:** keep in core for now. If admin/security want a unified
  view across services, extract to its own audit-svc later.
- **AI chat:** moving to realtime-svc since both AI chat and DM chat
  share the same channel pattern. Reconsider if AI usage grows enough to
  warrant its own service.
- **DM thread ownership:** thread metadata lives in realtime (decided in
  scoping). Core triggers creation via `realtime.EnsureThread` on
  booking accept; no booking → thread foreign key.
- **PCI scope:** PayBox handles tokenisation; payment-svc stores only
  brand + last4 + token references. Real PCI assessment is out of scope
  for the diploma but the structure permits it.
- **Distributed tracing:** add OpenTelemetry interceptors to
  `pkg/grpcutil` + `pkg/httpx` in a separate increment.

## How to resume

`git checkout main && cd backend`. Read this file. Pick a phase. The
auth-svc layout is the template for everything else — copy its structure
when in doubt.

Open questions to settle before phase 4:
1. Keep AI chat in core or move to realtime? (recommend: realtime)
2. Should the gateway pre-verify JWT and forward headers, or let each
   service verify? (recommend: gateway verifies, services trust headers
   except where they re-check role)
3. testcontainers compose vs. ad-hoc compose file for e2e? (recommend:
   ad-hoc compose so dev can `docker compose up` and the same file
   drives CI tests via testcontainers)
