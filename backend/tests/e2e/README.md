# Microservice E2E

Brings up the full stack (4 Postgres + auth + gateway minimum) and exercises
critical flows through the gateway. Requires Docker.

```bash
cd backend
go test -tags=e2e ./tests/e2e -v
```

The harness in `harness_test.go` is intentionally minimal — it covers the
smoke path (signup → login → /api/me) end-to-end. Extend with additional
tests as the suite grows; helpers `dsnFor`, `startBin`, `waitFor`, and
`postJSON` are reusable.

Add `core`, `payment`, `realtime` via `startBin(...)` when a test exercises
their flows. Each service needs `AUTH_GRPC_ADDR` plus its own DB DSN and
HTTP/gRPC ports.
