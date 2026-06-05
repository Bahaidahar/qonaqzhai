#!/usr/bin/env bash
# Reset the local databases to a clean, jury-ready demo state:
#   - wipes all transactional data (bookings, chats, reviews, cards, threads)
#   - removes e2e test users (…@e2e.test)
#   - keeps only the curated catalog (one approved vendor per category)
#   - re-seeds the curated vendors + vendor1's profile if missing
#
# Run after rehearsals or e2e runs:  ./demo-reset.sh
set -euo pipefail
cd "$(dirname "$0")"

PG="psql -h localhost -p 5433 -U qonaqzhai"
export PGPASSWORD=qonaqzhai
G=${GATEWAY:-http://localhost:8080}

KEEP="'Demo Vendor — Asyl Mereke','Asyl Toi Mekeni','Dastarkhan Catering','DJ Astana Pro','Aspan Studio','Gul Decor','Tatti Cakes'"

echo "→ wiping transactional data + test users…"
$PG -d qonaqzhai_core >/dev/null <<SQL
BEGIN;
DELETE FROM bookings; DELETE FROM reviews; DELETE FROM notifications;
DELETE FROM chat_messages; DELETE FROM chats;
DELETE FROM services WHERE vendor_id NOT IN (SELECT id FROM vendors WHERE name IN ($KEEP));
DELETE FROM photos   WHERE vendor_id NOT IN (SELECT id FROM vendors WHERE name IN ($KEEP));
DELETE FROM vendors  WHERE name NOT IN ($KEEP);
UPDATE vendors SET rating_avg=0, rating_count=0;
COMMIT;
SQL
$PG -d qonaqzhai_auth    -tAc "DELETE FROM users WHERE email LIKE '%@e2e.test';" >/dev/null
$PG -d qonaqzhai_payment -tAc "DELETE FROM payments; DELETE FROM cards;" >/dev/null
$PG -d qonaqzhai_realtime -tAc "DELETE FROM thread_messages; DELETE FROM threads;" >/dev/null

# Re-seed curated vendors if the catalog is short (e.g. fresh DB).
count=$($PG -d qonaqzhai_core -tAc "SELECT count(*) FROM vendors;")
if [ "$count" -lt 7 ]; then
  echo "→ catalog has $count vendors, re-seeding curated set…"
  ./seed-vendors.sh demo >/dev/null 2>&1 || true
  # vendor1 demo profile (approved Venue)
  VT=$(curl -s -X POST "$G/api/login" -H 'content-type: application/json' \
    -d '{"email":"vendor1@demo.kz","password":"demo12345"}' | sed -nE 's/.*"token":"([^"]+)".*/\1/p')
  AT=$(curl -s -X POST "$G/api/login" -H 'content-type: application/json' \
    -d '{"email":"admin@qonaqzhai.kz","password":"admin12345"}' | sed -nE 's/.*"token":"([^"]+)".*/\1/p')
  if [ -n "$VT" ] && [ -n "$AT" ]; then
    VID=$(curl -s -X POST "$G/api/me/vendor" -H "authorization: Bearer $VT" -H 'content-type: application/json' \
      -d '{"name":"Demo Vendor — Asyl Mereke","category":"Venue","city":"Almaty","description":"Банкетный зал в центре Алматы на 300 гостей.","priceFrom":250000}' \
      | sed -nE 's/.*"id":"([^"]+)".*/\1/p')
    curl -s -o /dev/null -X POST "$G/api/me/vendor/services" -H "authorization: Bearer $VT" -H 'content-type: application/json' \
      -d '{"name":"Вечерний банкет","description":"Полный вечер, до 300 гостей","price":400000,"unit":"day"}'
    curl -s -o /dev/null -X PATCH "$G/api/admin/vendors/$VID/status" -H "authorization: Bearer $AT" -H 'content-type: application/json' \
      -d '{"status":"approved"}'
  fi
fi

echo "→ done. catalog:"
$PG -d qonaqzhai_core -tAc "SELECT category || ' — ' || name FROM vendors ORDER BY category;"
echo "bookings/chats: $($PG -d qonaqzhai_core -tAc "SELECT count(*) FROM bookings") / $($PG -d qonaqzhai_core -tAc "SELECT count(*) FROM chats")"
