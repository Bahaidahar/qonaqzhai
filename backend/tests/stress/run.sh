#!/usr/bin/env bash
# One-shot runner that fires every scenario sequentially and writes a JSON
# summary per run to reports/. The deck reads the latest summary into the
# stress-test slide.
set -euo pipefail

HERE="$(cd "$(dirname "$0")" && pwd)"
BACKEND="${BASE_URL:-http://localhost:8080}"
TS=$(date +%Y%m%d-%H%M%S)
OUT="$HERE/reports/$TS"
mkdir -p "$OUT"

echo "› backend: $BACKEND"
echo "› report: $OUT"
echo

if ! curl -sf -o /dev/null "$BACKEND/api/vendors?limit=1"; then
  echo "✗ gateway not responding on $BACKEND" >&2
  exit 1
fi

echo "› seeding demo fixtures"
bash "$HERE/../../../mobile/.maestro/scripts/seed_fixtures.sh" >/dev/null
echo "  fixtures ready"
echo

run() {
  local name=$1
  echo "===> $name"
  BASE_URL="$BACKEND" k6 run \
    --summary-export "$OUT/$name.json" \
    --no-color \
    "$HERE/scenarios/$name.js" \
    | tee "$OUT/$name.log" || true
  echo
}

run vendors_search
run login_throughput
run booking_create
run chat_burst
run mixed

echo "✓ all scenarios done — see $OUT"
