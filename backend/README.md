# Qonaqzhai backend — microservices

Five Go services, four PostgreSQL databases, gRPC between services, HTTP
to the public.

```
┌──────────┐
│  client  │  mobile + web
└────┬─────┘
     │ HTTP
┌────▼──────────────────────────────────────────────────┐
│                       gateway :8080                   │
│   verifies JWT once (auth gRPC), routes by prefix,    │
│   forwards X-User-* headers to backends               │
└──┬──────────┬───────────────┬──────────────┬─────────┘
   │ HTTP     │ HTTP          │ HTTP         │ HTTP
   ▼          ▼               ▼              ▼
┌────────┐ ┌──────────┐ ┌────────────┐ ┌──────────────┐
│ auth   │ │  core    │ │  payment   │ │   realtime   │
│ :8081  │ │  :8082   │ │  :8083     │ │   :8084      │
│ :9081  │ │  :9082   │ │  :9083     │ │   :9084 grpc │
└───┬────┘ └──┬───────┘ └──┬─────────┘ └──┬───────────┘
    │         │           │              │
┌───▼───┐ ┌──▼─────┐  ┌───▼──────┐ ┌─────▼───────┐
│auth-db│ │core-db │  │payment-db│ │realtime-db  │
└───────┘ └────────┘  └──────────┘ └─────────────┘
```

gRPC edges between services:

- core → auth (`GetUser`, `GetUsersBatch`, verify)
- core → payment (`Charge`) — synchronous saga on booking pay
- core → realtime (`EnsureThread`, `PublishEvent`)
- payment → core (`MarkBookingPaid`) — saga callback
- realtime → auth (`GetUsersBatch`) — peer-name enrichment
- gateway → auth (verify) — once per inbound HTTP request

## Run

```bash
# All five services, four DBs, via compose:
cp deploy/.env.example deploy/.env && $EDITOR deploy/.env
docker compose -f deploy/docker-compose.yml --env-file deploy/.env up --build
```

Or run a single service against your own Postgres:

```bash
cd backend
make build       # builds every service binary into its own dir
make test        # runs go test ./... in every service module
go run ./services/auth/cmd/auth
```

## Module layout

```
backend/
├── proto/                  source .proto files
├── gen/proto/              generated Go (own module)
├── pkg/                    shared module (auth, errs, httpx, grpcutil, …)
├── services/
│   ├── auth/               users, JWT, password reset
│   ├── core/               vendors, bookings, reviews, photos, notifications
│   ├── payment/            cards, payments, PayBox
│   ├── realtime/           DM threads + WS hub
│   └── gateway/            edge reverse proxy
├── tests/e2e/              Docker-backed end-to-end (build tag: e2e)
├── go.work                 ties every module together for dev
├── Makefile
└── HANDOFF.md              historical scope doc
```

Each service is a separate Go module (own `go.mod`). They share `pkg/`
and `gen/proto/` via the workspace + `replace` directives so each
service can also be built in isolation.

## Configuration

Per-service env vars are documented in each service's `cmd/<name>/main.go`
plus `deploy/.env.example`. The non-negotiables in production:

| Var                | Service | Notes                                        |
|--------------------|---------|----------------------------------------------|
| `JWT_SECRET`       | auth    | Must persist across restarts                 |
| `ADMIN_EMAIL/PWD`  | auth    | Idempotent admin seed; was hardcoded before  |
| `AUTH_GRPC_ADDR`   | all     | Pointer to auth-svc gRPC port                |
| `*_DATABASE_URL`   | service | One Postgres per service, no cross-DB joins  |
| `PAYBOX_MERCHANT_ID/SECRET_KEY` | payment | Falls back to Mock gateway when unset |
| `SMTP_*`           | auth    | Password reset emails; skipped when unset    |

## Testing

```bash
cd backend
make test                           # per-service unit tests
go test -tags=e2e ./tests/e2e -v    # docker-backed E2E
```

## Architectural decisions

1. **One DB per service.** No `FK REFERENCES users(id)` across service
   boundaries — user ids are plain UUIDs. Cross-service joins are
   batched gRPC calls (`auth.GetUsersBatch`).
2. **gRPC only for service-to-service.** Public is HTTP/JSON via the
   gateway. JWT verification happens at the edge plus inside each
   backend service for defense-in-depth.
3. **Synchronous payment saga.** `core.booking.Pay` → `payment.Charge`
   → `core.MarkBookingPaid`. If the callback fails the payment row is
   still captured; reconciliation is a follow-up.
4. **Realtime owns chat threads.** Core triggers `EnsureThread` on
   accept; thread + messages live in the realtime DB. No FK back to
   bookings — booking_id is a unique constraint within realtime only.
5. **No distributed transactions.** Best-effort eventual consistency.
   This is appropriate for the diploma's scale; revisit when GMV
   reaches the millions.

## Where things used to live

`HANDOFF.md` documents the original split plan. The monolith deleted
in phase 9 lived under `cmd/qonaqzhai/`, `internal/`, and the old
`services/{auth,core,realtime}-svc` shells.
