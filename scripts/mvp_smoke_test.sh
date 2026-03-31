#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://127.0.0.1:8080}"
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
  echo "[mvp-smoke] $*"
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
    echo "[mvp-smoke][FAIL] ${msg} expected=${expected} actual=${actual}"
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
    go run ./cmd/server >/tmp/kubemanage-backend-smoke.log 2>&1
  ) &
  BACKEND_PID="$!"

  for _ in {1..30}; do
    if curl -sS "${BASE_URL}/api/v1/healthz" >/dev/null 2>&1; then
      log "backend startup ready"
      return
    fi
    sleep 1
  done

  echo "[mvp-smoke][FAIL] backend failed to start, see /tmp/kubemanage-backend-smoke.log"
  exit 1
}

main() {
  ensure_backend

  local body code ns_name
  ns_name="qa-smoke-$(date +%H%M%S)"

  body="${TMP_DIR}/health.json"
  code="$(request GET /api/v1/healthz viewer "" "${body}")"
  expect_status "${code}" "200" "healthz"

  body="${TMP_DIR}/clusters.json"
  code="$(request GET /api/v1/clusters viewer "" "${body}")"
  expect_status "${code}" "200" "cluster list"

  body="${TMP_DIR}/deny-switch.json"
  code="$(request POST /api/v1/clusters/switch viewer '{"name":"staging-cluster"}' "${body}")"
  expect_status "${code}" "403" "viewer switch cluster should be forbidden"

  body="${TMP_DIR}/switch.json"
  code="$(request POST /api/v1/clusters/switch operator '{"name":"staging-cluster"}' "${body}")"
  expect_status "${code}" "200" "operator switch cluster"

  body="${TMP_DIR}/ns-create.json"
  code="$(request POST /api/v1/namespaces operator "{\"name\":\"${ns_name}\"}" "${body}")"
  expect_status "${code}" "201" "create namespace"

  body="${TMP_DIR}/dep-update.txt"
  code="$(request PUT /api/v1/deployments/web-api/yaml operator '{"yaml":"apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: web-api\n"}' "${body}")"
  expect_status "${code}" "204" "update deployment yaml"

  body="${TMP_DIR}/pod-logs.txt"
  code="$(request GET /api/v1/pods/web-api-7bf59f6f9c-abcde/logs viewer "" "${body}")"
  expect_status "${code}" "200" "get pod logs"

  body="${TMP_DIR}/secret.json"
  code="$(request GET /api/v1/secrets viewer "" "${body}")"
  expect_status "${code}" "200" "list secrets"
  if ! grep -q '\*\*\*\*\*\*' "${body}"; then
    echo "[mvp-smoke][FAIL] secret list does not appear masked"
    exit 1
  fi

  body="${TMP_DIR}/audits.json"
  code="$(request GET /api/v1/audits admin "" "${body}")"
  expect_status "${code}" "200" "admin read audits"
  if ! grep -q '/api/v1/namespaces' "${body}"; then
    echo "[mvp-smoke][FAIL] audit log missing namespace operation"
    exit 1
  fi
  if ! grep -q '/api/v1/deployments/web-api/yaml' "${body}"; then
    echo "[mvp-smoke][FAIL] audit log missing deployment yaml update"
    exit 1
  fi

  body="${TMP_DIR}/ns-del.txt"
  code="$(request DELETE "/api/v1/namespaces/${ns_name}" operator "" "${body}")"
  expect_status "${code}" "204" "delete namespace"

  log "all smoke checks passed"
}

main "$@"
