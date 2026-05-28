#!/usr/bin/env bash
# Seeds a fixed set of users + vendor profiles for the Maestro suite. Idempotent —
# rerun any time. Credentials match `.maestro/.env` / `helpers/login.yaml`.
#
# Why not let each flow signup through the UI? Maestro's text input has subtle
# focus races when the keyboard animates; chaining a REST signup before every
# flow worked but slowed the suite ~25s per case. Fixed accounts let every
# flow open straight at the login screen and tap submit.
set -euo pipefail

BACKEND="${MAESTRO_BACKEND:-http://localhost:8080}"

post() {
  local path=$1 body=$2 token=${3:-}
  local hdr="content-type: application/json"
  if [ -n "$token" ]; then hdr="$hdr"$'\n'"authorization: Bearer $token"; fi
  curl -sS -X POST "$BACKEND$path" \
    -H "content-type: application/json" \
    ${token:+-H "authorization: Bearer $token"} \
    -d "$body"
}

# Signup (200) or fall back to login (on 409). Echoes the JWT.
signup_or_login() {
  local email=$1 password=$2 name=$3 role=$4
  local body
  body=$(curl -sS -X POST "$BACKEND/api/signup" \
    -H "content-type: application/json" \
    -d "{\"email\":\"$email\",\"password\":\"$password\",\"name\":\"$name\",\"role\":\"$role\"}")
  if echo "$body" | grep -q '"token"'; then
    echo "$body" | python3 -c 'import sys, json; print(json.load(sys.stdin)["token"])'
    return
  fi
  curl -sS -X POST "$BACKEND/api/login" \
    -H "content-type: application/json" \
    -d "{\"email\":\"$email\",\"password\":\"$password\"}" \
    | python3 -c 'import sys, json; print(json.load(sys.stdin)["token"])'
}

# --- Demo accounts the login screen's "Demo accounts" card auto-fills ---
echo "› seeding customer1@demo.kz"
signup_or_login "customer1@demo.kz" "demo12345" "Demo Customer" customer >/dev/null
echo "› seeding vendor1@demo.kz"
DEMO_VEND_TOKEN=$(signup_or_login "vendor1@demo.kz" "demo12345" "Demo Vendor" vendor)
curl -sS -X POST "$BACKEND/api/me/vendor" \
  -H "content-type: application/json" \
  -H "authorization: Bearer $DEMO_VEND_TOKEN" \
  -d '{"name":"Rixos Almaty Ballroom","category":"Venue","city":"Almaty","description":"Demo vendor used by Maestro and screenshots","priceFrom":500000}' >/dev/null
DEMO_VENDOR_ID=$(curl -sS "$BACKEND/api/me/vendor" \
  -H "authorization: Bearer $DEMO_VEND_TOKEN" \
  | python3 -c 'import sys, json; print(json.load(sys.stdin)["id"])')

# --- Fixed customer ---
echo "› seeding maestro_customer@test.kz"
CUST_TOKEN=$(signup_or_login "maestro_customer@test.kz" "password123" "Maestro Customer" customer)
echo "  token ${CUST_TOKEN:0:24}…"

# --- Approved vendor with a public listing ---
echo "› seeding maestro_vendor@test.kz"
VEND_TOKEN=$(signup_or_login "maestro_vendor@test.kz" "password123" "Maestro Vendor" vendor)
curl -sS -X POST "$BACKEND/api/me/vendor" \
  -H "content-type: application/json" \
  -H "authorization: Bearer $VEND_TOKEN" \
  -d '{"name":"Maestro Studio","category":"Venue","city":"Almaty","description":"Maestro test fixture vendor","priceFrom":450000}' >/dev/null
VENDOR_ID=$(curl -sS "$BACKEND/api/me/vendor" \
  -H "authorization: Bearer $VEND_TOKEN" \
  | python3 -c 'import sys, json; print(json.load(sys.stdin)["id"])')
echo "  vendor.id=$VENDOR_ID"

ADMIN_TOKEN=$(curl -sS -X POST "$BACKEND/api/login" \
  -H "content-type: application/json" \
  -d '{"email":"admin@qonaqzhai.kz","password":"admin12345"}' \
  | python3 -c 'import sys, json; print(json.load(sys.stdin)["token"])')

curl -sS -X PATCH "$BACKEND/api/admin/vendors/$VENDOR_ID/status" \
  -H "content-type: application/json" \
  -H "authorization: Bearer $ADMIN_TOKEN" \
  -d '{"status":"approved"}' >/dev/null
echo "  approved"

curl -sS -X PATCH "$BACKEND/api/admin/vendors/$DEMO_VENDOR_ID/status" \
  -H "content-type: application/json" \
  -H "authorization: Bearer $ADMIN_TOKEN" \
  -d '{"status":"approved"}' >/dev/null
echo "  demo vendor approved"

# --- Pending booking from a throwaway customer against the demo vendor ---
# (07_vendor_accept_decline.yaml logs in as the demo vendor, so its inbox
# only sees bookings whose vendorId is the demo vendor's id.)
echo "› seeding a pending booking against the demo vendor"
SEED_CUST_EMAIL="maestro_seed_$(date +%s)@test.kz"
SEED_CUST_TOKEN=$(signup_or_login "$SEED_CUST_EMAIL" "password123" "Seed" customer)
EVENT_DATE=$(date -v+30d +%Y-%m-%d 2>/dev/null || date -d "+30 days" +%Y-%m-%d)
curl -sS -X POST "$BACKEND/api/bookings" \
  -H "content-type: application/json" \
  -H "authorization: Bearer $SEED_CUST_TOKEN" \
  -d "{\"vendorId\":\"$DEMO_VENDOR_ID\",\"eventDate\":\"$EVENT_DATE\",\"guestCount\":80,\"amount\":500000}" >/dev/null
echo "  booking date=$EVENT_DATE"

# Also stash one against the maestro_vendor so a customer-side booking list
# shows something in 08.
DEMO_CUST_TOKEN=$(curl -sS -X POST "$BACKEND/api/login" \
  -H "content-type: application/json" \
  -d '{"email":"customer1@demo.kz","password":"demo12345"}' \
  | python3 -c 'import sys, json; print(json.load(sys.stdin)["token"])')
curl -sS -X POST "$BACKEND/api/bookings" \
  -H "content-type: application/json" \
  -H "authorization: Bearer $DEMO_CUST_TOKEN" \
  -d "{\"vendorId\":\"$DEMO_VENDOR_ID\",\"eventDate\":\"$EVENT_DATE\",\"guestCount\":40,\"amount\":250000}" >/dev/null
echo "  demo customer booking date=$EVENT_DATE"

echo "✓ fixtures seeded"
