#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://127.0.0.1:8080}"
KM_LISTEN_ADDR="${KM_LISTEN_ADDR:-:8080}"
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TMP_DIR="$(mktemp -d)"
BACKEND_PID=""

cleanup() {
  if [[ -n "${BACKEND_PID}" ]]; then
    kill "${BACKEND_PID}" >/dev/null 2>&1 || true
  fi
  rm -rf "${TMP_DIR}"
}
trap cleanup EXIT

log() {
  echo "[p203-smoke] $*"
}

request() {
  local method="$1"
  local path="$2"
  local role="$3"
  local outfile="$4"

  curl -sS -X "${method}" \
    -H "X-User: smoke-tester" \
    -H "X-User-Role: ${role}" \
    -o "${outfile}" \
    -w "%{http_code}" \
    "${BASE_URL}${path}"
}

expect_status() {
  local actual="$1"
  local expected="$2"
  local msg="$3"
  if [[ "${actual}" != "${expected}" ]]; then
    echo "[p203-smoke][FAIL] ${msg} expected=${expected} actual=${actual}"
    exit 1
  fi
}

ensure_backend() {
  if curl -sS "${BASE_URL}/api/v1/healthz" >/dev/null 2>&1; then
    log "backend already running"
    return
  fi

  log "backend not running, starting temporary backend process"
  (
    cd "${ROOT_DIR}/backend"
    KM_LISTEN_ADDR="${KM_LISTEN_ADDR}" GOCACHE=/tmp/go-build-cache go run ./cmd/server >/tmp/kubemanage-backend-p203-smoke.log 2>&1
  ) &
  BACKEND_PID="$!"

  for _ in {1..30}; do
    if curl -sS "${BASE_URL}/api/v1/healthz" >/dev/null 2>&1; then
      log "backend startup ready"
      return
    fi
    sleep 1
  done

  echo "[p203-smoke][FAIL] backend failed to start, see /tmp/kubemanage-backend-p203-smoke.log"
  exit 1
}

main() {
  ensure_backend

  local body code

  body="${TMP_DIR}/ingresses.json"
  code="$(request GET /api/v1/ingresses viewer "${body}")"
  expect_status "${code}" "200" "list ingresses"

  body="${TMP_DIR}/ingress-detail.json"
  code="$(request GET /api/v1/ingresses/web-api-ing viewer "${body}")"
  expect_status "${code}" "200" "get ingress detail"

  body="${TMP_DIR}/ingress-services.json"
  code="$(request GET /api/v1/ingresses/web-api-ing/services viewer "${body}")"
  expect_status "${code}" "200" "get ingress related services"

  if ! grep -q 'web-api-svc' "${body}"; then
    echo "[p203-smoke][FAIL] ingress relation missing web-api-svc"
    exit 1
  fi

  body="${TMP_DIR}/hpas.json"
  code="$(request GET /api/v1/hpas viewer "${body}")"
  expect_status "${code}" "200" "list hpas"

  body="${TMP_DIR}/hpa-detail.json"
  code="$(request GET /api/v1/hpas/web-api-hpa viewer "${body}")"
  expect_status "${code}" "200" "get hpa detail"

  body="${TMP_DIR}/hpa-target.json"
  code="$(request GET /api/v1/hpas/web-api-hpa/target viewer "${body}")"
  expect_status "${code}" "200" "get hpa target"

  if ! grep -q '"name":"web-api"' "${body}"; then
    echo "[p203-smoke][FAIL] hpa target missing web-api"
    exit 1
  fi

  log "all smoke checks passed"
}

main "$@"
