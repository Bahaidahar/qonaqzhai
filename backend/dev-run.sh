#!/usr/bin/env bash
# Local dev launcher for all five Go services.
# Loads backend/.env, sets cross-service gRPC addresses, starts each
# service in the background with logs under /tmp/qz-logs.
set -euo pipefail
cd "$(dirname "$0")"

LOGDIR=/tmp/qz-logs
mkdir -p "$LOGDIR"

# Load .env line-by-line so unquoted values with spaces survive
if [ -f .env ]; then
  while IFS= read -r line || [ -n "$line" ]; do
    case "$line" in ''|\#*) continue;; esac
    key=${line%%=*}
    val=${line#*=}
    export "$key=$val"
  done < .env
fi

# Cross-service wiring
export AUTH_GRPC_ADDR="localhost:9081"
export CORE_GRPC_ADDR="localhost:9082"
export PAYMENT_GRPC_ADDR="localhost:9083"
export REALTIME_GRPC_ADDR="localhost:9084"

# Gateway upstreams (override buggy default that pointed realtime at :8083)
export AUTH_URL="http://localhost:8081"
export CORE_URL="http://localhost:8082"
export REALTIME_URL="http://localhost:8084"

# auth also reads SMTP_PASSWORD (our .env uses SMTP_PASS)
export SMTP_PASSWORD="${SMTP_PASS:-}"

start() {
  local name=$1 path=$2
  echo "starting $name -> $LOGDIR/$name.log"
  ( go run "$path" >"$LOGDIR/$name.log" 2>&1 & echo $! >"$LOGDIR/$name.pid" )
}

start auth     ./services/auth/cmd/auth
sleep 4   # auth gRPC must be up before dependents dial it
start core     ./services/core/cmd/core
start payment  ./services/payment/cmd/payment
start realtime ./services/realtime/cmd/realtime
sleep 3
start gateway  ./services/gateway/cmd/main.go

echo "all launched. tail: tail -f $LOGDIR/*.log"
