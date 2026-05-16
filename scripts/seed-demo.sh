#!/usr/bin/env bash
# Seed the local backend with a rich demo dataset: customers, vendors across
# categories and cities, completed bookings, and reviews so ratings populate.
#
# Usage:
#   ./scripts/seed-demo.sh [API]
#
# API defaults to http://localhost:8080.

set -euo pipefail
API="${1:-http://localhost:8080}"

require() { command -v "$1" >/dev/null || { echo "missing $1"; exit 1; }; }
require curl
require jq

api() { curl -fsS "$@"; }

login() {
  local email="$1" pass="$2"
  api -X POST "$API/api/login" -H 'content-type: application/json' \
    -d "{\"email\":\"$email\",\"password\":\"$pass\"}" | jq -r '.token'
}

signup_or_login() {
  local email="$1" pass="$2" name="$3" role="$4"
  local res
  res=$(curl -s -X POST "$API/api/signup" -H 'content-type: application/json' \
    -d "{\"email\":\"$email\",\"password\":\"$pass\",\"name\":\"$name\",\"role\":\"$role\"}")
  if echo "$res" | jq -e .token >/dev/null 2>&1; then
    echo "$res" | jq -r .token
  else
    login "$email" "$pass"
  fi
}

upsert_vendor() {
  local tok="$1" name="$2" cat="$3" city="$4" price="$5" desc="$6"
  api -X POST "$API/api/vendor" -H "authorization: Bearer $tok" \
    -H 'content-type: application/json' \
    -d "$(jq -nc --arg n "$name" --arg c "$cat" --arg ct "$city" --arg d "$desc" --argjson p "$price" \
       '{name:$n,category:$c,city:$ct,priceFrom:$p,description:$d}')" \
    | jq -r .id
}

approve() {
  local atok="$1" vid="$2"
  api -X PATCH "$API/api/admin/vendors/$vid" -H "authorization: Bearer $atok" \
    -H 'content-type: application/json' -d '{"status":"approved"}' >/dev/null
}

create_booking() {
  local ctok="$1" vid="$2" date="$3" guests="$4" amount="$5"
  api -X POST "$API/api/bookings" -H "authorization: Bearer $ctok" \
    -H 'content-type: application/json' \
    -d "$(jq -nc --arg v "$vid" --arg d "$date" --argjson g "$guests" --argjson a "$amount" \
       '{vendorId:$v,eventDate:$d,guestCount:$g,amount:$a,note:"demo booking"}')" \
    | jq -r .id
}

vendor_transition() {
  local vtok="$1" bid="$2" status="$3"
  api -X PATCH "$API/api/bookings/$bid" -H "authorization: Bearer $vtok" \
    -H 'content-type: application/json' -d "{\"status\":\"$status\"}" >/dev/null
}

upload_photo() {
  local vtok="$1" seed="$2"
  local tmp
  tmp="$(mktemp -t qz_photo.XXXXXX)"
  if ! curl -fsSL "https://picsum.photos/seed/$seed/800/600" -o "$tmp" 2>/dev/null; then
    rm -f "$tmp"
    return 1
  fi
  curl -fsS -X POST "$API/api/vendor/photos" \
    -H "authorization: Bearer $vtok" \
    -F "photo=@${tmp};type=image/jpeg" >/dev/null || true
  rm -f "$tmp"
}

submit_review() {
  local ctok="$1" bid="$2" rating="$3" text="$4"
  api -X POST "$API/api/reviews" -H "authorization: Bearer $ctok" \
    -H 'content-type: application/json' \
    -d "$(jq -nc --arg b "$bid" --arg t "$text" --argjson r "$rating" \
       '{bookingId:$b,rating:$r,text:$t}')" >/dev/null
}

echo "==> Admin login"
ATOK=$(login "admin@qonaqzhai.kz" "admin12345")

echo "==> Customers"
CUSTOMER_TOKENS=()
for n in 1 2 3 4 5; do
  email="customer${n}@demo.kz"
  CUSTOMER_TOKENS+=("$(signup_or_login "$email" "demo12345" "Customer $n" "customer")")
done

