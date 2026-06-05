#!/usr/bin/env bash
# Completely remove the qonaqzhai deployment from the server — containers,
# volumes, locally-built images, network, and the uploaded source. Other
# projects on the host (e.g. rest-app) are untouched.
set -uo pipefail
cd "$(dirname "$0")"

echo "→ stopping + removing containers, volumes, images…"
docker-compose -p qonaqzhai down -v --rmi local 2>/dev/null || true

# Belt-and-suspenders: drop anything left behind by name.
docker rm -f $(docker ps -aq --filter name=qz-) 2>/dev/null || true
docker rmi qonaqzhai-backend qonaqzhai-web 2>/dev/null || true
docker network rm qonaqzhai-net 2>/dev/null || true
docker volume rm qonaqzhai-pgdata 2>/dev/null || true

echo "→ removing uploaded source…"
rm -rf "$HOME/qonaqzhai-tmp"

echo "→ done. remaining containers:"
docker ps --format "  {{.Names}} ({{.Status}})"
