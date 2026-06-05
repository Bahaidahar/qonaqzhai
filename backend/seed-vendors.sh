#!/usr/bin/env bash
# Seed approved demo vendors (one per catalog category) with a couple of
# services each. Idempotent-ish: re-running creates fresh emails by suffix.
set -euo pipefail
G=${GATEWAY:-http://localhost:8080}
SUFFIX=${1:-demo}

jqv() { sed -nE "s/.*\"$1\":\"?([^\",}]+)\"?.*/\1/p" | head -1; }

admin_token() {
  curl -s -X POST "$G/api/login" -H 'content-type: application/json' \
    -d '{"email":"admin@qonaqzhai.kz","password":"admin12345"}' | jqv token
}
ADMIN=$(admin_token)
[ -n "$ADMIN" ] || { echo "admin login failed"; exit 1; }

# vendor rows: email|name|category|city|priceFrom|description
VENDORS=(
"venue_$SUFFIX@demo.kz|Asyl Toi Mekeni|Venue|Almaty|250000|Banquet hall for 300 guests in central Almaty."
"catering_$SUFFIX@demo.kz|Dastarkhan Catering|Catering|Almaty|4000|National & European cuisine, full-service catering."
"music_$SUFFIX@demo.kz|DJ Almaty Pro|Music & DJ|Almaty|150000|Wedding & corporate DJ with light show."
"photo_$SUFFIX@demo.kz|Aspan Studio|Photo & Video|Almaty|80000|Photo and cinematic video for events."
"decor_$SUFFIX@demo.kz|Gul Decor|Decor & Florists|Almaty|120000|Floral arches, stage and table decor."
"cakes_$SUFFIX@demo.kz|Tatti Cakes|Cakes|Almaty|15000|Custom wedding cakes and dessert tables."
)

# bash 3.2 (macOS) has no associative arrays — services resolved via case
services_for() {
  case "$1" in
    "Venue")            printf '%s\n' 'Evening hall rental|400000|day|Full evening, up to 300 guests' 'Day package|250000|day|Daytime celebration package' ;;
    "Catering")         printf '%s\n' 'Banquet menu|8000|person|3-course banquet per guest' 'Coffee break|2500|person|Tea, coffee and pastries' ;;
    "Music & DJ")       printf '%s\n' 'DJ set 5h|180000|fixed|5 hours with equipment' 'MC + DJ|320000|fixed|Host plus DJ for the night' ;;
    "Photo & Video")    printf '%s\n' 'Photo day|120000|day|Full-day photo coverage' 'Photo + video|260000|day|Photo and edited video' ;;
    "Decor & Florists") printf '%s\n' 'Stage decor|200000|fixed|Backdrop and floral arch' 'Table decor|6000|item|Per-table floral setup' ;;
    "Cakes")            printf '%s\n' '3-tier cake|45000|item|Three-tier custom cake' 'Dessert table|90000|fixed|Assorted desserts for 100' ;;
  esac
}

for row in "${VENDORS[@]}"; do
  IFS='|' read -r email name category city price desc <<< "$row"
  echo "== $name ($category) =="

  # signup (or login if exists)
  resp=$(curl -s -X POST "$G/api/signup" -H 'content-type: application/json' \
    -d "{\"email\":\"$email\",\"password\":\"demo12345\",\"name\":\"$name\",\"role\":\"vendor\"}")
  tok=$(echo "$resp" | jqv token)
  if [ -z "$tok" ]; then
    tok=$(curl -s -X POST "$G/api/login" -H 'content-type: application/json' \
      -d "{\"email\":\"$email\",\"password\":\"demo12345\"}" | jqv token)
  fi
  [ -n "$tok" ] || { echo "  no token, skip"; continue; }

  # upsert vendor profile
  vresp=$(curl -s -X POST "$G/api/me/vendor" -H "authorization: Bearer $tok" \
    -H 'content-type: application/json' \
    -d "{\"name\":\"$name\",\"category\":\"$category\",\"city\":\"$city\",\"description\":\"$desc\",\"priceFrom\":$price}")
  vid=$(echo "$vresp" | jqv id)
  [ -n "$vid" ] || { echo "  vendor upsert failed: $vresp"; continue; }
  echo "  vendor id $vid"

  # add services
  while IFS='|' read -r sname sprice sunit sdesc; do
    [ -z "$sname" ] && continue
    curl -s -o /dev/null -X POST "$G/api/me/vendor/services" -H "authorization: Bearer $tok" \
      -H 'content-type: application/json' \
      -d "{\"name\":\"$sname\",\"description\":\"$sdesc\",\"price\":$sprice,\"unit\":\"$sunit\"}"
  done <<< "$(services_for "$category")"

  # approve as admin
  code=$(curl -s -o /dev/null -w '%{http_code}' -X PATCH "$G/api/admin/vendors/$vid/status" \
    -H "authorization: Bearer $ADMIN" -H 'content-type: application/json' \
    -d '{"status":"approved"}')
  echo "  approved -> $code"
done

echo "== done. catalog count: =="
curl -s "$G/api/vendors?status=" | sed -nE 's/.*"total":([0-9]+).*/total=\1/p'
