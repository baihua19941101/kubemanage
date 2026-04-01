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
  echo "[p205-smoke] $*"
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
    echo "[p205-smoke][FAIL] ${msg} expected=${expected} actual=${actual}"
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
    KM_LISTEN_ADDR="${KM_LISTEN_ADDR}" GOCACHE=/tmp/go-build-cache go run ./cmd/server >/tmp/kubemanage-backend-p205-smoke.log 2>&1
  ) &
  BACKEND_PID="$!"

  for _ in {1..30}; do
    if curl -sS "${BASE_URL}/api/v1/healthz" >/dev/null 2>&1; then
      log "backend startup ready"
      return
    fi
    sleep 1
  done

  echo "[p205-smoke][FAIL] backend failed to start, see /tmp/kubemanage-backend-p205-smoke.log"
  exit 1
}

main() {
  ensure_backend

  local body code

  body="${TMP_DIR}/pod-logs.txt"
  code="$(request GET '/api/v1/pods/web-api-7bf59f6f9c-abcde/logs?keyword=healthz&matchOnly=true' viewer "${body}")"
  expect_status "${code}" "200" "get filtered pod logs"
  if ! grep -q 'GET /healthz' "${body}"; then
    echo "[p205-smoke][FAIL] filtered log output missing expected line"
    exit 1
  fi
  if grep -q 'server started' "${body}"; then
    echo "[p205-smoke][FAIL] matchOnly logs should not contain unmatched lines"
    exit 1
  fi

  body="${TMP_DIR}/pod-follow-logs.txt"
  code="$(request GET '/api/v1/pods/web-api-7bf59f6f9c-abcde/logs?follow=true' viewer "${body}")"
  expect_status "${code}" "200" "get follow pod logs"
  if ! grep -q 'follow refresh tick=' "${body}"; then
    echo "[p205-smoke][FAIL] follow log output missing refresh marker"
    exit 1
  fi

  body="${TMP_DIR}/terminal-capabilities.json"
  code="$(request GET /api/v1/pods/web-api-7bf59f6f9c-abcde/terminal/capabilities viewer "${body}")"
  expect_status "${code}" "200" "get terminal capabilities"
  if ! grep -q '"enabled":false' "${body}"; then
    echo "[p205-smoke][FAIL] terminal capabilities should report disabled"
    exit 1
  fi

  body="${TMP_DIR}/terminal-session.json"
  code="$(request POST /api/v1/pods/web-api-7bf59f6f9c-abcde/terminal/sessions operator "${body}")"
  expect_status "${code}" "501" "create terminal session placeholder"
  if ! grep -q 'terminal gateway not enabled' "${body}"; then
    echo "[p205-smoke][FAIL] terminal placeholder message missing"
    exit 1
  fi

  log "all smoke checks passed"
}

main "$@"
