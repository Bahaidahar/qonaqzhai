# qonaqzhai-backend

Minimal Go API for MVP — auth + chat echo.

## Run

```bash
go run .
```

Default port `:8080`. Override with `ADDR=:9000`.

## Env

- `JWT_SECRET` — required in prod. Dev generates ephemeral.
- `CORS_ORIGIN` — defaults to `http://localhost:3000`.
- `ADDR` — defaults to `:8080`.

## Endpoints

| Method | Path | Auth | Body |
|--------|------|------|------|
| GET | `/api/health` | — | — |
| POST | `/api/signup` | — | `{ email, password, name? }` |
| POST | `/api/login` | — | `{ email, password }` |
| GET | `/api/me` | Bearer | — |
| POST | `/api/chat` | Bearer | `{ message }` |

Auth header: `Authorization: Bearer <jwt>`.

## Storage

In-memory map. Restart wipes users. SQLite/Postgres next iteration.

## Run tests

```bash
go test ./...           # full suite, ~7s
go test -cover ./...    # with coverage
go test -v ./...        # verbose
```

Covers: auth (validation, login, /me, suspension), RBAC (customer/vendor/admin), vendor upsert + photo upload/serve, booking lifecycle (pending→accept/decline/cancel), admin moderation + stats, chat with real vendor inject.

## Quick test

```bash
# signup
curl -s -X POST localhost:8080/api/signup \
  -H 'content-type: application/json' \
  -d '{"email":"a@b.kz","password":"password123","name":"Aigerim"}' | jq

# login
curl -s -X POST localhost:8080/api/login \
  -H 'content-type: application/json' \
  -d '{"email":"a@b.kz","password":"password123"}' | jq

# me (paste token)
TOKEN=...
curl -s localhost:8080/api/me -H "authorization: Bearer $TOKEN" | jq

# chat
curl -s -X POST localhost:8080/api/chat \
  -H "authorization: Bearer $TOKEN" \
  -H 'content-type: application/json' \
  -d '{"message":"give me a budget"}' | jq
```
