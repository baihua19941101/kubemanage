#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://127.0.0.1:18094}"
KM_LISTEN_ADDR="${KM_LISTEN_ADDR:-:18094}"
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
  echo "[p1001-node-smoke] $*"
}

expect_status() {
  local actual="$1"
  local expected="$2"
  local msg="$3"
  if [[ "${actual}" != "${expected}" ]]; then
    echo "[p1001-node-smoke][FAIL] ${msg} expected=${expected} actual=${actual}"
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
    KM_LISTEN_ADDR="${KM_LISTEN_ADDR}" KM_K8S_ADAPTER_MODE=mock GOCACHE=/tmp/go-build-cache go run ./cmd/server >/tmp/kubemanage-backend-p1001-node-smoke.log 2>&1
  ) &
  BACKEND_PID="$!"

  for _ in {1..45}; do
    if curl -sS "${BASE_URL}/api/v1/healthz" >/dev/null 2>&1; then
      log "backend startup ready"
      return
    fi
    sleep 1
  done

  echo "[p1001-node-smoke][FAIL] backend failed to start, see /tmp/kubemanage-backend-p1001-node-smoke.log"
  exit 1
}

main() {
  ensure_backend

  local code body node_name
  node_name="ip-10-10-1-21.ec2.internal"

  body="${TMP_DIR}/nodes-list.json"
  code="$(request GET /api/v1/nodes "${body}")"
  expect_status "${code}" "200" "list nodes"

  body="${TMP_DIR}/node-detail.json"
  code="$(request GET "/api/v1/nodes/${node_name}" "${body}")"
  expect_status "${code}" "200" "get node detail"

  body="${TMP_DIR}/node-yaml.yaml"
  code="$(request GET "/api/v1/nodes/${node_name}/yaml" "${body}")"
  expect_status "${code}" "200" "get node yaml"
  if ! rg -q "kind:\\s*Node" "${body}"; then
    echo "[p1001-node-smoke][FAIL] node yaml missing kind: Node"
    exit 1
  fi

  body="${TMP_DIR}/node-yaml-download.yaml"
  code="$(request GET "/api/v1/nodes/${node_name}/yaml/download" "${body}")"
  expect_status "${code}" "200" "download node yaml"
  if ! rg -q "metadata:" "${body}"; then
    echo "[p1001-node-smoke][FAIL] downloaded yaml missing metadata section"
    exit 1
  fi

  log "all p1001 node management smoke checks passed"
}

main "$@"

