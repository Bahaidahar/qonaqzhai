# Deploy

Single docker-compose file spins up the whole microservice stack: 4
Postgres instances + 5 Go services. Only the gateway is published on the
host.

## Quickstart

```bash
cp deploy/.env.example deploy/.env
$EDITOR deploy/.env        # set JWT_SECRET at minimum
docker compose -f deploy/docker-compose.yml --env-file deploy/.env up --build
```

Gateway listens on `localhost:8080`. Point the mobile app + frontend at
it (`API_BASE=http://localhost:8080`).

## Layout

| Container    | Port (internal)       | Exposed | Notes                          |
|--------------|-----------------------|---------|--------------------------------|
| gateway      | 8080                  | yes     | Public reverse proxy           |
| auth         | 8081 HTTP / 9081 gRPC | no      | Owns JWT + users               |
| core         | 8082 HTTP / 9082 gRPC | no      | Vendors / bookings / reviews   |
| payment      | 8083 HTTP / 9083 gRPC | no      | Cards / payments / PayBox      |
| realtime     | 8084 HTTP / 9084 gRPC | no      | DM threads + WebSocket hub     |
| auth-db      | 5432                  | no      | Dedicated Postgres for auth    |
| core-db      | 5432                  | no      | Dedicated Postgres for core    |
| payment-db   | 5432                  | no      | Dedicated Postgres for payment |
| realtime-db  | 5432                  | no      | Dedicated Postgres for realtime|

## Building one service image

`service.Dockerfile` is parametrised:

```bash
docker build -f deploy/service.Dockerfile --build-arg SERVICE=core -t qonaqzhai-core .
```

The build context must be the repo root (`.`) so the entire `backend/`
workspace is available.

## What's missing on purpose

- TLS termination: terminate at your edge (Cloudflare, nginx, ALB).
- Observability: add OTLP collector + a Prometheus scrape config if you
  want metrics. The Go services emit slog JSON to stdout already.
- Secrets manager: `.env` is convenient for diploma + dev. Replace with
  Doppler / Vault / k8s Secrets for production.
