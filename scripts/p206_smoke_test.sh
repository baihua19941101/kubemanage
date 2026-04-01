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
  echo "[p206-smoke] $*"
}

request() {
  local method="$1"
  local path="$2"
  local role="$3"
  local body="${4:-}"
  local outfile="$5"

  if [[ -n "${body}" ]]; then
    curl -sS -X "${method}" \
      -H "X-User: smoke-tester" \
      -H "X-User-Role: ${role}" \
      -H "Content-Type: application/json" \
      -d "${body}" \
      -o "${outfile}" \
      -w "%{http_code}" \
      "${BASE_URL}${path}"
  else
    curl -sS -X "${method}" \
      -H "X-User: smoke-tester" \
      -H "X-User-Role: ${role}" \
      -o "${outfile}" \
      -w "%{http_code}" \
      "${BASE_URL}${path}"
  fi
}

expect_status() {
  local actual="$1"
  local expected="$2"
  local msg="$3"
  if [[ "${actual}" != "${expected}" ]]; then
    echo "[p206-smoke][FAIL] ${msg} expected=${expected} actual=${actual}"
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
    KM_LISTEN_ADDR="${KM_LISTEN_ADDR}" GOCACHE=/tmp/go-build-cache go run ./cmd/server >/tmp/kubemanage-backend-p206-smoke.log 2>&1
  ) &
  BACKEND_PID="$!"

  for _ in {1..30}; do
    if curl -sS "${BASE_URL}/api/v1/healthz" >/dev/null 2>&1; then
      log "backend startup ready"
      return
    fi
    sleep 1
  done

  echo "[p206-smoke][FAIL] backend failed to start, see /tmp/kubemanage-backend-p206-smoke.log"
  exit 1
}

main() {
  ensure_backend

  local body code

  body="${TMP_DIR}/operator-dev.json"
  code="$(request PUT /api/v1/pods/task-worker-856ddcf69f-uvwxy/yaml operator '{"yaml":"apiVersion: v1\nkind: Pod\nmetadata:\n  name: task-worker-856ddcf69f-uvwxy\n"}' "${body}")"
  expect_status "${code}" "204" "operator update pod in dev"

  body="${TMP_DIR}/operator-default.json"
  code="$(request PUT /api/v1/deployments/web-api/yaml operator '{"yaml":"apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: web-api\n"}' "${body}")"
  expect_status "${code}" "403" "operator update deployment in default should be forbidden"

  body="${TMP_DIR}/auth-me.json"
  code="$(request GET /api/v1/auth/me operator '' "${body}")"
  expect_status "${code}" "200" "operator auth me"
  if ! grep -q 'dev' "${body}"; then
    echo "[p206-smoke][FAIL] auth/me missing allowed namespace dev"
    exit 1
  fi

  body="${TMP_DIR}/audit-filter.json"
  code="$(request GET '/api/v1/audits?path=/api/v1/deployments/web-api/yaml&statusCode=403&limit=5' admin '' "${body}")"
  expect_status "${code}" "200" "audit filtered query"
  if ! grep -q '/api/v1/deployments/web-api/yaml' "${body}"; then
    echo "[p206-smoke][FAIL] audit filter missing deployment path"
    exit 1
  fi

  log "all smoke checks passed"
}

main "$@"
