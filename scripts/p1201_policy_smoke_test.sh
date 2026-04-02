#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://127.0.0.1:18096}"
KM_LISTEN_ADDR="${KM_LISTEN_ADDR:-:18096}"
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
  echo "[p1201-policy-smoke] $*"
}

expect_status() {
  local actual="$1"
  local expected="$2"
  local msg="$3"
  if [[ "${actual}" != "${expected}" ]]; then
    echo "[p1201-policy-smoke][FAIL] ${msg} expected=${expected} actual=${actual}"
    exit 1
  fi
}

request() {
  local method="$1"
  local path="$2"
  local outfile="$3"
  curl -sS -X "${method}" -o "${outfile}" -w "%{http_code}" "${BASE_URL}${path}"
}

ensure_backend() {
  log "starting backend on ${KM_LISTEN_ADDR}"
  (
    cd "${ROOT_DIR}/backend"
    KM_LISTEN_ADDR="${KM_LISTEN_ADDR}" KM_K8S_ADAPTER_MODE=mock GOPROXY=https://goproxy.cn,direct GOCACHE=/tmp/go-build-cache go run ./cmd/server >/tmp/kubemanage-backend-p1201-policy-smoke.log 2>&1
  ) &
  BACKEND_PID="$!"

  for _ in {1..45}; do
    if curl -sS "${BASE_URL}/api/v1/healthz" >/dev/null 2>&1; then
      log "backend startup ready"
      return
    fi
    sleep 1
  done

  echo "[p1201-policy-smoke][FAIL] backend failed to start, see /tmp/kubemanage-backend-p1201-policy-smoke.log"
  exit 1
}

main() {
  ensure_backend

  local code body

  body="${TMP_DIR}/limitranges-list.json"
  code="$(request GET /api/v1/limitranges "${body}")"
  expect_status "${code}" "200" "list limitranges"

  body="${TMP_DIR}/limitrange-detail.json"
  code="$(request GET /api/v1/limitranges/compute-defaults "${body}")"
  expect_status "${code}" "200" "get limitrange detail"

  body="${TMP_DIR}/limitrange-yaml.yaml"
  code="$(request GET /api/v1/limitranges/compute-defaults/yaml "${body}")"
  expect_status "${code}" "200" "get limitrange yaml"
  if ! rg -q "kind:\\s*LimitRange" "${body}"; then
    echo "[p1201-policy-smoke][FAIL] limitrange yaml missing kind: LimitRange"
    exit 1
  fi

  body="${TMP_DIR}/resourcequotas-list.json"
  code="$(request GET /api/v1/resourcequotas "${body}")"
  expect_status "${code}" "200" "list resourcequotas"

  body="${TMP_DIR}/resourcequota-detail.json"
  code="$(request GET /api/v1/resourcequotas/compute-quota "${body}")"
  expect_status "${code}" "200" "get resourcequota detail"

  body="${TMP_DIR}/resourcequota-yaml.yaml"
  code="$(request GET /api/v1/resourcequotas/compute-quota/yaml "${body}")"
  expect_status "${code}" "200" "get resourcequota yaml"
  if ! rg -q "kind:\\s*ResourceQuota" "${body}"; then
    echo "[p1201-policy-smoke][FAIL] resourcequota yaml missing kind: ResourceQuota"
    exit 1
  fi

  body="${TMP_DIR}/networkpolicies-list.json"
  code="$(request GET /api/v1/networkpolicies "${body}")"
  expect_status "${code}" "200" "list networkpolicies"

  body="${TMP_DIR}/networkpolicy-detail.json"
  code="$(request GET /api/v1/networkpolicies/default-deny-all "${body}")"
  expect_status "${code}" "200" "get networkpolicy detail"

  body="${TMP_DIR}/networkpolicy-yaml.yaml"
  code="$(request GET /api/v1/networkpolicies/default-deny-all/yaml "${body}")"
  expect_status "${code}" "200" "get networkpolicy yaml"
  if ! rg -q "kind:\\s*NetworkPolicy" "${body}"; then
    echo "[p1201-policy-smoke][FAIL] networkpolicy yaml missing kind: NetworkPolicy"
    exit 1
  fi

  body="${TMP_DIR}/networkpolicy-yaml-download.yaml"
  code="$(request GET /api/v1/networkpolicies/default-deny-all/yaml/download "${body}")"
  expect_status "${code}" "200" "download networkpolicy yaml"

  log "all p1201 policy smoke checks passed"
}

main "$@"
