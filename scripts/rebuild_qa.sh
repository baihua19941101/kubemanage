#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

log() {
  echo "[rebuild-qa] $*"
}

log "run backend unit tests"
(
  cd "${ROOT_DIR}/backend"
  go test ./...
)

log "run frontend production build"
(
  cd "${ROOT_DIR}/frontend"
  npm run build
)

log "run MVP smoke checks"
bash "${ROOT_DIR}/scripts/mvp_smoke_test.sh"

log "all checks passed"
