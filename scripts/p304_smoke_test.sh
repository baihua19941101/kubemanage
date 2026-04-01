#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://127.0.0.1:18084}"
KM_LISTEN_ADDR="${KM_LISTEN_ADDR:-:18084}"
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
  echo "[p304-smoke] $*"
}

request() {
  local method="$1"
  local path="$2"
  local role="$3"
  local confirm="$4"
  local body="${5:-}"
  local outfile="$6"
  local headerfile="$7"

  local confirm_header=()
  if [[ "${confirm}" == "yes" ]]; then
    confirm_header=(-H "X-Action-Confirm: CONFIRM")
  fi

  if [[ -n "${body}" ]]; then
    curl -sS -X "${method}" \
      -D "${headerfile}" \
      -H "X-User: smoke-tester" \
      -H "X-User-Role: ${role}" \
      "${confirm_header[@]}" \
      -H "Content-Type: application/json" \
      -d "${body}" \
      -o "${outfile}" \
      -w "%{http_code}" \
      "${BASE_URL}${path}"
  else
    curl -sS -X "${method}" \
      -D "${headerfile}" \
      -H "X-User: smoke-tester" \
      -H "X-User-Role: ${role}" \
      "${confirm_header[@]}" \
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
    echo "[p304-smoke][FAIL] ${msg} expected=${expected} actual=${actual}"
    exit 1
  fi
}

ensure_backend() {
  log "starting temporary backend process on ${KM_LISTEN_ADDR}"
  (
    cd "${ROOT_DIR}/backend"
    KM_LISTEN_ADDR="${KM_LISTEN_ADDR}" KM_K8S_ADAPTER_MODE=mock GOCACHE=/tmp/go-build-cache go run ./cmd/server >/tmp/kubemanage-backend-p304-smoke.log 2>&1
  ) &
  BACKEND_PID="$!"

  for _ in {1..30}; do
    if curl -sS "${BASE_URL}/api/v1/healthz" >/dev/null 2>&1; then
      log "backend startup ready"
      return
    fi
    sleep 1
  done

  echo "[p304-smoke][FAIL] backend failed to start, see /tmp/kubemanage-backend-p304-smoke.log"
  exit 1
}

main() {
  ensure_backend

  local body headers code suffix
  suffix="$(date +%H%M%S)"

  body="${TMP_DIR}/deny.json"
  headers="${TMP_DIR}/deny.headers"
  code="$(request POST /api/v1/namespaces admin no "{\"name\":\"p304-${suffix}\"}" "${body}" "${headers}")"
  expect_status "${code}" "428" "create namespace without confirm header"
  if ! grep -q "X-Request-Id:" "${headers}"; then
    echo "[p304-smoke][FAIL] missing X-Request-Id header in failed response"
    exit 1
  fi
  if ! grep -q '"code":"confirmation_required"' "${body}"; then
    echo "[p304-smoke][FAIL] missing confirmation_required code in response body"
    exit 1
  fi

  body="${TMP_DIR}/create.json"
  headers="${TMP_DIR}/create.headers"
  code="$(request POST /api/v1/namespaces admin yes "{\"name\":\"p304-${suffix}\"}" "${body}" "${headers}")"
  expect_status "${code}" "201" "create namespace with confirm header"

  body="${TMP_DIR}/audit.json"
  headers="${TMP_DIR}/audit.headers"
  code="$(request GET "/api/v1/audits?method=POST&path=/api/v1/namespaces&limit=5" admin no "" "${body}" "${headers}")"
  expect_status "${code}" "200" "list audits"
  if ! grep -q '"requestId"' "${body}"; then
    echo "[p304-smoke][FAIL] audit records missing requestId field"
    exit 1
  fi
  if ! grep -q '"error"' "${body}"; then
    echo "[p304-smoke][FAIL] audit records missing error field"
    exit 1
  fi

  log "all smoke checks passed"
}

main "$@"
