#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://127.0.0.1:18097}"
KM_LISTEN_ADDR="${KM_LISTEN_ADDR:-:18097}"
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
  echo "[p1202-policy-write-smoke] $*"
}

expect_status() {
  local actual="$1"
  local expected="$2"
  local msg="$3"
  if [[ "${actual}" != "${expected}" ]]; then
    echo "[p1202-policy-write-smoke][FAIL] ${msg} expected=${expected} actual=${actual}"
    exit 1
  fi
}

request_yaml_update() {
  local path="$1"
  local body="$2"
  local role="$3"
  local confirm="${4:-yes}"
  local outfile="$5"

  local confirm_header=()
  if [[ "${confirm}" == "yes" ]]; then
    confirm_header=(-H "X-Action-Confirm: CONFIRM")
  fi

  curl -sS -X PUT \
    -H "Content-Type: application/json" \
    -H "X-User: smoke-tester" \
    -H "X-User-Role: ${role}" \
    "${confirm_header[@]}" \
    -d "${body}" \
    -o "${outfile}" \
    -w "%{http_code}" \
    "${BASE_URL}${path}"
}

ensure_backend() {
  log "starting backend on ${KM_LISTEN_ADDR}"
  (
    cd "${ROOT_DIR}/backend"
    KM_LISTEN_ADDR="${KM_LISTEN_ADDR}" KM_K8S_ADAPTER_MODE=mock GOPROXY=https://goproxy.cn,direct GOCACHE=/tmp/go-build-cache go run ./cmd/server >/tmp/kubemanage-backend-p1202-policy-write-smoke.log 2>&1
  ) &
  BACKEND_PID="$!"

  for _ in {1..45}; do
    if curl -sS "${BASE_URL}/api/v1/healthz" >/dev/null 2>&1; then
      log "backend startup ready"
      return
    fi
    sleep 1
  done

  echo "[p1202-policy-write-smoke][FAIL] backend failed to start, see /tmp/kubemanage-backend-p1202-policy-write-smoke.log"
  exit 1
}

main() {
  ensure_backend

  local code body

  body="${TMP_DIR}/limitrange-dev-update.json"
  code="$(request_yaml_update "/api/v1/limitranges/dev-container-limits/yaml" '{"yaml":"apiVersion: v1\nkind: LimitRange\nmetadata:\n  name: dev-container-limits\n  namespace: dev\n"}' "operator" "yes" "${body}")"
  expect_status "${code}" "204" "operator update dev limitrange"

  body="${TMP_DIR}/resourcequota-dev-update.json"
  code="$(request_yaml_update "/api/v1/resourcequotas/dev-quota/yaml" '{"yaml":"apiVersion: v1\nkind: ResourceQuota\nmetadata:\n  name: dev-quota\n  namespace: dev\n"}' "operator" "yes" "${body}")"
  expect_status "${code}" "204" "operator update dev resourcequota"

  body="${TMP_DIR}/networkpolicy-dev-update.json"
  code="$(request_yaml_update "/api/v1/networkpolicies/allow-web-to-api/yaml" '{"yaml":"apiVersion: networking.k8s.io/v1\nkind: NetworkPolicy\nmetadata:\n  name: allow-web-to-api\n  namespace: dev\n"}' "operator" "yes" "${body}")"
  expect_status "${code}" "204" "operator update dev networkpolicy"

  body="${TMP_DIR}/limitrange-default-forbidden.json"
  code="$(request_yaml_update "/api/v1/limitranges/compute-defaults/yaml" '{"yaml":"apiVersion: v1\nkind: LimitRange\nmetadata:\n  name: compute-defaults\n  namespace: default\n"}' "operator" "yes" "${body}")"
  expect_status "${code}" "403" "operator update default limitrange should be forbidden"

  body="${TMP_DIR}/missing-confirm.json"
  code="$(request_yaml_update "/api/v1/networkpolicies/allow-web-to-api/yaml" '{"yaml":"apiVersion: networking.k8s.io/v1\nkind: NetworkPolicy\nmetadata:\n  name: allow-web-to-api\n  namespace: dev\n"}' "admin" "no" "${body}")"
  expect_status "${code}" "428" "yaml update without confirm header should fail"

  log "all p1202 policy write smoke checks passed"
}

main "$@"
