#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://127.0.0.1:18099}"
KM_LISTEN_ADDR="${KM_LISTEN_ADDR:-:18099}"
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
  echo "[p1204-policy-create-smoke] $*"
}

expect_status() {
  local actual="$1"
  local expected="$2"
  local msg="$3"
  if [[ "${actual}" != "${expected}" ]]; then
    echo "[p1204-policy-create-smoke][FAIL] ${msg} expected=${expected} actual=${actual}"
    exit 1
  fi
}

request_create() {
  local path="$1"
  local body="$2"
  local role="$3"
  local confirm="${4:-yes}"
  local outfile="$5"

  local confirm_header=()
  if [[ "${confirm}" == "yes" ]]; then
    confirm_header=(-H "X-Action-Confirm: CONFIRM")
  fi

  curl -sS -X POST \
    -H "Content-Type: application/json" \
    -H "X-User: smoke-tester" \
    -H "X-User-Role: ${role}" \
    "${confirm_header[@]}" \
    -d "${body}" \
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
    KM_LISTEN_ADDR="${KM_LISTEN_ADDR}" KM_K8S_ADAPTER_MODE=mock GOPROXY=https://goproxy.cn,direct GOCACHE=/tmp/go-build-cache go run ./cmd/server >/tmp/kubemanage-backend-p1204-policy-create-smoke.log 2>&1
  ) &
  BACKEND_PID="$!"

  for _ in {1..45}; do
    if curl -sS "${BASE_URL}/api/v1/healthz" >/dev/null 2>&1; then
      log "backend startup ready"
      return
    fi
    sleep 1
  done

  echo "[p1204-policy-create-smoke][FAIL] backend failed to start, see /tmp/kubemanage-backend-p1204-policy-create-smoke.log"
  exit 1
}

main() {
  ensure_backend

  local code body

  body="${TMP_DIR}/create-lr-default-deny.json"
  code="$(request_create "/api/v1/limitranges" '{"namespace":"default","yaml":"apiVersion: v1\nkind: LimitRange\nmetadata:\n  name: p1204-default-deny\n  namespace: default\n"}' "operator" "yes" "${body}")"
  expect_status "${code}" "403" "operator create default limitrange should be forbidden"

  body="${TMP_DIR}/create-lr-dev.json"
  code="$(request_create "/api/v1/limitranges" '{"namespace":"dev","yaml":"apiVersion: v1\nkind: LimitRange\nmetadata:\n  name: p1204-dev-lr\n  namespace: dev\n"}' "operator" "yes" "${body}")"
  expect_status "${code}" "201" "operator create dev limitrange"

  body="${TMP_DIR}/create-rq-dev.json"
  code="$(request_create "/api/v1/resourcequotas" '{"namespace":"dev","yaml":"apiVersion: v1\nkind: ResourceQuota\nmetadata:\n  name: p1204-dev-rq\n  namespace: dev\nspec:\n  hard:\n    pods: \"10\"\n"}' "operator" "yes" "${body}")"
  expect_status "${code}" "201" "operator create dev resourcequota"

  body="${TMP_DIR}/create-np-dev.json"
  code="$(request_create "/api/v1/networkpolicies" '{"namespace":"dev","yaml":"apiVersion: networking.k8s.io/v1\nkind: NetworkPolicy\nmetadata:\n  name: p1204-dev-np\n  namespace: dev\nspec:\n  podSelector: {}\n  policyTypes:\n  - Ingress\n"}' "operator" "yes" "${body}")"
  expect_status "${code}" "201" "operator create dev networkpolicy"

  body="${TMP_DIR}/create-without-confirm.json"
  code="$(request_create "/api/v1/networkpolicies" '{"namespace":"dev","yaml":"apiVersion: networking.k8s.io/v1\nkind: NetworkPolicy\nmetadata:\n  name: p1204-no-confirm\n  namespace: dev\n"}' "admin" "no" "${body}")"
  expect_status "${code}" "428" "create without confirm header should fail"

  body="${TMP_DIR}/verify-created.json"
  code="$(request_get "/api/v1/networkpolicies/p1204-dev-np" "${body}")"
  expect_status "${code}" "200" "created networkpolicy should be queryable"

  log "all p1204 policy create smoke checks passed"
}

main "$@"
