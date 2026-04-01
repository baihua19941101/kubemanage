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
  echo "[p301-smoke] $*"
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
    echo "[p301-smoke][FAIL] ${msg} expected=${expected} actual=${actual}"
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
    KM_LISTEN_ADDR="${KM_LISTEN_ADDR}" KM_K8S_ADAPTER_MODE=mock GOCACHE=/tmp/go-build-cache go run ./cmd/server >/tmp/kubemanage-backend-p301-smoke.log 2>&1
  ) &
  BACKEND_PID="$!"

  for _ in {1..30}; do
    if curl -sS "${BASE_URL}/api/v1/healthz" >/dev/null 2>&1; then
      log "backend startup ready"
      return
    fi
    sleep 1
  done

  echo "[p301-smoke][FAIL] backend failed to start, see /tmp/kubemanage-backend-p301-smoke.log"
  exit 1
}

main() {
  ensure_backend

  local body code kubeconfig_id token_id suffix
  suffix="$(date +%H%M%S)"

  body="${TMP_DIR}/import-kubeconfig.json"
  code="$(request POST /api/v1/clusters/connections/import/kubeconfig admin "{\"name\":\"demo-kcfg-${suffix}\",\"kubeconfigContent\":\"apiVersion: v1\\nkind: Config\\nclusters:\\n- cluster:\\n    server: https://demo.example.local\\n  name: demo\\ncontexts:\\n- context:\\n    cluster: demo\\n    user: demo\\n  name: demo\\ncurrent-context: demo\\nusers:\\n- name: demo\\n  user:\\n    token: demo-token\\n\"}" "${body}")"
  expect_status "${code}" "201" "import kubeconfig"
  kubeconfig_id="$(grep -o '"id":[0-9]*' "${body}" | head -n1 | cut -d: -f2)"
  if [[ -z "${kubeconfig_id}" ]]; then
    echo "[p301-smoke][FAIL] failed to parse kubeconfig connection id"
    exit 1
  fi

  body="${TMP_DIR}/import-token.json"
  code="$(request POST /api/v1/clusters/connections/import/token admin "{\"name\":\"demo-token-${suffix}\",\"apiServer\":\"https://token.example.local\",\"bearerToken\":\"token-123\",\"caCert\":\"ca-data\",\"skipTlsVerify\":true}" "${body}")"
  expect_status "${code}" "201" "import token cluster"
  token_id="$(grep -o '"id":[0-9]*' "${body}" | head -n1 | cut -d: -f2)"
  if [[ -z "${token_id}" ]]; then
    echo "[p301-smoke][FAIL] failed to parse token connection id"
    exit 1
  fi

  body="${TMP_DIR}/test-connection.json"
  code="$(request POST /api/v1/clusters/connections/test admin '{"mode":"token","apiServer":"https://token.example.local","bearerToken":"token-123","caCert":"ca-data","skipTlsVerify":true}' "${body}")"
  expect_status "${code}" "200" "test connection"
  if ! grep -q 'connection ok' "${body}"; then
    echo "[p301-smoke][FAIL] connection test result missing success message"
    exit 1
  fi

  body="${TMP_DIR}/connections.json"
  code="$(request GET /api/v1/clusters/connections admin '' "${body}")"
  expect_status "${code}" "200" "list connections"
  if ! grep -q "demo-kcfg-${suffix}" "${body}"; then
    echo "[p301-smoke][FAIL] cluster connections missing imported kubeconfig connection"
    exit 1
  fi

  body="${TMP_DIR}/activate.json"
  code="$(request POST "/api/v1/clusters/connections/${kubeconfig_id}/activate" admin '' "${body}")"
  expect_status "${code}" "204" "activate first connection"

  body="${TMP_DIR}/live-cluster.json"
  code="$(request GET /api/v1/clusters/live viewer '' "${body}")"
  expect_status "${code}" "200" "get live cluster"
  if ! grep -q "demo-kcfg-${suffix}" "${body}"; then
    echo "[p301-smoke][FAIL] live cluster missing activated connection name"
    exit 1
  fi

  body="${TMP_DIR}/live-namespaces.json"
  code="$(request GET /api/v1/namespaces/live viewer '' "${body}")"
  expect_status "${code}" "200" "get live namespaces"
  if ! grep -q 'default' "${body}"; then
    echo "[p301-smoke][FAIL] live namespaces missing default"
    exit 1
  fi

  body="${TMP_DIR}/deny.json"
  code="$(request POST /api/v1/clusters/connections/import/token operator '{"name":"deny","apiServer":"https://deny.example.local","bearerToken":"token"}' "${body}")"
  expect_status "${code}" "403" "operator import token should be forbidden"

  log "all smoke checks passed"
}

main "$@"
