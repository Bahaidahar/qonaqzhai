#!/usr/bin/env bash
# Run the whole stack (Postgres + 5 Go services + Next.js web) locally in
# Docker. Same compose file as the server, but the public URLs point at
# localhost. Secrets come from deploy/.env.
#
#   ./local.sh            # build + start in background
#   ./local.sh logs       # follow logs
#   ./local.sh down       # stop + remove (keeps the DB volume)
#   ./local.sh nuke       # stop + remove everything incl. DB volume
#
# Web:     http://localhost:3001
# Gateway: http://localhost:8090
set -euo pipefail
cd "$(dirname "$0")"

# docker compose v2 plugin or the standalone binary
DC="docker compose"
if ! docker compose version >/dev/null 2>&1; then DC="docker-compose"; fi

PROJECT=qonaqzhai-local
export PUBLIC_API_URL="http://localhost:8090"
export WEB_ORIGIN="http://localhost:3001"

case "${1:-up}" in
  up|"")  $DC -p "$PROJECT" --env-file .env up -d --build
          echo "✓ web  → http://localhost:3001"
          echo "✓ api  → http://localhost:8090"
          echo "  demo: customer1@demo.kz / vendor1@demo.kz / admin@qonaqzhai.kz (pw demo12345 / admin12345)" ;;
  logs)   $DC -p "$PROJECT" logs -f ;;
  ps)     $DC -p "$PROJECT" ps ;;
  down)   $DC -p "$PROJECT" down ;;
  nuke)   $DC -p "$PROJECT" down -v --rmi local ;;
  *)      echo "usage: ./local.sh [up|logs|ps|down|nuke]"; exit 1 ;;
esac
