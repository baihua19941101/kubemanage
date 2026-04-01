#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://127.0.0.1:18082}"
KM_LISTEN_ADDR="${KM_LISTEN_ADDR:-:18082}"
KM_SECRET_KEY="${KM_SECRET_KEY:-p302-smoke-secret}"
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
  echo "[p302-smoke] $*"
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
    echo "[p302-smoke][FAIL] ${msg} expected=${expected} actual=${actual}"
    exit 1
  fi
}

ensure_backend() {
  log "starting temporary backend process on ${KM_LISTEN_ADDR}"
  (
    cd "${ROOT_DIR}/backend"
    KM_LISTEN_ADDR="${KM_LISTEN_ADDR}" KM_SECRET_KEY="${KM_SECRET_KEY}" KM_K8S_ADAPTER_MODE=mock GOCACHE=/tmp/go-build-cache go run ./cmd/server >/tmp/kubemanage-backend-p302-smoke.log 2>&1
  ) &
  BACKEND_PID="$!"

  for _ in {1..30}; do
    if curl -sS "${BASE_URL}/api/v1/healthz" >/dev/null 2>&1; then
      log "backend startup ready"
      return
    fi
    sleep 1
  done

  echo "[p302-smoke][FAIL] backend failed to start, see /tmp/kubemanage-backend-p302-smoke.log"
  exit 1
}

main() {
  ensure_backend

  local body code conn_id suffix
  local token_secret kubeconfig_secret
  suffix="$(date +%H%M%S)"
  token_secret="super-secret-token-p302-${suffix}"
  kubeconfig_secret="super-secret-kcfg-token-p302-${suffix}"

  body="${TMP_DIR}/import-token.json"
  code="$(request POST /api/v1/clusters/connections/import/token admin "{\"name\":\"p302-token-${suffix}\",\"apiServer\":\"https://token.example.local\",\"bearerToken\":\"${token_secret}\",\"caCert\":\"ca-data\",\"skipTlsVerify\":true}" "${body}")"
  expect_status "${code}" "201" "import token cluster"
  conn_id="$(grep -o '"id":[0-9]*' "${body}" | head -n1 | cut -d: -f2)"
  if [[ -z "${conn_id}" ]]; then
    echo "[p302-smoke][FAIL] failed to parse token connection id"
    exit 1
  fi

  body="${TMP_DIR}/import-kubeconfig.json"
  code="$(request POST /api/v1/clusters/connections/import/kubeconfig admin "{\"name\":\"p302-kcfg-${suffix}\",\"kubeconfigContent\":\"apiVersion: v1\\nkind: Config\\nclusters:\\n- cluster:\\n    server: https://demo.example.local\\n  name: demo\\ncontexts:\\n- context:\\n    cluster: demo\\n    user: demo\\n  name: demo\\ncurrent-context: demo\\nusers:\\n- name: demo\\n  user:\\n    token: ${kubeconfig_secret}\\n\"}" "${body}")"
  expect_status "${code}" "201" "import kubeconfig"

  body="${TMP_DIR}/connections.json"
  code="$(request GET /api/v1/clusters/connections admin '' "${body}")"
  expect_status "${code}" "200" "list connections"
  if grep -q "${token_secret}" "${body}"; then
    echo "[p302-smoke][FAIL] token secret leaked in list connections response"
    exit 1
  fi
  if grep -q "${kubeconfig_secret}" "${body}"; then
    echo "[p302-smoke][FAIL] kubeconfig secret leaked in list connections response"
    exit 1
  fi
  if ! grep -q '"hasBearerToken":true' "${body}"; then
    echo "[p302-smoke][FAIL] list connections missing hasBearerToken marker"
    exit 1
  fi
  if ! grep -q '"kubeconfigPreview":"\*\*\*"' "${body}"; then
    echo "[p302-smoke][FAIL] list connections missing kubeconfig masked preview"
    exit 1
  fi

  body="${TMP_DIR}/activate.json"
  code="$(request POST "/api/v1/clusters/connections/${conn_id}/activate" admin '' "${body}")"
  expect_status "${code}" "204" "activate cluster connection"

  body="${TMP_DIR}/live-cluster.json"
  code="$(request GET /api/v1/clusters/live viewer '' "${body}")"
  expect_status "${code}" "200" "get live cluster"

  body="${TMP_DIR}/live-namespaces.json"
  code="$(request GET /api/v1/namespaces/live viewer '' "${body}")"
  expect_status "${code}" "200" "get live namespaces"

  log "all smoke checks passed"
}

main "$@"
