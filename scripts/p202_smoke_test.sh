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
  echo "[p202-smoke] $*"
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
    echo "[p202-smoke][FAIL] ${msg} expected=${expected} actual=${actual}"
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
    KM_LISTEN_ADDR="${KM_LISTEN_ADDR}" GOCACHE=/tmp/go-build-cache go run ./cmd/server >/tmp/kubemanage-backend-p202-smoke.log 2>&1
  ) &
  BACKEND_PID="$!"

  for _ in {1..30}; do
    if curl -sS "${BASE_URL}/api/v1/healthz" >/dev/null 2>&1; then
      log "backend startup ready"
      return
    fi
    sleep 1
  done

  echo "[p202-smoke][FAIL] backend failed to start, see /tmp/kubemanage-backend-p202-smoke.log"
  exit 1
}

check_read_flow() {
  local kind="$1"
  local name="$2"
  local body code

  body="${TMP_DIR}/${kind}-list.json"
  code="$(request GET "/api/v1/${kind}" viewer "" "${body}")"
  expect_status "${code}" "200" "list ${kind}"

  body="${TMP_DIR}/${kind}-detail.json"
  code="$(request GET "/api/v1/${kind}/${name}" viewer "" "${body}")"
  expect_status "${code}" "200" "get ${kind}/${name}"

  body="${TMP_DIR}/${kind}-yaml.txt"
  code="$(request GET "/api/v1/${kind}/${name}/yaml" viewer "" "${body}")"
  expect_status "${code}" "200" "get ${kind}/${name} yaml"
}

check_write_flow() {
  local kind="$1"
  local name="$2"
  local payload="$3"
  local writer_role="$4"
  local body code

  body="${TMP_DIR}/${kind}-deny.json"
  code="$(request PUT "/api/v1/${kind}/${name}/yaml" viewer "${payload}" "${body}")"
  expect_status "${code}" "403" "viewer update ${kind}/${name} should be forbidden"

  body="${TMP_DIR}/${kind}-update.json"
  code="$(request PUT "/api/v1/${kind}/${name}/yaml" "${writer_role}" "${payload}" "${body}")"
  expect_status "${code}" "204" "${writer_role} update ${kind}/${name} yaml"
}

main() {
  ensure_backend

  check_read_flow "statefulsets" "mysql"
  check_read_flow "daemonsets" "node-exporter"
  check_read_flow "jobs" "db-migrate-20260401"
  check_read_flow "cronjobs" "cleanup"

  check_write_flow "statefulsets" "mysql" '{"yaml":"apiVersion: apps/v1\nkind: StatefulSet\nmetadata:\n  name: mysql\n"}' admin
  check_write_flow "daemonsets" "node-exporter" '{"yaml":"apiVersion: apps/v1\nkind: DaemonSet\nmetadata:\n  name: node-exporter\n"}' admin
  check_write_flow "jobs" "db-migrate-20260401" '{"yaml":"apiVersion: batch/v1\nkind: Job\nmetadata:\n  name: db-migrate-20260401\n"}' admin
  check_write_flow "cronjobs" "cleanup" '{"yaml":"apiVersion: batch/v1\nkind: CronJob\nmetadata:\n  name: cleanup\n"}' admin

  local body code
  body="${TMP_DIR}/audits.json"
  code="$(request GET /api/v1/audits admin "" "${body}")"
  expect_status "${code}" "200" "read audits"
  for p in \
    '/api/v1/statefulsets/mysql/yaml' \
    '/api/v1/daemonsets/node-exporter/yaml' \
    '/api/v1/jobs/db-migrate-20260401/yaml' \
    '/api/v1/cronjobs/cleanup/yaml'; do
    if ! grep -q "${p}" "${body}"; then
      echo "[p202-smoke][FAIL] audit missing path ${p}"
      exit 1
    fi
  done

  log "all smoke checks passed"
}

main "$@"
