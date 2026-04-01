#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BASE_PORT="${BASE_PORT:-18100}"

log() {
  echo "[p306-qa] $*"
}

run_step() {
  local name="$1"
  shift
  log "run ${name}"
  "$@"
}

main() {
  run_step "backend tests" bash -lc "cd '${ROOT_DIR}/backend' && GOCACHE=/tmp/go-build-cache go test ./..."
  run_step "frontend build" bash -lc "cd '${ROOT_DIR}/frontend' && npm run build"

  run_step "p304 smoke" bash -lc "cd '${ROOT_DIR}' && BASE_URL=http://127.0.0.1:$((BASE_PORT + 4)) KM_LISTEN_ADDR=:$((BASE_PORT + 4)) ./scripts/p304_smoke_test.sh"
  run_step "p305 smoke" bash -lc "cd '${ROOT_DIR}' && BASE_URL=http://127.0.0.1:$((BASE_PORT + 5)) KM_LISTEN_ADDR=:$((BASE_PORT + 5)) ./scripts/p305_smoke_test.sh"

  log "stage3 qa passed"
}

main "$@"
