#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://127.0.0.1:18085}"
KM_LISTEN_ADDR="${KM_LISTEN_ADDR:-:18085}"
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
  echo "[p305-smoke] $*"
}

request() {
  local method="$1"
  local path="$2"
  local role="$3"
  local confirm="$4"
  local body="${5:-}"
  local outfile="$6"

  local confirm_header=()
  if [[ "${confirm}" == "yes" ]]; then
    confirm_header=(-H "X-Action-Confirm: CONFIRM")
  fi

  if [[ -n "${body}" ]]; then
    curl -sS -X "${method}" \
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
    echo "[p305-smoke][FAIL] ${msg} expected=${expected} actual=${actual}"
    exit 1
  fi
}

ensure_backend() {
  log "starting temporary backend process on ${KM_LISTEN_ADDR}"
  (
    cd "${ROOT_DIR}/backend"
    KM_LISTEN_ADDR="${KM_LISTEN_ADDR}" KM_K8S_ADAPTER_MODE=mock GOCACHE=/tmp/go-build-cache go run ./cmd/server >/tmp/kubemanage-backend-p305-smoke.log 2>&1
  ) &
  BACKEND_PID="$!"

  for _ in {1..30}; do
    if curl -sS "${BASE_URL}/api/v1/healthz" >/dev/null 2>&1; then
      log "backend startup ready"
      return
    fi
    sleep 1
  done

  echo "[p305-smoke][FAIL] backend failed to start, see /tmp/kubemanage-backend-p305-smoke.log"
  exit 1
}

main() {
  ensure_backend

  local body code

  body="${TMP_DIR}/logs.json"
  code="$(request GET "/api/v1/pods/web-api-7bf59f6f9c-abcde/logs?container=web-api&keyword=healthz&matchOnly=true" viewer no "" "${body}")"
  expect_status "${code}" "200" "get pod logs with container filter"
  if ! grep -q "healthz" "${body}"; then
    echo "[p305-smoke][FAIL] filtered pod logs missing expected content"
    exit 1
  fi

  body="${TMP_DIR}/caps.json"
  code="$(request GET "/api/v1/pods/web-api-7bf59f6f9c-abcde/terminal/capabilities" viewer no "" "${body}")"
  expect_status "${code}" "200" "get terminal capabilities"
  if ! grep -q '"containers"' "${body}"; then
    echo "[p305-smoke][FAIL] terminal capabilities missing containers field"
    exit 1
  fi

  body="${TMP_DIR}/session-deny.json"
  code="$(request POST "/api/v1/pods/web-api-7bf59f6f9c-abcde/terminal/sessions" admin no '{"container":"web-api"}' "${body}")"
  expect_status "${code}" "428" "create terminal session should require confirm header"

  body="${TMP_DIR}/session-ok.json"
  code="$(request POST "/api/v1/pods/web-api-7bf59f6f9c-abcde/terminal/sessions" admin yes '{"container":"web-api"}' "${body}")"
  expect_status "${code}" "501" "create terminal session placeholder"

  log "all smoke checks passed"
}

main "$@"