echo "==> Vendors"
declare -a VENDORS_CATS=(Venue Venue Venue Catering Catering Photo Video Decor Music Cakes Photo Video)
declare -a VENDORS_CITIES=(Almaty Almaty Astana Almaty Shymkent Almaty Astana Almaty Shymkent Almaty Astana Almaty)
declare -a VENDORS_NAMES=(
  "Rixos Almaty Ballroom"
  "Esentai Ballroom"
  "St Regis Astana"
  "Aizada Catering"
  "Bayan Catering"
  "Studio Aitu Photo"
  "Daulet Video"
  "Almaty Floral Studio"
  "DJ Bahytzhan"
  "Cakes by Anel"
  "Tengri Photo"
  "Wedding Reels"
)
declare -a VENDORS_PRICES=(
  1500000
  1300000
  1800000
  1200000
  900000
  450000
  600000
  350000
  250000
  120000
  500000
  700000
)
declare -a VENDORS_DESCS=(
  "Premier ballroom in the heart of Almaty — capacity 400, stage, sound system included."
  "Elegant Esentai venue with rooftop terrace and curated catering partners."
  "Five-star ballroom in Astana with custom menus and floral arches."
  "Aigerim Saparbekova's catering — national dishes, plov, beshbarmak, vegetarian options."
  "Modern catering for corporate events and weddings, halal-certified."
  "Award-winning wedding and event photography studio."
  "Cinematic video, drone footage, same-day highlight reels."
  "Bespoke floral installations: ceremony arches, table runners, bridal bouquets."
  "Live DJ + MC, sound, light, smoke. Toi and corporate events."
  "Custom cake design, fondant, gluten-free options."
  "Documentary-style wedding photography across Kazakhstan."
  "Cinematic wedding films and editorial reels."
)

declare -a VENDOR_IDS
declare -a VENDOR_TOKENS
n=${#VENDORS_NAMES[@]}
for ((i=0; i<n; i++)); do
  email="vendor$((i+1))@demo.kz"
  echo "    - ${VENDORS_NAMES[$i]} (${VENDORS_CATS[$i]}, ${VENDORS_CITIES[$i]})"
  VTOK=$(signup_or_login "$email" "demo12345" "${VENDORS_NAMES[$i]}" "vendor")
  VID=$(upsert_vendor "$VTOK" "${VENDORS_NAMES[$i]}" "${VENDORS_CATS[$i]}" "${VENDORS_CITIES[$i]}" "${VENDORS_PRICES[$i]}" "${VENDORS_DESCS[$i]}")
  approve "$ATOK" "$VID"
  # 2 photos per vendor unless they already have some (idempotent re-run).
  existing=$(api "$API/api/vendors/$VID" | jq -r '.photoIds | length')
  if [ "$existing" -lt 2 ]; then
    upload_photo "$VTOK" "qz-${i}-a" || true
    upload_photo "$VTOK" "qz-${i}-b" || true
  fi
  VENDOR_IDS+=("$VID")
  VENDOR_TOKENS+=("$VTOK")
done

# Leave one vendor pending for admin moderation demo.
echo "==> Pending vendor (for moderation demo)"
PTOK=$(signup_or_login "vendor_pending@demo.kz" "demo12345" "Pending Co" "vendor")
upsert_vendor "$PTOK" "Aizhana Pending Co" "Decor" "Almaty" 200000 "Awaiting moderation." >/dev/null

echo "==> Bookings + reviews"
DATES=("2026-08-12" "2026-09-04" "2026-10-22" "2026-11-15" "2026-12-01")
RATINGS=(5 4 5 4 5 3 5 4)
TEXTS=(
  "Все прошло на высшем уровне!"
  "Отличная работа, рекомендую."
  "Спасибо, всё было прекрасно."
  "Хорошо, но были небольшие задержки."
  "Лучший опыт, спасибо!"
  "Нормально, есть что улучшить."
  "Топ команда, всё идеально."
  "Очень понравилось, ещё закажу."
)

# Create bookings: each customer books a few vendors; cycle through dates / amounts.
booking_counter=0
for c_idx in 0 1 2 3 4; do
  CTOK="${CUSTOMER_TOKENS[$c_idx]}"
  for ((j=0; j<3; j++)); do
    vendor_idx=$(( (c_idx * 3 + j) % n ))
    VID="${VENDOR_IDS[$vendor_idx]}"
    VTOK="${VENDOR_TOKENS[$vendor_idx]}"
    DATE=${DATES[$((booking_counter % ${#DATES[@]}))]}
    AMOUNT=${VENDORS_PRICES[$vendor_idx]}
    BID=$(create_booking "$CTOK" "$VID" "$DATE" "$((50 + booking_counter * 25))" "$AMOUNT")

    # Half of bookings go through full lifecycle with a review.
    if (( booking_counter % 2 == 0 )); then
      vendor_transition "$VTOK" "$BID" "accepted"
      vendor_transition "$VTOK" "$BID" "completed"
      RATING=${RATINGS[$((booking_counter % ${#RATINGS[@]}))]}
      TEXT=${TEXTS[$((booking_counter % ${#TEXTS[@]}))]}
      submit_review "$CTOK" "$BID" "$RATING" "$TEXT"
    elif (( booking_counter % 3 == 0 )); then
      vendor_transition "$VTOK" "$BID" "accepted"
    fi
    booking_counter=$((booking_counter+1))
  done
done

echo "==> Done. $booking_counter bookings, $n vendors approved, 1 pending."
echo
echo "Demo accounts:"
echo "  admin       admin@qonaqzhai.kz / admin12345"
echo "  customer1-5 customerN@demo.kz   / demo12345"
echo "  vendor1-12  vendorN@demo.kz     / demo12345"
echo "  pending     vendor_pending@demo.kz / demo12345"
