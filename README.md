# Qonaqzhai

Event services marketplace for Kazakhstan. Customers plan events with an AI assistant, browse local vendors (venues, catering, photo/video, decor, music, traditional services), book and pay online, leave reviews. Vendors manage bookings. Admins moderate the platform.

## Stack

- **Backend** — Go, Clean Architecture, SQLite + golang-migrate, JWT + refresh rotation, Gemini AI, Gmail SMTP, Firebase Cloud Messaging, Freedom Pay / PayBox
- **Frontend** — Next.js 16, Feature-Sliced Design, TypeScript, Tailwind
- **Mobile** — Flutter, Riverpod + MVVM, feature-first Clean Architecture
- **DevOps** — Docker Compose, nginx + Let's Encrypt, GitHub Actions CI/CD
- **Docs** — OpenAPI 3 (`backend/docs/openapi.yaml`), Swagger UI at `/api/docs`

## Repo layout

```
backend/     Go API server (cmd/qonaqzhai)
frontend/    Next.js 16 web client
mobile/      Flutter mobile client
deploy/      Docker compose + nginx + certbot
scripts/     dev helpers (seed-demo.sh)
.github/     CI/CD workflows
docs/        DIPLOMA_*.md, PLAN.md
```

## Quick start

```bash
# 1. backend
cd backend
cp .env.example .env   # fill JWT_SECRET, GEMINI_API_KEY, SMTP_*
go run ./cmd/qonaqzhai

# 2. frontend
cd ../frontend
pnpm install
pnpm dev

# 3. seed demo data (12 vendors, customers, bookings, reviews)
cd ../scripts
RATE_LIMIT_DISABLED=true ./seed-demo.sh
```

Open `http://localhost:3000`. Demo accounts in `ACCOUNTS.md`.

## Mobile

```bash
cd mobile
flutter pub get
flutter run --dart-define=API_BASE_URL=http://localhost:8080
```

## Tests

```bash
cd backend && go test -race ./...
cd frontend && pnpm build
```

## API docs

Backend serves Swagger UI at `http://localhost:8080/api/docs` and raw spec at `/api/docs/openapi.yaml`.

## Diploma

See `DIPLOMA_BRIEF.md` / `DIPLOMA_DOC.md` / `PLAN.md`.
