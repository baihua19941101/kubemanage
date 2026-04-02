#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://127.0.0.1:18098}"
KM_LISTEN_ADDR="${KM_LISTEN_ADDR:-:18098}"
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TMP_DIR="$(mktemp -d)"
BACKEND_PID=""

cleanup() {
  if [[ -n "${BACKEND_PID}" ]]; then
    kill "${BACKEND_PID}" >/dev/null 2>&1 || true
    wait "${BACKEND_PID}" 2>/dev/null || true
  fi
  rm -rf "${TMP_DIR}"
}
trap cleanup EXIT

log() {
  echo "[p1203-policy-delete-smoke] $*"
}

expect_status() {
  local actual="$1"
  local expected="$2"
  local msg="$3"
  if [[ "${actual}" != "${expected}" ]]; then
    echo "[p1203-policy-delete-smoke][FAIL] ${msg} expected=${expected} actual=${actual}"
    exit 1
  fi
}

request_delete() {
  local path="$1"
  local role="$2"
  local confirm="${3:-yes}"
  local outfile="$4"

  local confirm_header=()
  if [[ "${confirm}" == "yes" ]]; then
    confirm_header=(-H "X-Action-Confirm: CONFIRM")
  fi

  curl -sS -X DELETE \
    -H "X-User: smoke-tester" \
    -H "X-User-Role: ${role}" \
    "${confirm_header[@]}" \
    -o "${outfile}" \
    -w "%{http_code}" \
    "${BASE_URL}${path}"
}

request_get() {
  local path="$1"
  local outfile="$2"
  curl -sS -X GET -o "${outfile}" -w "%{http_code}" "${BASE_URL}${path}"
}

ensure_backend() {
  log "starting backend on ${KM_LISTEN_ADDR}"
  (
    cd "${ROOT_DIR}/backend"
    KM_LISTEN_ADDR="${KM_LISTEN_ADDR}" KM_K8S_ADAPTER_MODE=mock GOPROXY=https://goproxy.cn,direct GOCACHE=/tmp/go-build-cache go run ./cmd/server >/tmp/kubemanage-backend-p1203-policy-delete-smoke.log 2>&1
  ) &
  BACKEND_PID="$!"

  for _ in {1..45}; do
    if curl -sS "${BASE_URL}/api/v1/healthz" >/dev/null 2>&1; then
      log "backend startup ready"
      return
    fi
    sleep 1
  done

  echo "[p1203-policy-delete-smoke][FAIL] backend failed to start, see /tmp/kubemanage-backend-p1203-policy-delete-smoke.log"
  exit 1
}

main() {
  ensure_backend

  local code body

  body="${TMP_DIR}/deny-default-limitrange.json"
  code="$(request_delete "/api/v1/limitranges/compute-defaults" "operator" "yes" "${body}")"
  expect_status "${code}" "403" "operator delete default limitrange should be forbidden"

  body="${TMP_DIR}/delete-dev-limitrange.json"
  code="$(request_delete "/api/v1/limitranges/dev-container-limits" "operator" "yes" "${body}")"
  expect_status "${code}" "204" "operator delete dev limitrange"

  body="${TMP_DIR}/verify-limitrange-deleted.json"
  code="$(request_get "/api/v1/limitranges/dev-container-limits" "${body}")"
  expect_status "${code}" "404" "deleted limitrange should not be found"

  body="${TMP_DIR}/delete-dev-resourcequota.json"
  code="$(request_delete "/api/v1/resourcequotas/dev-quota" "operator" "yes" "${body}")"
  expect_status "${code}" "204" "operator delete dev resourcequota"

  body="${TMP_DIR}/delete-dev-networkpolicy.json"
  code="$(request_delete "/api/v1/networkpolicies/allow-web-to-api" "operator" "yes" "${body}")"
  expect_status "${code}" "204" "operator delete dev networkpolicy"

  body="${TMP_DIR}/missing-confirm.json"
  code="$(request_delete "/api/v1/networkpolicies/default-deny-all" "admin" "no" "${body}")"
  expect_status "${code}" "428" "delete without confirm header should fail"

  log "all p1203 policy delete smoke checks passed"
}

main "$@"
