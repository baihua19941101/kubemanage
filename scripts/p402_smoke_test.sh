#!/usr/bin/env bash
set -euo pipefail

MOCK_BASE_URL="${MOCK_BASE_URL:-http://127.0.0.1:18086}"
LIVE_BASE_URL="${LIVE_BASE_URL:-http://127.0.0.1:18087}"
MOCK_LISTEN_ADDR="${MOCK_LISTEN_ADDR:-:18086}"
LIVE_LISTEN_ADDR="${LIVE_LISTEN_ADDR:-:18087}"
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TMP_DIR="$(mktemp -d)"
MOCK_PID=""
LIVE_PID=""

cleanup() {
  if [[ -n "${MOCK_PID}" ]]; then
    kill "${MOCK_PID}" >/dev/null 2>&1 || true
  fi
  if [[ -n "${LIVE_PID}" ]]; then
    kill "${LIVE_PID}" >/dev/null 2>&1 || true
  fi
  rm -rf "${TMP_DIR}"
}
trap cleanup EXIT

log() {
  echo "[p402-smoke] $*"
}

start_backend() {
  local mode="$1"
  local listen_addr="$2"
  local base_url="$3"
  local logfile="$4"
  log "starting backend mode=${mode} on ${listen_addr}"
  (
    cd "${ROOT_DIR}/backend"
    KM_LISTEN_ADDR="${listen_addr}" KM_K8S_ADAPTER_MODE="${mode}" GOCACHE=/tmp/go-build-cache go run ./cmd/server >"${logfile}" 2>&1
  ) &
  local pid="$!"
  if [[ "${mode}" == "mock" ]]; then
    MOCK_PID="${pid}"
  else
    LIVE_PID="${pid}"
  fi

  for _ in {1..30}; do
    if curl -sS "${base_url}/api/v1/healthz" >/dev/null 2>&1; then
      log "backend ready mode=${mode}"
      return
    fi
    sleep 1
  done

  echo "[p402-smoke][FAIL] backend failed to start mode=${mode}, see ${logfile}"
  exit 1
}

expect_status() {
  local actual="$1"
  local expected="$2"
  local msg="$3"
  if [[ "${actual}" != "${expected}" ]]; then
    echo "[p402-smoke][FAIL] ${msg} expected=${expected} actual=${actual}"
    exit 1
  fi
}

request_ws_http() {
  local base_url="$1"
  local path="$2"
  local outfile="$3"
  curl -sS -X GET \
    -H "X-User: smoke-tester" \
    -H "X-User-Role: admin" \
    -o "${outfile}" \
    -w "%{http_code}" \
    "${base_url}${path}"
}

main() {
  local body code

  start_backend "mock" "${MOCK_LISTEN_ADDR}" "${MOCK_BASE_URL}" "/tmp/kubemanage-backend-p402-smoke-mock.log"
  body="${TMP_DIR}/mock-ws-disabled.json"
  code="$(request_ws_http "${MOCK_BASE_URL}" "/api/v1/pods/web-api-7bf59f6f9c-abcde/terminal/ws?sessionId=any" "${body}")"
  expect_status "${code}" "501" "terminal ws should be disabled in mock mode"

  start_backend "live" "${LIVE_LISTEN_ADDR}" "${LIVE_BASE_URL}" "/tmp/kubemanage-backend-p402-smoke-live.log"
  body="${TMP_DIR}/live-session-required.json"
  code="$(request_ws_http "${LIVE_BASE_URL}" "/api/v1/pods/web-api-7bf59f6f9c-abcde/terminal/ws" "${body}")"
  expect_status "${code}" "400" "terminal ws should require sessionId"

  body="${TMP_DIR}/live-invalid-session.json"
  code="$(request_ws_http "${LIVE_BASE_URL}" "/api/v1/pods/web-api-7bf59f6f9c-abcde/terminal/ws?sessionId=invalid" "${body}")"
  expect_status "${code}" "401" "terminal ws should reject invalid session"

  log "all smoke checks passed"
}

main "$@"
