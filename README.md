# Qonaqzhai

Event services marketplace for Kazakhstan. Customers plan events with an AI
assistant, browse local vendors (venues, catering, photo/video, decor, music,
traditional services), book and pay online, leave reviews. Vendors manage
bookings. Admins moderate the platform.

## Stack

- **Backend** — Go microservices (auth / core / payment / realtime / gateway),
  PostgreSQL per service, gRPC between services, HTTP at the edge.
  See [`backend/README.md`](backend/README.md).
- **Frontend** — Next.js 16, Feature-Sliced Design, TypeScript, Tailwind.
- **Mobile** — Flutter, Riverpod + MVVM, feature-first Clean Architecture.

## Repo layout

```
backend/   Go microservices + shared pkg + gen/proto + workspace
frontend/  Next.js 16 web client
mobile/    Flutter mobile client
```

## Run

```bash
# Backend (one service at a time during dev — see backend/README.md for full stack)
cd backend
go run ./services/auth/cmd/auth     # :8081 HTTP, :9081 gRPC

# Frontend
cd frontend
pnpm install && pnpm dev            # :3000

# Mobile
cd mobile
flutter pub get
flutter run --dart-define=API_BASE_URL=http://localhost:8080
```

Each backend service needs its own Postgres database; defaults in each
service's `cmd/<name>/main.go` point at `localhost:5433` with predictable
db names (`qonaqzhai_auth`, `qonaqzhai_core`, `qonaqzhai_payment`,
`qonaqzhai_realtime`).

## Tests

```bash
cd backend && make test             # per-service Go tests
cd frontend && pnpm test            # if configured
```
